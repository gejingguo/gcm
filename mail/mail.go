package main

import (
	"github.com/gejingguo/gcm/gcm"
	"fmt"
)

func main() {

	mail := gcm.NewTextMail("测试标题", "测试内地反对声浪")
	mail.SetAccount("smtp.yeah.net:25", "gejingguo@yeah.net", "abc123")
	mail.AddReceiver("404318634@qq.com")
	mail.AddCopier("343586350@qq.com")
	mail.AddFile("D:\\gejingguo\\go\\src\\github.com\\gejingguo\\gcm\\log\\test_pro\\1\\build.log")
	mail.AddFile("D:\\迅雷下载\\FTABCISSetup.exe")
	err := mail.Send()
	if err != nil {
		fmt.Println(err)
		return
	}
}
