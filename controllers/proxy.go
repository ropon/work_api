/*
Author:Ropon
Date:  2020-12-11
*/
package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/ropon/logger"
	"github.com/ropon/work_api/conf"
	"io/ioutil"
	"net/http"
	"time"
)

var client *http.Client

func HttpProxy(c *gin.Context) {
	server := c.Params.ByName("server")
	proxyUrl := *(c.Request.URL)
	proxyUrl.Scheme = "http"
	if host, ok := conf.Cfg.External[server]; ok {
		proxyUrl.Host = host
	} else {
		logger.Error("服务%s不在转发配置中", server)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	proxyURL := proxyUrl.String()
	requestData, err := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(requestData))
	proxyRequest, err := http.NewRequest(c.Request.Method, proxyURL, c.Request.Body)
	if err != nil {
		logger.Error("新建转发请求失败，错误信息：%s", err.Error())
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	proxyRequest.Header = c.Request.Header
	if client == nil {
		client = &http.Client{Timeout: time.Duration(5*60) * time.Second}
	}
	resp, err := client.Do(proxyRequest)
	logger.Debug("proxy: [%s] resp %+v err %v", proxyRequest.URL.String(), resp, err)
	if err != nil {
		logger.Error("转发请求失败，错误信息：%s", err.Error())
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	logger.Info("%s %s [url:%s,body:%s]", c.Request.Method, proxyURL, proxyUrl.RawQuery, requestData)
	dataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("读取转发请求返回数据失败，错误信息：%s", err.Error())
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	for respHeaderKey, _ := range resp.Header {
		c.Writer.Header().Add(respHeaderKey, resp.Header.Get(respHeaderKey))
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), dataBytes)
}
