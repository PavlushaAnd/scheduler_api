package controllers

import (
	"encoding/json"
	"fmt"
	"scheduler_api/core"
	"scheduler_api/logger"
	"scheduler_api/models"
	"scheduler_api/utils"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Operations about Users
type CoreController struct {
	core.Core
}

type LoginParam struct {
	UserCode string `json:"user_code"`
	Password string `json:"password"`
}

// write comment for swagger
// @Title Login
// @tags users
// @Description login with user_code and password
// @Param   loginParam		body		LoginParam  true		"loginParam"
// @Success 200 {object} utils.JSONStruct{data=core.UserDetailsWithPwd}
// @Failure 400
// @Failure 500
// @router /user/login [post]
func (c *CoreController) Login() {
	userCode := ""
	ipAddr := core.ReadClientIP(c.Ctx.Request)
	loginTime := time.Now()

	loginParam := LoginParam{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &loginParam)
	if err != nil {
		c.WriteLoginLog(userCode, ipAddr, loginTime)

		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	currentUser, err := c.GetUserDetailsWithPwd(loginParam.UserCode)
	if err != nil {
		c.WriteLoginLog(userCode, ipAddr, loginTime)

		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on getting user %s - %s", loginParam.UserCode, err.Error())}
		c.ServeJSON()
		return
	}

	if currentUser == nil {
		c.WriteLoginLog(userCode, ipAddr, loginTime)

		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: fmt.Sprintf("error - user %s doesn't exist", loginParam.UserCode)}
		c.ServeJSON()
		return
	}

	if currentUser.Inactive {
		c.WriteLoginLog(userCode, ipAddr, loginTime)

		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: fmt.Sprintf("error - user %s is inactivated", loginParam.UserCode)}
		c.ServeJSON()
		return
	}

	userCode = loginParam.UserCode
	if utils.GetMd5StrWithSalt(loginParam.Password, loginParam.UserCode) != currentUser.Password {
		c.WriteLoginLog(userCode, ipAddr, loginTime)

		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "error - user_code or password is not correct", Data: nil}
		c.ServeJSON()
		return
	}

	c.WriteLoginLog(userCode, ipAddr, loginTime)

	token, expireTime, err := core.CreateUserToken(userCode, true)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: fmt.Sprintf("error on creating token - %s", err.Error())}
		c.ServeJSON()
		return
	}

	err = c.UpdateUserLoginInfo(userCode, ipAddr, loginTime, token)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: fmt.Sprintf("error on updating login information - %s", err.Error())}
		c.ServeJSON()
		return
	}

	if core.TOKENCACHE != nil {
		core.TOKENCACHE.Put(c.Ctx.Request.Context(), token, expireTime, time.Duration(core.TokenAutoTime)*time.Second)
	}

	tokenstring := fmt.Sprintf("Bearer %s", token)
	c.Ctx.ResponseWriter.Header().Set("Authorization", tokenstring)
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin,X-Requested-With,Authorization,Access-Control-Request-Method,Access-Control-Request-Headers,Host,Content-Type,Accept,if-modified-since,soapaction")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Credentials,Access-Control-Expose-Headers,Content-Type,Authorization")
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: currentUser}
	c.ServeJSON()
}

// write comment for swagger
// @Title Get user list
// @tags users
// @Description get user list
// @Param	pageindex	query	int	true	"page index"
// @Param	pagesize	query	int	true	"page size"
// @Success 200 {object} utils.JSONStruct{data=core.UserPage}
// @Failure 400
// @router /user/userlist [get]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
func (c *CoreController) GetUserList() {
	c.RequireLogin()
	if c.CurrentUserDetail.Role != "admin" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "error - permission denied"}
		c.ServeJSON()
		return
	}

	pageIndex, err := c.GetInt("pageindex")
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	pageSize, err := c.GetInt("pagesize")
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	o := orm.NewOrmUsingDB("default")

	userList, cnt, err := models.ListUser("", pageIndex, pageSize, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	usrArr := make([]core.UserDetails, 0)
	for _, usr := range userList {
		var isOnline bool
		lastLoginToken, err := c.GetUserTokenFromDB(usr.UserCode)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
			c.ServeJSON()
			return
		}
		if lastLoginToken != "" {
			exptime, err := c.GetUserTokenExpireTimeFromCache(lastLoginToken)
			if err != nil {
				c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
				c.ServeJSON()
				return
			}
			if diff := exptime - time.Now().Unix(); diff > int64(2700) {
				isOnline = true
			} else {
				isOnline = false
			}
		}
		usrArr = append(usrArr, core.UserDetails{
			Id:                usr.Id,
			UserCode:          usr.UserCode,
			PositionCode:      usr.PositionCode,
			UserName:          usr.UserName,
			IsOnline:          isOnline,
			EmailAddress:      usr.EmailAddress,
			PhoneNo:           usr.PhoneNo,
			HasUploadedPage:   usr.HasUploadedPage,
			HasRecognisedPage: usr.HasRecognisedPage,
			HasConfirmedPage:  usr.HasConfirmedPage,
			HasPostedPage:     usr.HasPostedPage,
			Role:              usr.Role,
			Inactive:          usr.Inactive,
			ColorBackground:   usr.ColorBackground,
			ColorText:         usr.ColorText,
		})
	}

	usrpage := core.UserPage{}
	usrpage.CurrentPage = pageIndex
	l := cnt % pageSize
	if l == 0 {
		usrpage.TotalPages = cnt / pageSize
	} else {
		usrpage.TotalPages = (cnt / pageSize) + 1
	}
	usrpage.TotalUsers = cnt
	usrpage.UsersInPage = usrArr

	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success", Data: usrpage}
	c.ServeJSON()
}

