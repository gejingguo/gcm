package gcm

import (
	"errors"
	//"time"
	"time"
	"path"
	"os"
	"strconv"
	"fmt"
)

type ObjectHandle struct {
	ID int `json:"id"`			// 编号
	Name string `json:"name"`	// 名称
}


// 工程，一个完整任务集合
type Project struct {
	Name string 				// 工程名字
	Parallel bool 				// 是否支持并行

	lastResult ProjectResult	// 上次构建结果
	Mail MailInfo				// 邮件配置
	Running bool				// 是否正在构建

	// 主机列表
	hosts []*Host
	hostIdMap map[int]*Host
	hostNameMap map[string]*Host

	// 脚本列表
	scripts []*Script
	scriptIdMap map[int]*Script
	scriptNameMap map[string]*Script

	// 任务组列表
	groups []*TaskGroup
	groupIdMap map[int]*TaskGroup
	groupNameMap map[string]*TaskGroup

	// 任务列表
	tasks []*Task
	taskIdMap map[int]*Task
	taskNameMap map[string]*Task
}

func (pro* Project) GetLastBuildResult() *ProjectResult {
	return &pro.lastResult
}

// 获取主机
func (tg *Project) GetHostByID(id int) *Host {
	h, ok := tg.hostIdMap[id]
	if ok {
		return h
	}
	return nil
}

func (tg *Project) GetHostByName(name string) *Host {
	h, ok := tg.hostNameMap[name]
	if ok {
		return h
	}
	return nil
}

func (pro *Project) GetHostByHandle(handle ObjectHandle) *Host {
	if handle.ID > 0 {
		h, _ := pro.hostIdMap[handle.ID]
		return h
	}
	if handle.Name != "" {
		h, _ := pro.hostNameMap[handle.Name]
		return h
	}
	return nil
}

// 添加主机
func (tg *Project) AddHost(host *Host) error {
	if host == nil || host.ID == 0 || host.Name == "" {
		return errors.New("host param is nil")
	}
	if tg.GetHostByID(host.ID) != nil {
		return errors.New("host id has existed")
	}
	if tg.GetHostByName(host.Name) != nil {
		return errors.New("host name has existed")
	}

	tg.hosts = append(tg.hosts, host)
	tg.hostIdMap[host.ID] = host
	tg.hostNameMap[host.Name] = host

	return nil
}

// 获取所有主机列表
func (p *Project) GetAllHost() []*Host {
	return p.hosts
}

// 获取脚本
func (tg *Project) GetScriptByID(id int) *Script {
	h, ok := tg.scriptIdMap[id]
	if ok {
		return h
	}
	return nil
}

func (tg *Project) GetScriptByName(name string) *Script {
	h, ok := tg.scriptNameMap[name]
	if ok {
		return h
	}
	return nil
}

func (pro *Project) GetScriptByHandle(handle ObjectHandle) *Script {
	if handle.ID > 0 {
		s, _ := pro.scriptIdMap[handle.ID]
		return s
	}
	if handle.Name != "" {
		s, _ := pro.scriptNameMap[handle.Name]
		return s
	}
	return nil
}

// 添加主机
func (tg *Project) AddScript(script *Script) error {
	if script == nil || script.ID == 0 || script.Name == "" {
		return errors.New("script param is nil")
	}
	if tg.GetScriptByID(script.ID) != nil {
		return errors.New("script id has existed")
	}
	if tg.GetScriptByName(script.Name) != nil {
		return errors.New("script name has existed")
	}

	tg.scripts = append(tg.scripts, script)
	tg.scriptIdMap[script.ID] = script
	tg.scriptNameMap[script.Name] = script

	return nil
}

// 获取所有脚本列表
func (p *Project) GetAllScript() []*Script {
	return p.scripts
}

// 获取任务
func (tg *Project) GetTaskByID(id int) *Task {
	h, ok := tg.taskIdMap[id]
	if ok {
		return h
	}
	return nil
}

func (tg *Project) GetTaskByName(name string) *Task {
	h, ok := tg.taskNameMap[name]
	if ok {
		return h
	}
	return nil
}

func (pro *Project) GetTaskByHandle(handle ObjectHandle) *Task {
	if handle.ID > 0 {
		s, _ := pro.taskIdMap[handle.ID]
		return s
	}
	if handle.Name != "" {
		s, _ := pro.taskNameMap[handle.Name]
		return s
	}
	return nil
}

