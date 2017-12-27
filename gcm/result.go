package gcm

import (
	"time"
	"fmt"
	"encoding/json"
	"os"
	"bytes"
)

// 任务结果
type TaskResult struct {
	ID int `json:"id"`								//
	Name string `json:"name"`
	UsedSeconds float64 `json:"used_seconds"`		// 耗时
	err string `json:"err"`							// 出错信息
}

// 组结果
type GroupResult struct {
	ID int `json:"id"`
	Name string `json:"name"`
	UsedSeconds float64 `json:"used_seconds"`
	ErrTaskNum int `json:"err_task_num"`
}

// 构建结果报告
type ProjectResult struct {
	BuildID int	`json:"build_id"`									// 构建编号
	BuildTime string `json:"build_time"`							// 构建时间
	UsedTime float64 `json:"used_time"`								// 耗时
	SuccessedTasks []TaskResult `json:"successed_tasks"`			// 成功任务列表
	FailedTasks []TaskResult `json:"failed_tasks"`					// 失败任务列表
	SuccessedGroups []GroupResult `json:"successed_groups"`			// 成功的任务组列表
	FailedGroups []GroupResult `json:"failed_groups"`				// 失败的任务组列表
}

// 初始化
func (r *ProjectResult) InitWithProject(pro* Project, buildId int, buildTime time.Time, usedTime time.Duration) {
	if pro == nil {
		return
	}
	r.BuildID = buildId
	r.BuildTime = buildTime.Format("2006-01-02 15:04:05")
	r.UsedTime = usedTime.Seconds()
	r.SuccessedGroups = nil
	r.SuccessedTasks = nil
	r.FailedGroups = nil
	r.FailedTasks = nil

	for _, t := range pro.tasks {
		tr := TaskResult{}
		tr.ID = t.ID
		tr.Name = t.Name
		tr.UsedSeconds = t.lastUsedTime.Seconds()
		if t.lastError != nil {
			tr.err = fmt.Sprintf("%v", t.lastError)
			if r.FailedTasks == nil {
				r.FailedTasks = make([]TaskResult, 0, 10)
			}
			r.FailedTasks = append(r.FailedTasks, tr)
		} else {
			if r.SuccessedTasks == nil {
				r.SuccessedTasks = make([]TaskResult, 0, 10)
			}
			r.SuccessedTasks = append(r.SuccessedTasks, tr)
		}
	}

	// 筛选任务组
	for _, g := range pro.groups {
		gr := GroupResult{}
		gr.ID = g.ID
		gr.Name = g.Name
		gr.UsedSeconds = g.LastUsedTime.Seconds()
		gr.ErrTaskNum = g.LastErrTaskNum
		if g.LastErrTaskNum == 0 {
			if r.SuccessedGroups == nil {
				r.SuccessedGroups = make([]GroupResult, 0, 10)
			}
			r.SuccessedGroups = append(r.SuccessedGroups, gr)
		} else {
			if r.FailedGroups == nil {
				r.FailedGroups = make([]GroupResult, 0, 10)
			}
			r.FailedGroups = append(r.FailedGroups, gr)
		}
	}
}

// 结果输出到文件中
func (r *ProjectResult) OutToFile(file string) error {
	writer, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 666)
	if err != nil {
		return err
	}
	defer writer.Close()
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (r *ProjectResult) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("构建结果概要\n")
	buf.WriteString(fmt.Sprintf("构建编号: %d\n", r.BuildID))
	buf.WriteString(fmt.Sprintf("构建时间: %s\n", r.BuildTime))
	buf.WriteString(fmt.Sprintf("构建耗时: %f(S)\n", float32(r.UsedTime)))
	buf.WriteString(fmt.Sprintf("构建任务组数量:%d, 成功数量:%d, 失败数量:%d\n", len(r.SuccessedGroups) + len(r.FailedGroups), len(r.SuccessedGroups), len(r.FailedGroups)))
	buf.WriteString(fmt.Sprintf("构建任务数量:%d, 成功数量:%d, 失败数量:%d\n", len(r.SuccessedTasks) + len(r.FailedTasks), len(r.SuccessedTasks), len(r.FailedTasks)))
	return string(buf.Bytes())
}

// 邮件文本输出
func (r *ProjectResult) OutToMailText() string {
	return r.String()
}