package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SimpleRequest1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ping",
	})
}

func SimpleRequest2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func SimpleRequest3(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "Hello %s", name)
}

func SimpleRequest4(c *gin.Context) {
	name := c.Param("name")
	action := c.Param("action")
	message := name + " is " + action

	c.String(http.StatusOK, message)
}

func SimpleRequest5(c *gin.Context) {
	b := c.FullPath() == "/user/:name/*action" // true
	c.String(http.StatusOK, "%t", b)
}
