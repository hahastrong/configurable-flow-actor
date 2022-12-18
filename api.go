package configurable_flow_actor

import (
	"encoding/json"
	"github.com/configurable-flow-actor/context"
	"log"
	"time"
)

func CFARun(lg *log.Logger, flowName string, request string) ([]byte, error) {
	ctx := &context.Context{
		Logger: lg,
	}

	if err := ctx.Init(request); err != nil {
		// todo logger
		return nil, err
	}

	var testFlow string = `{"start":{"id":"start","type":"start","next":"getHTTPRSP"},"getHTTPRSP":{"id":"getHTTPRSP","type":"task","next":"databuilder","task":{"id":"getHTTPRSP","task_type":"HttpRequest","method":"GET","path":"http://ip-api.com/json/138.199.21.138?lang=zh-CN"}},"databuilder":{"id":"databuilder","type":"task","next":"end","task":{"id":"databuilder","task_type":"DataBuilder","method":"GET","response":{".":{"data":"__getHTTPRSP:RSP__","action":"expr"}}}}}`
	var flow Flow
	err := json.Unmarshal([]byte(testFlow), &flow)
	if err != nil {
		return nil, err
	}

	ch := make(chan error, 1)
	// asyc run flow
	go func() {
		err := flow.Run(ctx)
		ch <- err
	}()

	// add timeout
	select {
	case err := <- ch:
		if err != nil {
			// 补充 获取返回数据
			return nil, err
		}
	case <- time.After(time.Millisecond * 200000):
			// timeout

	}

	rsp, err := ctx.MarshalResponse()
	if err != nil {
		return nil, err
	}
	return rsp, nil
}
