package king_orm

import (
	"os"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	ROOT,_ := os.Getwd()
	path := ROOT+"/demo/mapper"
	v,_ := ReaderConfigBuilder(path)


	m,_ := v.GetMethodSql("demo.mapper.user.getuser")
	r := strings.NewReader(m.Sql)
	node  := parse(r)
	params := map[string]interface{}{
		"username": "majun",
	}
	sql  := CreateParamsSql(params,node.Elements...)
	sql1,sqlParams,_ :=ReadSQLParamsBySQL1(sql,params)
	println(sql)
	println(sql1)
	println(sqlParams)


}

