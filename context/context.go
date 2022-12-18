package context

import (
	"errors"
	"github.com/valyala/fastjson"
	"log"
	"strings"
)

type Context struct {
	Logger *log.Logger
	request *fastjson.Value
	response *fastjson.Value
	taskResult map[string]*TaskResult
	localVar *fastjson.Value
	A *fastjson.Arena

	// flow process result
	Exit bool
}

type TaskResult struct {
	request *fastjson.Value
	response *fastjson.Value
	headers map[string]string // 这个字段还需要后续在修改下
	responseByte []byte
}

func (c *Context) Init(request string) error {
	c.request, _ = fastjson.Parse(request)

	c.A = new(fastjson.Arena)

	c.response = c.A.NewObject()

	c.localVar = c.A.NewObject()

	c.taskResult = make(map[string]*TaskResult)

	c.Exit = true

	return nil
}

func (c *Context) NewTaskResult(id string) {
	c.taskResult[id] = &TaskResult{
		response: c.A.NewObject(),
		request: c.A.NewObject(),
		headers: make(map[string]string),
	}
}

func (c *Context) SetRsp(id string, rsp []byte) error {
	rspValue, err := fastjson.Parse(string(rsp))
	if err != nil {
		return err
	}
	if _, ok := c.taskResult[id]; !ok {
		return errors.New("failed to get taskResult value")
	}
	c.taskResult[id].response = rspValue
	return nil
}

func (c *Context) SetHeaders(id string, headers map[string]string) error {
	if _, ok := c.taskResult[id]; !ok {
		return errors.New("failed to get taskResult value")
	}
	c.taskResult[id].headers = make(map[string]string)
	c.taskResult[id].headers = headers
	return nil
}

func (c *Context) MarshalResponse() ([]byte, error) {
	var rspByte []byte
	rspByte = c.response.MarshalTo(rspByte)

	//var rsp = make(map[string]interface{})
	//err := json.Unmarshal(rspByte, &rsp)
	//if err != nil {
	//	return nil, err
	//}

	return rspByte, nil
}

func IsTaskTsp(source string) bool {
	return strings.Contains(source, ":RSP__")
}

func (c *Context) SetValue(dst string, value *fastjson.Value) error {
	//if IsTaskTsp(source) {
	//	id := getTaskId(source)
	//	v := c.taskResult[id].response
	//	// 补充 expr 表达式解析，
	//	return getValue(v, "")
	//}
	//return nil, errors.New("failed to get value")
	return nil
}

func (c *Context) SetResponse(v *fastjson.Value) {
	c.response = v
}

func (c *Context) GetValue(source string) (*fastjson.Value, error) {
	if IsTaskTsp(source) {
		id := getTaskId(source)
		v := c.taskResult[id].response
		// 补充 expr 表达式解析，
		return getValue(v, "")
	}
	return nil, errors.New("failed to get value")
}

func getTaskId(source string) string {
	idList := strings.Split(source, "__")
	if len(idList) < 3 {
		return ""
	}
	return strings.Replace(idList[1], ":RSP", "", 1)
}

func getValue(v *fastjson.Value, dst string) (*fastjson.Value, error) {
	if len(dst) == 0 {
		return v, nil
	}
	return nil, errors.New("failed to get value")
}
