package badge

import (
	"fmt"
	"strings"
)

var DefaultBadgeVariant = VariantDefault

type BadgeVariant string

const (
	VariantDefault   BadgeVariant = "default"
	VariantTextBased BadgeVariant = "text-based"
)

func BadgeVariantFromName(name string) (BadgeVariant, error) {
	switch strings.ToLower(name) {
	case "default":
		return VariantDefault, nil
	case "text-based":
		return VariantTextBased, nil
	}

	err := fmt.Errorf("invalid badge variant: %s (must be: default, text-based)", name)
	return DefaultBadgeVariant, err
}
