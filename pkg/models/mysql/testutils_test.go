package mysql

import (
	"database/sql"
	"io/ioutil"
	"testing"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestDB(t *testing.T) (*gorm.DB, func()) {
	db, err := sql.Open("mysql", "test_web:pass@/test_servente?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	script, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	gormDB, err := gorm.Open(gormMysql.Open("test_web:pass@/test_servente?charset=utf8mb4&parseTime=true"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}

	gormsqlDB, err := gormDB.DB()
	if err != nil {
		t.Fatal(err)
	}

	return gormDB, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		gormsqlDB.Close()
	}
}
