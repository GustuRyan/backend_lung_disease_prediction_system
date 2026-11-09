package config

import (
	"reflect"
	"testing"
	"time"
)

func TestSplitCSV(t *testing.T) {
	t.Parallel()
	want := []string{"http://localhost:8501", "https://app.example.com"}
	got := splitCSV(" http://localhost:8501, ,https://app.example.com ")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitCSV() = %#v, want %#v", got, want)
	}
}

func TestDurationEnv(t *testing.T) {
	t.Setenv("TEST_DURATION", "15s")
	if got := durationEnv("TEST_DURATION", time.Second); got != 15*time.Second {
		t.Fatalf("durationEnv() = %s, want 15s", got)
	}
}

func TestLoadRequiresDatabaseDSN(t *testing.T) {
	t.Setenv("DATABASE_DSN", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want missing DATABASE_DSN error")
	}
}
