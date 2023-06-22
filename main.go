package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// データベース接続設定
const (
	DBHost     = "localhost"
	DBPort     = 5432
	DBUser     = "your_username"
	DBPassword = "your_password"
	DBName     = "your_database_name"
)

// モデル定義
type User struct {
	gorm.Model
	Name  string
	Email string
}

// ルーティングの設定
func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// ユーザー作成
	router.HandleFunc("/users", createUser).Methods("POST")
	// ユーザー取得
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	// ユーザー更新
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	// ユーザー削除
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	return router
}

// ユーザー作成
func createUser(w http.ResponseWriter, r *http.Request) {
	// リクエストからデータを取得
	name := r.FormValue("name")
	email := r.FormValue("email")

	// データベースに接続
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ユーザー作成
	user := User{Name: name, Email: email}
	db.Create(&user)

	// レスポンスを返す
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}

// ユーザー取得
func getUser(w http.ResponseWriter, r *http.Request) {
	// パスパラメータからIDを取得
	vars := mux.Vars(r)
	id := vars["id"]

	// データベースに接続
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ユーザー取得
	var user User
	db.First(&user, id)

	// ユーザーが見つからない場合はエラーレスポンスを返す
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	// レスポンスを返す
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Name: %s, Email: %s", user.Name, user.Email)
}

// ユーザー更新
func updateUser(w http.ResponseWriter, r *http.Request) {
	// パスパラメータからIDを取得
	vars := mux.Vars(r)
	id := vars["id"]

	// リクエストからデータを取得
	name := r.FormValue("name")
	email := r.FormValue("email")

	// データベースに接続
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ユーザー取得
	var user User
	db.First(&user, id)

	// ユーザーが見つからない場合はエラーレスポンスを返す
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	// ユーザー更新
	user.Name = name
	user.Email = email
	db.Save(&user)

	// レスポンスを返す
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User updated successfully")
}

// ユーザー削除
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// パスパラメータからIDを取得
	vars := mux.Vars(r)
	id := vars["id"]

	// データベースに接続
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ユーザー取得
	var user User
	db.First(&user, id)

	// ユーザーが見つからない場合はエラーレスポンスを返す
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	// ユーザー削除
	db.Delete(&user)

	// レスポンスを返す
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User deleted successfully")
}

func main() {
	// ルーターのセットアップ
	router := setupRouter()

	// サーバーの起動
	log.Println("Server started on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
