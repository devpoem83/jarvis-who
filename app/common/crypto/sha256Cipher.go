package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
)

func Sha256Test() {
	var compareText = "89cdf965f08d1331362568c8ce92a41ea464b55a1e5503acd0e0d2679eaca2ce"
	var plainText = "faith83!"
	common.Logger.Debug("plainText : ", plainText)

	encryptedData := EncryptSha256(plainText)
	common.Logger.Debug("encryptedData : ", encryptedData)
	common.Logger.Debug("compareText   : ", compareText)
}

func EncryptSha256(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	md := hash.Sum(nil)
	text = hex.EncodeToString(md)
	return text
}