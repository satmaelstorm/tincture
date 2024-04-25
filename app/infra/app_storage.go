package infra

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/satmaelstorm/tincture/app/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
)

type TinctureDB struct {
	db *gorm.DB
}

func (t *TinctureDB) InitDb(uri fyne.URI) error {
	file, err := t.getStorage(uri)
	if err != nil {
		return fmt.Errorf("InitDb: %w", err)
	}
	exists, err := storage.Exists(file)
	if err != nil {
		return fmt.Errorf("InitDB: %w", err)
	}
	sqlDb, err := gorm.Open(sqlite.Open(file.Path()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	t.db = sqlDb
	if err != nil {
		return fmt.Errorf("InitDb: %w", err)
	}
	log.Printf("Opened DB: %s", file.Path())
	err = t.migrate()
	if err != nil {
		return fmt.Errorf("InitDb: %w", err)
	}
	if !exists {
		log.Println("Db not exists, initialize")
		initialData(t.db)
	}
	return nil
}

func (t *TinctureDB) GetReceipts() []domain.Receipt {
	var receipts []domain.Receipt
	result := t.db.Preload(clause.Associations).Find(&receipts)
	log.Printf("Loaded %d receipts\n", result.RowsAffected)
	return receipts
}

func (t *TinctureDB) GetTinctures() []domain.Tincture {
	var tinctures []domain.Tincture
	result := t.db.Find(&tinctures)
	log.Printf("Loaded %d tinctures\n", result.RowsAffected)
	return tinctures
}

func (t *TinctureDB) getStorage(uri fyne.URI) (fyne.URI, error) {
	db, err := storage.Child(uri, "tincture.db")
	if err != nil {
		return uri, fmt.Errorf("getStorage: %w", err)
	}
	return db, nil
}

func (t *TinctureDB) migrate() error {
	err := t.db.AutoMigrate(
		&domain.Receipt{},
		&domain.ReceiptItem{},
		&domain.Tincture{},
	)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
