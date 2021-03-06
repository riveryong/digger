///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/common"
	"digger/models"
	"digger/services/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取配置
func GetConfigs(c *gin.Context) {
	configs, err := service.ConfigService().ListConfigs()
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	if configs["admin_user"] == "" {
		configs["admin_user"] = DefaultUser.Username
	}
	if configs["admin_password"] == "" {
		configs["admin_password"] = DefaultUser.Password
	}
	if configs["secret"] == "" {
		configs["secret"] = common.DefaultSecret
	}
	c.JSON(http.StatusOK, Success(configs))
}

// 获取配置
func UpdateConfig(c *gin.Context) {
	// 绑定请求数据
	reqData := &models.Config{}
	if err := c.ShouldBindJSON(reqData); err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	err := service.ConfigService().UpdateConfig(reqData.Key, reqData.Value)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}
