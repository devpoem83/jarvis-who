package handler

import (
	"gitlab.eduwill.net/dev_team/jarvis-who/app/constant"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common/parser"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
	services "gitlab.eduwill.net/dev_team/jarvis-who/app/services/common"
)

func DefaultHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// before Process
		t := time.Now()

		uri := c.Request.RequestURI
		if uri == "/health" || strings.Index(uri, ".js") > -1 || strings.Index(uri, ".css") > -1 || strings.Index(uri, ".ico") > -1 || strings.Index(uri, ".html") > -1 {
			return
		}
		baseHandler(c)
		c.Next()

		go services.SetAccessHistory(c)		// 프로세스에 영향을 주지않기 위해 병렬처리

		status := c.Writer.Status()
		common.Logger.Info("RESPONSE STATUS - [" + strconv.Itoa(status) + "]")

		latency := time.Since(t)
		common.Logger.Info("TOTAL LATENCY   - [", latency, "]")
	}
}

func baseHandler(c *gin.Context) {
	edwUser, err := parser.GetEdwUser(c)
	isLogin := false
	if err != nil {
		c.Set(constant.CONST_EDW_USER, nil)

	} else {
		c.Set(constant.CONST_EDW_USER, edwUser)
		c.Set(constant.CONST_USER_ID, edwUser.UserId)
		c.Set(constant.CONST_USER_NM, edwUser.UserNm)
		isLogin = true

		common.Logger.Info("edwUser.UserId : ",edwUser.UserId)
		common.Logger.Info("edwUser.UserNm : ",edwUser.UserNm)
		
	}
	c.Set("LOGIN", isLogin)
	common.Logger.Info("edwUser.isLogin : ",isLogin)
}

func LoginCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		edwUser, isExist := c.Get(constant.CONST_EDW_USER)
		if !isExist || edwUser == nil {
			common.Logger.Info("401 Unauthorized!")
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
		}else{
			userId, _ := c.Get(constant.CONST_USER_ID)
			if(userId == ""){
				common.Logger.Info("401 Unauthorized!")
				c.JSON(http.StatusUnauthorized, gin.H{})
				c.Abort()
			}
		}
	}
}
