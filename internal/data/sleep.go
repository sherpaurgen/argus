package data

import (
	"context"
	"database/sql"
	"math"
	"time"
)

type SleepDataModel struct {
	DB *sql.DB
}
type SleepData struct {
	SleepId      int64     `json:"sleep_id" db:"sleep_id"`
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

func (sdm SleepDataModel) GetSleepData(child_id string, user_id string, StartSleep time.Time, EndSleep time.Time, page int, limit int) (map[string]interface{}, error) {

	if StartSleep.Equal(time.Time{}) {
		//StartSleep.Equal(time.Time{})
		// Fetch latest 24 hours of metrics by default
		EndSleep = time.Now()

	}
	if EndSleep.Equal(time.Time{}) {
		StartSleep = EndSleep.Add(-24 * time.Hour)
	}
	//(page-1)xpagesize
	offset := (page - 1) * limit
	query := `
	SELECT sleep_id,child_id,user_id,sleep_start,sleep_end
	FROM sleep_patterns 
	WHERE user_id = $1 
	  AND child_id = $2 
	  AND sleep_start >= $3 
	  AND sleep_end <= $4 
	ORDER BY sleep_start ASC
	LIMIT $5 OFFSET $6
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := sdm.DB.QueryContext(ctx, query, user_id, child_id, StartSleep, EndSleep, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sleepData []*SleepData
	for rows.Next() {
		var sd SleepData
		//sleep_id,child_id,user_id,sleep_start,sleep_end
		err := rows.Scan(
			&sd.SleepId,
			&sd.ChildId,
			&sd.UserID,
			&sd.StartSleep,
			&sd.EndSleep,
		)
		if err != nil {
			return nil, err
		}
		sleepData = append(sleepData, &sd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	var total_rows int
	sqlstatement := `SELECT COUNT(*) FROM sleep_patterns 
	WHERE user_id = $1 
	  AND child_id = $2 
	  AND sleep_start >= $3 
	  AND sleep_end <= $4 ;`
	err = sdm.DB.QueryRow(sqlstatement, user_id, child_id, StartSleep, EndSleep).Scan(&total_rows)
	if err != nil {
		return nil, err
	}
	response := map[string]interface{}{
		"results":      sleepData,
		"current_page": page,
		"limit":        limit,
		"total_rows":   total_rows,
		"total_pages":  int(math.Ceil(float64(total_rows) / float64(limit))),
	}
	return response, nil
}
