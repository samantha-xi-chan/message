package cmd

import (
	"context"
	"fmt"
	"github.com/Clouditera/message/api"
	"github.com/Clouditera/message/internal/repo"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

func MainModeGateway() {
	repo.InitMongo()

	fetchSession := func(c *gin.Context) {
		//  todo： 增加分页查询的保护，避免单次请求数据量过大

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
	}

	fetchSessionV2 := func(c *gin.Context) {
		sessionID := c.Param("session_id")
		evtType := c.Query("evt_type")
		log.Printf("sessionID: %s, evtType: %s ", sessionID, evtType)

		code, palyload := repo.GetSessionStatus(context.TODO(), sessionID, evtType)
		if code != repo.ERR_OK {
			body := api.BodyV2{
				ErrorCode: code,
				Message:   "ERR_NOT_FOUND",
			}
			c.JSON(200, body)
			return
		}
		log.Println("palyload:", palyload)
		data := api.QueryStatusResp{Payload: palyload}
		body := api.BodyV2{
			ErrorCode: 0,
			Message:   "ok",
			Data:      data,
		}
		c.JSON(200, body)
	}

	router := gin.Default()
	v1 := router.Group("/api/v1/sessions")
	{
		v1.GET("/:session_id", fetchSession) // const SESSION_ID = "ID10006666"  // curl localhost:8080/api/v1/session/ID10006666
	}
	v2 := router.Group("/api/v2/sessions")
	{
		v2.GET("/:session_id", fetchSessionV2)
	}
	router.Run(":18081")
}
