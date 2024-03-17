package controller

import (
	"net/http"
	"net/url"

	"github.com/akshaypatil3096/url-shortener/internal/model"
	"github.com/akshaypatil3096/url-shortener/internal/utils"

	"github.com/gin-gonic/gin"
)

func ShortenerURL(c *gin.Context) {
	var req model.Request

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	_, err := url.ParseRequestURI(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	utils.EnforceHTTP(req.URL)

}
