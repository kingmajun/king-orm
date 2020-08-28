package king_orm

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	ROOT,_ := os.Getwd()
	path := ROOT+"/demo/mapper"
	v,_ := ReaderConfigBuilder(path)

	//获取sql
	m,_ := v.GetMethodSql("demo.mapper.user.getuser")

	//把sql转成可以执行的sql
	params := map[string]interface{}{
		"username": "majun",
	}
	sql1,sqlParams,_ := GetExecSqlInfo(m.Sql,params)
	println(sql1)
	println(sqlParams)


}

