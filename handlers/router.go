package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Index is the handler for the path "/"
func Index(c *gin.Context) {
	fmt.Fprintf(c.Writer, "Hello World Paxos Server\n")
}

// Logger is the middleware to
// log the incoming request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"path":   r.URL,
			"method": r.Method,
		}).Info("incoming request")

		next.ServeHTTP(w, r)
	})
}

// Router returns a mux router
func Router() *gin.Engine {
	router := gin.Default()

	router.GET("/", Index)
	router.GET("/store/get/:key", GetValue)
	router.POST("/store/set/", SetValue)
	router.GET("/log/print/*path", PrintLog)
	router.GET("/store/nop", SendNop)
	router.GET("/crash", GetCrash)
	// router.POST("learn/:index", Learn)
	return router
}
