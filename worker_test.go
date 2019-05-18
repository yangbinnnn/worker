package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type Job struct {
	Num int
	Url string
}
type Result struct {
	Msg     string
	Success bool
	Job     *Job
}

type MyWorker struct {
	num int // worker context
}

func (this *MyWorker) HandlerJob(ijob interface{}) (result interface{}) {
	job := ijob.(*Job)
	result0 := &Result{
		Job:     job,
		Success: true,
	}

	defer func() {
		result = result0
	}()

	resp, err := http.Get(job.Url)
	if err != nil {
		result0.Msg = err.Error()
		result0.Success = false
		return
	}
	defer resp.Body.Close()

	respModel := struct {
		Unixtime int    `json:"unixtime"`
		Datetime string `json:"datetime"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&respModel)
	if err != nil {
		result0.Msg = err.Error()
		result0.Success = false
		return
	}

	result0.Msg = fmt.Sprintf("%s %d", respModel.Datetime, respModel.Unixtime)
	this.num++
	return
}

func (this *MyWorker) HandlerResult(iresult interface{}) {
	result := iresult.(*Result)
	fmt.Printf("job %d, fetch %s: result %v %s\n", result.Job.Num, result.Job.Url, result.Success, result.Msg)
	return
}

func TestWorkerPool(t *testing.T) {
	myWorker := &MyWorker{}
	pool := NewWorkerPool(5, 100, myWorker)
	pool.Start()
	fmt.Println("start worker...")
	fmt.Println("worker running:", pool.IsRunning())

	for i := 0; i < 10; i++ {
		job := Job{
			Num: i,
			Url: "https://worldtimeapi.org/api/timezone/Asia/Shanghai",
		}
		pool.Push(&job)
	}
	fmt.Println("wait job done...")
	pool.WaitAndClose()
	fmt.Println("worker running:", pool.IsRunning())
	fmt.Println("all job done:", myWorker.num)
}
