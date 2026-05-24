package gomod

import (
	"strings"
	"testing"
)

func sampleNotifyPayload() NotifyPayload {
	return NotifyPayload{
		Branch: "feature/upgrade-deps",
		Summary: Summary{Added: 2, Removed: 1, Changed: 3},
		Score:   BreakingScore{Score: 72.5, Level: "high"},
		Alerts: []Alert{
			{Level: "critical", Message: "major version bump: github.com/foo/bar v1 -> v2"},
			{Level: "warning", Message: "dependency removed: github.com/old/pkg"},
		},
	}
}

func TestDefaultNotifyConfig(t *testing.T) {
	cfg := DefaultNotifyConfig()
	if cfg.Channel != ChannelStdout {
		t.Errorf("expected stdout channel, got %s", cfg.Channel)
	}
	if cfg.MinRisk != "medium" {
		t.Errorf("expected medium min risk, got %s", cfg.MinRisk)
	}
}

func TestFormatNotifyText_ContainsBranch(t *testing.T) {
	p := sampleNotifyPayload()
	out := FormatNotifyText(p)
	if !strings.Contains(out, "feature/upgrade-deps") {
		t.Error("expected branch name in text output")
	}
}

func TestFormatNotifyText_ContainsSummary(t *testing.T) {
	p := sampleNotifyPayload()
	out := FormatNotifyText(p)
	if !strings.Contains(out, "Added: 2") {
		t.Error("expected added count in text output")
	}
	if !strings.Contains(out, "Removed: 1") {
		t.Error("expected removed count in text output")
	}
}

func TestFormatNotifyText_ContainsAlerts(t *testing.T) {
	p := sampleNotifyPayload()
	out := FormatNotifyText(p)
	if !strings.Contains(out, "critical") {
		t.Error("expected critical alert level in text output")
	}
	if !strings.Contains(out, "major version bump") {
		t.Error("expected alert message in text output")
	}
}

func TestFormatNotifyMarkdown_ContainsBranch(t *testing.T) {
	p := sampleNotifyPayload()
	out := FormatNotifyMarkdown(p)
	if !strings.Contains(out, "feature/upgrade-deps") {
		t.Error("expected branch name in markdown output")
	}
}

func TestFormatNotifyMarkdown_ContainsScore(t *testing.T) {
	p := sampleNotifyPayload()
	out := FormatNotifyMarkdown(p)
	if !strings.Contains(out, "72.5") {
		t.Error("expected score in markdown output")
	}
	if !strings.Contains(out, "high") {
		t.Error("expected risk level in markdown output")
	}
}

func TestFormatNotifyMarkdown_NoAlerts(t *testing.T) {
	p := sampleNotifyPayload()
	p.Alerts = nil
	out := FormatNotifyMarkdown(p)
	if strings.Contains(out, "### Alerts") {
		t.Error("expected no alerts section when alerts are empty")
	}
}
