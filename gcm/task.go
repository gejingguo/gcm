package gcm

import (
	"os"
	//"path/filepath"
	//"fmt"
	//"io"
	"fmt"
	"path/filepath"
	"time"
	"errors"
	"path"
)

type TaskRunTimeParam struct {
	ConsoleOut bool 		// 是否输出到终端
	LogPath string			// 日志目录
	FilePath string 		// 传输文件，文件路径
	BuildID int				// 构建执行编号
}

// 远程传输文件
type FileHandle struct {
	File string 		// 原始文件，路径相对工程目录
	Target string 		// 目标路径，包含文件名, 绝对路径
}

// 远程任务
type Task struct {
	ID int
	Name string					// 名称
	Pro *Project				// 工程
	Group ObjectHandle			// 组ID
	Script ObjectHandle			// 命令
	Host ObjectHandle			// 主机
	File FileHandle				// 文件

	lastError error				// 上次操作结果
	lastUsedTime time.Duration		// 上次操作耗时
}

func (st *Task) Execute(runParam TaskRunTimeParam) error {
	st.lastError = nil
	st.lastUsedTime = 0
	var err error = nil
	begTime := time.Now()
	for {
		host := st.Pro.GetHostByHandle(st.Host)
		if host == nil {
			err = errors.New("SSHTask gethost failed")
			break
		}

		fileWriter := os.Stdout
		if !runParam.ConsoleOut {
			file := fmt.Sprintf("task_%d.log", st.ID)
			fileWriter, err = os.OpenFile(filepath.Join(runParam.LogPath, file), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 666)
			if err != nil {
				break
			}
			defer fileWriter.Close()
		}

		// 文件
		if st.File.File != "" && st.File.Target != "" {
			file := path.Join(runParam.FilePath, st.File.File)
			err = ScpFile(file, st.File.Target, host.Host, host.User, host.Pass)
			fmt.Fprintf(fileWriter, "scpfile %s %s, ret:%v\n", file, st.File.Target, err)
			if err != nil {
				break
			}
		}

		// 指令
		if st.Script.ID != 0 || st.Script.Name != "" {
			script := st.Pro.GetScriptByHandle(st.Script)
			if script == nil {
				err = errors.New("SSHTask getscript failed")
				break
			}

			fmt.Fprintf(fileWriter, "===task(%d),time(%s) begin ...\n", st.ID, begTime.Format("2006-01-02 15:04:05"))
			err = SshCmd(host, script.Body, fileWriter)
		}
		break
	}

	st.lastError = err
	st.lastUsedTime = time.Since(begTime)

	return err
}