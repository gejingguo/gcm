package gcm

import "time"

// 任务组概念
type TaskGroup struct {
	ID int `json:"id"`					// 组ID
	Name string `json:"name"`			// 名称
	Parallel bool `json:"parallel"` 	// 是否支持并行
	Pro *Project						// 工程
	Tasks []*Task 						// 任务列表

	LastErrTaskNum int 					// 上次构建错误任务数量, 0=success
	LastUsedTime time.Duration			// 上次构建耗时
}

func (tg *TaskGroup) AddTask(task *Task) error {
	if tg.Tasks == nil {
		tg.Tasks = make([]*Task, 0, 10)
	}
	tg.Tasks = append(tg.Tasks, task)
	return nil
}

// 任务组也支持Execute接口
func (tg *TaskGroup) Execute(rp TaskRunTimeParam) error {
	tg.LastUsedTime = 0
	tg.LastErrTaskNum = 0
	begTime := time.Now()
	// 并发处理
	if tg.Parallel && len(tg.Tasks) > 1 {
		result := make([]bool, len(tg.Tasks))
		retChan := make(chan taskRunRet)
		//
		for i := 0; i < len(tg.Tasks); i++  {
			go func(task int, c chan taskRunRet, rp TaskRunTimeParam) {
				err := tg.Tasks[task].Execute(rp)
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
				if ret.err != nil {
					tg.LastErrTaskNum++
				}
			}
		}
	} else {
		for _, t := range tg.Tasks  {
			err := t.Execute(rp)
			if err != nil {
				tg.LastErrTaskNum++
			}
		}
	}
	tg.LastUsedTime = time.Since(begTime)

	return nil
}
