package controllers

import (
	"github.com/gin-gonic/gin"
)

// MenuItemsController ...
type MenuItemsController struct{}

// Retrieve ...
func (u MenuItemsController) Retrieve(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
	return
}
