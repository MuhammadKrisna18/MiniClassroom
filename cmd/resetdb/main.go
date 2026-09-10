package main

import (
	"fmt"
	"log"
	"siakad-pro/config"
	"siakad-pro/internal/modules/auth/domain"
	kelasDomain "siakad-pro/internal/modules/kelas/domain"
	mkDomain "siakad-pro/internal/modules/matakuliah/domain"
	psDomain "siakad-pro/internal/modules/programstudi/domain"
	semDomain "siakad-pro/internal/modules/semester/domain"
	"siakad-pro/internal/shared/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Dropping all tables...")
	err = db.Migrator().DropTable(
		&kelasDomain.PesertaKelas{},
		&kelasDomain.Absensi{},
		&kelasDomain.Pertemuan{},
		&semDomain.SemesterMataKuliah{},
		&semDomain.Semester{},
		&mkDomain.PengajuanMataKuliah{},
		&kelasDomain.PengajuanKelas{},
		&kelasDomain.Kelas{},
		&domain.EmailChangeRequest{},
		&psDomain.ProgramStudi{},
		&mkDomain.MataKuliah{},
		&domain.User{},
	)
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}
	log.Println("All tables dropped successfully.")

	log.Println("Re-initializing database (Migration & Seeding)...")
	_, err = database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	log.Println("Database reset successfully! All data is wiped except for the seeds.")
}
