package badge

import (
	"fmt"
	"strings"
)

var DefualtBadgeStyle = StyleSummary

// BadgeStyle represents the type of badge to generate
type BadgeStyle string

const (
	StyleSummary  BadgeStyle = "summary"  // 400x200 - Key metrics
	StyleCompact  BadgeStyle = "compact"  // 280x28 - Shields.io style
	StyleDetailed BadgeStyle = "detailed" // 400x320 - Full stats
	StyleMinimal  BadgeStyle = "minimal"  // 120x28 - Project count only
)

func BadgeStyleFromName(name string) (BadgeStyle, error) {
	switch strings.ToLower(name) {
	case "summary":
		return StyleSummary, nil
	case "compact":
		return StyleCompact, nil
	case "detailed":
		return StyleDetailed, nil
	case "minimal":
		return StyleMinimal, nil
	}
	err := fmt.Errorf("invalid badge style: %s (must be: summary, compact, detailed, minimal)", name)
	return DefualtBadgeStyle, err
}

func (s BadgeStyle) TemplateStr() (string, error) {
	switch s {
	case StyleSummary:
		return summaryTemplate, nil
	case StyleCompact:
		return compactTemplate, nil
	case StyleDetailed:
		return detailedTemplate, nil
	case StyleMinimal:
		return minimalTemplate, nil
	}

	return "", fmt.Errorf("unsupported badge style: %s", s)
}
