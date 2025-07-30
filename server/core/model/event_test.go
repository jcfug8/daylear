package model_test

import (
	"testing"
	"time"

	"github.com/jcfug8/daylear/server/core/model"
)

func TestEvent_GenerateInstances(t *testing.T) {
	now := time.Now()

	event := model.EventSet{
		EventData: model.EventData{
			StartTime: now.Add(12 * time.Hour),
		},
	}

	instances, err := event.GenerateInstances(now, now.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("failed to generate instances: %v", err)
	}

	if len(instances) != 1 {
		t.Fatalf("expected 1 instance, got %d", len(instances))
	}
}
