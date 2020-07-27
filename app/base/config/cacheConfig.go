package config

import (
	"io/ioutil"
	"strings"

	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
	cache "github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

var cacheMemory *cache.Cache

var root = "mapper"

func CacheInit() {
	cacheMemory = cache.New(0, 0)

	path := root
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileName := file.Name()

		if strings.Index(fileName, ".") > -1 {
			//common.Logger.Debug("path : ", path + "/" + fileName)
			setCache(path, fileName)
		} else {
			subPath := root + "/" + fileName
			subFiles, err := ioutil.ReadDir(subPath)
			if err != nil {
				panic(err)
			}

			for _, subFile := range subFiles {
				subFileName := subFile.Name()
				if strings.Index(subFileName, ".") > -1 {
					setCache(subPath, subFileName)
				}
			}
		}
	}
}

func setCache(path, fileName string) {
	name := strings.Split(fileName, ".")[0]
	ext := strings.Split(fileName, ".")[1]

	if ext == "yml" || ext == "yaml" {
		viper.SetConfigType(ext)
		viper.SetConfigName(name)
		viper.AddConfigPath("./" + path)

		err := viper.ReadInConfig()
		if err != nil {
			common.Logger.Error("error : ", err)
		} else {
			keys := viper.AllKeys()
			for _, key := range keys {
				SetCache(key, viper.GetString(key))
			}
		}
	}
}

func GetCache(key string) string {
	value, found := cacheMemory.Get(key)
	if found {
		return value.(string)
	} else {
		return ""
	}

}

func SetCache(key string, value string) {
	//common.Logger.Debug("key / value : ", key + " / " + value)
	cacheMemory.Set(key, value, cache.DefaultExpiration)
}
