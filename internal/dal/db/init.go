package db

import (
	"fmt"
	"github.com/Alf_Grindel/save/conf"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	dsn := fmt.Sprintf(constant.MySQLDefaultDsn, conf.DB.User, conf.DB.Password, conf.DB.Addr, conf.DB.Database)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		hlog.Fatal("can not connect database")
	}
}
