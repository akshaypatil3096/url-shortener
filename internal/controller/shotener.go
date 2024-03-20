package controller

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/akshaypatil3096/url-shortener/internal/dao"
	"github.com/akshaypatil3096/url-shortener/internal/model"
	"github.com/akshaypatil3096/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func ShortenerURL(c *gin.Context) {
	var req model.Request
	ctx := c.Request.Context()
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	rd := dao.CreateClient(1)
	defer rd.Close()

	value, err := rd.Get(ctx, c.ClientIP()).Result()
	if err == redis.Nil {
		rd.Set(ctx, c.ClientIP(), os.Getenv("API_QUOTA"), time.Minute*30)
	}

	valInt, _ := strconv.Atoi(value)
	if valInt <= 0 {
		limit, _ := rd.TTL(ctx, c.ClientIP()).Result()
		c.IndentedJSON(http.StatusServiceUnavailable, map[string]any{
			"error": "short url not found in the database",
			"limit": limit,
		})
	}

	_, err = url.ParseRequestURI(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	req.URL = utils.EnforceHTTP(req.URL)

	var id string
	if req.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = req.CustomShort
	}

	val, _ := rd.Get(ctx, id).Result()
	if val != "" {
		c.IndentedJSON(http.StatusForbidden, map[string]any{
			"error": "short url already exits ",
		})
	}

	if req.Expiry == 0 {
		req.Expiry = 24
	}

	if err := rd.Set(ctx, id, req.URL, req.Expiry).Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, map[string]any{
			"error": "failed to set url",
		})
	}

	resp := model.Response{
		URL:             req.URL,
		CustomShort:     "",
		Expiry:          req.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	rd.Decr(ctx, c.ClientIP())

	val1, _ := rd.Get(ctx, c.ClientIP()).Result()

	resp.XRateRemaining, _ = strconv.Atoi(val1)

	ttl, _ := rd.TTL(ctx, c.ClientIP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	c.IndentedJSON(http.StatusOK, map[string]any{
		"body": resp,
	})

}
