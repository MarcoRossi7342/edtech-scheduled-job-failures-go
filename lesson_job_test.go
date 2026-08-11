package main

import (
	"errors"
	"testing"
)

func TestRunLessonSyncCapturesByJob(t *testing.T) {
	var captured map[string]any
	err := runLessonSync(func(payload map[string]any) error {
		captured = payload
		return nil
	}, "lesson-sync", func() error { return errors.New("database check failed") })
	if err == nil || captured == nil {
		t.Fatal("expected job error and capture payload")
	}
	if captured["message"] != "scheduled job lesson-sync failed: database check failed" {
		t.Fatalf("unexpected message: %v", captured["message"])
	}
	fingerprint := captured["fingerprint"].([]string)
	if fingerprint[0] != "edtech" || fingerprint[1] != "lesson-sync" {
		t.Fatalf("unexpected fingerprint: %v", fingerprint)
	}
}
