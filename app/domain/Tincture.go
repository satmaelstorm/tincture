package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"strconv"
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

func (t *Tincture) ReadyToBottle(now time.Time) float64 {
	return readiness(t.CreatedAt, t.NeedBottledAt, now)
}

func (t *Tincture) ReadyToDrink(now time.Time) float64 {
	return readiness(t.NeedBottledAt, t.ReadyAt, now)
}

func (t *Tincture) IsExpire(now time.Time) bool {
	return t.ExpiredAt.Before(now)
}

func (t *Tincture) IsNearExpire(now time.Time) bool {
	return t.ExpiredAt.Add(time.Hour * 24 * 14).Before(now)
}

func (t *Tincture) NeedToBottled(now time.Time) bool {
	return !t.BottledAt.Valid && t.NeedBottledAt.Before(now)
}

func (t *Tincture) IsBottled() bool {
	return t.BottledAt.Valid
}

func (t *Tincture) IsReady(now time.Time) bool {
	return t.ReadyAt.Before(now)
}

func (t *Tincture) Bottled(now time.Time) {
	t.BottledAt.Valid = true
	t.BottledAt.Time = now
}

func NewTincture(
	name string,
	createdAt string,
	daysToBottle string,
	daysToRest string,
	daysExpire string,
) Tincture {
	created, err := time.Parse(time.DateOnly, createdAt)
	if err != nil {
		created = time.Now()
	}
	dtb, _ := strconv.Atoi(daysToBottle)
	dtr, _ := strconv.Atoi(daysToRest)
	de, _ := strconv.Atoi(daysExpire)

	toBottle := created.Add(time.Hour * 24 * time.Duration(dtb))
	toReady := toBottle.Add(time.Hour * 24 * time.Duration(dtr))
	toExpire := toBottle.Add(time.Hour * 24 * time.Duration(de))

	return Tincture{
		Uuid:          uuid.Must(uuid.NewV7()),
		Name:          name,
		CreatedAt:     created,
		NeedBottledAt: toBottle,
		ReadyAt:       toReady,
		ExpiredAt:     toExpire,
	}
}
