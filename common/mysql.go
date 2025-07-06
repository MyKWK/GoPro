package common

import (
	"awesomeProject/datamodels"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

func NewMysqlConn() (*gorm.DB, error) {
	username := "root"
	password := "thms1368"
	hostname := "localhost"
	port := "3306"
	database := "imooc"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, hostname, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&datamodels.User{}, &datamodels.Product{}, &datamodels.Order{}); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MySQL via GORM (AutoMigrate applied)")
	return db, nil
}

func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	//返回所有列
	columns, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(columns))
	//这里表示一行填充数据
	scans := make([]interface{}, len(columns))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := columns[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result
}
