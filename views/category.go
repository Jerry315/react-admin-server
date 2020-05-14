package views

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"test/admin-server/models"
	"time"
)

type CategoryListResponse struct {
	Status int                `json:"status"`
	Data   []*models.Category `json:"data"`
}

type RegCategoryResponse struct {
	Status int              `json:"status"`
	Data   *models.Category `json:"data"`
}

type CategoryResponse struct {
	RegCategoryResponse
}



func AddCategory(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("categoryName")
	filter := bson.M{"name": name}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("category")
	var category models.Category
	err := collection.FindOne(ctx, filter).Decode(&category)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Println(name)
		category.Name = name
		_, err = collection.InsertOne(ctx, &category)
		if err != nil {
			fmt.Printf("插入数据失败：%v\n", err)
			errResponse.Msg = "插入数据失败"
		} else {
			collection.FindOne(ctx, filter).Decode(&category)
			var regCategoryResponse RegCategoryResponse
			regCategoryResponse.Status = 0
			regCategoryResponse.Data = &category
			json.NewEncoder(w).Encode(regCategoryResponse)
			return
		}
	} else {
		fmt.Printf("插入数据成功")
		errResponse.Msg = "数据已存在"
	}
	json.NewEncoder(w).Encode(errResponse)
}

func GetCategoryList(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("category")
	cur, err := collection.Find(ctx, bson.M{})
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Printf("获取分类失败: %v\n", err)
		errResponse.Msg = "获取分类失败"
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	defer cur.Close(ctx)
	var categories []*models.Category
	for cur.Next(ctx) {
		var category *models.Category
		err = cur.Decode(&category)
		if err != nil {
			fmt.Printf("获取用户失败: %v\n", err)
			continue
		}
		categories = append(categories, category)
	}

	var categoryListResponse CategoryListResponse
	categoryListResponse.Status = 0
	categoryListResponse.Data = categories
	json.NewEncoder(w).Encode(categoryListResponse)

}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	cid :=r.FormValue("categoryId")
	objId, _ := primitive.ObjectIDFromHex(cid)
	var category models.Category
	category.Id = objId
	filter := bson.M{"_id": objId}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("category")
	err := collection.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		var errResponse ErrResponse
		errResponse.Status = 1
		errResponse.Msg = "查询失败"
		fmt.Printf("查询失败：%v\n", err)
		json.NewEncoder(w).Encode(errResponse)
	} else {
		var categoryResponse CategoryResponse
		categoryResponse.Status = 0
		categoryResponse.Data = &category
		json.NewEncoder(w).Encode(categoryResponse)

	}
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("categoryName")
	cid := r.FormValue("categoryId")
	objId, _ := primitive.ObjectIDFromHex(cid)
	var category models.Category
	category.Id = objId
	filter := bson.M{"_id": objId}
	updateOption := bson.M{"$set": bson.M{"name": name}}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("category")
	err := collection.FindOne(ctx, filter).Decode(&category)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Println("更新失败，暂未该记录")
		errResponse.Msg = "更新失败，暂未该记录"
	} else {
		category.Name = name
		_, err = collection.UpdateOne(ctx, filter, updateOption)
		errResponse.Status = 0
		errResponse.Msg = "更新成功"
	}
	json.NewEncoder(w).Encode(errResponse)

}
