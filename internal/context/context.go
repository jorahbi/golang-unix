package context

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"socket/pkg/logger"
	"socket/pkg/socket"
)

var modules = getModules()

const (
	M_DATABASE = "database"
)

type Processor interface {
	Execute () (int, string, map[interface{}]interface{})
}

func getProcessor (request *socket.Request) (Processor, error) {
	if module, isOk := modules[request.GetModuleName()]; isOk {
		return module(request.GetBody())
	}
	return nil, errors.New("not defined execute module")
}

func getModules () map[string] func(data map[string]interface{}) (Processor, error) {

	m := make(map[string]func(data map[string]interface{}) (Processor, error))
	m[M_DATABASE] = func(data map[string]interface{}) (Processor, error) {
		return NewStore(data)
	}
	return m
}

func Handle (s *socket.Socket) {
	defer s.Closer()
	for {
		request, err := s.GetRequest()
		if err != nil {
			logger.Handle.Panic("socket get request error",
				zap.Int("socket id", s.GetFd()),
				zap.String("request data", request.OriginData()),
				zap.String("error info ", err.Error()))
			body := make(map[interface{}]interface{})
			err := s.SetRespBody(1001, "", body).Send()
			logger.Handle.Panic("send response to socket",
				zap.Int("socket id", s.GetFd()),
				zap.String("response", s.RespToString()),
				zap.String("error info ", err.Error()))
			break
		}
		processor, err := getProcessor(request)
		if err != nil {
			body := make(map[interface{}]interface{})
			err := s.SetRespBody(1001, "", body).Send()
			fmt.Println("error1", err)
			break
		}
		code, msg, data := processor.Execute()
		err = s.SetRespBody(code, msg, data).Send()
		fmt.Println("error2", err)
	}
}
