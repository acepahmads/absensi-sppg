package handler

import (
	"absensi-sppg/pkg/utils"

	"github.com/gin-gonic/gin"
)

func RenderDashboard(c *gin.Context, desktopTpl, mobileTpl string) {
	device, _ := c.Get("device_type")
	claims, ok := utils.GetClaims(c)
	if !ok {
		return
	}

	data := gin.H{
		"email":  claims.Email,
		"device": device,
	}

	if device == utils.DeviceMobile {
		c.HTML(200, mobileTpl, data)
		return
	}

	c.HTML(200, desktopTpl, data)
}
