package databases

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jwt/models"
	"log"
)

var MysqlDb *gorm.DB
var err error

func init() {
	MysqlDb, err = gorm.Open(mysql.Open("root:root@/golang"))
	if err != nil {
		log.Fatal("connect to db failed : ", err)
	}

	//drop table
	//MysqlDb.Migrator().DropTable(&models.User{})
	MysqlDb.AutoMigrate(&models.User{})
	fmt.Println("connected to mysql successfully")
}
