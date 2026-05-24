package gomod

import (
	"fmt"
	"strings"
)

// NotifyChannel represents a supported notification output channel.
type NotifyChannel string

const (
	ChannelStdout  NotifyChannel = "stdout"
	ChannelSlack   NotifyChannel = "slack"
	ChannelWebhook NotifyChannel = "webhook"
)

// NotifyConfig holds configuration for sending diff notifications.
type NotifyConfig struct {
	Channel    NotifyChannel
	WebhookURL string
	MinRisk    string // "low", "medium", "high", "critical"
	Branch     string
}

// DefaultNotifyConfig returns a sensible default notification config.
func DefaultNotifyConfig() NotifyConfig {
	return NotifyConfig{
		Channel: ChannelStdout,
		MinRisk: "medium",
	}
}

// NotifyPayload is the structured message sent to a channel.
type NotifyPayload struct {
	Branch  string
	Summary Summary
	Score   BreakingScore
	Alerts  []Alert
}

// FormatNotifyText renders a plain-text notification message.
func FormatNotifyText(p NotifyPayload) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[diffdep] Dependency report for branch: %s\n", p.Branch))
	sb.WriteString(fmt.Sprintf("  Added: %d | Removed: %d | Changed: %d\n",
		p.Summary.Added, p.Summary.Removed, p.Summary.Changed))
	sb.WriteString(fmt.Sprintf("  Risk Score: %.1f (%s)\n", p.Score.Score, p.Score.Level))
	if len(p.Alerts) > 0 {
		sb.WriteString(fmt.Sprintf("  Alerts (%d):\n", len(p.Alerts)))
		for _, a := range p.Alerts {
			sb.WriteString(fmt.Sprintf("    [%s] %s\n", a.Level, a.Message))
		}
	}
	return sb.String()
}

// FormatNotifyMarkdown renders a Markdown-formatted notification message.
func FormatNotifyMarkdown(p NotifyPayload) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## diffdep report — `%s`\n", p.Branch))
	sb.WriteString(fmt.Sprintf("**Added:** %d | **Removed:** %d | **Changed:** %d\n\n",
		p.Summary.Added, p.Summary.Removed, p.Summary.Changed))
	sb.WriteString(fmt.Sprintf("**Risk Score:** `%.1f` — _%s_\n\n", p.Score.Score, p.Score.Level))
	if len(p.Alerts) > 0 {
		sb.WriteString(fmt.Sprintf("### Alerts (%d)\n", len(p.Alerts)))
		for _, a := range p.Alerts {
			sb.WriteString(fmt.Sprintf("- **[%s]** %s\n", a.Level, a.Message))
		}
	}
	return sb.String()
}