// write comment for swagger
// @Title add or update user
// @tags users
// @Description add or update user
// @Param	userDetail		body		core.UserDetailsWithPwd	true		"user detail"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /user/addorupd [post]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string false "With the bearer in front"
func (c *CoreController) AddOrUpdateUser() {
	c.RequireLogin()

	userDetailStr := string(c.Ctx.Input.RequestBody)
	logger.D("docdetail json:", userDetailStr)
	d := &core.UserDetailsWithPwd{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, d)
	if err != nil {
		logger.E("json.Unmarshal failed, err:", err)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: "Request body is not a valid json"}
		c.ServeJSON()
		return
	}

	if d.UserCode == "" {
		logger.E("User code cannot be empty")
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "User code cannot be empty"}
		c.ServeJSON()
		return
	}

	if c.CurrentUserDetail.Role != "admin" {
		//only admin can update the following fields
		if d.HasConfirmedPage != c.CurrentUserDetail.HasConfirmedPage ||
			d.HasPostedPage != c.CurrentUserDetail.HasPostedPage ||
			d.HasRecognisedPage != c.CurrentUserDetail.HasRecognisedPage ||
			d.HasUploadedPage != c.CurrentUserDetail.HasUploadedPage ||
			d.Inactive != c.CurrentUserDetail.Inactive ||
			d.UserCode != c.CurrentUserDetail.UserCode ||
			d.Role != c.CurrentUserDetail.Role ||
			d.PositionCode != c.CurrentUserDetail.PositionCode ||
			d.ColorBackground != c.CurrentUserDetail.ColorBackground ||
			d.ColorText != c.CurrentUserDetail.ColorText {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "error - permission denied"}
			c.ServeJSON()
			return
		}
	}

	o := orm.NewOrmUsingDB("default")

	//base64.StdEncoding.EncodeToString([]byte(d.Password))
	if d.Id == 0 {
		user := models.User{
			UserCode:          d.UserCode,
			UserName:          d.UserName,
			PositionCode:      d.PositionCode,
			EmailAddress:      d.EmailAddress,
			PhoneNo:           d.PhoneNo,
			HasUploadedPage:   d.HasUploadedPage,
			HasRecognisedPage: d.HasRecognisedPage,
			HasConfirmedPage:  d.HasConfirmedPage,
			HasPostedPage:     d.HasPostedPage,
			Inactive:          d.Inactive,
			Role:              d.Role,
			ColorText:         d.ColorText,
			ColorBackground:   d.ColorBackground,
			Password:          utils.GetMd5StrWithSalt(d.Password, d.UserCode),
			CreatedAt:         time.Now(),
			LastModified:      time.Now(),
			CreatorCode:       c.CurrentUserDetail.Role,
		}
		//add to db
		if d.Password == "" {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "Password cannot be empty"}
			c.ServeJSON()
			return
		}

		err = models.InsertUser(&user, o)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
			c.ServeJSON()
			return
		}
	} else {
		user := models.User{
			UserCode:          d.UserCode,
			UserName:          d.UserName,
			PositionCode:      d.PositionCode,
			EmailAddress:      d.EmailAddress,
			PhoneNo:           d.PhoneNo,
			HasUploadedPage:   d.HasUploadedPage,
			HasRecognisedPage: d.HasRecognisedPage,
			HasConfirmedPage:  d.HasConfirmedPage,
			HasPostedPage:     d.HasPostedPage,
			Inactive:          d.Inactive,
			Role:              d.Role,
			ColorText:         d.ColorText,
			ColorBackground:   d.ColorBackground,
			Password:          utils.GetMd5StrWithSalt(d.Password, d.UserCode),
			LastModified:      time.Now(),
			EditorCode:        c.CurrentUserDetail.Role,
		}
		//update to db
		user.Id = d.Id
		err = models.UpdateUserWithoutPwd(&user, o)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
			c.ServeJSON()
			return
		}
	}

	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success", Data: nil}
	c.ServeJSON()
}

