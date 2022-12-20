package context

import (
	"errors"
	"fmt"
	"github.com/valyala/fastjson"
	"strings"
)

type Expr struct {
	expr string
	Type int // localvar
	Token string
	next *Expr
	head *Expr
}

const (
	ExprTypeHead = iota
	ExprTypeObject
	ExprTypeArray
)

func (e *Expr) getValue(v *fastjson.Value) ([]*fastjson.Value, error) {
	p := e.next
	pva := []*fastjson.Value{v}
	var err error

	for p != nil {
		pv := make([]*fastjson.Value, 0)
		for _, value := range pva {
			if p.Type == ExprTypeObject {
				tmp := value.Get(p.Token)
				if tmp == nil {
					err = errors.New(fmt.Sprintf("%s is not exist", p.head.expr))
					goto flag
				}
				pv = append(pv, tmp)
			} else if p.Type == ExprTypeArray {
				tmp := value.GetArray(p.expr)
				if tmp == nil {
					err = errors.New(fmt.Sprintf("%s is not exist", p.head.expr))
					goto flag
				}
				pv = append(pv, tmp...)
			} else {
				err = errors.New("")
				goto flag
			}
			pva = pv
			p = p.next
		}
	}
flag:
	return pva, err
}

func (e *Expr) IsTaskTsp() bool {
	return strings.Contains(e.Token, ":RSP__")
}

func ExpressionParse(expr string) (*Expr, error) {
	items := strings.Split(expr, ".")
	if len(items) == 0 {
		return nil, errors.New("invalid expr")
	}

	head := &Expr{
		Type: ExprTypeHead,
		expr: expr,
		Token: items[0],
		next: nil,
	}
	head.head = head

	p := head

	for i:=1; i<len(items); i++ {
		var tmp *Expr
		if IsObject(items[i]) {
			tmp = &Expr{
				Type:  ExprTypeObject,
				expr:  expr,
				Token: items[i],
				head:  head,
				next: nil,
			}
		} else if IsArray(expr) {
			tmp = &Expr{
				Type: ExprTypeArray,
				expr: expr,
				Token: items[i],
				head: head,
				next: nil,
			}
		} else {
			return nil, errors.New(fmt.Sprintf("failed to parse the expr, %s is illeagal", items[i]))
		}

		p.next = tmp
		p = p.next
	}

	return head, nil
}

func IsObject(expr string) bool {
	return !strings.Contains(expr, "[") && !strings.Contains(expr, "]")
}

func IsArray(expr string) bool {
	return !IsObject(expr) && (strings.Index(expr, "]")  - strings.Index(expr, "[")) == 1
}

func (c *Context) GetValue(source string) ([]*fastjson.Value, error) {
	ee, err := ExpressionParse(source)
	if err != nil {
		return nil, err
	}

	var value []*fastjson.Value

	if ee.IsTaskTsp() {
		id := getTaskId(source)
		v, err := c.getActionResponse(id)
		if err != nil {
			return nil, err
		}
		// 补充 expr 表达式解析，
		value, err = ee.getValue(v)
	} else {
		return nil, errors.New("failed to get value")
	}
	return value, nil
}
