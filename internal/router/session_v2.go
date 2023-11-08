package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"message/apiv2"
	"message/internal/repo"
	"message/package/util_struct"
	"net/http"
)

func GetSessionV1Status(c *gin.Context) { // to be deprecated
	ctx := context.Background()
	sessionId := c.Param("session_id")

	if sessionId == "" {
		c.JSON(http.StatusOK, apiv2.HttpRespBody{
			Code: apiv2.ERR_URL_ID,
			Msg:  "ERR_URL_ID",
		})
		return
	}

	var query apiv2.QueryGetSessionReq
	if err := c.BindJSON(&query); err != nil {
		c.JSON(http.StatusOK, apiv2.HttpRespBody{
			Code: apiv2.ERR_FORMAT,
			Msg:  "ERR_FORMAT",
		})
		return
	}

	log.Printf("query: %#v \n", query)

	// todo: check if it really exists
	// repo.GetRedisMgr().Exists()

	elem, total, e := repo.GetRedisMgr().Query(ctx, true, sessionId, query.TimeAsc, query.PageId, query.PageSize)
	if e != nil {
		c.JSON(http.StatusOK, apiv2.HttpRespBody{
			Code: apiv2.ERR_OTHER,
			Msg:  "ERR_OTHER",
		})
		return
	}

	interfaces, e := util_struct.MultiConvertJsonStr2Interface(elem)
	if e != nil {
		c.JSON(http.StatusOK, apiv2.HttpRespBody{
			Code: apiv2.ERR_OTHER,
			Msg:  "ERR_OTHER",
		})
		return
	}

	log.Printf("interfaces: %#v \n", interfaces)

	c.JSON(http.StatusOK, apiv2.HttpRespBody{
		Code: 0,
		Msg:  apiv2.WARNING_DEPRECATE,
		Data: apiv2.QueryGetSessionResp{
			Obj:   interfaces,
			Total: total,
		},
	})
	return
}

func GetSessionV1Log(c *gin.Context) {
	ctx := context.Background()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusOK, apiv2.HttpRespBody{
			Code: apiv2.ERR_URL_ID,
			Msg:  "ERR_URL_ID",
		})
		return
	}

	// Upsert Redis KeyValue

	var query apiv2.QueryGetSessionReqForm
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusOK, apiv2.HttpRespBody{
			Code: apiv2.ERR_FORMAT,
			Msg:  "ERR_FORMAT",
		})
		return
	}

	log.Printf("query: %#v \n", query)

	// todo: check if it really exists
	// repo.GetRedisMgr().Exists()

	repo.GetRedisMgr().TouchKeyHot(ctx, id, "", 0)

	elem, total, e := repo.GetRedisMgr().Query(ctx, true, id, query.TimeAsc, query.PageId, query.PageSize)
	if e != nil {
		c.JSON(http.StatusOK, apiv2.HttpRespBody{
			Code: apiv2.ERR_OTHER,
			Msg:  "ERR_OTHER",
		})
		return
	}

	interfaces, e := util_struct.MultiConvertJsonStr2Interface(elem)
	if e != nil {
		c.JSON(http.StatusOK, apiv2.HttpRespBody{
			Code: apiv2.ERR_OTHER,
			Msg:  "ERR_OTHER",
		})
		return
	}

	log.Printf("interfaces: %#v \n", interfaces)

	c.JSON(http.StatusOK, apiv2.HttpRespBody{
		Code: 0,
		Msg:  "",
		Data: apiv2.QueryGetSessionResp{
			Obj:   interfaces,
			Total: total,
		},
	})
	return
}
