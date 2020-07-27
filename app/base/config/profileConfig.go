package config

import (
	"fmt"
	"flag"
	"log"
	"os"

	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
	"github.com/spf13/viper"
)

type Profile struct {
	Env                             string
	ServerHost                      string
	ServerDomain                    string
	ServerPort                      int
	GanagosiDbDriver                string
	GanagosiDbServer                string
	GanagosiDbName                  string
	GanagosiDbUser                  string
	GanagosiDbPassword              string
	GanagosiDbOptionsPort           int
	GanagosiDbOptionsReadOnlyIntent bool
	GanagosiDbOptionsEncrypt        bool
	LogDbDriver                     string
	LogDbServer                     string
	LogDbName                       string
	LogDbUser                       string
	LogDbPassword                   string
	LogDbOptionsPort                int
	LogDbOptionsReadOnlyIntent      bool
	LogDbOptionsEncrypt             bool
	LoggerLevel                     string
}

var profile Profile

func ProfileInit() Profile {
	SetProfile(&profile)
	return profile
}

func SetEnv() string {

	env := "local"
	osEnv := os.Getenv("env")
	argEnv := flag.String("env", "", "a string")
	flag.Parse()

	if *argEnv != "" {
		env = *argEnv
	} else {
		if osEnv != "" {
			env = osEnv
		}
	}

	fmt.Println("profile : ", env)
	return env
}

func SetProfile(p *Profile) {
	env := SetEnv()

	// 프로파일 설정
	viper.SetConfigName("common")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./profile")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln(err)
	}

	viper.SetConfigName(env)
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.MergeInConfig()

	p.Env = env
	p.ServerHost = viper.GetString("server.host")
	p.ServerDomain = viper.GetString("server.domain")
	p.ServerPort = viper.GetInt("server.port")

	p.GanagosiDbDriver = viper.GetString("datasource.ganagosi.driver")
	p.GanagosiDbServer = viper.GetString("datasource.ganagosi.server")
	p.GanagosiDbName = viper.GetString("datasource.ganagosi.name")
	p.GanagosiDbUser = viper.GetString("datasource.ganagosi.user")
	p.GanagosiDbPassword = viper.GetString("datasource.ganagosi.password")
	p.GanagosiDbOptionsPort = viper.GetInt("datasource.ganagosi.options.port")
	p.GanagosiDbOptionsReadOnlyIntent = viper.GetBool("datasource.ganagosi.options.readOnlyIntent")
	p.GanagosiDbOptionsEncrypt = viper.GetBool("datasource.ganagosi.options.encrypt")

	p.LogDbDriver = viper.GetString("datasource.log.driver")
	p.LogDbServer = viper.GetString("datasource.log.server")
	p.LogDbName = viper.GetString("datasource.log.name")
	p.LogDbUser = viper.GetString("datasource.log.user")
	p.LogDbPassword = viper.GetString("datasource.log.password")
	p.LogDbOptionsPort = viper.GetInt("datasource.log.options.port")
	p.LogDbOptionsReadOnlyIntent = viper.GetBool("datasource.log.options.readOnlyIntent")
	p.LogDbOptionsEncrypt = viper.GetBool("datasource.log.options.encrypt")

	p.LoggerLevel = viper.GetString("logger.level")

	common.Logger.Info("p.ServerDomain : ", p.ServerDomain)
	common.Logger.Info("Profile configuration Complete!")
}

func GetProfile() Profile {
	return profile
}
