package mqtt

import (
	"container/list"
	"sync"
)

type Queue struct {
	queue *list.List
	lock sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{
		queue: list.New(),
		lock: sync.RWMutex{},
	}
}

func (q *Queue) EnqueueMessage(msg interface{}) {
	q.lock.Lock()
	//defer q.lock.Unlock()
	q.queue.PushBack(msg)
	q.lock.Unlock()
}
func (q *Queue) DequeueMessage() interface{} {
	q.lock.RLock()
	//defer q.lock.Unlock()
	if q.queue.Len() > 0 {
		q.lock.RUnlock()
		q.lock.Lock()
		el := q.queue.Front()
		v := q.queue.Remove(el)
		q.lock.Unlock()
		return v
	} else {
		q.lock.RUnlock()
		return nil
	}
}
func (q *Queue) Size() int {
	q.lock.RLock()
	//defer q.lock.RUnlock()
	//return q.queue.Len()
	ret := q.queue.Len()
	q.lock.RUnlock()
	return ret
}
