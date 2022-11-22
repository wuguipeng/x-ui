package controller

import (
	"encoding/base64"
	"x-ui/logger"
	"x-ui/web/job"
	"x-ui/web/service"

	"github.com/gin-gonic/gin"
)

type SubscribeController struct {
	subscribeService service.SubscribeService
}

func NewSubscribeController(g *gin.RouterGroup) *SubscribeController {
	a := &SubscribeController{}
	a.initRouter(g)
	return a
}

func (a *SubscribeController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/subscribe")
	g.GET("/v2", a.vmess)
	g.GET("/clash", a.clash)
}

func (a *SubscribeController) vmess(c *gin.Context) {
	text, _, _ := a.subscribeService.Publish()
	// 再次编码，返回
	_, err := c.Writer.WriteString(base64.StdEncoding.EncodeToString([]byte(text)))
	if err != nil {
		logger.Debug("返回失败")
	}
}

func (a *SubscribeController) clash(c *gin.Context) {
	text, newPort, oldPort := a.subscribeService.Clash()
	_, err := c.Writer.WriteString(text)
	if err != nil {
		logger.Debug("返回失败")
	}
	job.NewStatsNotifyJob().UpdatePortNotify(newPort, oldPort)
}
