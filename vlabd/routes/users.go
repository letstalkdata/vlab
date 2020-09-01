package routes

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/letstalkdata/vlab/vlabd/models"
)

func UsersLogin(c *gin.Context){
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err = user.IsAuthenticated(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := user.GetAuthToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "There was an error authenticating",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
	
	
	
	
}

func UsersRegister(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err = user.Register(&conn)
	if err != nil {
		fmt.Println("Error in user.Register()")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := user.GetAuthToken()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}