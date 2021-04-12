package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func GetCrash(c *gin.Context) {
	os.Exit(1000)
}
