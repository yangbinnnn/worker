package worker

import (
	"sync"
)

type Worker interface {
	HandlerJob(job interface{}) (result interface{})
	HandlerResult(result interface{})
}

type WorkerPool struct {
	WorkNum  int
	QueueNum int
	worker   Worker
	jobChan  chan interface{}
	wg       sync.WaitGroup
	running  bool
}

func NewWorkerPool(workNum, queueNum int, worker Worker) *WorkerPool {
	wp := &WorkerPool{
		WorkNum:  workNum,
		QueueNum: queueNum,
		worker:   worker,
		jobChan:  make(chan interface{}, queueNum),
	}
	return wp
}

func (this *WorkerPool) Push(job interface{}) {
	this.jobChan <- job
	return
}

func (this *WorkerPool) Start() (err error) {
	this.running = true
	for i := 0; i < this.WorkNum; i++ {
		this.wg.Add(1)
		go func() {
			defer this.wg.Done()
			for this.running {
				job := <-this.jobChan
				if job == nil {
					return
				}
				this.worker.HandlerResult(this.worker.HandlerJob(job))
			}
		}()
	}
	return
}

func (this *WorkerPool) IsRunning() bool {
	return this.running
}

func (this *WorkerPool) WaitAndClose() (err error) {
	close(this.jobChan)
	this.wg.Wait()
	this.running = false
	return
}
