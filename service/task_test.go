package service

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTaskService_SendRateLimitTask(t *testing.T) {
	fakeFetchFunc := func(request any) Result {
		time.Sleep(time.Second * 1)
		return Result{
			Data: request,
			Err:  nil,
		}
	}

	const (
		count       = 30
		millisecond = 100
	)

	wg := sync.WaitGroup{}
	taskService := NewTaskService(5, 100, time.Millisecond*millisecond)
	start := time.Now()
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(index int) {
			result := <-taskService.SendRateLimitTask(func() Result {
				return fakeFetchFunc(index)
			})
			require.Equal(t, index, result.Data)
			require.Nil(t, result.Err)
			wg.Done()
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	if time.Duration(duration.Milliseconds()) < time.Duration(millisecond)*count {
		t.Errorf("expect greater than %d , but got %d",
			time.Duration(millisecond)*count,
			duration.Milliseconds())
	}
}

func TestTaskService_SendPriorityTask(t *testing.T) {
	fakeFetchFunc := func(request any) Result {
		time.Sleep(time.Second * 1)
		return Result{
			Data: request,
			Err:  nil,
		}
	}

	const (
		count       = 30
		millisecond = 100
	)

	wg := sync.WaitGroup{}
	taskService := NewTaskService(5, 100, time.Millisecond*millisecond)
	start := time.Now()
	wg.Add(1)

	go func() {
		time.Sleep(time.Second * 2)
		result := <-taskService.SendPriorityTask(func() Result {
			duration := time.Since(start)
			return fakeFetchFunc(duration)
		})
		end := result.Data.(time.Duration)
		if time.Duration(end.Milliseconds()) > time.Duration(millisecond)*count {
			t.Errorf("expect less than %d , but got %d",
				time.Duration(millisecond)*count,
				end.Milliseconds())
		}
		wg.Done()
	}()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(index int) {
			result := <-taskService.SendRateLimitTask(func() Result {
				return fakeFetchFunc(index)
			})
			require.Equal(t, index, result.Data)
			require.Nil(t, result.Err)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
