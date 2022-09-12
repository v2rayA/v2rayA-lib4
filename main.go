package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/v2rayA/v2rayA-lib4/conf/report"
	_ "github.com/v2rayA/v2rayA-lib4/pkg/plugin/pingtunnel"
	_ "github.com/v2rayA/v2rayA-lib4/pkg/plugin/simpleobfs"
	_ "github.com/v2rayA/v2rayA-lib4/pkg/plugin/socks5"
	_ "github.com/v2rayA/v2rayA-lib4/pkg/plugin/ss"
	_ "github.com/v2rayA/v2rayA-lib4/pkg/plugin/ssr"
	_ "github.com/v2rayA/v2rayA-lib4/pkg/plugin/tcp"
	_ "github.com/v2rayA/v2rayA-lib4/pkg/plugin/tls"
	_ "github.com/v2rayA/v2rayA-lib4/pkg/plugin/trojanc"
	_ "github.com/v2rayA/v2rayA-lib4/pkg/plugin/ws"
	"github.com/v2rayA/v2rayA-lib4/pkg/util/log"
	"runtime"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	checkEnvironment()
	if runtime.GOOS == "linux" {
		checkTProxySupportability()
	}
	initConfigure()
	checkUpdate()
	hello()
	if err := run(); err != nil {
		log.Fatal("main: %v", err)
	}
}
