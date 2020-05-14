package main

import (
	"github.com/gorilla/mux"
	"github.com/micro/go-log"
	"net/http"
	"test/admin-server/views"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", views.Login).Methods("POST")
	router.HandleFunc("/manage/user/register", views.Register).Methods("POST")
	router.HandleFunc("/manage/user/update", views.Update).Methods("POST")
	router.HandleFunc("/manage/user/delete", views.Delete).Methods("POST")
	router.HandleFunc("/manage/user/list", views.GetAllUsers).Methods("GET")
	router.HandleFunc("/manage/category/list", views.GetCategoryList).Methods("GET")
	router.HandleFunc("/manage/category/add", views.AddCategory).Methods("POST")
	router.HandleFunc("/manage/category/update", views.UpdateCategory).Methods("POST")
	router.HandleFunc("/manage/category/info", views.GetCategory).Methods("GET")
	router.HandleFunc("/manage/product/list", views.GetProductList).Methods("GET")
	router.HandleFunc("/manage/product/search", views.GetProductList).Methods("GET")
	router.HandleFunc("/manage/product/info", views.GetProduct).Methods("GET")
	router.HandleFunc("/manage/product/add", views.AddProduct).Methods("POST")
	router.HandleFunc("/manage/product/update", views.UpdateProduct).Methods("POST")
	router.HandleFunc("/manage/product/updateStatus", views.UpdateProductStatus).Methods("POST")
	router.HandleFunc("/manage/img/upload", views.UploadImg).Methods("POST")
	router.HandleFunc("/manage/img/delete", views.DeleteImg).Methods("POST")
	router.HandleFunc("/manage/role/add", views.AddRole).Methods("POST")
	router.HandleFunc("/manage/role/list", views.GetRoleList).Methods("GET")
	router.HandleFunc("/manage/role/update", views.UpdateRole).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
