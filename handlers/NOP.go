package handlers

import (
	"log"
	"net/http"

	"github.com/fudute/GoPaxos/paxos"
	"github.com/fudute/GoPaxos/request"
	"github.com/gin-gonic/gin"
)

func Nop(c *gin.Context) {
	req := request.Nop()

	paxos.GetBatcherInstance().In <- req

	err := <-req.Done()
	if err != nil {
		log.Printf("NOP error: %v", err)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
	c.JSON(http.StatusOK, nil)
}
