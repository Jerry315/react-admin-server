package views

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"test/admin-server/models"
	"time"
)

type ImgResponse struct {
	Status int         `json:"status"`
	Data   *models.Img `json:"data"`
}

func GetImg(w http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.RequestURI, "/")
	img := "./upload/" + tmp[len(tmp)-1]
	fmt.Println(img)
	http.ServeFile(w, r, img)
}

func UploadImg(w http.ResponseWriter, r *http.Request) {
	uploadFile, handle, err := r.FormFile("image")
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Printf("获取上传图片信息失败, %v\n", err)
		errResponse.Msg = "获取上传图片信息失败"
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	// 检查图片后缀
	ext := strings.ToLower(path.Ext(handle.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		fmt.Printf("上传图片后缀须为.jpg、.jpeg、.png, %v\n", err)
		errResponse.Msg = "上传图片后缀须为.jpg、.jpeg、.png"
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	// 保存图片

	var img models.Img
	imgName := handle.Filename
	img.Name = imgName
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("img")
	filter := bson.M{"name": imgName}
	fmt.Println(filter)
	err = collection.FindOne(ctx, bson.M{"name": imgName}).Decode(&img)
	fmt.Printf("查询结果，%v\n", err)
	if err != nil {

		saveFile, _ := os.OpenFile("./upload/"+imgName, os.O_WRONLY|os.O_CREATE, 0666)
		io.Copy(saveFile, uploadFile)

		defer uploadFile.Close()
		collection.InsertOne(ctx, &img)
		collection.FindOne(ctx, bson.M{"name": imgName}).Decode(&img)
		var imgResponse ImgResponse
		imgResponse.Status = 0
		imgResponse.Data = &img
		json.NewEncoder(w).Encode(imgResponse)
	} else {
		fmt.Printf("图片已存在，请勿重复上传,%v\n", err)
		errResponse.Status = 0
		errResponse.Msg = "图片已存在，请勿重复上传"
		json.NewEncoder(w).Encode(errResponse)
		return
	}

}

func DeleteImg(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	decoder.Decode(&params)
	name := params["name"]
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := models.Client.Database("admin_db").Collection("img")
	filter := bson.M{"name": name}
	_, err := collection.DeleteOne(ctx, filter)
	var errResponse ErrResponse
	errResponse.Status = 1
	if err != nil {
		fmt.Printf("删除%s图片失败, %v\n", name, err)
		errResponse.Msg = "删除图片失败"
	} else {
		errResponse.Status = 0
		errResponse.Msg = "删除图片成功"
	}
	json.NewEncoder(w).Encode(errResponse)
}
