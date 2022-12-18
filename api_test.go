package configurable_flow_actor

import (
	"fmt"
	"testing"
)

func TestApi(t *testing.T) {
	rsp, _ := CFARun(nil, "helo", "")
	fmt.Println(string(rsp))
}
