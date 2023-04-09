package main

import (
	"github.com/gin-gonic/gin"
	_ "gorm.io/driver/mysql"
	"jwt/controllers"
	"jwt/middlewares"
	"log"
)

type Hi struct {
	Name *string
}

func main() {

	r := gin.Default()

	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.SignUp)
	r.GET("users/all", middlewares.Authorize, controllers.GetUsers)

	log.Fatal(r.Run())

}
