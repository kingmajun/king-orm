package king_orm

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
	"time"
)

//加载sqlMapper.xml配置文件
//xmlconfigBuilder是加载xml文件入口

//select、insert、update、delete标签内容结构体
type methodObj struct{
	Id string `xml:"id,attr"`
	SqlXml string `xml:",innerxml"`//读取该标签先所有内容，包含xml标签内容
}


//xml内容全部
type sqlMapper struct {
	XMLName xml.Name `xml:"sqlMapper"`
	Namespace string `xml:"namespace,attr"`
	Selects  []methodObj `xml:"select"`
	Inserts  []methodObj `xml:"insert"`
	Updates  []methodObj `xml:"update"`
	Deletes  []methodObj `xml:"delete"`
}

type Method struct {
	Namespace string
	Id string
	Sql string
}

type Osm struct {
	sqlMapperMap map[string]*Method
}



func xmlBuilder(xmlPath string,ch chan sqlMapper) {
	byte,err := ioutil.ReadFile(xmlPath)
	if err !=nil{
		logrus.Error("not file error:",err.Error())
		return
	}
	s := &sqlMapper{}
	err = xml.Unmarshal(byte, &s)
	if err != nil {
		logrus.Error("reader xml error:",err)
		return
	}
	ch <- *s
	return
}

//读取pathDir目录下所有SQLMapper路径
func readerMapperPath(pathDir string)(xmlPath []string,err error){
	files, err := ioutil.ReadDir(pathDir)
	if err!=nil{
		fmt.Printf("pathdir read error \"%v\"", err)
		return
	}
	for _,file := range files{
		if !file.IsDir()  {
			idex := strings.Index(file.Name(),"Mapper")
			filesuffix := path.Ext(file.Name())
			if (idex >0 && filesuffix==".xml"){
				xmlPath = append(xmlPath,pathDir+"/"+file.Name())
			}
		}
	}
	return
}

//加载xml
func ReaderConfigBuilder(pathDir string)(osm *Osm,err error){

	osm  = new(Osm)

	xmlPath,err := readerMapperPath(pathDir)
	if err !=nil{
		fmt.Printf("not found xml  file")
		return
	}

	chs := make([] chan sqlMapper, len(xmlPath))
	defer func() {
		for _, c := range chs {
			if c != nil {
				close(c)
			}
		}
	}()

	for i:=0;i<len(xmlPath);i++{
		chs[i] = make(chan sqlMapper)
		go xmlBuilder(xmlPath[i],chs[i])
	}

	sqlMapperMap := make(map[string]*Method)
	// 获取结果
	for _, ch := range chs {
		mapper := <-ch

		for _,v := range mapper.Selects{
			key := fmt.Sprintf("%v.%v",mapper.Namespace,v.Id)
			v1 := new(Method)
			v1.Sql = fmt.Sprintf("<select id=\"%v\">%v</select>",v.Id,v.SqlXml)
			v1.Namespace = mapper.Namespace
			v1.Id = v.Id
			sqlMapperMap[key] = v1
		}
		for _,v := range mapper.Inserts{
			key := fmt.Sprintf("%v.%v",mapper.Namespace,v.Id)
			v1 := new(Method)
			v1.Sql = fmt.Sprintf("<insert id=\"%v\">%v</insert>",v.Id,v.SqlXml)
			v1.Namespace = mapper.Namespace
			v1.Id = v.Id
			sqlMapperMap[key] = v1
		}

		for _,v := range mapper.Deletes{
			key := fmt.Sprintf("%v.%v",mapper.Namespace,v.Id)
			v1 := new(Method)
			v1.Sql =  fmt.Sprintf("<delete id=\"%v\">%v</delete>",v.Id,v.SqlXml)
			v1.Namespace = mapper.Namespace
			v1.Id = v.Id
			sqlMapperMap[key] = v1
		}

		for _,v := range mapper.Updates{
			key := fmt.Sprintf("%v.%v",mapper.Namespace,v.Id)
			v1 := new(Method)
			v1.Sql =  fmt.Sprintf("<update id=\"%v\">%v</update>",v.Id,v.SqlXml)
			v1.Namespace = mapper.Namespace
			v1.Id = v.Id
			sqlMapperMap[key] = v1
		}
	}

	osm.sqlMapperMap = sqlMapperMap

	return osm,nil
}

