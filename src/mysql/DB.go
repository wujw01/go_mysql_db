package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
)
//Db数据库连接池
var DB *sql.DB

//链接数据库
func Init(path string)*sql.DB  {
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)

	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil{
		fmt.Println("opon database fail")
		os.Exit(404)
	}
	//defer DB.Close()//关闭数据库资源
	return DB

}

//查询一条
func FetchRow(table string,id int)map[int]map[string]string{
	//查询数据，取所有字段
	sql := "select * from "+table+" where id ="+strconv.Itoa(id)

	rows2, _ := DB.Query(sql);
	//返回所有列
	cols, _ := rows2.Columns();//返回字段名

	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols));//存一个空的数组，有多少字段，空数组就有多少子数组[[],[],[]]

	//这里表示一行填充数据
	scans := make([]interface{}, len(cols));//定义接口存值

	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k];//把地址存在scans里面
	}

	i := 0;
	result := make(map[int]map[string]string);
	for rows2.Next() {
		//填充数据
		rows2.Scan(scans...);
		//每行数据
		row := make(map[string]string);
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k];
			//这里把[]byte数据转成string
			row[key] = string(v);
		}
		//放入结果集
		result[i] = row;
		i++;
	}
	//fmt.Println(reflect.TypeOf(result))
	//fmt.Println(result)
	return result
}

//查询所有
func FetchAll(table string)map[int]map[string]string {
	//查询数据，取所有字段
	sql := "select * from "+table

	rows2, _ := DB.Query(sql);
	//返回所有列
	cols, _ := rows2.Columns();//返回字段名

	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols));//存一个空的数组，有多少字段，空数组就有多少子数组[[],[],[]]

	//这里表示一行填充数据
	scans := make([]interface{}, len(cols));//定义接口存值

	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k];//把地址存在scans里面
	}

	i := 0;
	result := make(map[int]map[string]string);
	for rows2.Next() {
		//填充数据
		rows2.Scan(scans...);
		//每行数据
		row := make(map[string]string);
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k];
			//这里把[]byte数据转成string
			row[key] = string(v);
		}
		//放入结果集
		result[i] = row;
		i++;
	}
	return result
}

//添加一条数据
func Create(table string,value map[string]string)(int64,error) {

	//接受数据构造sql
	//循环数据
	//声明两个字符串，存储K,v值
	var ks string
	var vs  string
	for k,v:=range value{
		ks+=k+","
		vs+="'"+v+"'"+","
	}
	ks = strings.Trim(ks,",")//去除头部和尾部的“，”
	vs = strings.Trim(vs,",")

	sql :="insert into "+table+"("+ks+")"+"values("+vs+")"
	//开启事务
	tx, err := DB.Begin()
	if err != nil{
		fmt.Println("tx fail")
	}

	//准备sql语句
	stmt, err := tx.Prepare(sql)//sql语句写入

	if err != nil{
		fmt.Println("Prepare fail")
		os.Exit(1)
	}
	res, err := stmt.Exec()//执行sql
	if err != nil{
		fmt.Println("Exec fail")
		os.Exit(2)
	}

	//提交事务
	tx.Commit()
	//获得上一个insert的id
	//fmt.Println(res.LastInsertId())
	//UserId,err :=res.LastInsertId()
	return res.LastInsertId()
}

//删除一条数据
func Delete(table string,id int)(int64,error)  {
	//开启事务
	tx, err := DB.Begin()
	if err != nil{
		fmt.Println("tx fail")
	}
	sql :="delete from "+table+" where id="+strconv.Itoa(id)
	//准备sql语句
	stmt, err := tx.Prepare(sql)//sql语句写入
	if err != nil{
		fmt.Println("Prepare fail")
		os.Exit(1)
	}
	res, err := stmt.Exec(id)//执行sql
	if err != nil{
		fmt.Println("Exec fail")
		os.Exit(2)
	}

	//提交事务
	tx.Commit()
	//获得上一个insert的id
	return res.RowsAffected()//影响记录的行数

}

//编辑一条记录
func Update(table string,id int,value map[string]string)(int,int64)  {

	//循环数据
	//声明1个字符串，存储sql更新字符串
	var ks string

	for k,v:=range value{
		ks+= k+" = "+"'"+v+"'"+","
	}
	ks = strings.Trim(ks,",")//去除头部和尾部的“，”


	sql :="update "+table+" set "+ks+" where id = ?"
	/*fmt.Println(sql)
	os.Exit(1)*/
	//开启事务
	tx, err := DB.Begin()
	if err != nil{
		fmt.Println("tx fail")
	}
	//准备sql语句
	stmt, err := tx.Prepare(sql)//sql语句写入
	if err != nil{
		fmt.Println("Prepare fail")
		os.Exit(1)
	}
	res, err := stmt.Exec(id)//执行sql
	if err != nil{
		fmt.Println("Exec fail")
		os.Exit(2)
	}

	//提交事务
	tx.Commit()
	//获得上一个insert的id
	num,err:=res.RowsAffected()
	return id,num
}

