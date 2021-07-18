package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home ...
func (server *Server) Home(c *gin.Context) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"msg": "Welcome To This Awesome API",
	})

}
