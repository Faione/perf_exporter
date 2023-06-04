package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addPerfEventCollector(c *gin.Context) {
	var collector PerfEventConfig
	var resp Response

	if err := c.BindJSON(&collector); err != nil {
		resp = Response{Message: err.Error(), Data: nil}
	} else if err := AddCgroupPerfEventCollector(&collector); err != nil {
		resp = Response{Message: "add collector failed", Data: collector}
	} else {
		resp = Response{Message: "OK", Data: collector}
	}

	c.JSON(http.StatusOK, &resp)

}

func delPerfEventCollector(c *gin.Context) {
	var collector PerfEventConfig
	var resp Response

	if err := c.BindJSON(&collector); err != nil {
		resp = Response{Message: err.Error(), Data: nil}
	} else if err := DelCgroupPerfEventCollector(&collector); err != nil {
		resp = Response{Message: "delete collector failed", Data: collector}
	} else {
		resp = Response{Message: "OK", Data: collector}
	}

	c.JSON(http.StatusOK, &resp)
}
