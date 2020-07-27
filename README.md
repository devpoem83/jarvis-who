# Token 기반의 인증 시스템

## 백엔드 인증 서비스

### 도메인

-   운영 : http://jarvis-who.eduwill.net
-   스테이지 : http://s-jarvis-who.eduwill.net
-   개발 : http://d-jarvis-who.eduwill.net
-   로컬 : http://l-jarvis-who.eduwill.net

### 개발언어 & 프레임워크

-   Language : Go v1.14
-   Framework : Gin v1.3

### 형상관리

-   Gitlab : http://gitlab.eduwill.net/dev_team/land-api

### 빌드 & 배포관리

-   빌드 : Jenkins
-   배포 : Marathon
-   알림 : Hangout Chat

### 개발환경 구축

-   Go 설치
    https://golang.org
-   VSCode 설치 https://code.visualstudio.com/
-   Git 연동 https://git-scm.com/
-   소스 clone

```
git clone http://gitlab.eduwill.net/dev_team/jarvis-who.git
```

-   호스트 설정

```
127.0.0.1   l-jarvis-who.eduwill.net
```

-   디펜던시

```
$ go mod download
```

-   빌드

```
$ go build
```

-   실행

```
$ jarvis-who.exe // windows os
$ jarvis-who // linux
```

-   테스트

```
$ go test -v -cover ./...
```

## 개발 가이드

### 프로파일 설정

-   구동시 환경 옵션설정

```
$ jarvis-who.exe --env=live // windows
$ jarvis-who --env=live // linux
```

### 로거 설정 및 사용법

-   포맷 및 레벨 설정 /app/base/common/logger.go

```
func LoggerInit(level string) {
	logLevel := logrus.DebugLevel		// Default level
	timeFormat := "2006-01-02 15:04:05"	// Date format
	logFormat := "[%lvl%] %time% --> %msg% \n"	// Log print format

	if level == "info" {
		logLevel = logrus.InfoLevel
	} else if level == "warn" {
		logLevel = logrus.WarnLevel
	} else if level == "error" {
		logLevel = logrus.ErrorLevel
	}

	Logger = logrus.Logger{
		Out:   os.Stdout,
		Level: logLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: timeFormat,
			LogFormat:       logFormat,
		},
	}
}
```

-   사용법

```
import (
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
)
common.Logger.Debug("Debug Log")
common.Logger.Info("Info Log")
common.Logger.Warn("Warn Log")
common.Logger.Error("Error Log")
```

### 라우팅

-   Basic route

```
func main() {
	r := gin.Default()
	r.GET("/sample", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "sample"})
	})
}
```

-   Grouping routes

```
func main() {
	r := gin.Default()

	// Simple group: v1
	v1 := r.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// Simple group: v2
	v2 := r.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}
}

```

### 데이터 바인딩

-   Querystring Parameter
    (http://jarvis-who.eduwill.net?firstname=gildong&lastname=hong)

```
lastname := c.Query("lastname")
firstname := c.DefaultQuery("firstname", "Guest")
```

-   Parameter in path
    (http://jarvis-who.eduwill.net/sample/:name)

```
name := c.Param("name")
```

-   Multipart/Urlencoded Form

```
message := c.PostForm("message")
nick := c.DefaultPostForm("nick", "anonymous")
```

### 미들웨어 & 핸들러

-   Using Middleware

```
func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// Per route middleware, you can add as many as you desire.
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)

		// nested group
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}
}
```

-   Custom Middleware

```
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func main() {
	r := gin.New()
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)

		// it would print: "12345"
		log.Println(example)
	})
}
```

### Redirect

```
r.GET("/sample", func(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/sample/list")
})
```

### 회원정보

```
import (
	"gitlab.eduwill.net/dev_team/jarvis-who/app/constant"
)
userId, _ := c.Get(constant.CONST_USER_ID)
userNm, _ := c.Get(constant.CONST_USER_NM)
```

### 인증 / 권한

-   개별 라우팅 인증

```
func SampleRouter(r *gin.Engine) {
	sample := r.Group("sample")
	{
		sample.GET("/list", handler.LoginCheckHandler(), Samples)
	}
}
```

-   그룹 라우팅 인증

```
func SampleRouter(r *gin.Engine) {
	sample := r.Group("sample")
	{
		sample.Use(handler.LoginCheckHandler())
		sample.GET("/list", Samples)
	}
}
```

### 파일 업로드

-   Single file

```
router.POST("/upload", func(c *gin.Context) {
    // single file
    file, _ := c.FormFile("file")
    log.Println(file.Filename)

    // Upload the file to specific dst.
    // c.SaveUploadedFile(file, dst)

    c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
})
```

-   Multiple files

```
router.POST("/upload", func(c *gin.Context) {
    // Multipart form
    form, _ := c.MultipartForm()
    files := form.File["upload[]"]

    for _, file := range files {
        log.Println(file.Filename)

        // Upload the file to specific dst.
        // c.SaveUploadedFile(file, dst)
    }
    c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
})
```

### 데이터 응답

-   Text Response

```
func HealthRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "SUCCESS")
	})
}
```

-   Json Response

```
func HealthRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"result":  "SUCCESS"})
	})
}
```

-   Json Pretty Response

```
func HealthRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.Render(http.StatusOK, render.IndentedJSON{Data: "SUCCESS"})
	})
}
```

### 트랜잭션

```
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
```

### 고루틴

```
go func() {
    // 처리내용
}()
```

### 테스팅

-   Target

```
package main

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
```

-   Test

```
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
```

### CRUD Sample

-   아래 파일 참고

```
/app/services/sample/SampleRouter.go 	// 라우팅(GET, POST, PUT, DELETE)
/app/services/sample/SampleService.go	// 서비스(CRUD)
```

### 푸시 & 빌드 & 배포 자동화
