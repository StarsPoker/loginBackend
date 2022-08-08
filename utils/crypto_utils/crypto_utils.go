package crypto_utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
)

func GetMd5(input string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetToken(min int, max int) string {
	token := strconv.Itoa(rand.Intn(max-min) + min)
	return token
}
