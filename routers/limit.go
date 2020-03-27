package routers

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"workApi/config"
	"workApi/dblayer"
	"workApi/models"
)

//检查频率
func checkRateLimit(ip, url, method string) bool {
	current := int(time.Now().Unix())
	currentStr := strconv.Itoa(current)
	limit, timeSet := getRateLimitConfig()
	allowanceStr, timestampStr := loadAllowance(ip, url, method)
	allowance, _ := strconv.Atoi(allowanceStr)
	timestamp, _ := strconv.Atoi(timestampStr)
	allowance += int(current-timestamp) * limit / timeSet
	if allowance > limit {
		allowance = limit
	}
	if allowance < 1 {
		saveAllowance(ip, url, method, "0", currentStr)
		//返回true 代表速率超过,进行错误输出
		return true
	} else {
		allowanceStr = strconv.Itoa(allowance - 1)
		saveAllowance(ip, url, method, allowanceStr, currentStr)
		//返回false 代表速率未超过
		return false
	}
}

//读取频率
func loadAllowance(ip, url, method string) (allowance, timestamp string) {
	res, err := dblayer.RedisGet(ip + url + method)
	if err != nil {
		currentStr := string(time.Now().Unix())
		defaultLimitInt, _ := getRateLimitConfig()
		defaultLimitStr := strconv.Itoa(defaultLimitInt)
		allowance, timestamp = defaultLimitStr, currentStr
	} else {
		kv := strings.Split(res.(string), "-")
		allowance, timestamp = kv[0], kv[1]
	}
	return

}

//保存
func saveAllowance(ip, url, method, allowance, current string) {
	err := dblayer.RedisSet(ip+url+method, allowance+"-"+current, 60)
	if err != nil {
		return
	}
}

func getRateLimitConfig() (limit, timeSet int) {
	limit = config.Cfg.Count
	timeSet = config.Cfg.Second
	return
}

func getDefaultIp(c *gin.Context) string {
	remoteAddr := c.Request.RemoteAddr
	if ip := c.Request.Header.Get("X-Real-Ip"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func limitRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		realIp := getDefaultIp(c)
		reqUrl := c.Request.URL.Path
		method := c.Request.Method
		if checkRateLimit(realIp, reqUrl, method) {
			//请求过于频繁
			res := models.NewDefaultResult()
			res.Code = 4002
			res.Msg = config.CodeType[res.Code]
			c.JSON(http.StatusForbidden, res)
			c.AbortWithStatus(http.StatusForbidden)
		}
		c.Next()
	}
}
