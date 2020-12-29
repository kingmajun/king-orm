package king_orm

import (
	"encoding/xml"
	"os"
	"testing"
)

func Test(t *testing.T) {
	ROOT, _ := os.Getwd()
	path := ROOT + "/demo/mapper"
	v, _ := ReaderConfigBuilder(path)

	//获取sql
	m, err := v.GetMethodSql("demo.mapper.user.getuser")
	if err != nil {
		println(err.Error())
		return
	}
	println(m.Sql)
	//把sql转成可以执行的sql
	params := map[string]interface{}{
		"pwd": "222", "name": "",
	}
	sql1, sqlParams, err := GetExecSqlInfo(m.Sql, params)
	if err != nil {
		println(err.Error())
		return
	}
	println(sql1)
	println(sqlParams)

}

//验证xml
func Test_Validate(t *testing.T) {
	ROOT, _ := os.Getwd()
	path := ROOT + "/demo/mapper/aaMapper.xml"
	file, err := os.Open(path)
	if nil != err {
		panic("Open mapper config: " + path + " err:" + err.Error())
	}
	defer file.Close()
	parser := xml.NewDecoder(file)
	startTag := make([]string, 1)
	endTag := make([]string, 1)
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement: //tag start
			elmt1 := xml.StartElement(t)
			name := elmt1.Name.Local
			if name != "sqlMapper" && name != "update" && name != "select" && name != "delete" && name != "insert" {
				panic("sql mapper tag error,please check your xml  " + path)
			} else {
				startTag = append(startTag, name)
			}
		case xml.EndElement: //tag end
			elmt2 := xml.EndElement(t)
			name := elmt2.Name.Local
			if name != "sqlMapper" && name != "update" && name != "select" && name != "delete" && name != "insert" {
				panic("sql mapper tag error,please check your xml " + path)
			} else {
				endTag = append(endTag, name)
			}
		}
	}
	if len(startTag) != len(endTag) {
		panic("xml error, there is tag no close, please check your xml config!")
	}
}
