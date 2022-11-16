package controller

import (
	"github.com/gin-gonic/gin"
	"x-ui/logger"
	"x-ui/web/service"
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
    g.GET("/update",a.subscribe)

}

func (a *SubscribeController) subscribe(c *gin.Context){
    text := a.subscribeService.Publish()
    _, err := c.Writer.WriteString(text)
    if err != nil {
        logger.Debug("返回失败")
    }
}
