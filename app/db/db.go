package db

import (
	"log"

	"fmt"

	"github.com/elBroom/meteo/app/config"
	"github.com/elBroom/meteo/app/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var sql_connect (*gorm.DB) = nil

func Sql_connect() *gorm.DB {
	if sql_connect != nil {
		return sql_connect
	}

	cfg := config.GetSql()
	conn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password)
	var err error
	sql_connect, err = gorm.Open("postgres", conn)
	//defer connect.Close() in main  =(

	if err != nil {
		log.Fatal(err)
	}

	sql_connect.LogMode(true)

	sql_connect.AutoMigrate(&model.Disignation{}, &model.Indication{})

	return sql_connect
}
