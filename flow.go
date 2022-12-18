package configurable_flow_actor

import (
	"errors"
	"github.com/configurable-flow-actor/context"
)

const (
	GATEWAY = "gateway"
	START = "start"
	END = "end"
	TASK = "task"
)

type Flow map[string]*Node

type Node struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"` // start, end, gateway, task
	Expr    string    `json:"expr"`
	Next    string    `json:"next"`
	Default string    `json:"default"`
	Task    *TaskParam `json:"task"`
}

type RunNode struct {
	n *Node
	tp Task
	flow Flow
}

func (rn *RunNode) getDefault(ctx *context.Context) (*RunNode, error) {
	defaultId := rn.n.Default
	for k, _ := range rn.flow {
		if k == defaultId {
			ctx.NewTaskResult(k)
			return &RunNode{
				n: rn.flow[k],
				flow: rn.flow,
				tp: NewTask(rn.flow[k].Task, rn.flow[k].Task.TaskType),
			}, nil

		}
	}
	return nil, errors.New("failed to get default node")
}

func (rn *RunNode) getNext(ctx *context.Context) (*RunNode, error) {
	defaultId := rn.n.Next
	for k, _ := range rn.flow {
		if k == defaultId {

			ctx.NewTaskResult(k)
			return &RunNode{
				n: rn.flow[k],
				flow: rn.flow,
				tp: NewTask(rn.flow[k].Task, rn.flow[k].Task.TaskType),
			}, nil

		}
	}
	return nil, errors.New("failed to get next node")
}

func (rn *RunNode) eval(ctx *context.Context) (bool, error) {
	if rn.n.Type != GATEWAY {
		return true, nil
	}
	// add gateway condition
	return true, nil
}

func (f Flow) getStart() (*RunNode, error) {
	for k, _ := range f {
		if k == START {
			return &RunNode{
				n: f[k],
				tp: NewTask(f[k].Task, START),
				flow: f,
			}, nil

		}
	}

	return nil, errors.New("failed to get start node")
}

func (f Flow) Run(ctx *context.Context) error {
	n, err := f.getStart()
	if err != nil {
		return err
	}

	for n != nil {
		// gateway
		rv, err := n.eval(ctx)
		if err != nil {
			return err
		}

		if err := n.tp.DoTask(ctx); err != nil {
			return err
		}

		if rv {
			n, _ = n.getNext(ctx)
		} else {
			n, _ = n.getDefault(ctx)
		}
	}

	return nil
}
