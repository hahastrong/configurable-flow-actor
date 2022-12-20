package configurable_flow_actor

import (
	"fmt"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	start := time.Now().UnixMilli()
	rsp, _ := CFARun(nil, "helo", "")
	fmt.Println(time.Now().UnixMilli() - start, string(rsp))
}
