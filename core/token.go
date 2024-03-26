package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/cache"
	beego "github.com/beego/beego/v2/server/web"
)

// ***important: key  length must be 32 24 16 bit
const Key = "8f95b53c3c0f598266b801a4b621423d"
const LenghtOfUserCode = 8

var (
	TOKENCACHE    cache.Cache
	TokenAutoTime int64
)

func init() {
	TOKENCACHE, _ = cache.NewCache("memory", `{"interval":60}`)

	if TOKENCACHE == nil {
		fmt.Println("Init token memcache failed, create file cache!")
		TOKENCACHE, _ = cache.NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":2,"EmbedExpiry":0}`)
	}

	//auto login expires time
	TokenAutoTime, _ = beego.AppConfig.Int64("TokenAutoTime")

	if TokenAutoTime <= 0 {
		fmt.Println("Read token expire time error")
		//set default value, one hour - 3600 seconds
		TokenAutoTime = 3600
	}
}

/*
Get user code and expire time from token. If token is valid, need to check need updating.
returned: 1. User Code; 2. ExpireTime; 3. Error
*/
func GetUserExpireTimeFromToken(token *string) (string, int64, error) {

	//Decrypt token string
	str, err := Decrypt(token)
	if err != nil {
		return "", 0, err
	}

	rs := []rune(str)

	//get user id
	userCode := string(rs[:LenghtOfUserCode])
	userCode = strings.Replace(userCode, " ", "", -1)

	//get Token expired time
	var expiretime int64
	strTokenTime := string(rs[LenghtOfUserCode:])
	expiretime, err = strconv.ParseInt(strTokenTime, 10, 64)
	if err != nil {
		return "", 0, err
	} else {
		return userCode, expiretime, err
	}
}

// AES Decrypt function
func Decrypt(decryptString *string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(*decryptString)
	if err != nil {
		return "", err
	}

	keyByte := []byte(Key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	// get blockMode
	blockMode := cipher.NewCBCDecrypter(block, keyByte[:blockSize])

	origData := make([]byte, len(b))

	// pre-process the data to avoid panic to make sure the server keep running
	if len(b)%blockSize != 0 || len(origData) < len(b) || len(b) == 0 {
		return "", errors.New("incorrect token")
	}

	blockMode.CryptBlocks(origData, b)
	origData = ZeroUnPadding(origData)
	// remove the blank character in bytes[]
	origData = bytes.Trim(origData, "\x00")
	return string(origData), nil
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func CreateUserToken(userCode string, auto bool) (string, int64, error) {
	if userCode == "" {
		return "", 0, fmt.Errorf("error - user code is blank")
	}
	//The length of user code could be any
	//In this project, user code is analysis code of JCCOMP, the length of this field is 8
	if len(userCode) > LenghtOfUserCode {
		return "", 0, fmt.Errorf("error - the length of user code is %d, more than 8", len(userCode))
	}

	var str string
	var tokentime int64
	//is auto login
	if auto {
		//Generate token with user id and auto login expires time
		tokentime = time.Now().Unix() + TokenAutoTime
	} else {
		//make the token expired arfter 24 hours
		dd, _ := time.ParseDuration("24h")
		tokentime = time.Now().Add(dd).Unix()
	}
	//If length of user code is less than 8, padding space on the right and making its length to 8
	str = fmt.Sprintf("%-*s", LenghtOfUserCode, userCode) + strconv.FormatInt(tokentime, 10)

	//encrypt string to token
	token, err := Encryption(str)
	return token, tokentime, err
}

// AES Encrypt function
func Encryption(encryptionString string) (string, error) {
	//key := beego.AppConfig.String("JWTSecret")
	plaintext := []byte(encryptionString)
	keyByte := []byte(Key)

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	plaintext = ZeroPadding(plaintext, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, keyByte[:blockSize])

	crypted := make([]byte, len(plaintext))

	blockMode.CryptBlocks(crypted, plaintext)

	// base 64 after encrypt
	return base64.URLEncoding.EncodeToString(crypted), err
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}
