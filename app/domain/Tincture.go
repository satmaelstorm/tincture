package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Tincture struct {
	Uuid          uuid.UUID `gorm:"primarykey"`
	Name          string
	CreatedAt     time.Time
	NeedBottledAt time.Time
	BottledAt     sql.NullTime
	ReadyAt       time.Time
	ExpiredAt     time.Time
}

func (t Tincture) ReadyToBottle(now time.Time) float64 {
	return readyness(t.CreatedAt, t.NeedBottledAt, now)
}

func (t Tincture) ReadyToDrink(now time.Time) float64 {
	return readyness(t.NeedBottledAt, t.ReadyAt, now)
}

func (t Tincture) IsExpire(now time.Time) bool {
	return t.ExpiredAt.Before(now)
}

func (t Tincture) NeedToBottled(now time.Time) bool {
	return !t.BottledAt.Valid && t.NeedBottledAt.Before(now)
}
