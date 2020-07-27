package common

import (
	"github.com/gin-gonic/gin"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/constant"
)

var mapperPrefix = "common."

func SetAccessHistory(c *gin.Context) {
	common.Logger.Info("==============================================================================================")
	common.Logger.Debug("접근이력 기록 Begin...")
	
	userId, _ := c.Get(constant.CONST_USER_ID)
	userNm, _ := c.Get(constant.CONST_USER_NM)

	uri := c.Request.RequestURI
	ip := c.ClientIP()
	agent := c.Request.UserAgent()
	referer := c.Request.Referer()
	method := c.Request.Method

	common.Logger.Info("USER_ID       : ", userId)
	common.Logger.Info("USER_NM       : ", userNm)
	common.Logger.Info("IP       : ", ip)
	common.Logger.Info("AGENT    : ", agent)
	common.Logger.Info("URI      : ", uri)
	common.Logger.Info("METHOD   : ", method)
	common.Logger.Info("REFERER  : ", referer)
	common.Logger.Debug("접근이력 기록 End...")
	common.Logger.Info("==============================================================================================")
	
}