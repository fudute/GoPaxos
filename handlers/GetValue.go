package handlers

import (
	"net/http"

	"github.com/fudute/GoPaxos/sm"
	"github.com/gin-gonic/gin"
)

// GetValue  is the HTTP handler to process incoming KV Get requests
// It gets the value from the in memory KV store
func GetValue(c *gin.Context) {
	key := c.Param("key")

	value, err := sm.GetKVStatMachineInstance().Execute("GET " + key)
	if err != nil {
		c.AbortWithError(http.StatusOK, err)
	}
	c.String(http.StatusOK, value)
}
