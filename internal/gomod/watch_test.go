package gomod

import (
	"sync"
	"testing"
	"time"
)

func sampleDepsForWatch(ver string) []Dependency {
	return []Dependency{
		{Module: "github.com/foo/bar", Version: ver},
		{Module: "github.com/baz/qux", Version: "v1.0.0"},
	}
}

func TestDefaultWatchConfig(t *testing.T) {
	cfg := DefaultWatchConfig("main")
	if cfg.Branch != "main" {
		t.Errorf("expected branch main, got %s", cfg.Branch)
	}
	if cfg.IntervalSecs != 60 {
		t.Errorf("expected interval 60, got %d", cfg.IntervalSecs)
	}
}

func TestRunWatch_InvalidInterval(t *testing.T) {
	cfg := DefaultWatchConfig("main")
	cfg.IntervalSecs = 0
	err := RunWatch(cfg, func(string) ([]Dependency, error) { return nil, nil }, make(chan struct{}))
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestRunWatch_NilLoader(t *testing.T) {
	cfg := DefaultWatchConfig("main")
	err := RunWatch(cfg, nil, make(chan struct{}))
	if err == nil {
		t.Fatal("expected error for nil loader")
	}
}

func TestRunWatch_StopsOnSignal(t *testing.T) {
	stop := make(chan struct{})
	cfg := DefaultWatchConfig("watch-stop-test")
	cfg.IntervalSecs = 1

	loader := func(branch string) ([]Dependency, error) {
		return sampleDepsForWatch("v1.0.0"), nil
	}

	done := make(chan error, 1)
	go func() {
		done <- RunWatch(cfg, loader, stop)
	}()

	time.Sleep(50 * time.Millisecond)
	close(stop)

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("watch did not stop in time")
	}
}

func TestRunWatch_FiresOnChange(t *testing.T) {
	stop := make(chan struct{})
	cfg := DefaultWatchConfig("watch-change-test")
	cfg.IntervalSecs = 1

	call := 0
	var mu sync.Mutex
	var events []WatchEvent

	cfg.OnChange = func(e WatchEvent) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	}

	loader := func(branch string) ([]Dependency, error) {
		mu.Lock()
		n := call
		call++
		mu.Unlock()
		if n == 0 {
			return sampleDepsForWatch("v1.0.0"), nil
		}
		return sampleDepsForWatch("v2.0.0"), nil
	}

	go func() {
		_ = RunWatch(cfg, loader, stop)
	}()

	time.Sleep(1600 * time.Millisecond)
	close(stop)

	mu.Lock()
	got := len(events)
	mu.Unlock()

	if got == 0 {
		t.Error("expected at least one change event")
	}
}
