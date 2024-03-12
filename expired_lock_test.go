package expiredlock

import "testing"

func Test_ExpiredLock(t *testing.T) {
	lock := NewExpiredLock()
	lock.Lock(0)
	if err := lock.Unlock(); err != nil {
		t.Error(err)
	}
}

func Test_ExpiredLockNotSameUser(t *testing.T) {
	lock := NewExpiredLock()
	lock.Lock(0)
	ch := make(chan struct{})
	go func() {
		if err := lock.Unlock(); err == nil {
			t.Error(err)
		}
		close(ch)
	}()
	<-ch
}

func Test_ExpiredLockExp(t *testing.T) {
	lock := NewExpiredLock()
	lock.Lock(1)
	lock.Lock(0)
	if err := lock.Unlock(); err != nil {
		t.Error(err)
	}
}