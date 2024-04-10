package core

import (
	"fmt"
	"scheduler_api/models"
	"scheduler_api/utils"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type UserDetailsWithPwd struct {
	Id                int    `json:"id"`
	UserCode          string `json:"user_code"`
	PositionCode      string `json:"position_code"`
	UserName          string `json:"user_name"`
	Inactive          bool   `json:"inactive"`
	PhoneNo           string `json:"phone_no"`
	EmailAddress      string `json:"email_address"`
	HasUploadedPage   bool   `json:"has_uploaded_page"`
	HasRecognisedPage bool   `json:"has_recognised_page"`
	HasConfirmedPage  bool   `json:"has_confirmed_page"`
	HasPostedPage     bool   `json:"has_posted_page"`
	Password          string `json:"password"`
	Role              string `json:"role"`
	ColorText         string `json:"color_text"`
	ColorBackground   string `json:"color_background"`
}

type UserDetails struct {
	Id                int    `json:"id"`
	UserCode          string `json:"user_code"`
	PositionCode      string `json:"position_code"`
	UserName          string `json:"user_name"`
	Inactive          bool   `json:"inactive"`
	PhoneNo           string `json:"phone_no"`
	EmailAddress      string `json:"email_address"`
	HasUploadedPage   bool   `json:"has_uploaded_page"`
	HasRecognisedPage bool   `json:"has_recognised_page"`
	HasConfirmedPage  bool   `json:"has_confirmed_page"`
	HasPostedPage     bool   `json:"has_posted_page"`
	Role              string `json:"role"`
	ColorText         string `json:"color_text"`
	ColorBackground   string `json:"color_background"`
}

type ModifyPwd struct {
	OldPwd   string `json:"old_pwd"`
	NewPwd   string `json:"new_pwd"`
	UserCode string `json:"user_code"`
}

type UserPage struct {
	TotalUsers  int           `json:"ToltalUsers" example:"1" format:"int"`
	TotalPages  int           `json:"ToltalPages" example:"1" format:"int"`
	CurrentPage int           `json:"CurrentPage" example:"1" format:"int"`
	UsersInPage []UserDetails `json:"Users"`
}

type Core struct {
	CurrentUserDetail UserDetails
	beego.Controller
}

func (c *Core) WriteLoginLog(loginUser, loginIp string, loginTime time.Time) error {
	o := orm.NewOrmUsingDB("default")

	loginLog := models.LoginLog{}
	loginLog.LoginUser = loginUser
	loginLog.LoginIp = loginIp
	loginLog.LoginTime = loginTime

	return models.InsertLoginLog(&loginLog, o)
}

func (c *Core) UpdateUserLoginInfo(loginUser, loginIp string, loginTime time.Time, token string) error {
	o := orm.NewOrmUsingDB("default")

	loginInfo, err := models.GetUserLoginInfo(loginUser, o)
	if err != nil {
		return fmt.Errorf("error on getting user login info - %s", err.Error())
	}

	if loginInfo == nil {
		newLoginInfo := models.UserLoginInfo{}
		newLoginInfo.UserCode = loginUser
		newLoginInfo.LastLoginIp = loginIp
		newLoginInfo.LastLoginTime = loginTime
		newLoginInfo.LastLoginToken = token

		err = models.InsertUserLoginInfo(&newLoginInfo, o)
		if err != nil {
			return fmt.Errorf("error on inserting user login info - %s", err.Error())
		}
		return nil
	}

	loginInfo.LastLoginIp = loginIp
	loginInfo.LastLoginTime = loginTime
	loginInfo.LastLoginToken = token

	err = models.UpdateUserLoginInfo(loginInfo, o)
	if err != nil {
		return fmt.Errorf("error on updating user login info - %s", err.Error())
	}

	return nil
}

func (c *Core) GetTokenFromHttpRequest() (string, error) {
	authorization := c.Ctx.Input.Header("Authorization")
	if len(authorization) == 0 {
		return "", fmt.Errorf("error - authorization on http header is null")
	}

	tokenSlice := strings.Split(authorization, " ")
	if len(tokenSlice) != 2 {
		return "", fmt.Errorf("error - authorization on http header is wrong")
	}

	if strings.TrimSpace(tokenSlice[0]) != "Bearer" {
		return "", fmt.Errorf("error - authorization type on http header is wrong")
	}

	if len(strings.TrimSpace(tokenSlice[1])) == 0 {
		return "", fmt.Errorf("error - bearer token is empty")
	}

	return tokenSlice[1], nil
}

func (c *Core) GetUserDetailsWithPwd(userCode string) (*UserDetailsWithPwd, error) {
	o := orm.NewOrmUsingDB("default")

	userDB, err := models.GetUser(userCode, o)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error on getting user details - %s", err.Error())
	}

	if userDB == nil {
		return nil, nil
	}

	UserDetails := UserDetailsWithPwd{}
	UserDetails.Id = userDB.Id
	UserDetails.UserCode = userDB.UserCode
	UserDetails.UserName = userDB.UserName
	UserDetails.PositionCode = userDB.PositionCode
	UserDetails.Inactive = userDB.Inactive
	UserDetails.PhoneNo = userDB.PhoneNo
	UserDetails.EmailAddress = userDB.EmailAddress
	UserDetails.HasUploadedPage = userDB.HasUploadedPage
	UserDetails.HasRecognisedPage = userDB.HasRecognisedPage
	UserDetails.HasConfirmedPage = userDB.HasConfirmedPage
	UserDetails.HasPostedPage = userDB.HasPostedPage
	UserDetails.Password = userDB.Password
	UserDetails.Role = userDB.Role
	UserDetails.ColorBackground = userDB.ColorBackground
	UserDetails.ColorText = userDB.ColorText
	return &UserDetails, nil
}

