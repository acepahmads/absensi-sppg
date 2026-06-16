package handler

import (
	"absensi-sppg/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ComplainHandler struct {
	complainService service.ComplainService
}

func NewComplainHandler(service service.ComplainService) *ComplainHandler {
	return &ComplainHandler{
		complainService: service,
	}
}

func (h *ComplainHandler) GetAll(c *gin.Context) {
	fmt.Println("ComplainHandler GetAll")
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	status := c.Query("status")                 // optional
	types := c.Query("type")                    // optional
	parameter_flat := c.Query("parameter_flat") // optional
	parameter_zero := c.Query("parameter_zero") // optional
	ispu := c.Query("ispu")
	mutu_air := c.Query("mutu_air")                     // optional
	validate := c.Query("validate")                     // optional
	accuweather_values := c.Query("accuweather_values") // optional
	perPageStr := c.Query("per_page")
	per_page, err := strconv.Atoi(perPageStr)
	if err != nil {
		per_page = 0 // default value
	}
	all_error := c.Query("all_error")

	result, err := h.complainService.GetComplains(page, status, types, parameter_flat, parameter_zero, ispu, mutu_air, validate, accuweather_values, per_page, all_error)
	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed load complain data",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
