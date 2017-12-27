package main

import (
	"github.com/gejingguo/gcm/gcm"
	//"os"
	"fmt"
	"os"
	"time"
)

func main() {
	/*
	hosts := [...]gcm.HostInfo {
		{"111.231.64.47:22", "ubuntu", "Gejingguo@004", "/home/ubuntu" },
		{"111.231.64.47:22", "root", "Gejingguo@004", "/root" },
	}
	//task.SshCmd(host, "echo 'hello'; whoami; pwd; uname -a", os.Stdout)
	fmt.Println(hosts)
	task := gcm.NewSshTask(1, "task1", false, "whoami;pwd;uname -a", hosts[:])
	runParam := gcm.TaskRunTimeParam{"D:\\gejingguo\\go\\src\\github.com\\gejingguo\\gcm", 1}
	err := task.Execute(runParam)
	if err != nil {
		fmt.Println("task exec err:", err)
		os.Exit(1)
	}
	fmt.Println("task exec ok.")
	*/
	pro, err := gcm.LoadProjectFromJson(".\\project.json")
	if err != nil {
		fmt.Println("pro load failed, err:", err)
		os.Exit(1)
	}

	err = pro.Run(1, "D:\\gejingguo\\go\\src\\github.com\\gejingguo\\gcm")
	if err != nil {
		fmt.Println("pro run failed, err:", err)
		os.Exit(1)
	}

	fmt.Println("pro run ok.", pro.GetLastBuildResult())
	time.Now().Nanosecond()
}
