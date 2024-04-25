package infra

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/satmaelstorm/tincture/app/domain"
	"gorm.io/gorm"
	"time"
)

func initialData(db *gorm.DB) {
	r := receipts()
	db.CreateInBatches(r, len(r))
	t := times()
	db.CreateInBatches(t, len(t))
}

func times() []domain.Tincture {
	return []domain.Tincture{
		{
			Uuid:          uuid.Must(uuid.NewV7()),
			Name:          "Хвойная черная смородина",
			CreatedAt:     time.Now().Add(-(time.Hour * 24)),
			NeedBottledAt: time.Now().Add(time.Hour * 24 * 6),
			BottledAt: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
			ReadyAt:   time.Now().Add(time.Hour * 24 * 11),
			ExpiredAt: time.Now().Add(time.Hour * 24 * (365 + 11)),
		},
		{
			Uuid:          uuid.Must(uuid.NewV7()),
			Name:          "Хвойная черная смородина",
			CreatedAt:     time.Now().Add(-(time.Hour * 24 * 10)),
			NeedBottledAt: time.Now().Add(-(time.Hour * 24 * 4)),
			BottledAt: sql.NullTime{
				Time:  time.Now().Add(-(time.Hour * 24 * 4)),
				Valid: true,
			},
			ReadyAt:   time.Now().Add(time.Hour * 24 * 2),
			ExpiredAt: time.Now().Add(time.Hour * 24 * (365 + 2)),
		},
	}
}

func receipts() []domain.Receipt {
	return []domain.Receipt{
		makeReceipt(
			"Хвойная Черная Смородина",
			"Ягодный вкус с хвойными нотками. \n "+
				"Приготовление: Ягоды разморозить в холодильнике. \n"+
				"Ингриденты смешать, залить 1л крепкого алкоголя (40гр) и настаивать неделю, после перелить в бутылку и дать постоять еще неделю.",
			[]domain.ReceiptItem{
				{
					Name:     "Водка",
					Quantity: "1 л",
				},
				{
					Name:     "Черная смородина замороженная",
					Quantity: "0.5 кг",
				},
				{
					Name:     "Розмарин",
					Quantity: "5 г",
				},
				{
					Name:     "Сахар",
					Quantity: "150 г",
				},
			},
		),
		makeReceipt(
			"Перцовка",
			"Острая перцовка",
			[]domain.ReceiptItem{
				{
					Name:     "Водка",
					Quantity: "1 л",
				},
				{
					Name:     "Смесь перцев",
					Quantity: "7 г",
				},
				{
					Name:     "Сахар",
					Quantity: "3 г",
				},
			},
		),
	}
}

func makeReceipt(
	title string,
	desc string,
	items []domain.ReceiptItem,
) domain.Receipt {
	uid := uuid.Must(uuid.NewV7())
	for i, item := range items {
		item.ReceiptUuid = uid
		items[i] = item
	}
	return domain.Receipt{
		Uuid:        uid,
		Title:       title,
		Description: desc,
		Items:       items,
	}
}
