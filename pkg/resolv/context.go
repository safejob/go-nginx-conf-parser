// resove 包，该包包含了项目最基础的上下文相关对象，及相关方法及函数
// 创建者： ClessLi
// 创建时间：2020-1-17 11:14:15
package resolv

import (
	"fmt"
	"strings"
)

var INDENT string = "    "

// Context, 上下文接口对象，定义了上下文接口需实现的增、删、改等方法
type Context interface {
	Add(...Parser)
	Remove(...Parser)
	Modify(int, Parser) error
	//Filter(string, string) *Context
	//getReg() string
	//Dict() map[string]interface{}
	String() []string
	//BumpChildDepth(int)
	dump() ([]string, error)
}

// BasicContext, 上下文基础对象，定义了上下文类型的基本属性及基础方法
type BasicContext struct {
	Name     string   `json:"Name"`
	Value    string   `json:"value"`
	depth    int      `json:"-"`
	Children []Parser `json:"children"`
}

// Add, BasicContext 类新增子对象的方法， Context.Add(...interface{}) 的实现
func (c *BasicContext) Add(contents ...Parser) {
	for _, content := range contents {
		/*if _, isBC := content.(Context); isBC {
			content.(Context).BumpChildDepth(c.depth+1)
		}*/
		c.Children = append(c.Children, content)
	}
}

// Remove, BasicContext 类删除子对象的方法， Context.Remove(...interface{}) 的实现
func (c *BasicContext) Remove(contents ...Parser) {
	for _, content := range contents {
		for index, child := range c.Children {
			if content == child {
				c.remove(index)
			}
		}
	}
}

// Modify, BasicContext 类修改子对象的方法， Context.Modify(int, interface{}) error 的实现
func (c *BasicContext) Modify(index int, content Parser) error {
	switch content.(type) {
	case Context, *Comment, *Key:
		c.Children[index] = content
	default:
		return fmt.Errorf("conf format not supported with: %T", content)
	}
	return nil
}

// Filter, BasicContext 类生成过滤对象的方法， Context.Filter(string, string) []*Context 的实现
//TODO: 过滤器
//func (c *BasicContext) Filter(btype, name string) *Context {
//}

//func (c *BasicContext) getReg() string {
//	return c.Value
//}

//func (c *BasicContext) Dict() map[string]interface{} {
//}

func (c *BasicContext) String() []string {
	ret := make([]string, 0)

	contextTitle := c.getTitle()

	ret = append(ret, contextTitle)

	for _, child := range c.Children {
		switch child.(type) {
		case *Key:
			ret = append(ret, INDENT+child.String()[0])
		case *Comment:
			if child.(*Comment).Inline && len(ret) >= 1 {
				ret[len(ret)-1] = strings.TrimRight(ret[len(ret)-1], "\n") + "  " + child.String()[0]
			} else {
				ret = append(ret, INDENT+child.String()[0])
			}
		case Context:
			strs := child.String()
			//ret = append(ret, INDENT+strs[0])
			//for _, str := range strs[1:] {
			for _, str := range strs {
				ret = append(ret, INDENT+str)
			}
		default:
			str := child.String()
			if str != nil {
				ret = append(ret, str...)
			}
		}
	}
	ret[len(ret)-1] = RegEndWithCR.ReplaceAllString(ret[len(ret)-1], "}\n")
	ret = append(ret, "}\n\n")

	return ret
}

func (c *BasicContext) dump() ([]string, error) {
	ret := make([]string, 0)
	contextTitle := c.getTitle()
	ret = append(ret, contextTitle)

	for _, child := range c.Children {
		switch child.(type) {
		case *Key:
			ret = append(ret, INDENT+child.String()[0])
		case *Comment:
			if child.(*Comment).Inline && len(ret) >= 1 {
				ret[len(ret)-1] = strings.TrimRight(ret[len(ret)-1], "\n") + "  " + child.String()[0]
			} else {
				ret = append(ret, INDENT+child.String()[0])
			}
		case Context:
			strs, err := child.(Context).dump()
			if err != nil {
				return ret, err
			}

			for _, str := range strs {
				ret = append(ret, INDENT+str)
			}
		default:
			str := child.String()
			if str != nil {
				ret = append(ret, str...)
			}
		}
	}
	ret[len(ret)-1] = RegEndWithCR.ReplaceAllString(ret[len(ret)-1], "}\n")
	ret = append(ret, "}\n\n")

	return ret, nil
}

func (c *BasicContext) remove(index int) {
	c.Children = append(c.Children[:index], c.Children[index+1:]...)
}

func (c *BasicContext) getTitle() string {
	contextTitle := ""
	/*for i := 0; i < c.depth; i++ {
		contextTitle += INDENT
	}*/
	contextTitle += c.Name

	if c.Value != "" {
		contextTitle += " " + c.Value
	}

	contextTitle += " {\n"
	return contextTitle
}

//func (c *BasicContext) BumpChildDepth(depth int) {
//	for i := 0; i < len(c.Children); i++ {
//		if bc, isBC := c.Children[i].(*BasicContext); isBC {
//			bc.depth = depth + 1
//			c.BumpChildDepth(bc.depth)
//		}
//	}
//}
