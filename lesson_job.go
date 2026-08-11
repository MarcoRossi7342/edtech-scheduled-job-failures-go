package main

import (
	"fmt"
	"runtime/debug"
)

func runLessonSync(capture func(map[string]any) error, jobName string, work func() error) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%v", recovered)
		}
		if err == nil {
			return
		}
		captureErr := capture(map[string]any{
			"message":         fmt.Sprintf("scheduled job %s failed: %v", jobName, err),
			"level":           "error",
			"fingerprint":     []string{"edtech", jobName},
			"exception":       string(debug.Stack()),
			"context":         map[string]any{"job": jobName, "domain": "lesson-sync"},
			"idempotency_key": "edtech-lesson-sync-" + jobName,
		})
		if captureErr != nil {
			err = fmt.Errorf("job failed and capture failed: %w", captureErr)
		}
	}()
	return work()
}

func syncLessons() error {
	return fmt.Errorf("lesson roster validation failed")
}
