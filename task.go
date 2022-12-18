package configurable_flow_actor

import (
	"github.com/configurable-flow-actor/context"
	"io/ioutil"
	"net/http"
)

type Task interface {
	DoTask(ctx *context.Context) error
}

type TaskParam struct {
	ID       string            `json:"id"`
	TaskType string            `json:"task_type"`
	Path     string            `json:"path"`
	Method   string            `json:"method"`
	Request  map[string]*ParamNode `json:"request"` // 后续在更改为可以替换的
	Response map[string]*ParamNode `json:"response"`
}

type ParamNode struct {
	Data      string            `json:"data"`
	Type      string            `json:"type"`       // string, number, bool, array
	Action    string            `json:"action"`     // data parse method, expr, data, iNfunc
	TargetIdx map[string]string `json:"target_idx"` // to replace target dst item
	SourceIdx map[string]string `json:"source_idx"` // to replace source dst item
}

func (p *ParamNode) exec(ctx *context.Context, k string) error {
	if p.Action == "expr" {
		v, _ := ctx.GetValue(p.Data)
		// 直接复制给 respose
		ctx.SetResponse(v)
	}
	return nil
}

const (
	HTTPREQUEST = "HttpRequest"
	DATABUILDER = "DataBuilder"
)

type HttpRequest struct {
	tp *TaskParam
}

type DataBuilder struct {
	tp *TaskParam
}

type Start struct {
	tp *TaskParam
}

type End struct {
	tp *TaskParam
}

func (t *HttpRequest) DoTask(ctx *context.Context) error {
	if t.tp.Method == "GET" {
		// 需要自定义请求头
		// 添加header
		client := &http.Client{}
		request, err := http.NewRequest(http.MethodGet, t.tp.Path, nil)
		if err != nil {
			// 添加日志打印
			return err
		}
		rsp, err := client.Do(request)
		if err != nil {
			return err
		}

		rspBody, err := ioutil.ReadAll(rsp.Body)
		_ = rsp.Body.Close()

		err = ctx.SetRsp(t.tp.ID, rspBody)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *DataBuilder) DoTask(ctx *context.Context) error {
	for k, v := range t.tp.Response {
		err := v.exec(ctx, k)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Start) DoTask(ctx *context.Context) error {
	return nil
}

func (t *End) DoTask(ctx *context.Context) error {
	return nil
}

func NewTask(tp *TaskParam, taskType string) Task {
	switch taskType {
	case START:
		return &Start{tp: tp}
	case END:
		return &Start{tp: tp}
	case HTTPREQUEST:
		return &HttpRequest{tp: tp}
	case DATABUILDER:
		return &DataBuilder{tp: tp}
	default:
		return nil
	}
}