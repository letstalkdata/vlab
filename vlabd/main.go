package main

import (
	"fmt"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/letstalkdata/vlab/vlabd/routes"
	
	"github.com/gin-gonic/gin"

)

func connectDB() (c *pgx.Conn, err error){
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:vlab@localhost:5432/vlab")
	if err != nil {
		fmt.Println("Error Connecting to DB")
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleware(conn pgx.Conn) gin.HandlerFunc {
	return func (c *gin.Context) {
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
	}


	router.Run(":3000")
}