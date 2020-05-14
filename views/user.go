package views

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"test/admin-server/models"
	"time"
)

type ErrResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type RegResponse struct {
	Status int          `json:"status"`
	Data   *models.User `json:"data"`
}

type LoginResponse struct {
	RegResponse
}

type UpdateResponse struct {
	RegResponse
}

type UserResponseList struct {
	Status int `json:"status"`
	Data   struct {
		Users []*models.User
	} `json:"data"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("user")
	//name := r.FormValue("username")
	//password := r.FormValue("password")
	//phone := r.FormValue("phone")
	//email := r.FormValue("email")
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	decoder.Decode(&params)
	name := params["username"]
	password := params["password"]
	phone := params["phone"]
	email := params["email"]
	var user models.User
	filter := bson.M{"username": name}
	err := collection.FindOne(ctx, filter).Decode(&user)
	var errResponse ErrResponse
	if err != nil {
		fmt.Println("暂未查到此用户可以注册")
		user.PassWord = password
		user.UserName = name
		user.Phone = phone
		user.Email = email
		user.CreateTime = time.Now()
		_, err = collection.InsertOne(ctx, &user)
		if err != nil {
			fmt.Println("插入注册数据失败")
			errResponse.Status = 1
			errResponse.Msg = "插入注册数据失败"
			json.NewEncoder(w).Encode(errResponse)
		} else {
			collection.FindOne(ctx, filter).Decode(&user)
			fmt.Println("用户注册成功")
			var regResponse RegResponse
			regResponse.Status = 1
			regResponse.Data = &user
			json.NewEncoder(w).Encode(regResponse)
		}
	} else {
		fmt.Println("用户已存在")
		errResponse.Status = 1
		errResponse.Msg = "用户已存在"
		json.NewEncoder(w).Encode(errResponse)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	decoder.Decode(&params)
	name := params["username"]
	password := params["password"]
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var user *models.User
	filter := bson.M{"username": name}
	collection := models.Client.Database("admin_db").Collection("user")
	err := collection.FindOne(ctx, filter).Decode(&user)
	var errResponse ErrResponse
	if err != nil {
		fmt.Printf("登录失败，用户名密码错, err: %s\n", err)
		errResponse.Status = 1
		errResponse.Msg = "登录失败，用户名密码错"
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	if user.PassWord != password {
		fmt.Printf("登录失败，用户名密码错, err: %s\n", err)
		errResponse.Status = 1
		errResponse.Msg = "登录失败，用户名密码错"
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	var loginResponse LoginResponse
	loginResponse.Status = 0
	loginResponse.Data = user
	json.NewEncoder(w).Encode(loginResponse)
}

func Update(w http.ResponseWriter, r *http.Request) {
	//name := r.FormValue("username")
	//password := r.FormValue("password")
	//phone := r.FormValue("phone")
	//email := r.FormValue("email")
	//role_id := r.FormValue("role_id")
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	decoder.Decode(&params)
	name := params["username"]
	password := params["password"]
	phone := params["phone"]
	email := params["email"]
	role_id := params["role_id"]
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var user *models.User
	filter := bson.M{"username": name}
	collection := models.Client.Database("admin_db").Collection("user")
	err := collection.FindOne(ctx, filter).Decode(&user)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Printf("%s用户不存在\n", name)
		errResponse.Status = 1
		errResponse.Msg = "用户不存在"
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	if password == "" || phone == "" || email == "" {
		fmt.Println("无更新内容")

		errResponse.Msg = "无更新内容"
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	update := bson.M{}
	if password != "" {
		update["password"] = password
	}
	if phone != "" {
		update["phone"] = phone
	}
	if email != "" {
		update["email"] = email
	}
	if role_id != "" {
		update["role_id"] = role_id
	}
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		errResponse.Msg = "无更新内容"
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	collection.FindOne(ctx, filter).Decode(&user)
	var updateResponse UpdateResponse
	updateResponse.Status = 0
	updateResponse.Data = user
	json.NewEncoder(w).Encode(updateResponse)

}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []*models.User
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := models.Client.Database("admin_db").Collection("user")
	cur, err := collection.Find(ctx, bson.M{})
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Printf("查询全部用户失败: %v\n", err)
		errResponse.Msg = "查询全部用户失败"
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user *models.User
		err = cur.Decode(&user)
		if err != nil {
			fmt.Printf("获取用户失败: %v\n", err)
			continue
		}
		users = append(users, user)
	}

	var userResponseList UserResponseList
	userResponseList.Status = 0
	for _, item := range users {
		userResponseList.Data.Users = append(userResponseList.Data.Users, item)
	}
	json.NewEncoder(w).Encode(userResponseList)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("user")
	//uid := r.FormValue("uid")
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	decoder.Decode(&params)
	uid := params["uid"]
	objId, _ := primitive.ObjectIDFromHex(uid)
	filter := bson.M{"_id": objId}

	_, err := collection.DeleteOne(ctx, filter)
	var errResponse ErrResponse
	if err != nil {
		errResponse.Status = 1
		errResponse.Msg = "删除用户失败"
	} else {
		errResponse.Status = 0
		errResponse.Msg = "用户删除成功"
	}
	json.NewEncoder(w).Encode(errResponse)
}
