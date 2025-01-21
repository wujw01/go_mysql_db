package main

import (
	"mysql"
	"strings"
)

const (
	userName = "root"
	pass = "root"
	host = "127.0.0.1"
	port = "3306"
	dbName = "student"
)

func main()  {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName,":",pass,"@tcp(",host,":",port, ")/",dbName,"?charset=utf8"}, "")//拼接字符串
	mysql.Init(path)
	//mysql.FetchRow("student2",2)
	//values := make(map[string]string)
	//values["name"]="wan2"
	//values["class_id"]="3"
	//fmt.Println(reflect.TypeOf(student1))
	//fmt.Println(mysql.Update("student2",2,values))
	//fmt.Println(mysql.FetchAll("student"))
	/*rows:=mysql.FetchAll("student")
	for _,vs:=range rows{
		for k,v:=range vs{
			fmt.Printf("%v====>%v\n",k,v)
		}
	}*/
	/*values:=make(map[string]string)
	values["student_hight"]="150cm"
	values["student_xihao"]="跳舞"
	mysql.Update("student",24,values)*/

}

