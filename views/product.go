package views

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"test/admin-server/models"
	"time"
)

type ProductListResponse struct {
	Status int `json:"status"`

	Data struct {
		PageNum  int               `json:"pageNum"`
		Total    int               `json:"total"`
		Pages    int               `json:"pages"`
		PageSize int               `json:"pageSize"`
		List     []*models.Product `json:"list"`
	} `json:"data"`
}

type RegProductResponse struct {
	Status int             `json:"status"`
	Data   *models.Product `json:"data"`
}

type ProductResponse struct {
	RegProductResponse
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]map[string]interface{}
	// 解析参数 存入map
	decoder.Decode(&params)
	data := params["product"]
	filter := bson.M{"name": data["name"].(string)}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("product")
	var product models.Product
	err := collection.FindOne(ctx, filter).Decode(&product)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Println("插入数据")
		product.Name = data["name"].(string)
		product.Detail = data["detail"].(string)
		product.Desc = data["desc"].(string)
		product.Price = data["price"].(string)
		product.CategoryId = data["categoryId"].(string)
		imgs := data["imgs"].([]interface{})
		for _, img := range imgs {
			product.Imgs = append(product.Imgs, img.(string))
		}

		_, err = collection.InsertOne(ctx, &product)
		if err != nil {
			fmt.Printf("插入数据失败：%v\n", err)
			errResponse.Msg = "插入数据失败"
		} else {
			collection.FindOne(ctx, filter).Decode(&product)
			var productResponse ProductResponse
			productResponse.Status = 0
			productResponse.Data = &product
			json.NewEncoder(w).Encode(productResponse)
			return
		}
	} else {
		fmt.Printf("插入数据成功")
		errResponse.Msg = "数据已存在"
	}
	json.NewEncoder(w).Encode(errResponse)
}

func GetProductList(w http.ResponseWriter, r *http.Request) {
	pageNum, err := strconv.Atoi(r.URL.Query().Get("pageNum"))
	if err != nil {
		pageNum = 0
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 0
	}
	searchName := r.URL.Query().Get("searchName")
	searchType := r.URL.Query().Get("searchType")
	filter := bson.M{}
	if searchType != "" {
		if searchType == "productName" {
			filter["name"] = bson.M{"$regex": searchName, "$options": "i"}
		} else if searchType == "productDesc" {
			filter["desc"] = bson.M{"$regex": searchName, "$options": "i"}
		}
	}

	findOptions := options.Find()
	if pageNum != 0 && pageSize != 0 {
		findOptions.SetSkip(int64((pageNum - 1) * pageSize))
		findOptions.SetLimit(int64(pageSize))
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("product")
	cur, err := collection.Find(ctx, filter, findOptions)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Printf("获取商品分页列表: %v\n", err)
		errResponse.Msg = "获取商品分页列表"
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	defer cur.Close(ctx)
	var products []*models.Product
	cur.All(ctx, &products)
	allCur, _ := collection.Find(ctx, filter)
	var AllProducts []*models.Product
	allCur.All(ctx, &AllProducts)
	defer allCur.Close(ctx)
	var productListResponse ProductListResponse
	productListResponse.Status = 0
	productListResponse.Data.List = products
	productListResponse.Data.PageNum = pageNum
	productListResponse.Data.PageSize = pageSize
	productListResponse.Data.Total = len(AllProducts)
	json.NewEncoder(w).Encode(productListResponse)

}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	pid := r.URL.Query().Get("productId")
	objId, _ := primitive.ObjectIDFromHex(pid)
	filter := bson.M{"_id": objId}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("product")
	var procudt models.Product
	err := collection.FindOne(ctx, filter).Decode(&procudt)
	if err != nil {
		var errResponse ErrResponse
		errResponse.Status = 1
		errResponse.Msg = "查询失败"
		fmt.Printf("查询失败：%v\n", err)
		json.NewEncoder(w).Encode(errResponse)
	} else {
		var productResponse ProductResponse
		productResponse.Status = 0
		productResponse.Data = &procudt
		json.NewEncoder(w).Encode(productResponse)

	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]map[string]interface{}
	// 解析参数 存入map
	decoder.Decode(&params)
	data := params["product"]
	pid := data["_id"]
	name := data["name"]
	desc := data["desc"]
	price := data["price"]
	detail := data["detail"]
	categoryId := data["categoryId"]
	imgs := data["imgs"]
	objId, _ := primitive.ObjectIDFromHex(pid.(string))
	filter := bson.M{"_id": objId}
	var product models.Product

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("product")
	err := collection.FindOne(ctx, filter).Decode(&product)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Println("更新失败，暂未该记录")
		errResponse.Msg = "更新失败，暂未该记录"
	} else {
		updateOption := bson.M{"categoryId": categoryId, "name": name, "desc": desc, "detail": detail, "imgs": imgs, "price": price}

		_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": updateOption})
		errResponse.Status = 0
		errResponse.Msg = "更新成功"
	}
	json.NewEncoder(w).Encode(errResponse)

}

func UpdateProductStatus(w http.ResponseWriter, r *http.Request) {

	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]interface{}
	// 解析参数 存入map
	decoder.Decode(&params)
	pid := params["productId"].(string)

	status := int(params["status"].(float64))
	objId, _ := primitive.ObjectIDFromHex(pid)
	var product models.Product
	product.Id = objId
	filter := bson.M{"_id": objId}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("product")
	err := collection.FindOne(ctx, filter).Decode(&product)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Println("更新失败，暂未该记录")
		errResponse.Msg = "更新失败，暂未该记录"
	} else {
		updateOption := bson.M{"$set": bson.M{"status": status}}

		collection.UpdateOne(ctx, filter, updateOption)
		errResponse.Status = 0
		errResponse.Msg = "商品更新成功！"
	}
	json.NewEncoder(w).Encode(errResponse)
}
