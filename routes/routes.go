package routes

import (
	"go-sample/controllers"
	"go-sample/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes(
	authController *controllers.AuthControllerType,
	s3Controller *controllers.S3ControllerType,
	personController *controllers.PersonControllerType,
	userController *controllers.UserControllerType) *mux.Router {
	router := mux.NewRouter()

	jwt := middlewares.Jwt(authController)

	// Define protected routes

	// person routes
	router.Handle("/people", jwt.Handler(http.HandlerFunc(personController.Base.All))).Methods("GET")
	router.Handle("/people", jwt.Handler(http.HandlerFunc(personController.Base.Create))).Methods("POST")
	router.Handle("/people/{id}", jwt.Handler(http.HandlerFunc(personController.Base.Find))).Methods("GET")
	router.Handle("/people/{id}", jwt.Handler(http.HandlerFunc(personController.Base.Update))).Methods("PUT")
	router.Handle("/people/{id}", jwt.Handler(http.HandlerFunc(personController.Base.Delete))).Methods("DELETE")

	// user routes
	router.Handle("/users", jwt.Handler(http.HandlerFunc(userController.Base.All))).Methods("GET")
	router.Handle("/users/{id}", jwt.Handler(http.HandlerFunc(userController.Base.Find))).Methods("GET")
	router.Handle("/users/{id}", jwt.Handler(http.HandlerFunc(userController.Base.Update))).Methods("PUT")
	router.Handle("/users/{id}", jwt.Handler(http.HandlerFunc(userController.Base.Delete))).Methods("DELETE")

	// s3 routes
	router.Handle("/s3/upload", jwt.Handler(http.HandlerFunc(s3Controller.Upload))).Methods("POST")
	router.Handle("/s3/get/{bucket}/{object}", jwt.Handler(http.HandlerFunc(s3Controller.Get))).Methods("GET")
	router.Handle("/s3/list/{bucket}", jwt.Handler(http.HandlerFunc(s3Controller.List))).Methods("GET")
	router.Handle("/s3/delete/{bucket}/{object}", jwt.Handler(http.HandlerFunc(s3Controller.Delete))).Methods("DELETE")
	router.Handle("/s3/share/{bucket}/{object}", jwt.Handler(http.HandlerFunc(s3Controller.Share))).Methods("GET")
	router.Handle("/s3/rename/{bucket}/{object}", jwt.Handler(http.HandlerFunc(s3Controller.Rename))).Methods("PUT")

	// Define public routes
	router.HandleFunc("/auth", authController.Auth).Methods("POST")
	router.HandleFunc("/users", userController.Create).Methods("POST")

	return router
}