//组装sqlmapper
func conversSqlMapper(mapper *sqlMapper)(sqlMapperMap map[string]*Method,err error){
	sqlMapperMap = make(map[string]*Method)
	for _,v := range mapper.Selects{
		key := fmt.Sprintf("%v.%v",mapper.Namespace,v.Id)
		v1 := new(Method)
		v1.Sql = fmt.Sprintf("<select id=\"%v\">%v</select>",v.Id,v.SqlXml)
		v1.Namespace = mapper.Namespace
		v1.Id = v.Id
		sqlMapperMap[key] = v1
	}
	for _,v := range mapper.Inserts{
		key := fmt.Sprintf("%v.%v",mapper.Namespace,v.Id)
		v1 := new(Method)
		v1.Sql = fmt.Sprintf("<insert id=\"%v\">%v</insert>",v.Id,v.SqlXml)
		v1.Namespace = mapper.Namespace
		v1.Id = v.Id
		sqlMapperMap[key] = v1
	}

	for _,v := range mapper.Deletes{
		key := fmt.Sprintf("%v.%v",mapper.Namespace,v.Id)
		v1 := new(Method)
		v1.Sql =  fmt.Sprintf("<delete id=\"%v\">%v</delete>",v.Id,v.SqlXml)
		v1.Namespace = mapper.Namespace
		v1.Id = v.Id
		sqlMapperMap[key] = v1
	}

	for _,v := range mapper.Updates{
		key := fmt.Sprintf("%v.%v",mapper.Namespace,v.Id)
		v1 := new(Method)
		v1.Sql =  fmt.Sprintf("<update id=\"%v\">%v</update>",v.Id,v.SqlXml)
		v1.Namespace = mapper.Namespace
		v1.Id = v.Id
		sqlMapperMap[key] = v1
	}

	return
}




func (o *Osm)GetMethodSql(id string)(v *Method,err error){
	v,ok := o.sqlMapperMap[id]
	if !ok {
		//创建异常
		err = errors.New(fmt.Sprintf("%v id is not found \n",id))
 		return
	}
	return
}

type sqlFragment struct {
	content     string
	paramValue  interface{}
	paramValues []interface{}
	isParam     bool
	isIn        bool
}
//解析sql，返回sql和对应的参数值
//sqlOrg=select * from im_user where username = #{username}
func ReadSQLParamsBySQL1(sqlOrg string, params ...interface{}) (sql string, sqlParams []interface{}, err error) {
	var param interface{}
	paramsSize := len(params)
	if paramsSize > 0 {
		if paramsSize == 1 {
			param = params[0]
		} else {
			param = params
		}

		//sql start
		sqls := []*sqlFragment{}
		paramNames := []*sqlFragment{}

		sqlTemp := sqlOrg
		errorIndex := 0
		for strings.Contains(sqlTemp, "#{") {
			si := strings.Index(sqlTemp, "#{")
			lastSQLText := sqlTemp[0:si]
			sqls = append(sqls, &sqlFragment{
				content: lastSQLText,
			})
			sqlTemp = sqlTemp[si+2:]
			errorIndex += si + 2

			ei := strings.Index(sqlTemp, "}")
			if ei != -1 {
				pni := &sqlFragment{
					content: strings.TrimSpace(sqlTemp[0:ei]),
					isParam: true,
					isIn:    sqlIsIn(lastSQLText),
				}
				sqls = append(sqls, pni)
				paramNames = append(paramNames, pni)
				sqlTemp = sqlTemp[ei+1:]
				errorIndex += ei + 1
			} else {
				fmt.Printf("sql read error \"%v\"", sqlOrg)
				return
			}
		}
		sqls = append(sqls, &sqlFragment{
			content: sqlTemp,
		})
		//sql end

		v := reflect.ValueOf(param)

		kind := v.Kind()
		switch {
		case kind == reflect.Array || kind == reflect.Slice:
			if len(paramNames) == 1 && paramNames[0].isIn {
				setDataToParamName(paramNames[0], v)
			} else {
				for i := 0; i < v.Len() && i < len(paramNames); i++ {
					vv := v.Index(i)
					if vv.IsValid() {
						setDataToParamName(paramNames[i], v.Index(i))
					}
				}
			}
		case kind == reflect.Map:
			for _, paramName := range paramNames {
				vv := v.MapIndex(reflect.ValueOf(paramName.content))
				if vv.IsValid() {
					setDataToParamName(paramName, vv)
				} else {
					err = fmt.Errorf("sql '%s' error : Key '%s' no exist", sqlOrg, paramName.content)
					return
				}
			}
		case kind == reflect.Struct:
			for _, paramName := range paramNames {
				firstChar := paramName.content[0]
				if firstChar < 'A' || firstChar > 'Z' {
					err = fmt.Errorf("sql '%s' error : Field '%s' unexported", sqlOrg, paramName.content)
					return
				}
				vv := v.FieldByName(paramName.content)
				if vv.IsValid() {
					setDataToParamName(paramName, vv)
				} else {
					err = fmt.Errorf("sql '%s' error : Field '%s' no exist", sqlOrg, paramName.content)
					return
				}
			}
		case kind == reflect.Bool ||
			kind == reflect.Int ||
			kind == reflect.Int8 ||
			kind == reflect.Int16 ||
			kind == reflect.Int32 ||
			kind == reflect.Int64 ||
			kind == reflect.Uint ||
			kind == reflect.Uint8 ||
			kind == reflect.Uint16 ||
			kind == reflect.Uint32 ||
			kind == reflect.Uint64 ||
			kind == reflect.Uintptr ||
			kind == reflect.Float32 ||
			kind == reflect.Float64 ||
			kind == reflect.Complex64 ||
			kind == reflect.Complex128 ||
			kind == reflect.String:
			for _, paramName := range paramNames {
				setDataToParamName(paramName, v)
			}
		default:
		}

		var sqlTexts []string

		for _, sql := range sqls {
			if sql.isParam {
				if sql.isIn {
					sqlTexts = append(sqlTexts, "(")
					for index, pv := range sql.paramValues {
						if index > 0 {
							sqlTexts = append(sqlTexts, ",")
						}
						sqlTexts = append(sqlTexts, "?")
						sqlParams = append(sqlParams, pv)
					}
					sqlTexts = append(sqlTexts, ")")
				} else {
					sqlTexts = append(sqlTexts, "?")
					sqlParams = append(sqlParams, sql.paramValue)
				}
			} else {
				sqlTexts = append(sqlTexts, sql.content)
			}
		}

		sql = strings.Join(sqlTexts, "")
	} else {
		sql = sqlOrg

	}
	return
}