func (c *Core) GetUserDetails(userCode string) (*UserDetails, error) {
	userDetailsWithPwd, err := c.GetUserDetailsWithPwd(userCode)
	if err != nil {
		return nil, err
	}
	if userDetailsWithPwd == nil {
		return nil, nil
	}
	if userDetailsWithPwd.Inactive {
		return nil, fmt.Errorf("error - user is inactive")
	}

	UserDetails := UserDetails{}
	UserDetails.Id = userDetailsWithPwd.Id
	UserDetails.UserCode = userDetailsWithPwd.UserCode
	UserDetails.PositionCode = userDetailsWithPwd.PositionCode
	UserDetails.UserName = userDetailsWithPwd.UserName
	UserDetails.Inactive = userDetailsWithPwd.Inactive
	UserDetails.PhoneNo = userDetailsWithPwd.PhoneNo
	UserDetails.EmailAddress = userDetailsWithPwd.EmailAddress
	UserDetails.HasUploadedPage = userDetailsWithPwd.HasUploadedPage
	UserDetails.HasRecognisedPage = userDetailsWithPwd.HasRecognisedPage
	UserDetails.HasConfirmedPage = userDetailsWithPwd.HasConfirmedPage
	UserDetails.HasPostedPage = userDetailsWithPwd.HasPostedPage
	UserDetails.Role = userDetailsWithPwd.Role
	UserDetails.ColorBackground = userDetailsWithPwd.ColorBackground
	UserDetails.ColorText = userDetailsWithPwd.ColorText

	return &UserDetails, nil
}
func (c *Core) GetUserTokenFromDB(userCode string) (string, error) {
	o := orm.NewOrmUsingDB("default")

	loginInfo, err := models.GetUserLoginInfo(userCode, o)
	if err != nil {
		return "", fmt.Errorf("error on getting user %s login information from db - %s", userCode, err.Error())
	}

	if loginInfo == nil {
		return "", nil
	}

	return loginInfo.LastLoginToken, nil
}

func (c *Core) ValidateUserToken(userToken string) error {
	if len(userToken) == 0 {
		return fmt.Errorf("error - token on http header is null")
	}

	userCode, expiretime, err := GetUserExpireTimeFromToken(&userToken)
	if err != nil {
		if TOKENCACHE != nil {
			TOKENCACHE.Delete(c.Ctx.Request.Context(), userToken)
		}
		return fmt.Errorf("error on parse token - %s", err.Error())
	}

	//Just read user and client data
	currentUser, err := c.GetUserDetails(userCode)
	if err != nil {
		return fmt.Errorf("error on getting user information - %s", err.Error())
	}

	if currentUser == nil {
		return fmt.Errorf("error - user %s not found", userCode)
	}

	if currentUser.Inactive {
		return fmt.Errorf(("error - current user is inactive"))
	}

	lastLoginToken, err := c.GetUserTokenFromDB(userCode)
	if err != nil {
		return fmt.Errorf("error checking token in db - %s", err.Error())
	}
	if userToken != lastLoginToken {
		return fmt.Errorf("error - token is different from last login token")
	}

	if TOKENCACHE != nil {
		tmp, _ := TOKENCACHE.Get(c.Ctx.Request.Context(), userToken)
		if tmp != nil {
			expiretime = tmp.(int64)
		}
	}

	if expiretime <= 0 {
		return fmt.Errorf("error - token is expired, require to login")
	}

	nowTime := time.Now().Unix()
	if expiretime < nowTime {
		return fmt.Errorf("error - token is expired, require to login")
	}

	c.CurrentUserDetail = *currentUser

	return nil
}

func (c *Core) DeleteUserTokenInCache(userToken string) error {
	if TOKENCACHE != nil {
		TOKENCACHE.Delete(c.Ctx.Request.Context(), userToken)
	}
	return nil
}

func (c *Core) UpdateUserTokenExpireTimeInCache(userToken string) error {
	NewExpireTime := time.Now().Unix() + TokenAutoTime
	if TOKENCACHE != nil {
		TOKENCACHE.Put(c.Ctx.Request.Context(), userToken, NewExpireTime, time.Duration(TokenAutoTime)*time.Second)
	}

	return nil
}

// RequireLogin is used by majority of APIs
func (c *Core) RequireLogin() {
	userToken, err := c.GetTokenFromHttpRequest()
	if len(userToken) <= 0 {
		c.Ctx.Output.SetStatus(utils.Success)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: err.Error()}
		c.ServeJSON()
		c.StopRun()
	}

	err = c.ValidateUserToken(userToken)
	if err != nil {
		c.DeleteUserTokenInCache(userToken)

		c.Ctx.Output.SetStatus(utils.Success)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: err.Error()}
		c.ServeJSON()
		c.StopRun()
	}

	//If token is valid, regardless privileges is valid or not, update token's expiry time
	c.UpdateUserTokenExpireTimeInCache(userToken)

	tokenstring := fmt.Sprintf("Bearer %s", userToken)
	c.Ctx.ResponseWriter.Header().Set("Authorization", tokenstring)
}

// RequireLoginOnJWT is used by APIs for synchronization
func (c *Core) RequireLoginOnJWT() {
	token, err := c.GetTokenFromHttpRequest()
	if len(token) <= 0 {
		c.Ctx.Output.SetStatus(utils.Success)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: err.Error()}
		c.ServeJSON()
		c.StopRun()
	}
	defaultToken := "KYK6xK5lNLiT7IlzuBsGqXMpZbPtI5RKZTwzN45L27g="
	if defaultToken != token {
		c.Ctx.Output.SetStatus(utils.Success)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: err.Error()}
		c.ServeJSON()
		c.StopRun()
	}
}
