package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"siakad-pro/internal/modules/programstudi/domain"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	prodis := []domain.ProgramStudi{
		{Name: "Mata Kuliah Departemen", Code: "DEPT"},
		{Name: "Mata Kuliah Umum Bersama", Code: "MKUB"},
	}

	for _, p := range prodis {
		var existing domain.ProgramStudi
		err := db.WithContext(ctx).Where("code = ?", p.Code).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				p.ID = uuid.New().String()
				if err := db.WithContext(ctx).Create(&p).Error; err != nil {
					log.Fatal("Failed to create", p.Name, err)
				}
				fmt.Println("Inserted:", p.Name)
			} else {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Already exists:", p.Name)
		}
	}
}
