package badge

import (
	"fmt"
	"strings"
)

var DefaultSortBy = SortByPRs

// SortBy represents how contributions should be sorted in detailed view
type SortBy string

const (
	SortByPRs     SortBy = "prs"
	SortByStars   SortBy = "stars"
	SortByCommits SortBy = "commits"
)

func SortByFromName(name string) (SortBy, error) {
	switch strings.ToLower(name) {
	case "prs":
		return SortByPRs, nil
	case "stars":
		return SortByStars, nil
	case "commits":
		return SortByCommits, nil
	}
	err := fmt.Errorf("invalid badge sort: %s (must be: prs, stars, commits)", name)
	return DefaultSortBy, err
}
