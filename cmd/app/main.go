package main

import (
	"log"
	"net/http"
	"os"
	"time"

	myDB "github.com/AndreySmirnoffv/golang-auth-training/internal/adapter/db"
	myHttp "github.com/AndreySmirnoffv/golang-auth-training/internal/adapter/http"
	"github.com/AndreySmirnoffv/golang-auth-training/internal/adapter/jwt"
	usecases "github.com/AndreySmirnoffv/golang-auth-training/internal/usecases/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&myDB.UserModel{}, &myDB.PaymentModel{}); err != nil {
		log.Fatal(err)
	}

	jwtSrv := jwt.NewJWTService(
		os.Getenv("ACCESS_SECRET"),
		os.Getenv("REFRESH_SECRET"),
		time.Minute*15,
		time.Hour*24*7,
	)

	uRepo := myDB.NewUserRepoGORM(db)
	uuc := usecases.NewUserUseCase(uRepo, jwtSrv)
	uHandler := myHttp.NewUserHandler(*uuc)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Refresh-Token"},
		ExposeHeaders:    []string{"X-New-Access-Token", "X-New-Refresh-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	auth := api.Group("/auth")

	auth.POST("/register", uHandler.Register)
	auth.POST("/login", uHandler.Login)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Println("Server starting on :8080")
	log.Fatal(srv.ListenAndServe())
}
