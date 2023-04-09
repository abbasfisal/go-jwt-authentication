package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"jwt/databases"
	"jwt/models"
	"net/http"
	"time"
)

func Login(c *gin.Context) {

	var userReq models.User

	//bind user data request
	bErr := c.ShouldBindJSON(&userReq)
	if bErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": bErr,
			"text":    "failed while binding data",
			"data":    userReq,
		})
		return
	}

	//validate
	fmt.Println("validate ---------", userReq.Validate())
	if userReq.Validate() == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request not valid",
			"data":    userReq,
		})
		return
	}

	//get email from db
	var userDb models.User
	databases.MysqlDb.Where("email = ?", userReq.Email).Find(&userDb)
	if userDb.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "email | password is invalid",
			"data":    userReq,
		})
		return
	}

	//check hashed password
	if !userDb.CheckPassword(userReq.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "password is invalid",
		})
		return
	}

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userDb.Email,
		"exp": time.Now().Add(time.Hour * 3).Unix(),
	})
	tokenString, err := token.SignedString([]byte("54564s65d46s5d4f56sd4f56s4d4ew")) //TODO : use .env file for secret key
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//set cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	//response
	c.JSON(http.StatusOK, gin.H{
		"user":  userDb,
		"token": tokenString,
	})
}

func SignUp(c *gin.Context) {
	var userReq models.User

	err := c.ShouldBindJSON(&userReq)
	//bind user data
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "error while binding data to user",
		})
		return
	}

	//validate
	valRes := userReq.Validate()
	if valRes == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request not valid ",
		})
		return
	}

	//check user existence
	var dbUser models.User
	databases.MysqlDb.Where("email=?", userReq.Email).Find(&dbUser)
	if dbUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email was exist! login",
		})
		return
	}

	//hash password
	hErr := userReq.GeneratePassword()
	if hErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   hErr,
			"message": "hash password failed",
		})
		return

	}
	//save to db
	res := databases.MysqlDb.Create(&userReq)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
		"data":    userReq,
	})
}

func GetUsers(c *gin.Context) {

	var usersDb []models.User
	res := databases.MysqlDb.Find(&usersDb)
	if res.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "error while fetching data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": usersDb,
	})
}
