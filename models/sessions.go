package models

import (
	"context"
	"encoding/json"
	"fmt"
	"gymBackend/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/robfig/cron/v3"
)

type Session struct {
	UserID         int       `json:"user_id"`
	TrainerID      int       `json:"trainer_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	ScheduledTime  time.Time `json:"scheduled_time"`
	WorkoutDetails struct {
		Exercises []Exercise `json:"exercises"`
	} `json:"workout_details"`
	Status string `json:"status"`
}

type Exercise struct {
	Name            string `json:"name"`
	MuscleGroup     string `json:"muscle_group"`
	DurationMinutes int    `json:"duration_minutes"`
	Sets            int    `json:"sets"`
	RepsPerSet      int    `json:"reps_per_set"`
}

type WorkoutDetail struct {
	Exercises []Exercise `json:"exercises"`
}

var customTimeFormat = "2006-01-02 15:04:05"

func (session *Session) CreateSession(exercises []Exercise, UserId int) error {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	session.Status = "pending"
	// Convert exercises slice to JSON
	exercisesJSON, err := json.Marshal(WorkoutDetail{Exercises: exercises})
	if err != nil {
		return fmt.Errorf("error marshalling exercises: %v", err)
	}

	sqlQuery := `
        INSERT INTO sessions (user_id, trainer_id, start_time, end_time, workout_details, scheduled_time, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`
	err = conn.QueryRow(context.Background(), sqlQuery, UserId, session.TrainerID, session.StartTime, session.EndTime, exercisesJSON, session.ScheduledTime, session.Status).Scan(&UserId)
	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}
	return nil
}

func MarkSessionCompleted(sessionID int) error {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	sqlQuery := `UPDATE sessions SET status = 'completed' WHERE id = $1`
	_, err = conn.Exec(context.Background(), sqlQuery, sessionID)
	if err != nil {
		return fmt.Errorf("error marking session as completed: %v", err)
	}

	return nil
}

func MarkMissedSessions() error {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	sqlQuery := `UPDATE sessions SET status = 'missed' WHERE status = 'pending' AND scheduled_time < NOW() - INTERVAL '1 hour'`
	_, err = conn.Exec(context.Background(), sqlQuery)
	if err != nil {
		return fmt.Errorf("error marking missed sessions: %v", err)
	}
	return nil
}

func ScheduleNotification(session Session) {
	c := cron.New()

	// Calculate notification time: 30 minutes before scheduled time
	notificationTime := session.ScheduledTime.Add(-30 * time.Minute)

	// Format notificationTime for CRON expression
	cronSchedule := fmt.Sprintf("%d %d %d %d %d", notificationTime.Minute(), notificationTime.Hour(), notificationTime.Day(), notificationTime.Month(), notificationTime.Weekday())

	// Schedule notification
	c.AddFunc(cronSchedule, func() {
		fmt.Printf("Notification: Training session at %s\n", session.ScheduledTime.Format(customTimeFormat))
	})

	c.Start()
}

func GetSession(sessionID int) (*Session, error) {
	dbURL, _ := utils.CheckDbConnection()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return &Session{}, fmt.Errorf("UNABLE TO CONNECT TO DATABASE: %v", err)
	}
	defer conn.Close(context.Background())

	var session Session

	err = conn.QueryRow(context.Background(), `
		SELECT id, user_id, trainer_id, start_time, end_time, workout_details, scheduled_time, status
		FROM sessions
		WHERE id = $1`, sessionID).Scan(
		&session.TrainerID,
		&session.StartTime,
		&session.EndTime,
		&session.WorkoutDetails,
		&session.ScheduledTime,
		&session.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("error retrieving session: %v", err)
	}

	return &session, nil
}
