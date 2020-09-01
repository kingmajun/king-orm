package king_orm

import (
	"encoding/xml"
	"errors"
	"github.com/kingmajun/king-orm/model"
	"os"
)

//验证xml格式是否正确
func validateXml(xmlPath string) error  {
	file, err := os.Open(xmlPath)
	if nil != err {
		return errors.New("Open mapper config: " + xmlPath + " err:" + err.Error())
	}
	defer file.Close()
	parser := xml.NewDecoder(file)
	startTag := make([]string,1)
	endTag   := make([]string,1)
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement: //tag start
			startElmt := xml.StartElement(t)
			name := startElmt.Name.Local
			if isContain(model.XMLTAG,name) {
				startTag = append(startTag,name)
			}else{
				return errors.New("sql mapper tag error,please check your xml  "+xmlPath)
			}
		case xml.EndElement: //tag end
			endElmt := xml.EndElement(t)
			name := endElmt.Name.Local
			if isContain(model.XMLTAG,name) {
				endTag = append(endTag,name)
			}else{
				return errors.New("sql mapper tag error,please check your xml  "+xmlPath)
			}
		}
	}
	if len(startTag)!=len(endTag){
		return errors.New("xml error, there is tag no close, please check your xml config!")
	}
	return nil
}

//判断字符串是否在数组中
func isContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}