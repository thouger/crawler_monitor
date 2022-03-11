package dbs

import (
	"container/list"
	"database/sql"
	"fmt"
	mysql_config "spider/utils/config/mysql_config"
	_log "spider/utils/log"
	"strings"
	"time"

	"github.com/cevaris/ordered_map"
	_ "github.com/go-sql-driver/mysql"
)

var MysqlDb *sql.DB
var MysqlDbErr error

const (
	Host     = mysql_config.Host
	Username = mysql_config.Username
	Password = mysql_config.Password
	Port     = mysql_config.Port
	Database = mysql_config.Database
	HARSET   = mysql_config.HARSET
)

// 初始化链接
var log = _log.Init_log(_log.Params{})

func Connect() {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", Username, Password, Host, Port, Database, HARSET)

	// 打开连接失败
	MysqlDb, MysqlDbErr = sql.Open("mysql", dbDSN)
	if MysqlDbErr != nil {
		log.Info("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + MysqlDbErr.Error())
	}

	// 最大连接数
	MysqlDb.SetMaxOpenConns(100)
	// 闲置连接数
	MysqlDb.SetMaxIdleConns(20)
	// 最大连接周期
	MysqlDb.SetConnMaxLifetime(100 * time.Second)

	if MysqlDbErr = MysqlDb.Ping(); nil != MysqlDbErr {
		panic("数据库链接失败: " + MysqlDbErr.Error())
	}
}

func Close() {
	MysqlDb.Close()
}

// 用户表结构体
type User struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func Make_sql(sql_type string, table_name string, field_num int, field_value []interface{}) string {
	sql := fmt.Sprintf("%s into %s (", sql_type, table_name) + strings.Repeat("'%s'", field_num) + ")"
	new_sql := fmt.Sprintf(sql, field_value...)
	return new_sql
}

func Execute(sql string) int64 {
	res, err := MysqlDb.Exec(sql)
	if err != nil {
		fmt.Printf("执行sql错误,sql:%s\n,错误是%s\n", sql, err.Error())
	} else {
		id, err := res.LastInsertId()
		if err == nil {
			return id
		} else {
			println("Error:", err.Error())
		}
	}
	return -1
}

// 插入数据
func StructInsert() {

	ret, _ := MysqlDb.Exec("insert INTO users(name,age) values(?,?)", "小红", 23)

	//插入数据的主键id
	lastInsertID, _ := ret.LastInsertId()
	fmt.Println("LastInsertID:", lastInsertID)

	//影响行数
	rowsaffected, _ := ret.RowsAffected()
	fmt.Println("RowsAffected:", rowsaffected)

}

// 更新数据
func StructUpdate() {

	ret, _ := MysqlDb.Exec("UPDATE users set age=? where id=?", "100", 1)
	upd_nums, _ := ret.RowsAffected()

	fmt.Println("RowsAffected:", upd_nums)
}

// 删除数据
func StructDel() {

	ret, _ := MysqlDb.Exec("delete from users where id=?", 1)
	del_nums, _ := ret.RowsAffected()

	fmt.Println("RowsAffected:", del_nums)
}

// 事务处理,结合预处理
func StructTx() {

	//事务处理
	tx, _ := MysqlDb.Begin()

	// 新增
	userAddPre, _ := MysqlDb.Prepare("insert into users(name, age) values(?, ?)")
	addRet, _ := userAddPre.Exec("zhaoliu", 15)
	ins_nums, _ := addRet.RowsAffected()

	// 更新
	userUpdatePre1, _ := tx.Exec("update users set name = 'zhansan'  where name=?", "张三")
	upd_nums1, _ := userUpdatePre1.RowsAffected()
	userUpdatePre2, _ := tx.Exec("update users set name = 'lisi'  where name=?", "李四")
	upd_nums2, _ := userUpdatePre2.RowsAffected()

	fmt.Println(ins_nums)
	fmt.Println(upd_nums1)
	fmt.Println(upd_nums2)

	if ins_nums > 0 && upd_nums1 > 0 && upd_nums2 > 0 {
		tx.Commit()
	} else {
		tx.Rollback()
	}

}

// 查询数据，指定字段名,不采用结构体
func RawQueryField() {

	rows, _ := MysqlDb.Query("select id,name from users")
	if rows == nil {
		return
	}
	id := 0
	name := ""
	fmt.Println(rows)
	fmt.Println(rows)
	for rows.Next() {
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
}

// 查询数据,取所有字段,不采用结构体
func Select(sql string) *list.List {

	//查询数据，取所有字段
	rows2, err1 := MysqlDb.Query(sql)
	if err1 != nil {
		log.Fatalf("%s", err1)
	}

	//返回所有列
	cols, err2 := rows2.Columns()
	if err2 != nil {
		log.Fatalf("%s", err2)
	}

	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))

	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	result := list.New()
	for rows2.Next() {
		//填充数据
		rows2.Scan(scans...)
		//每行数据
		row := ordered_map.NewOrderedMap()

		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row.Set(key, string(v))
		}
		//放入结果集
		result.PushBack(row)
		i++
	}
	return result
}
