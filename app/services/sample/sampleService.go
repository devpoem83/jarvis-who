package sample

import (
	"net/http"
	"gitlab.eduwill.net/dev_team/landus-api/app/base/config"
	"gitlab.eduwill.net/dev_team/landus-api/app/common/dbTemplate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	//"gitlab.eduwill.net/dev_team/landus-api/app/common"
)
var mapperPrefix = "sample."

func Samples(c *gin.Context) {
	params := make(map[string]interface{})
	db := config.GetGanagosiDB()
	list, _ := dbTemplate.SelectList(db, mapperPrefix+"list", params)
	c.Render(http.StatusOK, render.IndentedJSON{Data: list})
}

func ReadSample(c *gin.Context) {
	obj := make(map[string]interface{})

	id := c.Param("id")
	if id != "" {
		params := make(map[string]interface{})
		params["id"] = id
		db := config.GetGanagosiDB()
		obj, _ = dbTemplate.SelectOne(db, mapperPrefix+"read", params)
	}

	c.Render(http.StatusOK, render.IndentedJSON{Data: obj})
}

func CreateSample(c *gin.Context) {
	
	var resultData = make(map[string]interface{})
	var result = false
	
	//userId, _ := c.Get("edwUserId")
	userId := "chris83"
	content := c.PostForm("content")
	
	if userId != "" && content != "" {
		params := make(map[string]interface{})
		params["userId"] = userId
		params["content"] = content
		
		db := config.GetGanagosiDB()
		if value, err := dbTemplate.SelectValue(db, mapperPrefix+"create", params); err == nil && value > 0 {
			result = true
			resultData["value"] = value
		}
	}

	resultData["result"] = result
	c.Render(http.StatusOK, render.IndentedJSON{Data: resultData})
}

func UpdateSample(c *gin.Context) {
	var resultData = make(map[string]interface{})
	var result = false
	
	//userId, _ := c.Get("edwUserId")
	userId := "chris83"
	id := c.Param("id")
	content := c.PostForm("content")
	
	if id != "" && userId != "" && content != "" {
		params := make(map[string]interface{})
		params["userId"] = userId
		params["id"] = id
		params["content"] = content
		
		db := config.GetGanagosiDB()
		if value, err := dbTemplate.SelectValue(db, mapperPrefix+"update", params); err == nil && value > 0 {
			result = true
		}
	}

	resultData["result"] = result
	c.Render(http.StatusOK, render.IndentedJSON{Data: resultData})
}

func DeleteSample(c *gin.Context) {
	var resultData = make(map[string]interface{})
	var result = false

	id := c.Param("id")
	//userId, _ := c.Get("edwUserId")
	userId := "chris83"
	
	if id != "" && userId != "" {
		params := make(map[string]interface{})
		params["id"] = id
		params["userId"] = userId
		
		db := config.GetGanagosiDB()
		if value, err := dbTemplate.SelectValue(db, mapperPrefix+"delete", params); err == nil && value > 0 {
			result = true
		}
	}

	resultData["result"] = result
	c.Render(http.StatusOK, render.IndentedJSON{Data: resultData})
}

func TxSamples(c *gin.Context) {
	result := false
	params := make(map[string]interface{})
	params["userId"] = "chris83"
	params["studyPlanNo"] = 1149
	
	db := config.GetGanagosiDB()

	// 트랜잭션 시작
	tx, err1 := dbTemplate.BeginTx(db)
	if err1 == nil {
		_, _, err2 := dbTemplate.ExecTx(tx, mapperPrefix+"update1", params)
		_, rowsAffected, err3 := dbTemplate.ExecTx(tx, mapperPrefix+"update2", params)

		if err2 != nil || err3 != nil || rowsAffected == 0 {
			dbTemplate.RollbackTx(tx)	// 트랜잭션 롤백
		} else {
			dbTemplate.CommitTx(tx)		// 트랜잭션 커밋
		}
	}
	
	c.Render(http.StatusOK, render.IndentedJSON{Data: result})
}
