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

func GetSessionV2(c *gin.Context) {
	ctx := context.Background()

	var query apiv2.QueryGetSessionReq
	if err := c.BindJSON(&query); err != nil {
		c.JSON(http.StatusOK, apiv2.HttpRespBody{
			Code: 1000,
			Msg:  "",
			Data: apiv2.QueryGetSessionResp{
				Obj:   nil,
				Total: 0,
			},
		})
		return
	}

	log.Printf("query: %#v \n", query)

	//repo.GetRedisMgr().Traversal(ctx, true, query.SessionId, true)
	elem, total, e := repo.GetRedisMgr().Query(ctx, true, query.SessionId, !query.TimeAsc, query.PageId, query.PageSize)
	if e != nil {

		return
	}

	interfaces, e := util_struct.MultiConvertJsonStr2Interface(elem)
	if e != nil {

		log.Println("e: ", e)
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
