package epoll

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net"
	"socket/pkg/logger"
	"socket/pkg/socket"
	"sync"
	"sync/atomic"
	"syscall"
)

var Epolls = newEventPoll()

type eventPoll struct {
	fd int
	sockets *sync.Map
	len int32
}

func newEventPoll () *eventPoll {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil
	}
	return &eventPoll{
		fd: fd,
		len: 100,
		sockets: &sync.Map{},
	}
}


func Start () {
	for {
		connects, err := Epolls.Wait()
		if err != nil {
			logger.Handle.Info("epoll wait error",
				zap.String("error info", err.Error()))
		}
		for _, conn := range connects {
			if conn.GetConn() == nil{
				continue
			}
			Works.AddTask(conn)
		}
	}
}

func (ep *eventPoll) Add (conn net.Conn) error {
	s := socket.NewSocket(conn)
	fd := s.GetFd()
	err := syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{Events: syscall.EPOLLIN | syscall.EPOLLHUP, Fd: int32(fd)})
	if err != nil{
		return errors.WithStack(err)
	}
	logger.Handle.Info("socket to add epoll success",
		zap.Int("socket id", s.GetFd()))
	ep.sockets.Store(fd, s)
	atomic.AddInt32(&ep.len, -1)
	if ep.len < 0 {
		logger.Handle.Info("socket connection number is full",
			zap.Int("connection total", 100))
		return errors.New("socket connection number is full")
	}
	return nil
}

func (ep *eventPoll) Remove (conn *socket.Socket) error{
	err := syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_DEL, conn.GetFd(), nil)
	if err != nil {
		return errors.WithStack(err)
	}
	ep.sockets.Delete(conn.GetFd())
	atomic.AddInt32(&ep.len, 1)
	return nil

}

func (ep *eventPoll) Wait () ([]*socket.Socket, error) {
	events := make([]syscall.EpollEvent, 100)
	n, err := syscall.EpollWait(ep.fd, events, 100)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var sockets []*socket.Socket
	for i := 0; i < n; i++ {
		if conn, isOk:= ep.sockets.Load(int(events[i].Fd)); isOk {
			sockets = append(sockets, conn.(*socket.Socket))
		}
	}
	return sockets, nil
}



