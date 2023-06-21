package service

import (
	"time"
)

type Result struct {
	Data any
	Err  error
}

type Job struct {
	resultChan chan Result
	taskFunc   func() Result
}

type TaskService struct {
	rateLimitJobChan chan Job
	priorityJobChan  chan Job
}

func NewTaskService(workerCount int, chanSize int, rateLimitTime time.Duration) *TaskService {
	rateLimitJobChan := make(chan Job, chanSize)
	priorityJobChan := make(chan Job, chanSize)

	limiter := time.Tick(rateLimitTime)

	for i := 0; i < workerCount; i++ {
		go doWork(rateLimitJobChan, priorityJobChan, limiter)
	}

	return &TaskService{
		rateLimitJobChan: rateLimitJobChan,
		priorityJobChan:  priorityJobChan,
	}
}

func doWork(rateLimitJobChan chan Job, priorityJobChan chan Job, limiter <-chan time.Time) {
	for {
		select {
		case priorityJob := <-priorityJobChan:
			priorityJob.resultChan <- priorityJob.taskFunc()
		default:
			job := <-rateLimitJobChan
			// 限速
			<-limiter
			job.resultChan <- job.taskFunc()
		}
	}
}

// 增加速率限制任務
func (ts *TaskService) SendRateLimitTask(taskFunc func() Result) chan Result {
	resultChan := make(chan Result, 1)
	go func() {
		ts.rateLimitJobChan <- Job{resultChan, taskFunc}
	}()
	return resultChan
}

// 增加急件任務
func (ts *TaskService) SendPriorityTask(taskFunc func() Result) chan Result {
	resultChan := make(chan Result, 1)
	go func() {
		ts.priorityJobChan <- Job{resultChan, taskFunc}
	}()
	return resultChan
}
