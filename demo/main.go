package main

import (
	king_orm "king.plugin/king-orm"
	"os"
)

type userInfo struct {
	username string
	id int
	ids string
}


func main()  {
	ROOT,_ := os.Getwd()
	path := ROOT+"/king-orm/demo/mapper/userMapper.xml"
	v,_ := king_orm.ReaderConfigBuilder(path)
	m,_ := v.GetMethodSql("demo.mapper.getuser")
	println(m.Sql)
	/*m,err  := v.GetMethodSql("demo.mapper.update")
	if err!=nil {
		fmt.Print(err)
		return
	}
	println(m.Sql)


	user := make(map[string]interface{})
	user["username"] = "majun"
	user["id"] = 1
	user["params"] = "1,2,3,3"
	newSql := king_orm.DollarTokenHandler(m.Sql,user)
	sql,sqlParams,err := king_orm.ReadSQLParamsBySQL1(newSql,user)
	if err!=nil {
		println(err)
		return
	}
	println(sql)
	println(len(sqlParams))*/
}