func sqlIsIn(lastSQLText string) bool {
	lastSQLText = strings.TrimSpace(lastSQLText)
	lenLastSQLText := len(lastSQLText)
	if lenLastSQLText > 3 {
		return strings.EqualFold(lastSQLText[lenLastSQLText-3:], " IN")
	}
	return false
}

func timeFormat(t time.Time, format string) string {
	return t.Format(format)
}

var formatDateTime = "2006-01-02 15:04:05"

func setDataToParamName(paramName *sqlFragment, v reflect.Value) {
	if paramName.isIn {
		v = reflect.ValueOf(v.Interface())
		kind := v.Kind()
		if kind == reflect.Array || kind == reflect.Slice {
			for j := 0; j < v.Len(); j++ {
				vv := v.Index(j)
				if vv.Type().String() == "time.Time" {
					paramName.paramValues = append(paramName.paramValues, timeFormat(vv.Interface().(time.Time), formatDateTime))
				} else {
					paramName.paramValues = append(paramName.paramValues, vv.Interface())
				}
			}
		} else {
			if v.Type().String() == "time.Time" {
				paramName.paramValues = append(paramName.paramValues, timeFormat(v.Interface().(time.Time), formatDateTime))
			} else {
				paramName.paramValues = append(paramName.paramValues, v.Interface())
			}
		}
	} else {
		if v.Type().String() == "time.Time" {
			paramName.paramValue = timeFormat(v.Interface().(time.Time), formatDateTime)
		} else {
			paramName.paramValue = v.Interface()
		}
	}
}


// ${xx}处理
func  DollarTokenHandler(sqlStr string,params map[string]interface{})(sql string) {
 	if strings.Index(sqlStr, "$") == -1 {
		return
	}

	finalSqlStr := ""
	itemStr := ""
	start := 0
	for i := 0; i < len(sqlStr); i++ {
		if start > 0 {
			itemStr += string(sqlStr[i])
		}

		if i != 0 && i < len(sqlStr) {
			if string([]byte{sqlStr[i-1], sqlStr[i]}) == "${" {
				start = i
			}
		}

		if start != 0 && i < len(sqlStr)-1 && sqlStr[i+1] == '}' {
			finalSqlStr += sqlStr[:start-1]
			sqlStr = sqlStr[i+2:]

			itemStr = strings.Trim(itemStr, " ")

			item, ok := params[itemStr]
			if !ok {

				panic("params:" + itemStr + " not found")
			}

			finalSqlStr += item.(string)

			i = 0
			start = 0
			itemStr = ""
		}
	}

	finalSqlStr += sqlStr
	finalSqlStr = strings.Trim(finalSqlStr, " ")
	return finalSqlStr
}