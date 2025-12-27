package badge

import (
	"fmt"
	"strings"
)

var DefaultBadgeStyle = StyleSummary

// BadgeStyle represents the type of badge to generate
type BadgeStyle string

const (
	StyleSummary  BadgeStyle = "summary"  // 400x200 - Key metrics
	StyleCompact  BadgeStyle = "compact"  // 280x28 - Shields.io style
	StyleDetailed BadgeStyle = "detailed" // 400x320 - Full stats
)

func BadgeStyleFromName(name string) (BadgeStyle, error) {
	switch strings.ToLower(name) {
	case "summary":
		return StyleSummary, nil
	case "compact":
		return StyleCompact, nil
	case "detailed":
		return StyleDetailed, nil
	}
	err := fmt.Errorf("invalid badge style: %s (must be: summary, compact, detailed)", name)
	return DefaultBadgeStyle, err
}
