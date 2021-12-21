package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ropon/logger"
	"net/http"
	"time"
)

const (
	StatusOk      = "OK"
	StatusErr     = "Error"
	MaxLogDataLen = 2048
)

const (
	ErrCodeParamInvalid = 2001
	ErrCodeGeneralFail  = 2002
)

func ColorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return logger.Green
	case code >= 300 && code <= 399:
		return logger.White
	case code >= 400 && code <= 499:
		return logger.Yellow
	default:
		return logger.Red
	}
}

func ColorForMethod(method string) string {
	switch {
	case method == "GET":
		return logger.Blue
	case method == "POST":
		return logger.Cyan
	case method == "PUT":
		return logger.Yellow
	case method == "DELETE":
		return logger.Red
	case method == "PATCH":
		return logger.Green
	case method == "HEAD":
		return logger.Magenta
	case method == "OPTIONS":
		return logger.White
	default:
		return logger.Reset
	}
}

func ColorForReset() string {
	return logger.Reset
}

func GetOffsetAndLimit(PageSize, PageNum int64) (int64, int64) {
	if PageSize < 1 {
		PageSize = 10
	}
	if PageNum < 1 {
		PageNum = 1
	}
	offset := PageSize * (PageNum - 1)
	limit := PageSize
	return offset, limit
}

func Cuts(s string, n int) string {
	if len(s) > n {
		return s[:n]
	} else {
		return s
	}
}

func GinOKRsp(c *gin.Context, data interface{}, desc interface{}) {
	GinRsp(c, http.StatusOK, gin.H{"status": StatusOk, "data": data, "description": desc})
}

func GinRsp(c *gin.Context, statusCode int, obj interface{}) {
	requestData := GetRequestData(c)
	objData := fmt.Sprintf("%+v", obj)
	clientIP := c.ClientIP()
	method := c.Request.Method
	statusColor := ColorForStatus(statusCode)
	methodColor := ColorForMethod(method)
	resetColor := ColorForReset()
	userEmail := c.Request.Header.Get("user_email")
	logger.Info("[GIN-RSP] %s%s%s %s%d%s %s [ip:%s] [user_email:%s] [rsp:%s]",
		methodColor, method, resetColor,
		statusColor, statusCode, resetColor,
		Cuts(requestData, MaxLogDataLen),
		clientIP,
		userEmail,
		Cuts(objData, MaxLogDataLen))
	c.JSON(statusCode, obj)
}

func GinErrRsp(c *gin.Context, errCode interface{}, errMsg interface{}) {
	GinRsp(c, http.StatusOK, gin.H{"status": StatusErr, "data": errCode, "description": errMsg})
}

func GetRequestData(c *gin.Context) string {
	var requestData string
	method := c.Request.Method
	if method == "GET" || method == "DELETE" {
		requestData = c.Request.RequestURI
	} else {
		_ = c.Request.ParseForm()
		requestData = fmt.Sprintf("%s [%s]", c.Request.RequestURI, c.Request.Form.Encode())
	}
	return requestData
}

func StrInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}

// Intersect 取两个切片的交集
func Intersect(slice1 []string, slice2 []string) []string {
	m := make(map[string]int)
	n := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}
	return n
}

// FormatTime 格式化时间为指定格式字符串
func FormatTime(t time.Time, args ...string) string {
	layout := "2006-01-02 15:04:05"
	if len(args) == 1 {
		layout = args[0]
	}
	return t.Format(layout)
}
