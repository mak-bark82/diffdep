package gomod

import (
	"strings"
	"testing"
)

func TestGenerateBadge_LowRisk(t *testing.T) {
	b := GenerateBadge(BreakingScore(5), BadgeStyleFlat)
	if b.Label != "dep-risk" {
		t.Errorf("expected label 'dep-risk', got %q", b.Label)
	}
	if b.Color != "4c1" {
		t.Errorf("expected green color for low risk, got %q", b.Color)
	}
	if b.Style != BadgeStyleFlat {
		t.Errorf("expected flat style, got %q", b.Style)
	}
}

func TestGenerateBadge_HighRisk(t *testing.T) {
	b := GenerateBadge(BreakingScore(80), BadgeStyleFlatSquare)
	if b.Color != "e05d44" {
		t.Errorf("expected red color for high risk, got %q", b.Color)
	}
}

func TestGenerateBadge_DefaultStyle(t *testing.T) {
	b := GenerateBadge(BreakingScore(0), "")
	if b.Style != BadgeStyleFlat {
		t.Errorf("expected default style to be flat, got %q", b.Style)
	}
}

func TestFormatBadgeMarkdown_ContainsURL(t *testing.T) {
	b := GenerateBadge(BreakingScore(10), BadgeStyleFlat)
	md := FormatBadgeMarkdown(b)
	if !strings.HasPrefix(md, "![dep-risk]") {
		t.Errorf("expected markdown badge prefix, got %q", md)
	}
	if !strings.Contains(md, "shields.io") {
		t.Errorf("expected shields.io URL in badge markdown")
	}
}

func TestFormatBadgeSVG_ContainsSVGTag(t *testing.T) {
	b := GenerateBadge(BreakingScore(50), BadgeStyleFlat)
	svg := FormatBadgeSVG(b)
	if !strings.Contains(svg, "<svg") {
		t.Errorf("expected SVG tag in output")
	}
	if !strings.Contains(svg, b.Label) {
		t.Errorf("expected label %q in SVG output", b.Label)
	}
	if !strings.Contains(svg, "</svg>") {
		t.Errorf("expected closing SVG tag")
	}
}

func TestFormatBadgeMarkdown_SpacesEncoded(t *testing.T) {
	b := GenerateBadge(BreakingScore(60), BadgeStyleFlat)
	md := FormatBadgeMarkdown(b)
	if strings.Contains(md, " ") {
		t.Errorf("badge markdown URL should not contain raw spaces")
	}
}
