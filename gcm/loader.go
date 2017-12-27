package gcm


import (
	"os"
	"io/ioutil"
	"encoding/json"
	//"errors"
	"fmt"
)

type TaskInfo struct {
	ID int `json:"id"`
	Name string `json:"name"`
	GroupID int `json:"group_id"`
	GroupName string `json:"group_name"`
	HostID int `json:"host_id"`
	HostName string `json:"host_name"`
	ScriptID int `json:"script_id"`
	ScriptName string `json:"script_name"`
	SrcFile string `json:"src_file"`
	DstFile string `json:"dst_file"`
}

type MailInfo struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	To []string `json:"to"`
	CC []string `json:"cc"`
	Subject string `json:"subject"`
	Body string `json:"body"`
}

type ProjectInfo struct {
	Name string `json:"name"`
	Parallel bool `json:"parallel"`

	Hosts []Host `json:"hosts"`
	Scripts []Script `json:"scripts"`
	Groups []TaskGroup `json:"groups"`
	Tasks []TaskInfo `json:"tasks"`

	Mail MailInfo `json:"mail"`
}

// 记载工程
func LoadProjectFromJson(file string) (*Project, error) {
	data, err := loadFile(file)
	if err != nil {
		return nil,err
	}

	info := &ProjectInfo{}
	err = json.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	pro := &Project{}
	pro.Name = info.Name
	pro.Parallel = info.Parallel
	pro.hosts = make([]*Host, 0, 10)
	pro.hostIdMap = make(map[int]*Host)
	pro.hostNameMap = make(map[string]*Host)
	pro.scripts = make([]*Script, 0, 10)
	pro.scriptIdMap = make(map[int]*Script)
	pro.scriptNameMap = make(map[string]*Script)
	pro.groups = make([]*TaskGroup, 0, 10)
	pro.groupIdMap = make(map[int]*TaskGroup)
	pro.groupNameMap = make(map[string]*TaskGroup)
	pro.tasks = make([]*Task, 0, 10)
	pro.taskIdMap = make(map[int]*Task)
	pro.taskNameMap = make(map[string]*Task)

	for _, h := range info.Hosts {
		host := &Host{}
		*host = h
		err := pro.AddHost(host)
		if err != nil {
			return nil, err
		}
		fmt.Printf("pro addhost %d.\n", host.ID)
	}

	for _, s := range info.Scripts {
		script := &Script{}
		*script = s
		err := pro.AddScript(script)
		if err != nil {
			return nil, err
		}
		fmt.Printf("pro addscript %d.\n", script.ID)
	}

	for _, g := range info.Groups {
		group := &TaskGroup{}
		*group = g
		err := pro.AddGroup(group)
		if err != nil {
			return nil, err
		}
		fmt.Printf("pro addgroup %d.\n", group.ID)
	}

	for _, t := range info.Tasks {
		task := &Task{}
		task.ID = t.ID
		task.Name = t.Name
		task.Pro = pro
		task.Host.ID = t.HostID
		task.Host.Name = t.HostName
		task.Group.ID = t.GroupID
		task.Group.Name = t.GroupName
		task.Script.ID = t.ScriptID
		task.Script.Name = t.ScriptName
		task.File.File = t.SrcFile
		task.File.Target = t.DstFile
		err := pro.AddTask(task)
		if err != nil {
			return nil, err
		}
		fmt.Printf("pro addtask %d.\n", task.ID)
	}

	pro.Mail = info.Mail

	return pro, nil
}

func loadFile(fileName string) ([]byte, error) {
	//定义变量
	var (
		open               *os.File
		file_data          []byte
		open_err, read_err error
	)
	//打开文件
	open, open_err = os.Open(fileName)
	if open_err != nil {
		return nil, open_err
	}
	//关闭资源
	defer open.Close()

	//读取所有文件内容
	file_data, read_err = ioutil.ReadAll(open)
	if read_err != nil {
		return nil, read_err
	}
	return file_data, nil
}