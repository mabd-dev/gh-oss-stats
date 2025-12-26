package badge

// BadgeOptions contains all configuration for badge generation
type BadgeOptions struct {
	Style   BadgeStyle
	Variant BadgeVariant
	Theme   BadgeTheme
	SortBy  SortBy // For detailed badge - how to sort contributions (default: prs)
	Limit   int    // For detailed badge - max contributions to show (default: 5)
}
