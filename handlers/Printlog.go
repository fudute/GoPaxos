package handlers

import (
	"log"
	"net/http"

	"github.com/fudute/GoPaxos/paxos"
	"github.com/gin-gonic/gin"
)

const (
	defaultFileName = "default.log"
)

type PrintLogResp struct {
	FileName string `json:"FileName"`
}

func PrintLog(c *gin.Context) {
	fileName := c.Param("path")
	if fileName == "/" {
		fileName = defaultFileName
	} else {
		fileName = fileName[1:]
	}
	log.Printf("try to write logs to file %v\n", fileName)
	paxos.DB.PrintLog(fileName)
	resp := PrintLogResp{FileName: fileName}
	c.JSON(http.StatusOK, resp)
}
