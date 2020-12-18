package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var DB *sql.DB

// 查询
func QueryRows(sentence string, args ...interface{}) ([]map[string]string, error) {
	var ret []map[string]string

	stmt, err := DB.Prepare(sentence)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 字段
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}

		var value string
		instance := make(map[string]string)
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			instance[columns[i]] = value
		}

		ret = append(ret, instance)
	}
	return ret, nil
}

// 插入/删除, 返回最后插入的ID和错误
func Execute(sentence string, args ...interface{}) (int64, error) {
	stmt, err := DB.Prepare(sentence)
	if err != nil {
		log.Fatal(err)
		return 0, nil
	}
	defer stmt.Close()

	rows, err := stmt.Exec(args...)
	if err != nil {
		log.Fatal(err)
		return 0, nil
	}
	return rows.LastInsertId()
}

func init() {
	var err error

	viper.SetConfigName("configure")
	viper.SetConfigType("json")
	viper.AddConfigPath("$GOPATH/src/github.com/service-computing-2020/bbs_backend/config/")
	viper.AddConfigPath("config/")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	dbstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", viper.GetString(`db.username`), viper.GetString(`db.password`), viper.GetString(`db.host`), viper.GetString(`db.port`), viper.GetString(`db.db`))
	DB, err = sql.Open("mysql", dbstr)
	fmt.Println(dbstr)
	if err != nil {
		log.Fatal("sql open error")
	}

	DB.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超时的连接就close
	DB.SetMaxOpenConns(100)                  //设置最大连接数
}
