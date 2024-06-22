package data

import (
	"context"
	"database/sql"
	"time"
)

type SleepDataModel struct {
	DB *sql.DB
}
type SleepData struct {
	StartSleep   time.Time `json:"start_sleep" db:"start_sleep"`
	EndSleep     time.Time `json:"end_sleep" db:"end_sleep"`
	DeviceID     string    `json:"device_id" db:"device_id"`
	UserID       string    `json:"user_id" db:"user_id"`
	ChildId      string    `json:"child_id" db:"child_id"`
	SleepQuality int       `json:"sleep_quality" db:"sleep_quality"`
}

func (sdm SleepDataModel) Insert(sd *SleepData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		INSERT INTO sleep_patterns (sleep_start, sleep_end, device_id, child_id,user_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := sdm.DB.ExecContext(ctx, query, sd.StartSleep, sd.EndSleep, sd.DeviceID, sd.ChildId, sd.UserID)
	if err != nil {
		return err
	}
	return nil

}

func (sdm SleepDataModel) listSleepData(child *Child, StartSleep time.Time, EndSleep time.Time) {

}
