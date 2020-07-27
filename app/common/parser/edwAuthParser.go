package parser

import (
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common/crypto"
	"github.com/gin-gonic/gin"
	"strings"
)

type EdwUser struct {
	UserNo   string
	UserNm   string
	UserId   string
	Progress string
	BrnchNo  string
	PrgrNo   string
	InstNo   string
	LoginDt  string
	SystemCd string
	LoginIP  string
	HisSeq   string
	Edwauth  string
}

func GetEdwUser(c *gin.Context) (EdwUser, error) {
	edwUser := EdwUser{"", "", "", "", "", "", "", "", "", "", "", ""}
	cookie, err := c.Cookie("EDWAUTH01")

	//usrNo=1072306|usrNm=안치순|usrId=chris83|progress=G|brnchNo=100|prgrNo=0|instNo=0|loginDt=20190905134723|systemCd=|loginIP=10.10.13.16|hisSeq=115535250
	if err != nil {
		return edwUser, err
	}
	if cookie != "" {
		cookie = strings.ReplaceAll(cookie, " ", "+")
		authCookie, _ := crypto.DecryptAes256(cookie)
		items := strings.Split(authCookie, "|")
		for _, obj := range items {
			data := strings.Split(obj, "=")
			key := data[0]
			value := data[1]

			if key == "usrNo" {
				edwUser.UserNo = value
			} else if key == "usrNm" {
				edwUser.UserNm = value
			} else if key == "usrId" {
				edwUser.UserId = value
			} else if key == "progress" {
				edwUser.Progress = value
			} else if key == "brnchNo" {
				edwUser.BrnchNo = value
			} else if key == "prgrNo" {
				edwUser.PrgrNo = value
			} else if key == "instNo" {
				edwUser.InstNo = value
			} else if key == "loginDt" {
				edwUser.LoginDt = value
			} else if key == "systemCd" {
				edwUser.SystemCd = value
			} else if key == "loginIP" {
				edwUser.LoginIP = value
			} else if key == "hisSeq" {
				edwUser.HisSeq = value
			}
		}
	}

	return edwUser, nil
}
