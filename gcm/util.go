package gcm

import (
	"fmt"
	"net"
	"golang.org/x/crypto/ssh"
	"io"
	//"os"
	//"path/filepath"
	"github.com/pkg/sftp"
	"os"
)

// 远程执行ssh命令
func SshCmd(host *Host, cmd string, out io.Writer) error {
	client, err := ssh.Dial("tcp", host.Host, &ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{ssh.Password(host.Pass)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		fmt.Println("ssh dial err,", err)
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("client newsession err,", err)
		return err
	}
	defer session.Close()

	session.Stdout = out
	session.Stderr = out
	//session.Stdin = os.Stdin

	err = session.Run(cmd)
	if err != nil {
		fmt.Println("session run cmd, err:", err)
		return err
	}
	return nil
}

type taskRunRet struct {
	task int
	err error
}

type executor interface {
	Execute(rp TaskRunTimeParam) error
}

// 并发执行任务
func RunParallelTask(tasks []executor, rp TaskRunTimeParam) error {
	result := make([]bool, len(tasks))
	retChan := make(chan taskRunRet)
	//
	for i := 0; i < len(tasks); i++  {
		go func(task int, c chan taskRunRet, rp TaskRunTimeParam) {
			err := tasks[i].Execute(rp)
			eret := taskRunRet{task, err}
			c <- eret
		}(i, retChan, rp)
	}
	// 等待结果
	for {
		finish := true
		for _, v := range result  {
			if !v {
				finish = false
				break
			}
		}
		if finish {
			break
		}
		//var ret execResult
		select {
		case ret := <- retChan:
			//if ret.err != nil {
			//	return ret.err
			//}
			result[ret.task] = true
		}
	}
	return nil
}

// 传输文件, 路径使用绝对路径
func ScpFile(srcFile string, dstFile string, host string, user string, passwd string) error {
	//addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", host, &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(passwd)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		return err
	}
	defer client.Close()
	// create sftp client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}

	sf, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := sftpClient.Create(dstFile)
	if err != nil {
		return err
	}
	defer df.Close()

	buf := make([]byte, 1024)
	for {
		n, err := sf.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		_, err = df.Write(buf[0:n])
		if err != nil {
			return err
		}
	}
	return nil
}