package controller

import (
	"net/http"

	"github.com/akshaypatil3096/url-shortener/internal/dao"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func ResolveURL(c *gin.Context) {
	ctx := c.Request.Context()
	url := c.Param("url")

	dConn := dao.CreateClient(0)
	defer dConn.Close()

	value, err := dConn.Get(ctx, url).Result()
	if err == redis.Nil {
		c.IndentedJSON(http.StatusNotFound, map[string]any{
			"error": "short url not found in the database",
		})
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, map[string]any{
			"error": err,
		})
	}

	_ = dConn.Incr(ctx, "counter")

	c.Redirect(301, value)
}
