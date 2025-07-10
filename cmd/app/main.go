package main

import (
	"log"
	"net/http"
	"time"

	myDB "github.com/AndreySmirnoffv/golang-auth-training/internal/db"
	"github.com/AndreySmirnoffv/golang-auth-training/internal/usecases"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=... user=... password=... dbname=... port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&myDB.UserModel{}); err != nil {
		log.Fatal(err)
	}

	uRepo := myDB.NewUserRepoGORM(db)
	uuc := usecases.NewUserUseCase(uRepo)
	uHandler := myHttp.newUserHandler(uuc)

	r := gin.Default()
	r.POST("/register", uHandler.Register)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Println("Server starting on :8080")
	log.Fatal(srv.ListenAndServe())
}
