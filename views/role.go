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

type RoleResponse struct {
	Status int          `json:"status"`
	Data   *models.Role `json:"data"`
}

type RoleListResponse struct {
	Status int            `json:"status"`
	Data   []*models.Role `json:"data"`
}

func AddRole(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	decoder.Decode(&params)
	name := params["roleName"]
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("role")
	filter := bson.M{"name": name}
	var role models.Role
	err := collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		var roleResponse RoleResponse
		role.Name = name
		role.CreateTime = time.Now()
		collection.InsertOne(ctx, &role)
		collection.FindOne(ctx, filter).Decode(&role)
		roleResponse.Status = 0
		roleResponse.Data = &role
		json.NewEncoder(w).Encode(roleResponse)
	} else {
		var errResponse ErrResponse
		fmt.Println("改角色已存在")
		errResponse.Status = 1
		errResponse.Msg = "改角色已存在"
		json.NewEncoder(w).Encode(errResponse)
	}
}

func GetRoleList(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("role")
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("暂无角色数据")
	} else {
		var roles []*models.Role
		var roleListResponse RoleListResponse
		for cur.Next(ctx) {
			var role models.Role
			cur.Decode(&role)
			roles = append(roles, &role)
		}
		roleListResponse.Status = 0
		roleListResponse.Data = roles
		json.NewEncoder(w).Encode(roleListResponse)
	}
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]map[string]interface{}
	// 解析参数 存入map
	decoder.Decode(&params)
	data := params["role"]
	rid := data["_id"].(string)
	objId, _ := primitive.ObjectIDFromHex(rid)
	filter := bson.M{"_id": objId}
	menus := data["menus"]
	auth_time := time.Now().Local()
	auth_name := data["auth_name"]
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("role")
	var role models.Role
	err := collection.FindOne(ctx, filter).Decode(&role)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Printf("角色不存在，%#v\n", err)
		errResponse.Msg = "角色不存在"
	} else {
		updateOption := bson.M{"$set": bson.M{"menus": menus, "auth_time": auth_time, "auth_name": auth_name}}
		collection.UpdateOne(ctx, filter, updateOption)
		collection.FindOne(ctx, filter).Decode(&role)
		var roleResponse RoleResponse
		roleResponse.Status = 0
		roleResponse.Data = &role
		json.NewEncoder(w).Encode(roleResponse)
	}
}
