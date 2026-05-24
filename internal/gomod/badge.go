package gomod

import (
	"fmt"
	"strings"
)

// BadgeStyle controls the output format of the badge.
type BadgeStyle string

const (
	BadgeStyleFlat        BadgeStyle = "flat"
	BadgeStyleFlatSquare  BadgeStyle = "flat-square"
	BadgeStylePlastic     BadgeStyle = "plastic"
)

// BadgeData holds the computed values for a dependency badge.
type BadgeData struct {
	Label   string
	Message string
	Color   string
	Style   BadgeStyle
}

// GenerateBadge produces a BadgeData summary from a scored diff.
func GenerateBadge(score BreakingScore, style BadgeStyle) BadgeData {
	if style == "" {
		style = BadgeStyleFlat
	}

	level := RiskLevel(score)
	color := badgeColor(level)

	return BadgeData{
		Label:   "dep-risk",
		Message: fmt.Sprintf("%s (%d)", strings.ToLower(string(level)), int(score)),
		Color:   color,
		Style:   style,
	}
}

// FormatBadgeMarkdown returns a shields.io markdown badge string.
func FormatBadgeMarkdown(b BadgeData) string {
	msg := strings.ReplaceAll(b.Message, " ", "%20")
	label := strings.ReplaceAll(b.Label, " ", "%20")
	url := fmt.Sprintf(
		"https://img.shields.io/badge/%s-%s-%s?style=%s",
		label, msg, b.Color, string(b.Style),
	)
	return fmt.Sprintf("![dep-risk](%s)", url)
}

// FormatBadgeSVG returns a minimal inline SVG badge.
func FormatBadgeSVG(b BadgeData) string {
	var sb strings.Builder
	sb.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="160" height="20">`)
	sb.WriteString(fmt.Sprintf(
		`<rect width="80" height="20" fill="#555"/>`+
			`<rect x="80" width="80" height="20" fill="#%s"/>`+
			`<text x="40" y="14" fill="#fff" font-size="11" text-anchor="middle">%s</text>`+
			`<text x="120" y="14" fill="#fff" font-size="11" text-anchor="middle">%s</text>`,
		b.Color, b.Label, b.Message,
	))
	sb.WriteString(`</svg>`)
	return sb.String()
}

func badgeColor(level RiskLabel) string {
	switch level {
	case RiskLow:
		return "4c1"
	case RiskMedium:
		return "fe7d37"
	case RiskHigh:
		return "e05d44"
	case RiskCritical:
		return "9b0000"
	default:
		return "9f9f9f"
	}
}
