package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type User struct {
	ID uint `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	Email string `json:"email"`
	Age int `json:"age"`
}

func main() {
	dsn := "host=localhost user=postgres password=19LynnP-Marie60 dbname=theatre-kid-db port=5432 sslmode=disable"

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	defer closeDBConnection(db)

	db.AutoMigrate(&User{})

	router := gin.Default()

	router.GET("/users", getUsers)
	router.POST("/users", createUser)

	router.Run(":8080")
}

func getUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func closeDBConnection(db *gorm.DB) {
	dbInstance, _ := db.DB()
	_ = dbInstance.Close()
}