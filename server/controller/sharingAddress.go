package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"github.com/v2rayA/v2rayA-lib4/common"
	"github.com/v2rayA/v2rayA-lib4/db/configure"
	"github.com/v2rayA/v2rayA-lib4/server/service"
)

func GetSharingAddress(ctx *gin.Context) {
	var w configure.Which
	err := jsoniter.Unmarshal([]byte(ctx.Query("touch")), &w)
	if err != nil {
		common.ResponseError(ctx, logError("bad request"))
		return
	}
	addr, err := service.GetSharingAddress(&w)
	if err != nil {
		common.ResponseError(ctx, logError(err))
		return
	}
	common.ResponseSuccess(ctx, gin.H{"sharingAddress": addr})
}
