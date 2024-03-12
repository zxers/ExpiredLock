package expiredlock

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type ExpiredLock struct {
	mutex sync.Mutex
	processMutex sync.Mutex
	owner string
	stop context.CancelFunc // 异步 goroutine 生命周期中止控制器
}

func NewExpiredLock() *ExpiredLock {
	return &ExpiredLock{}
}

func (e *ExpiredLock) Lock(expireSeconds int) {
	e.mutex.Lock()
	
	e.processMutex.Lock()
	defer e.processMutex.Unlock()
	token := GetCurrentProcessAndGoroutineIDStr()
	fmt.Println("Lock token: ", token)
	e.owner = token

	if expireSeconds <= 0 {
		return
	} 

	ctx, cancel := context.WithCancel(context.Background())
	e.stop = cancel

	go func ()  {
		select {
		case <-time.After(time.Duration(expireSeconds) * time.Second):
			e.unlock(token)
		case <-ctx.Done():

		}
	}()
}

func (e *ExpiredLock) Unlock() error {
	token := GetCurrentProcessAndGoroutineIDStr()
	return e.unlock(token)
}

func (e *ExpiredLock) unlock(token string) error {
	e.processMutex.Lock()
	defer e.processMutex.Unlock()
	fmt.Println("unLock token: ", token)
	if token != e.owner {
		return errors.New("not your lock")
	}

	e.owner = ""
	// 中止异步 goroutine 生命周期
	if e.stop != nil {
		e.stop()
	}
	e.mutex.Unlock()
	return nil
}