package main

import (
	"strings"
	"fmt"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/letstalkdata/vlab/vlabd/routes"
	"github.com/letstalkdata/vlab/vlabd/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func connectDB() (c *pgx.Conn, err error){
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:vlab@localhost:5432/vlab")
	if err != nil || conn == nil {
		fmt.Println("Error Connecting to DB")
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleware(conn pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}

func main() {

	conn, err := connectDB()
	if err != nil {
		return
	}

	router := gin.Default()

	router.Use(dbMiddleware(*conn))

	usersGroup := router.Group("users") 
	{
		usersGroup.POST("register", routes.UsersRegister)
		usersGroup.POST("login", routes.UsersLogin)
	}
	//router.POST("/users/register", routes.UsersRegister)
	//router.POST("/users/login", routes.UsersLogin)
	itemsGroup := router.Group("items")
	{
		itemsGroup.GET("index", routes.ItemIndex)
		itemsGroup.POST("create", authMiddleWare(), routes.ItemsCreate)
		itemsGroup.GET("sold_by_user", authMiddleWare(), routes.ItemsForSaleByCurrentUser)
	}


	router.Run(":3000")
}

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		split := strings.Split(bearer, "Bearer ")
		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}
		token := split[1]
		isValid, userID := models.IsTokenValid(token)
		if isValid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
		} else {
			c.Set("user_id", userID)
			c.Next()
		}
	}
}