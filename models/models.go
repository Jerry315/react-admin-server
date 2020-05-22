package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Client *mongo.Client

type User struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserName   string             `json:"username" bson:"username" form:"username"`
	PassWord   string             `json:"password" bson:"password" form:"password"`
	Phone      string             `json:"phone" bson:"phone"`
	Email      string             `json:"email" bson:"email"`
	CreateTime time.Time          `json:"create_time" bson:"create_time"`
	RoleId     string             `json:"role_id" bson:"role_id"`
	Level      int                `json:"level" bson:"level"`
}

type Category struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	ParentId string             `json:"parentId" bson:"parentId"`
}

type Product struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	CategoryId string             `json:"categoryId" bson:"categoryId"`
	Name       string             `json:"name" bson:"name"`
	Price      string             `json:"price" bson:"price"`
	Desc       string             `json:"desc" bson:"desc"`
	Status     int                `json:"status" bson:"status"`
	Imgs       []string           `json:"imgs" bson:"imgs"`
	Detail     string             `json:"detail" bson:"detail"`
}

type Role struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	AuthName   string             `json:"auth_name" bson:"auth_name"`
	CreateTime time.Time          `json:"create_time" bson:"create_time"`
	AuthTime   time.Time          `json:"auth_time" bson:"auth_time"`
	Menus      []string           `json:"menus" bson:"menus"`
}

type Img struct {
	Name string `json:"name" bson:"name"`
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("连接数据库失败")
	}

	err = c.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("数据库连接超时")
	}
	Client = c

}
