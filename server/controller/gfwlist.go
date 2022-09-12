package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/v2rayA/v2rayA-lib4/common"
	"github.com/v2rayA/v2rayA-lib4/core/v2ray/asset/dat"
)

func PutGFWList(ctx *gin.Context) {
	localGFWListVersion, err := dat.CheckAndUpdateGFWList()
	if err != nil {
		common.ResponseError(ctx, logError(err))
		return
	}
	common.ResponseSuccess(ctx, gin.H{
		"localGFWListVersion": localGFWListVersion,
	})
}
