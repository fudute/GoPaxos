package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index is the handler for the path "/"
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World Paxos Server")
}

// Router returns a mux router
func Router() *gin.Engine {
	router := gin.Default()

	router.GET("/", Index)
	router.GET("/store/get/:key", GetValue)
	router.POST("/store/set/", SetValue)
	router.GET("/log/print/*path", PrintLog)

	router.GET("/store/nop", Nop)
	router.GET("/crash", GetCrash)
	// router.POST("learn/:index", Learn)
	return router
}
