package socket

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Socket struct {
	conn net.Conn
	fd int
	*Request
	response
}

type Request struct {
	module string
	data map[string]interface{}
	originData []byte
}

type response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data map[interface{}]interface{} `json:"data"`
	toByte []byte
}

func NewSocket (conn net.Conn) *Socket{
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return &Socket{conn: conn, fd: int(pfdVal.FieldByName("Sysfd").Int()), Request: &Request{}, response: response{}}
}

func (s *Socket) GetConn () net.Conn {
	return s.conn
}

func (s *Socket) GetFd () int {
	return s.fd
}

func (s *Socket) SetReadDeadline () error{
	return s.GetConn().SetReadDeadline(time.Now().Add(5000 * time.Millisecond))
}

func (s *Socket) Closer () {
	if err := s.GetConn().Close(); err != nil {
		fmt.Println("close conn error", err)
	}
}

func (s *Socket) Send () error {
	_, err := s.GetConn().Write(s.toByte)
	return err
}

func (s *Socket) Error(code int) {

}

func (s *Socket) GetRequest() (*Request, error) {

	request, err := unpack(s.GetConn())
	s.Request = request
	return s.Request, err
}


func (s *Socket) SetRespBody(code int, msg string, data map[interface{}]interface{}) *Socket {
	s.response.Msg = msg
	s.response.Code = code
	s.response.Data = data
	s.response.toByte = pack(s.response)
	return s
}

func (s *Socket) RespToString () string {
	return string(s.response.toByte)
}

func (r *Request) GetModuleName () string {
	return r.module
}

func (r *Request) OriginData () string {
	return string(r.originData)
}

func (r *Request) GetBody () map[string]interface{} {
	return r.data
}


func unpack (conn net.Conn) (*Request, error) {

	request := &Request{}
	rstBody := make(map[string]interface{})
	buf := make([]byte, 20)
	if _, err := conn.Read(buf); err != nil {
		fmt.Println("unpack", err)
		return request, errors.Wrap(err, "read socket stream error")
	}

	indexOf := strings.IndexAny(string(buf), "\n")
	length := strings.Split(string(buf[0:indexOf]), ":")
	index, _ := strconv.Atoi(length[len(length)-1])
	preCont := make([]byte, index-len(buf[indexOf+9:]))

	if _, err := conn.Read(preCont); err != nil {
		return request, errors.Wrap(err, "read socket stream error")
	}
	request.originData = append(buf[indexOf+9:], preCont...)
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(request.originData, &rstBody); err != nil {
		return request, errors.Wrap(err, "json decode error")
	}

	if module, isOk := rstBody["module"]; isOk {
		request.module = module.(string)
	} else {
		return nil, errors.New("not defined module")
	}
	if data, isOk := rstBody["data"]; isOk {
		request.data = data.(map[string]interface{})
	} else {
		return nil, errors.New("not defined data")
	}

	return request, nil
}

func pack(body response) []byte{
	if resp, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(body); err == nil {
		return resp
	}else {
		return []byte{}
	}
}
