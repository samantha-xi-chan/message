package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"message/api"
	"message/internal/repo"
	"strconv"
	"strings"
)

func FetchSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	pageID := c.Query("page_id")
	pageSize := c.Query("page_size")
	mode := c.Query("mode")
	log.Printf("sessionID %s, page_id %s, size %s, mode %s", sessionID, pageID, pageSize, mode)
	iPageID, _ := strconv.Atoi(pageID) // todo: check valid
	iPageSize, _ := strconv.Atoi(pageSize)
	detail := true
	if strings.Compare("simple", mode) == 0 {
		detail = false
	}

	err, cnt, lines := repo.GetSessionProfile(context.TODO(), sessionID, iPageID, iPageSize, detail)
	fmt.Println(err, cnt)
	if err == repo.ERR_OK {
		data := api.QuerySessionResp{cnt, lines}
		data = api.QuerySessionResp{Count: cnt, Array: lines}
		body := api.Body{
			ErrorCode: 0,
			Message:   "ok",
			Data:      data,
		}
		c.JSON(200, body)
	} else {
		data := api.QuerySessionResp{0, nil}
		body := api.Body{
			ErrorCode: repo.ERR_INVALID_PARA,
			Message:   "invalid para",
			Data:      data,
		}
		c.JSON(200, body)
	}

	return
}

func FetchSessionV2(c *gin.Context) {
	sessionID := c.Param("session_id")
	evtType := c.Query("evt_type")
	log.Printf("sessionID: %s, evtType: %s ", sessionID, evtType)

	code, payload := repo.GetSessionStatus(context.TODO(), sessionID, evtType)
	if code != repo.ERR_OK {
		body := api.BodyV2{
			ErrorCode: code,
			Message:   "ERR_NOT_FOUND",
		}
		c.JSON(200, body)
		return
	}
	log.Println("payload:", payload)
	data := api.QueryStatusResp{Payload: payload}
	body := api.BodyV2{
		ErrorCode: 0,
		Message:   "ok",
		Data:      data,
	}
	c.JSON(200, body)
}
