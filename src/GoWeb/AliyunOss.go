package main

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
)

func initOss() *oss.Client {
	accessKeyId := strings.ReplaceAll(rc.getValue("accessKeyId"), "\"", "")
	accessKeySecret := strings.ReplaceAll(rc.getValue("accessKeySecret"), "\"", "")
	endpoint := strings.ReplaceAll(rc.getValue("endpoint"), "\"", "")
	if accessKeyId == "" || accessKeySecret == "" || endpoint == "" {
		fmt.Println("加载阿里云配置信息失败，请检查redis缓存")
	}
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("初始化阿里云OSS失败", err)
	}
	return client
}

var ossClient = initOss()
