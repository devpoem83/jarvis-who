package auth

import (
	"net/http"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/base/config"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common/dbTemplate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/dgrijalva/jwt-go"
	"time"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
	"strings"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common/crypto"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/constant"
)
var mapperPrefix = "auth."

var jwtKey = []byte("my_secret_key")

var users = map[string]string {
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Password string `form:"password" json:"password"`
	Username string `form:"username" json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Headers struct {
	Authorization string `header: "Authorization"`
}

type Response struct {
	Status bool
	Code	int
	Message string
}

// 인증 및 Token 발급
func Signin(c *gin.Context) {
	var creds Credentials
	var accessToken = ""
	var refreshToken = ""
	var makeType = "1"
	var userId = ""

	// Header Data Bind
	if err := c.ShouldBind(&creds); err != nil {
		common.Logger.Error(err)
		c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err})
		return
	}

	// 로그인 기반처리
	if creds.Username != "" && creds.Password != "" {
		common.Logger.Debug("creds : ", creds)
		
		// 인증결과 조회
		ok := validSign(creds.Username, creds.Password)
	
		// 회원 정보 확인
		if !ok {
			common.Logger.Debug("password is valid")
			c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: ""})
			return
		}

		common.Logger.Info("Credential for Username and Password")
		userId = creds.Username

	} else {
		value, _ := c.Get(constant.CONST_USER_ID)		// 인증쿠키에서 회원정보 조회
		userId := value.(string)
		makeType = "2"

		if userId != "" {
			common.Logger.Info("Credential for Authenticated Cookie")	
		}
	}

	if userId != "" {

		var err1 error
		var err2 error

		// Token 생성
		accessToken, err1 = makeToken(userId, time.Now().Add(5 * time.Minute))
		//refreshToken, err2 = makeToken(userId, time.Now().Add(14 * 24 * time.Hour))
		refreshToken, err2 = makeToken(userId, time.Now().Add(60 * time.Minute))

		if err1 != nil || err2 != nil {
			// Token 생성시 에러 발생
			if err1 != nil {
				common.Logger.Error(err1)
				c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err1})
			}

			if err2 != nil {
				common.Logger.Error(err2)
				c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err2})
			}
			return
		}

		// Token정보를 DB에 저장한다.
		_, err3 := registToken(userId, accessToken, refreshToken, makeType)
		if err3 != nil {
			c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err3})
		}

	}else{
		// 인증할 수 있는 수단이 없을경우
		c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: nil})
		return
	}

	c.Writer.Header().Set("Authorization", accessToken)
	
	resultData := make(map[string]interface{})
	resultData["RefreshToken"] = refreshToken
	c.Render(http.StatusOK, render.IndentedJSON{Data: resultData})
}

// 유효한 아이디, 비밀번호인지 확인
func validSign(username string, password string) bool {
	params := make(map[string]interface{})
	params["userId"] = username

	db := config.GetGanagosiDB()
	obj, err := dbTemplate.SelectOne(db, mapperPrefix+"user-info", params)

	if err != nil {
		common.Logger.Error(err)
		return false
	}

	if obj == nil {
		return false
	}
	
	// 입력한 비밀번호를 Sha256으로 암호화 문자열 생성
	ecriptedPassword := crypto.EncryptSha256(password)
	common.Logger.Debug("ecriptedPassword : " + ecriptedPassword)
	
	// DB에 Sha256으로 암호화 되어있는 비밀번호 조회
	dbPassword := obj["password"].(string)
	common.Logger.Debug("dbPassword : " + dbPassword)

	// 암호화된 비밀번호 비교
	if dbPassword == ecriptedPassword {
		common.Logger.Debug("[" + username + "] Username and Password is valid")
		return true
	}

	return false
}

// Token 생성처리
func makeToken(username string, expirationTime time.Time) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Token정보 DB저장처리
func registToken(userId string, accessToken string, refreshToken string, makeType string) (int, error) {
	params := make(map[string]interface{})
	params["userId"] = userId
	params["accessToken"] = accessToken
	params["refreshToken"] = refreshToken
	params["makeType"] = makeType

	db := config.GetGanagosiDB()
	return dbTemplate.SelectValue(db, mapperPrefix+"add-token", params)
}

func updateToken(userId string, oldAccessToken string, newAccessToken string, refreshToken string) (int, error) {
	params := make(map[string]interface{})
	params["userId"] = userId
	params["oldAccessToken"] = oldAccessToken
	params["newAccessToken"] = newAccessToken
	params["refreshToken"] = refreshToken

	db := config.GetGanagosiDB()
	return dbTemplate.SelectValue(db, mapperPrefix+"update-token", params)
}

// Token 유효성 확인
func validToken(token string) (bool, error) {
	if token != "" {
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil {
			if tkn.Valid {
				return true, nil
			}
		}else{
			return false, err
		}
	}
	return false, nil
}

// AccessToken으로 서비스요청
func Welcome(c *gin.Context) {
	headers := Headers{}
	if err := c.ShouldBindHeader(&headers); err != nil{
		common.Logger.Error(err)
		c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err})
		return
	}
	accessToken := strings.Replace(headers.Authorization, "Bearer ", "", 1)

	if accessToken == "" {
		c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: nil})
		return
	}

	// Refresh Token 유효여부 확인
	ok, err := validToken(accessToken)

	if err != nil {
		c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: err})
		return
	}
	
	common.Logger.Debug("Valid AccessToken : ", ok)
	if !ok {
		c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: nil})
		return
	}

	c.Render(http.StatusOK, render.IndentedJSON{Data: nil})
}

// Refresh Token을 통한 새로운 Token 발급
func Refresh(c *gin.Context) {
	
	headers := Headers{}
	if err := c.ShouldBindHeader(&headers); err != nil{
		common.Logger.Error(err)
		c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err})
		return
	}
	accessToken := strings.Replace(headers.Authorization, "Bearer ", "", 1)
	newAccessToken := ""
	refreshToken := c.PostForm("refreshToken")

	common.Logger.Debug("accessToken : " + accessToken)
	common.Logger.Debug("refreshToken : " + refreshToken)

	if accessToken == "" || refreshToken == "" {
		c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: nil})
		return
	}
	
	// Refresh Token 유효여부 확인
	ok, err := validToken(refreshToken)

	if err != nil {
		common.Logger.Error(err)
		c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: err})
		return
	}
	
	common.Logger.Debug("Valid RefreshToken : ", ok)
	if !ok {
		c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: nil})
		return
	}

	params := make(map[string]interface{})
	params["accessToken"] = accessToken
	params["refreshToken"] = refreshToken

	db := config.GetGanagosiDB()
	obj, err := dbTemplate.SelectOne(db, mapperPrefix+"read-token", params)

	if err != nil {
		c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err})
		return
	}

	if obj == nil {
		c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: nil})
		return
	}

	userId := obj["userId"].(string)
	if userId == "" {
		c.Render(http.StatusUnauthorized, render.IndentedJSON{Data: nil})
		return
	}

	var err1 error
	var err2 error
	newAccessToken, err1 = makeToken(userId, time.Now().Add(5 * time.Minute))

	if err1 != nil || err2 != nil {
		// Token 생성시 에러 발생
		if err1 != nil {
			common.Logger.Error(err1)
			c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err1})
		}

		if err2 != nil {
			common.Logger.Error(err2)
			c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err2})
		}
		return
	}

	// Token정보를 DB에 저장한다.
	_, err3 := updateToken(userId, accessToken, newAccessToken, refreshToken)
	if err3 != nil {
		c.Render(http.StatusInternalServerError, render.IndentedJSON{Data: err3})
		return
	}

	c.Writer.Header().Set("Authorization", newAccessToken)
	c.Render(http.StatusOK, render.IndentedJSON{Data: nil})

}
