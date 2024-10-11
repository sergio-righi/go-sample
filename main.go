package main

import (
	"fmt"
	"go-sample/controllers"
	"go-sample/db"
	"go-sample/integrations"
	"go-sample/routes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongo := db.MongoDB()
	// Connect to MongoDB
	mongo.Connect()

	minio := integrations.Minio()
	// Connect to MinIO
	minio.Connect()

	authController := controllers.AuthController()
	s3Controller := controllers.S3Controller(minio)
	userController := controllers.UserController(mongo.GetCollection("users"))

	router := routes.InitRoutes(authController, s3Controller, userController)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
