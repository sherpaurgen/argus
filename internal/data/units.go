package data

import "time"

type Unit struct {
	UnitID      int       `json:"unit_id" db:"unit_id"`
	BuildingID  string    `json:"building_id,string" db:"building_id"`
	OwnerID     int       `json:"owner_id" db:"owner_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	FloorNumber int       `json:"floor_number" db:"floor_number"`
	PricePerDay float64   `json:"price_per_day,omitempty" db:"price_per_day"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
