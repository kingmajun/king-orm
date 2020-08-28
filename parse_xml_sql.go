package king_orm

import (
	"encoding/xml"
	"github.com/antonmedv/expr"
	"github.com/kingmajun/king-orm/model"
	"io"
	"strings"
)

/*
*把xml格式的sql转成对应的结构体，方便操作
* <select id="getuser">
*	select * from user2
*	<if test="username!=nil and username!=''">
* 		where username = #{username}
* 	</if>
*</select>
*/




type node struct {
	Id       string
	TagName  string
	Attrs    map[string]xml.Attr
	Elements []element
}

type element struct {
	ElementType model.ElemType
	Val         interface{}
}

//解析xml格式的sql
func parse(r io.Reader) *node {
	parser := xml.NewDecoder(r)
	var root node

	st := NewStack()
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement: //标签 开始
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			attr := elmt.Attr
			attrMap := make(map[string]xml.Attr)
			for _, val := range attr {
				attrMap[val.Name.Local] = val
			}
			node := node{
				TagName:     name,
				Attrs:    attrMap,
				Elements: make([]element, 0),
			}
			for _, val := range attr {
				if val.Name.Local == "id" {
					node.Id = val.Value
				}
			}
			st.Push(node)

		case xml.EndElement: //标签 结束
			if st.Len() > 0 {
				//cur node
				n := st.Pop().(node)
				if st.Len() > 0 { //if the root node then append to element
					e := element{
						ElementType: model.ELEMTYPNODE,
						Val:         n,
					}

					pn := st.Pop().(node)
					els := pn.Elements
					els = append(els, e)
					pn.Elements = els
					st.Push(pn)
				} else { //else root = n
					root = n
				}
			}
		case xml.CharData: //标签 内容
			if st.Len() > 0 {
				n := st.Pop().(node)

				bytes := xml.CharData(t)
				content := strings.TrimSpace(string(bytes))
				if content != "" {
					e := element{
						ElementType: model.ELEMTYPTEXT,
						Val:         content,
					}
					els := n.Elements
					els = append(els, e)
					n.Elements = els
				}

				st.Push(n)
			}

		case xml.Comment: //标签 注释内容
		case xml.ProcInst:
		case xml.Directive:
		default:
		}
	}

	if st.Len() != 0 {
		panic("Parse xml error, there is tag no close, please check your xml config!")
	}

	return &root
}

//创建一个带有参数的sql
func createParamsSql(params map[string]interface{},elements ...element)(sql string)  {
	sql = ""
	if len(elements)==1 {
		elem := elements[0]
		if elem.ElementType == model.ELEMTYPTEXT{
			sql += elem.Val.(string)
		}else if elem.ElementType == model.ELEMTYPNODE{
			node := elem.Val.(node)
			if node.TagName == "if" {
				testVal := node.Attrs["test"].Value
				ok,err := expr.Eval(testVal, params)
				if err !=nil{
					println(err)
					return
				}
				if ok.(bool) {
					sql += " "+createParamsSql(params,node.Elements...)+" "
				}

			}
		}
		return
	}
	for _, elem := range elements {
		sql += createParamsSql(params,elem)
	}
	return
}

