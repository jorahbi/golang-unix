package epoll

import (
	"go.uber.org/zap"
	"socket/internal/context"
	"socket/pkg/logger"
	"socket/pkg/socket"
	"sync"
	"time"
)

var Works = newWorkerPool(10, 100)

type workerPool struct {
	workers int
	maxTasks int
	closed bool
	done chan struct{}
	taskQueue chan *socket.Socket
	mu *sync.RWMutex
}

func newWorkerPool (workers int, tasks int) *workerPool{
	return &workerPool{
		workers: workers,
		maxTasks: tasks,
		closed: false,
		done: make(chan struct{}),
		taskQueue: make(chan *socket.Socket, tasks),
		mu: &sync.RWMutex{},
	}
}


func (wp *workerPool) Closer (){
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.closed = true
	close(wp.taskQueue)
	close(wp.done)

}

func (wp *workerPool) Start () {
	for i := 0; i < wp.workers; i ++ {
		go wp.StartWorker()
	}
}

func (wp *workerPool) StartWorker () {
	for {
		select {
		case <-wp.done:
			return
		case s := <-wp.taskQueue:
			if s.GetConn() == nil {
				return
			}
			context.Handle(s)
		case <-time.NewTicker(2000 * time.Microsecond).C:

		}
	}
}

func (wp *workerPool) AddTask (socket *socket.Socket) {
	wp.mu.Lock()
	wp.mu.Unlock()
	if wp.closed {
		wp.mu.Unlock()
		return
	}
	err := Epolls.Remove(socket)
	logger.Handle.Info("remove socket to epoll error",
		zap.String("error info", err.Error()))
	wp.taskQueue <- socket
}