// 添加主机
func (tg *Project) AddTask(task *Task) error {
	if task == nil || task.ID == 0 || task.Name == "" {
		return errors.New("task param is nil")
	}
	if tg.GetTaskByID(task.ID) != nil {
		return errors.New("task id has existed")
	}
	if tg.GetTaskByName(task.Name) != nil {
		return errors.New("task name has existed")
	}

	group := tg.GetGroupByHandle(task.Group)
	if group == nil {
		return errors.New("task add group not found.")
	}

	group.AddTask(task)

	tg.tasks = append(tg.tasks, task)
	tg.taskIdMap[task.ID] = task
	tg.taskNameMap[task.Name] = task

	return nil
}

// 获取所有脚本列表
func (p *Project) GetAllTasks() []*Task {
	return p.tasks
}

// 获取任务
func (tg *Project) GetGroupByID(id int) *TaskGroup {
	h, ok := tg.groupIdMap[id]
	if ok {
		return h
	}
	return nil
}

func (tg *Project) GetGroupByName(name string) *TaskGroup {
	h, ok := tg.groupNameMap[name]
	if ok {
		return h
	}
	return nil
}

func (pro *Project) GetGroupByHandle(handle ObjectHandle) *TaskGroup {
	if handle.ID > 0 {
		s, _ := pro.groupIdMap[handle.ID]
		return s
	}
	if handle.Name != "" {
		s, _ := pro.groupNameMap[handle.Name]
		return s
	}
	return nil
}

// 添加主机
func (tg *Project) AddGroup(group *TaskGroup) error {
	if group == nil || group.ID == 0 || group.Name == "" {
		return errors.New("group param is nil")
	}
	if tg.GetGroupByID(group.ID) != nil {
		return errors.New("group id has existed")
	}
	if tg.GetGroupByName(group.Name) != nil {
		return errors.New("group name has existed")
	}

	tg.groups = append(tg.groups, group)
	tg.groupIdMap[group.ID] = group
	tg.groupNameMap[group.Name] = group

	return nil
}

// 获取所有脚本列表
func (p *Project) GetAllGroups() []*TaskGroup {
	return p.groups
}

//
func (pro *Project) Run(buildId int, proPath string) error {
	if pro.Running {
		return errors.New("project is running")
	}
	pro.Running = true
	defer func() {
		pro.Running = false
	}()

	buildLogPath := path.Join(proPath, pro.Name, "logs", strconv.Itoa(buildId))
	buildResultLog := path.Join(buildLogPath, "build.log")
	//fmt.Println(buildResultLog)
	err := os.MkdirAll(buildLogPath, 666)
	if err != nil {
		return err
	}
	buildFilePath := path.Join(proPath, pro.Name, "files")
	err = os.MkdirAll(buildFilePath, 666)
	if err != nil {
		return err
	}

	begTime := time.Now()
	rp := TaskRunTimeParam{false,buildLogPath, buildFilePath,buildId}
	// 并发处理
	if pro.Parallel && len(pro.groups) > 1 {
		result := make([]bool, len(pro.groups))
		retChan := make(chan taskRunRet)
		//
		for i := 0; i < len(pro.groups); i++  {
			go func(task int, c chan taskRunRet, rp TaskRunTimeParam) {
				err := pro.groups[i].Execute(rp)
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
				result[ret.task] = true
			}
		}
	} else {
		for _, t := range pro.groups  {
			t.Execute(rp)
		}
	}

	usedTime := time.Since(begTime)
	pro.lastResult.InitWithProject(pro, buildId, begTime, usedTime)
	err = pro.lastResult.OutToFile(buildResultLog)
	if err != nil {
		return err
	}

	err = pro.SendResultMail(pro.lastResult)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("send mail ret:", err)
	return nil
}

func (pro* Project) SendResultMail(result ProjectResult) error {
	if pro.Mail.User == "" {
		return errors.New("project mail setting failed")
	}
	content := result.OutToMailText()
	mail := NewTextMail(pro.Mail.Subject, pro.Mail.Body + "\n" +content)
	mail.SetAccount(pro.Mail.Host, pro.Mail.User, pro.Mail.Pass)
	mail.To = pro.Mail.To
	mail.CC = pro.Mail.CC
	err := mail.Send()
	if err != nil {
		return err
	}
	return nil
}