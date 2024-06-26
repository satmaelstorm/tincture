package infra

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/google/uuid"
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

func (t *TinctureDB) GetReceipt(uid uuid.UUID) (domain.Receipt, bool) {
	var receipt domain.Receipt
	result := t.db.Preload(clause.Associations).Take(&receipt, "uuid = ?", uid.String())
	if result.RowsAffected < 1 {
		return receipt, false
	}
	return receipt, true
}

func (t *TinctureDB) DeleteReceipt(receipt domain.Receipt) bool {
	result := t.db.Delete(receipt)
	return result.RowsAffected > 0
}

func (t *TinctureDB) SaveReceipt(receipt *domain.Receipt) {
	//так как мы зачищаем все ингриденты для простоты, когда берем их из формы, то и тут мы делаем тоже самое
	t.db.Delete(&domain.ReceiptItem{}, "receipt_uuid = ?", receipt.Uuid.String())
	t.db.Save(receipt)
}

func (t *TinctureDB) CreateReceipt(receipt *domain.Receipt) {
	t.db.Create(receipt)
}

func (t *TinctureDB) GetPreparingTinctures() []domain.Tincture {
	var tinctures []domain.Tincture
	result := t.db.Order("need_bottled_at asc, created_at asc").Find(&tinctures, "bottled_at IS NULL")
	log.Printf("Loaded %d tinctures\n", result.RowsAffected)
	return tinctures
}

func (t *TinctureDB) GetReadyTinctures() []domain.Tincture {
	var tinctures []domain.Tincture
	result := t.db.Order("ready_at asc, expired_at asc").Find(&tinctures, "bottled_at NOT NULL")
	log.Printf("Loaded %d tinctures\n", result.RowsAffected)
	return tinctures
}

func (t *TinctureDB) SaveTincture(tincture *domain.Tincture) {
	t.db.Save(tincture)
}

func (t *TinctureDB) CreateTincture(tincture *domain.Tincture) {
	t.db.Create(tincture)
}

func (t *TinctureDB) DeleteTincture(tincture *domain.Tincture) {
	t.db.Delete(tincture)
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