// write comment for swagger
// @Title  delete user
// @tags users
// @Description delete user
// @Param	userCode		query		string	true		"user code"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /user/delete [delete]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
func (c *CoreController) DeleteUser() {
	c.RequireLogin()
	if c.CurrentUserDetail.Role != "admin" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "error - permission denied"}
		c.ServeJSON()
		return
	}

	delUser := c.GetString("userCode")
	if delUser == "" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "User code cannot be empty"}
		c.ServeJSON()
		return
	}

	if c.CurrentUserDetail.UserCode == delUser {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "error - cannot delete current user"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrmUsingDB("default")

	dbUser, err := models.GetUser(delUser, o)
	if (err != nil) || (dbUser == nil) {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot find user %s, err: - %s", delUser, err.Error())}
		c.ServeJSON()
		return
	}

	err = models.DelUser(dbUser, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot delete user %s, err: - %s", delUser, err.Error())}
		c.ServeJSON()
		return
	}

	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Delete User Success"}
	c.ServeJSON()
}

// write comment for swagger
// @Title  Modify password
// @tags users
// @Description Modify password
// @Param	ModifyPwd		body		core.ModifyPwd	true	"modify user password"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /user/updpasswd [post]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
func (c *CoreController) ModifyPassword() {
	c.RequireLogin()

	modifyPwdStr := string(c.Ctx.Input.RequestBody)
	logger.D("docdetail json:", modifyPwdStr)
	d := &core.ModifyPwd{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, d)
	if err != nil {
		logger.E("json.Unmarshal failed, err:", err)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: "Request body is not a valid json"}
		c.ServeJSON()
		return
	}

	//userCode := c.GetString("userCode")
	if d.UserCode == "" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "User code cannot be empty"}
		c.ServeJSON()
		return
	}

	//oldPassword := c.GetString("oldPassword")
	if d.OldPwd == "" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "Old password cannot be empty"}
		c.ServeJSON()
		return
	}

	//newPassword := c.GetString("newPassword")
	if d.NewPwd == "" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "New password cannot be empty"}
		c.ServeJSON()
		return
	}

	if c.CurrentUserDetail.UserCode != d.UserCode {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "error - cannot modify other user's password"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	dbUser, err := models.GetUser(d.UserCode, o)
	if err != nil || dbUser.Inactive {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot find user %s, err: - %s", d.UserCode, err.Error())}
		c.ServeJSON()
		return
	}

	if utils.GetMd5StrWithSalt(d.OldPwd, d.UserCode) != dbUser.Password {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "error - old password is not correct"}
		c.ServeJSON()
		return
	}

	dbUser.Password = utils.GetMd5StrWithSalt(d.NewPwd, d.UserCode) //base64.StdEncoding.EncodeToString([]byte(d.NewPwd))
	dbUser.EditorCode = c.CurrentUserDetail.Role
	dbUser.LastModified = time.Now()
	err = models.UpdateUserPwd(dbUser, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
		c.ServeJSON()
		return
	}

	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Modify Password Success"}
	c.ServeJSON()
}

// write comment for swagger
// @Title  reset password
// @Description reset password
// @tags users
// @Param	ModifyPwd		body		core.ModifyPwd	true	"reset user password"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /user/rstpasswd [post]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
func (c *CoreController) ResetPassword() {
	c.RequireLogin()
	modifyPwdStr := string(c.Ctx.Input.RequestBody)
	logger.D("docdetail json:", modifyPwdStr)
	d := &core.ModifyPwd{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, d)
	if err != nil {
		logger.E("json.Unmarshal failed, err:", err)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: "Request body is not a valid json"}
		c.ServeJSON()
	}

	if c.CurrentUserDetail.Role != "admin" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "error - permission denied"}
		c.ServeJSON()
		return
	}

	userCode := d.UserCode //c.GetString("userCode")
	if userCode == "" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "User code cannot be empty"}
		c.ServeJSON()
		return
	}

	newPassword := d.NewPwd //c.GetString("newPassword")
	if newPassword == "" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "New password cannot be empty"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrmUsingDB("default")

	dbUser, err := models.GetUser(userCode, o)
	if err != nil || dbUser.Inactive {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot find user %s, err: - %s", userCode, err.Error())}
		c.ServeJSON()
		return
	}

	dbUser.Password = utils.GetMd5StrWithSalt(newPassword, dbUser.UserCode) //base64.StdEncoding.EncodeToString([]byte(newPassword))
	dbUser.EditorCode = c.CurrentUserDetail.Role
	dbUser.LastModified = time.Now()
	err = models.UpdateUserPwd(dbUser, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
		c.ServeJSON()
		return
	}

	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Reset Password Success"}
	c.ServeJSON()
}
