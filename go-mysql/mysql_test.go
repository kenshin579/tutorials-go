package go_mysql

import (
	"fmt"
	"testing"

	db "github.com/kenshin579/tutorials-go/go-mysql/common/database"
	"github.com/kenshin579/tutorials-go/go-mysql/config"
	"github.com/kenshin579/tutorials-go/go-mysql/model"
	"gorm.io/gorm"
)

func setup() *gorm.DB {
	cfg, err := config.New("config/config.yaml")
	if err != nil {
		panic(err)
	}

	return db.NewMysqlDB(cfg)
}

func Test_Create(t *testing.T) {
	db := setup()

	for i := 6; i <= 1000; i++ {
		sql := fmt.Sprintf("INSERT INTO `locations_flat`(`name`, `position`) "+
			"VALUES ('point_%d', ST_GeomFromText('POINT( %d %d )', 0));", i, i, i)
		db.Exec(sql)
	}
	db.Commit()
}

func Test_Create_4329(t *testing.T) {
	db := setup()

	for i := 1; i <= 1000; i++ {
		sql := fmt.Sprintf("INSERT INTO `locations_earth`(`name`, `position`) "+
			"VALUES ('point_%d', ST_GeomFromText('POINT( %d %d )', 4326));", i, (90/1000)*i, (180/1000)*i)
		db.Exec(sql)
	}
	db.Commit()
}

func Test_Select(t *testing.T) {
	db := setup()
	rows, err := db.Raw("SELECT * from `study_db`.`locations_flat`;").Rows()
	fmt.Println(err)
	fmt.Println(rows)
	defer rows.Close()

	var result model.Location
	for rows.Next() {
		db.ScanRows(rows, &result)
		fmt.Printf("result: %+v\n", result)
	}

}

func Test_Select2(t *testing.T) {
	db := setup()
	var result model.Location

	db.Raw("SELECT * from `study_db`.`locations_flat`;").Scan(&result)
	fmt.Println(result)
}
