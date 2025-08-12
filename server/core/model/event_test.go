package model_test

import (
	"testing"
	"time"

	"github.com/jcfug8/daylear/server/core/model"
)

var fixedNow = time.Date(2025, time.August, 10, 0, 0, 0, 0, time.UTC)

func TestEvent_GenerateClones_NoRRule(t *testing.T) {
	now := fixedNow

	event := model.Event{
		StartTime: now.Add(12 * time.Hour),
	}

	clones, err := event.GenerateClones(now, now.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 0 {
		t.Fatalf("expected 0 clone, got %d", len(clones))
	}
}

func TestEvent_GenerateClones_DailyRepeatCount2(t *testing.T) {
	now := fixedNow

	expectedStartTimes := []time.Time{
		now.Add(12 * time.Hour * 3),
	}

	event := model.Event{
		StartTime:      now.Add(12 * time.Hour),
		RecurrenceRule: &[]string{"FREQ=DAILY;COUNT=2"}[0],
	}

	clones, err := event.GenerateClones(now, now.Add(24*time.Hour*365))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 1 {
		t.Fatalf("expected 1 clone, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_DailyRepeatCount10_1ExcludedDates(t *testing.T) {
	now := fixedNow

	expectedStartTimes := []time.Time{
		// now.Add(12 * time.Hour*3),
		now.Add(12 * time.Hour * 5),
		now.Add(12 * time.Hour * 7),
		now.Add(12 * time.Hour * 9),
		now.Add(12 * time.Hour * 11),
		now.Add(12 * time.Hour * 13),
		now.Add(12 * time.Hour * 15),
		now.Add(12 * time.Hour * 17),
		now.Add(12 * time.Hour * 19),
	}

	event := model.Event{
		StartTime:      now.Add(12 * time.Hour),
		RecurrenceRule: &[]string{"FREQ=DAILY;COUNT=10"}[0],
		ExcludedDates:  []time.Time{now.Add(36 * time.Hour)},
	}

	clones, err := event.GenerateClones(now, now.Add(24*time.Hour*365))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 8 {
		t.Fatalf("expected 8 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_DailyRepeatCount10_1AdditionalDates(t *testing.T) {
	now := fixedNow

	expectedStartTimes := []time.Time{
		now.Add(12 * time.Hour * 3),
		now.Add(12 * time.Hour * 4),
		now.Add(12 * time.Hour * 5),
		now.Add(12 * time.Hour * 7),
		now.Add(12 * time.Hour * 9),
		now.Add(12 * time.Hour * 11),
		now.Add(12 * time.Hour * 13),
		now.Add(12 * time.Hour * 15),
		now.Add(12 * time.Hour * 17),
		now.Add(12 * time.Hour * 19),
	}

	event := model.Event{
		StartTime:       now.Add(12 * time.Hour),
		RecurrenceRule:  &[]string{"FREQ=DAILY;COUNT=10"}[0],
		AdditionalDates: []time.Time{now.Add(48 * time.Hour)},
	}

	clones, err := event.GenerateClones(now, now.Add(24*time.Hour*365))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 10 {
		t.Fatalf("expected 10 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_WeeklyRepeat(t *testing.T) {
	now := fixedNow
	// Start on a Monday
	startTime := now.AddDate(0, 0, -int(now.Weekday())+1) // Monday
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 9, 0, 0, 0, startTime.Location())

	expectedStartTimes := []time.Time{
		startTime.AddDate(0, 0, 7),  // Next Monday
		startTime.AddDate(0, 0, 14), // Monday after next
	}

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=WEEKLY;COUNT=3"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_MonthlyRepeat(t *testing.T) {
	now := fixedNow
	// Start on the 15th of current month
	startTime := time.Date(now.Year(), now.Month(), 15, 10, 0, 0, 0, now.Location())

	expectedStartTimes := []time.Time{
		startTime.AddDate(0, 1, 0), // Next month 15th
		startTime.AddDate(0, 2, 0), // Month after next 15th
	}

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=MONTHLY;COUNT=3"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 6, 0))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_YearlyRepeat(t *testing.T) {
	now := fixedNow
	// Start on January 1st
	startTime := time.Date(now.Year(), 1, 1, 12, 0, 0, 0, now.Location())

	expectedStartTimes := []time.Time{
		startTime.AddDate(1, 0, 0), // Next year January 1st
		startTime.AddDate(2, 0, 0), // Year after next January 1st
	}

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=YEARLY;COUNT=3"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(5, 0, 0))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_WithInterval(t *testing.T) {
	now := fixedNow
	startTime := now.Add(12 * time.Hour)

	expectedStartTimes := []time.Time{
		startTime.AddDate(0, 0, 2), // Every 2 days
		startTime.AddDate(0, 0, 4),
		startTime.AddDate(0, 0, 6),
	}

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=DAILY;INTERVAL=2;COUNT=4"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 10))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 3 {
		t.Fatalf("expected 3 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_WithByDay(t *testing.T) {
	now := fixedNow
	// Start on a Monday
	startTime := now.AddDate(0, 0, -int(now.Weekday())+1) // Monday
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 9, 0, 0, 0, startTime.Location())

	expectedStartTimes := []time.Time{
		startTime.AddDate(0, 0, 7),  // Next Monday
		startTime.AddDate(0, 0, 14), // Monday after next
	}

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=WEEKLY;BYDAY=MO;COUNT=3"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_InvalidRRule(t *testing.T) {
	now := fixedNow

	event := model.Event{
		StartTime:      now.Add(12 * time.Hour),
		RecurrenceRule: &[]string{"INVALID_RRULE_FORMAT"}[0],
	}

	_, err := event.GenerateClones(now, now.AddDate(0, 0, 30))
	if err == nil {
		t.Fatalf("expected error for invalid RRULE, got nil")
	}

	expectedError := "failed to parse recurrence rule"
	if !contains(err.Error(), expectedError) {
		t.Fatalf("expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEvent_GenerateClones_EmptyRRule(t *testing.T) {
	now := fixedNow

	event := model.Event{
		StartTime:      now.Add(12 * time.Hour),
		RecurrenceRule: &[]string{""}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 0 {
		t.Fatalf("expected 0 clones for empty RRULE, got %d", len(clones))
	}
}

func TestEvent_GenerateClones_EventBeforeRange(t *testing.T) {
	now := fixedNow
	startTime := now.Add(-24 * time.Hour) // Event started yesterday

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=DAILY;COUNT=5"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	// Should get 4 clones (COUNT=5 total - 1 original, original is before search range)
	if len(clones) != 4 {
		t.Fatalf("expected 4 clones, got %d", len(clones))
	}

	// First clone should be today
	expectedFirstClone := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), startTime.Nanosecond(), startTime.Location())
	if clones[0].StartTime.Unix() != expectedFirstClone.Unix() {
		t.Fatalf("expected first clone at %v, got %v", expectedFirstClone, clones[0].StartTime)
	}
}

func TestEvent_GenerateClones_EventAfterRange(t *testing.T) {
	now := fixedNow
	startTime := now.AddDate(0, 0, 10) // Event starts in 10 days

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=DAILY;COUNT=5"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	// Should get no clones since event starts after the search range
	if len(clones) != 0 {
		t.Fatalf("expected 0 clones, got %d", len(clones))
	}
}

func TestEvent_GenerateClones_MultipleExcludedDates(t *testing.T) {
	now := fixedNow
	startTime := now.Add(12 * time.Hour)

	excludedDates := []time.Time{
		now.Add(36 * time.Hour), // Day 3 (first clone)
		now.Add(60 * time.Hour), // Day 5 (second clone)
		now.Add(84 * time.Hour), // Day 7 (third clone)
	}

	expectedStartTimes := []time.Time{
		now.Add(12 * time.Hour * 9),  // Day 9
		now.Add(12 * time.Hour * 11), // Day 11
		now.Add(12 * time.Hour * 13), // Day 13
		now.Add(12 * time.Hour * 15), // Day 15
		now.Add(12 * time.Hour * 17), // Day 17
		now.Add(12 * time.Hour * 19), // Day 19
	}

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=DAILY;COUNT=10"}[0],
		ExcludedDates:  excludedDates,
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 365))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	// Should get 6 clones (10 total - 1 original - 3 excluded)
	if len(clones) != 6 {
		t.Fatalf("expected 6 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_MultipleAdditionalDates(t *testing.T) {
	now := fixedNow
	startTime := now.Add(12 * time.Hour)

	additionalDates := []time.Time{
		now.Add(48 * time.Hour), // 2 days from now at 00:00
		now.Add(72 * time.Hour), // 3 days from now at 00:00
		now.Add(96 * time.Hour), // 4 days from now at 00:00
	}

	expectedStartTimes := []time.Time{
		now.Add(36 * time.Hour),  // RRULE: Day 2 at 12:00
		now.Add(48 * time.Hour),  // RDATE: 2d at 00:00
		now.Add(60 * time.Hour),  // RRULE: Day 3 at 12:00
		now.Add(72 * time.Hour),  // RDATE: 3d at 00:00
		now.Add(84 * time.Hour),  // RRULE: Day 4 at 12:00
		now.Add(96 * time.Hour),  // RDATE: 4d at 00:00
		now.Add(108 * time.Hour), // RRULE: Day 5 at 12:00
		now.Add(132 * time.Hour), // RRULE: Day 6 at 12:00
		now.Add(156 * time.Hour), // RRULE: Day 7 at 12:00
		now.Add(180 * time.Hour), // RRULE: Day 8 at 12:00
		now.Add(204 * time.Hour), // RRULE: Day 9 at 12:00
		now.Add(228 * time.Hour), // RRULE: Day 10 at 12:00
	}

	event := model.Event{
		StartTime:       startTime,
		RecurrenceRule:  &[]string{"FREQ=DAILY;COUNT=10"}[0],
		AdditionalDates: additionalDates,
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 365))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	// Should get 12 clones (9 from RRULE + 3 additional)
	if len(clones) != 12 {
		t.Fatalf("expected 12 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_ComplexRRule(t *testing.T) {
	now := fixedNow
	// Start on a Monday
	startTime := now.AddDate(0, 0, -int(now.Weekday())+1) // Monday
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 9, 0, 0, 0, startTime.Location())

	expectedStartTimes := []time.Time{
		startTime.AddDate(0, 0, 14), // Every 2 weeks, Monday
		startTime.AddDate(0, 0, 28), // Every 2 weeks, Monday
	}

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=WEEKLY;INTERVAL=2;BYDAY=MO;COUNT=3"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 60))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(clones))
	}

	for i, clone := range clones {
		if clone.StartTime.Unix() != expectedStartTimes[i].Unix() {
			t.Fatalf("expected start time %v, got %v", expectedStartTimes[i], clone.StartTime)
		}
	}
}

func TestEvent_GenerateClones_PreserveDuration(t *testing.T) {
	now := fixedNow
	startTime := now.Add(12 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	expectedDuration := 2 * time.Hour

	event := model.Event{
		StartTime:      startTime,
		EndTime:        &endTime,
		RecurrenceRule: &[]string{"FREQ=DAILY;COUNT=3"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(clones))
	}

	for _, clone := range clones {
		if clone.EndTime == nil {
			t.Fatalf("expected clone to have end time")
		}

		cloneDuration := clone.EndTime.Sub(clone.StartTime)
		if cloneDuration != expectedDuration {
			t.Fatalf("expected duration %v, got %v", expectedDuration, cloneDuration)
		}
	}
}

func TestEvent_GenerateClones_PreserveFields(t *testing.T) {
	now := fixedNow
	startTime := now.Add(12 * time.Hour)

	parent := model.EventParent{
		UserId:     123,
		CircleId:   456,
		CalendarId: 789,
	}
	parentEventId := int64(999)

	event := model.Event{
		StartTime:      startTime,
		Parent:         parent,
		ParentEventId:  &parentEventId,
		RecurrenceRule: &[]string{"FREQ=DAILY;COUNT=3"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(clones))
	}

	for _, clone := range clones {
		// Check Parent field is copied
		if clone.Parent != parent {
			t.Fatalf("expected Parent %v, got %v", parent, clone.Parent)
		}

		// Check ParentEventId is copied
		if clone.ParentEventId == nil {
			t.Fatalf("expected ParentEventId to be copied")
		}
		if *clone.ParentEventId != parentEventId {
			t.Fatalf("expected ParentEventId %d, got %d", parentEventId, *clone.ParentEventId)
		}
	}
}

func TestEvent_GenerateClones_NoEndTime(t *testing.T) {
	now := fixedNow
	startTime := now.Add(12 * time.Hour)

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=DAILY;COUNT=3"}[0],
		// No EndTime set
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(clones))
	}

	for _, clone := range clones {
		// Clones should not have EndTime when original doesn't
		if clone.EndTime != nil {
			t.Fatalf("expected clone to not have end time, got %v", clone.EndTime)
		}
	}
}

func TestEvent_GenerateClones_OriginalEventExcluded(t *testing.T) {
	now := fixedNow
	startTime := now.Add(12 * time.Hour)

	event := model.Event{
		StartTime:      startTime,
		RecurrenceRule: &[]string{"FREQ=DAILY;COUNT=3"}[0],
	}

	clones, err := event.GenerateClones(now, now.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("failed to generate clones: %v", err)
	}

	if len(clones) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(clones))
	}

	// Verify that the original event start time is never included
	for _, clone := range clones {
		if clone.StartTime.Unix() == startTime.Unix() {
			t.Fatalf("clone should not have the same start time as original event")
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || contains(s[1:], substr))))
}
