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
	"strings"
	"test/admin-server/models"
	"time"
)

type ProductListResponse struct {
	Status int `json:"status"`

	Data struct {
		PageNum  int `json:"pageNum"`
		Total    int `json:"total"`
		Pages    int `json:"pages"`
		PageSize int `json:"pageSize"`
		List     []*models.Product
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
	name := r.FormValue("name")
	cid := r.FormValue("categoryId")
	desc := r.FormValue("desc")
	price := r.FormValue("price")
	detail := r.FormValue("detail")
	imgs := strings.Split(r.FormValue("imgs"), ",")
	filter := bson.M{"name": name}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("product")
	var product models.Product
	err := collection.FindOne(ctx, filter).Decode(&product)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Println(name)
		product.Name = name
		product.CategoryId = cid
		product.Desc = desc
		product.Price = price
		product.Detail = detail
		product.Imgs = imgs
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
	pageNum, _ := strconv.Atoi(r.FormValue("pageNum"))
	pageSize, _ := strconv.Atoi(r.FormValue("pageSize"))
	productName := r.FormValue("productName")
	productDesc := r.FormValue("productDesc")
	filter := bson.M{}
	if productName != "" {
		filter["name"] = productName
	}
	if productDesc != "" {
		filter["desc"] = productDesc
	}
	findOptions := options.Find()
	findOptions.SetSkip(int64(pageNum - 1))
	findOptions.SetLimit(int64(pageSize))

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
	for cur.Next(ctx) {
		var product *models.Product
		err = cur.Decode(&product)
		if err != nil {
			fmt.Printf("获取用户失败: %v\n", err)
			continue
		}
		products = append(products, product)
	}

	var productListResponse ProductListResponse
	productListResponse.Status = 0
	productListResponse.Data.List = products
	productListResponse.Data.PageNum = pageNum
	productListResponse.Data.PageSize = pageNum
	productListResponse.Data.Total = len(products)
	json.NewEncoder(w).Encode(productListResponse)

}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	pid := r.FormValue("productId")
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
	pid := r.FormValue("_id")
	name := r.FormValue("name")
	cid := r.FormValue("categoryId")
	desc := r.FormValue("desc")
	price := r.FormValue("price")
	detail := r.FormValue("detail")
	imgs := strings.Split(r.FormValue("imgs"), ",")
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
		updateOption := bson.M{"categoryId": cid}
		if desc != "" {
			updateOption["desc"] = desc
		}
		if price != "" {
			updateOption["price"] = price
		}

		if detail != "" {
			updateOption["detail"] = detail
		}
		if len(imgs) > 0 {
			updateOption["imgs"] = imgs
		}
		product.Name = name
		_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": updateOption})
		errResponse.Status = 0
		errResponse.Msg = "更新成功"
	}
	json.NewEncoder(w).Encode(errResponse)

}

func UpdateProductStatus(w http.ResponseWriter, r *http.Request) {
	pid := r.FormValue("productId")
	status, _ := strconv.Atoi(r.FormValue("status"))
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
		errResponse.Msg = "商品下架成功！"
	}
	json.NewEncoder(w).Encode(errResponse)
}
