package utils

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

const (
	Success               = 200
	SuccessExisting       = 210
	SuccessNoContent      = 204
	ErrorParameter        = 400
	ErrorUnLogin          = 401
	ErrorForbidden        = 403
	ErrorNodata           = 404
	ErrorLogic            = 409
	ErrorExpire           = 410
	ErrorParseJson        = 422
	ErrorParseXML         = 423
	ErrorFrequency        = 490 // do something too frequency
	ErrorWebsocketRead    = 4100
	ErrorWebsocketWrite   = 4101
	ErrorService          = 500
	ErrorWebsocketUpgrade = 570 // error
	ErrorCache            = 800
	ErrorDB               = 900
	ErrorNetWork          = 990 // sending http request error
	ErrorHttp             = 991 // sending http request error

	ErrorPackage           = 1000
	ErrorThirdParty        = 1100
	PartErrorOfPostAutoPay = 1200

	//for last transaction
	DoNothing           = 205
	NeedWebSocket       = 206
	NeedPayAgain        = 207
	ShowTransactionInfo = 208
	GoToTransaction     = 209

	ErrorMWLogger = 2000
	ErrorMWEmail  = 2001

	ErrorLoyaltyCreation = 3001

	BeegoNoData = "<QuerySeter> no row found"
)

// use the beginning byte of file to determine the file type
func GetFileType(file []byte) string {
	var filetype string
	if file[0] == 0xff && file[1] == 0xd8 {
		filetype = "jpg"
	} else if file[0] == 0x89 && file[1] == 0x50 && file[2] == 0x4e && file[3] == 0x47 {
		filetype = "png"
	} else if file[0] == 0x25 && file[1] == 0x50 && file[2] == 0x44 && file[3] == 0x46 {
		filetype = "pdf"
	} else {
		filetype = "unknown"
	}
	return filetype
}

// to check if the file exist in local
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// to check if the directory exist in local
func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func GetCurrentDir() string {
	curdir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return curdir
}

func GetUploadDir() string {
	uploadDir := beego.AppConfig.String("UploadDir")
	if uploadDir[0] == '.' {
		curPath := GetCurrentDir()
		pathAfter, find := strings.CutPrefix(uploadDir, ".")
		if find {
			uploadDir = curPath + pathAfter
		} else {
			uploadDir = curPath + "/Uploads/"
		}
	}

	if uploadDir != "" && !DirExists(uploadDir) {
		err := os.Mkdir(uploadDir, 0777)
		if err != nil {
			return ""
		}
	}
	return uploadDir
}

func Unixtimetodatetime(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

func GetMd5StrWithSalt(pwd string, salt string) string {
	a := md5.Sum([]byte(pwd))
	b := []byte(strings.ToLower(salt))
	c := append(a[:], b[:]...)
	finalStr := fmt.Sprintf("%x", md5.Sum(c))
	return finalStr
}

type JSONStruct struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type OptionsItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
