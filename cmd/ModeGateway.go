package cmd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"message/api"
	"message/internal/config"
	"message/internal/repo"
	"message/internal/router"
	"message/package/util_debug"
	"strconv"
	"strings"
)

func MainModeGateway() {
	debugOn, e := config.GetDebugMode()
	if e != nil {
		log.Fatal("GetDebugMode: ", e)
	}

	if debugOn {
		addr, e := config.GetDebugPprofGateway()
		if e != nil {
			log.Fatal("config e: ", e)
		}

		log.Println("GetDebugPprofGateway addr: ", addr)
		go util_debug.InitPProf(addr)
	}

	v, _ := config.GetDependMongo()
	repo.InitMongo(v)

	redisDsn, e := config.GetDependRedisDsn()
	if e != nil {
		log.Fatal("GetDependRedisDsn: ", e)
	}
	log.Println("redisDsn: ", redisDsn)
	e = repo.InitRedis(context.Background(), redisDsn, 0, 0)
	if e != nil {
		log.Fatal("InitRedis: ", e)
	}

	fetchSession := func(c *gin.Context) {
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

	r := gin.Default()
	v1 := r.Group("/api/v1/sessions")
	{
		v1.GET("/:session_id", fetchSession) // const SESSION_ID = "ID10006666"  // curl localhost:8080/api/v1/session/ID10006666
	}
	v2 := r.Group("/api/v2/sessions")
	{
		v2.GET("/:session_id", fetchSessionV2)
	}
	v1r := r.Group("/api/v1r/session")
	{
		//v1r.GET("/:session_id", router.GetSessionV2)
		v1r.GET("", router.GetSessionV2)
	}

	val, _ := config.GetGwPortHttp()
	r.Run(val)
}
