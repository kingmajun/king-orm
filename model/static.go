package model

type ElemType string

const (
	ELEMTYPTEXT ElemType = "text" // 静态文本节点
	ELEMTYPNODE ElemType = "node" // 节点子节点
)

//定义xml中标签符号
var XMLTAG = []string{"sqlMapper","select","update","insert","delete","if"}