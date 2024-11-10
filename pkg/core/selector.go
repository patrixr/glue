package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/patrixr/q"
)

const RootLevel = "root"
const Wildcard = "*"
const GroupSeparator = "."
const SelectorFilterSeparator = ","
const NegationRune = '~'

type Selector struct {
	filters []string
	prefix  []string
}

func NewSelector(selector string) Selector {
	return NewSelectorWithPrefix(selector, []string{})
}

func NewSelectorWithPrefix(selector string, prefix []string) Selector {
	filters := q.Filter(strings.Split(selector, SelectorFilterSeparator), func(str string) bool {
		return len(str) > 0
	})

	return Selector{
		filters: filters,
		prefix:  prefix,
	}
}

func (selector Selector) Test(levels []string) (bool, error) {
	if len(selector.filters) == 0 {
		return true, nil
	}

	var traverse func(level int, filter []string) (bool, error)

	// Recursive deep-dive
	traverse = func(level int, userFilter []string) (bool, error) {
		if level >= len(userFilter) {
			// User filter doesn't specify anything for lower levels, we run all the subgroups
			return true, nil
		}

		if level >= len(levels) {
			return false, nil
		}

		currentFilterLevel := userFilter[level]

		pass := nameMatch(currentFilterLevel, levels[level])

		last := level == len(userFilter)-1

		if !pass {
			return false, nil
		}

		if pass && last {
			return true, nil
		}

		return traverse(level+1, userFilter)
	}

	res := false

	hasPositiveFilters, _, _ := q.Find(selector.filters, func(filter string, _ int) bool {
		return filter[0] != NegationRune
	})

	// Go through each filter given to us by the user
	for _, filter := range selector.filters {
		if !ValidSelectorString(filter) {
			return false, errors.New(fmt.Sprintf("Invalid selector '%s'. Bad format or invalid characters", filter))
		}

		negation := len(filter) > 0 && filter[0] == NegationRune

		if negation {
			filter = filter[1:]
		}

		path := append([]string{}, selector.prefix...)
		path = append(path, strings.Split(filter, GroupSeparator)...)

		pass, error := traverse(0, path)

		if error != nil {
			return false, error
		}

		if pass {
			if negation {
				// explicit rejection, intercept and reject immediatly
				return false, nil
			}

			res = true
		} else if negation && !hasPositiveFilters {
			res = true
		}
	}

	return res, nil
}

func nameMatch(filter string, level string) bool {
	if strings.Contains(filter, Wildcard) {
		pattern := strings.ReplaceAll(filter, Wildcard, ".*")
		rexp := regexp.MustCompile(pattern)
		match := rexp.MatchString(level)
		return match
	}

	return strings.EqualFold(filter, level)
}

func ValidSelectorString(selector string) bool {
	c := "[a-zA-Z0-9 _\\-\\:\\$#\\*]+"
	filterPattern := fmt.Sprintf("%c?%s(\\.%s)*", NegationRune, c, c)
	pattern := fmt.Sprintf("^%[1]v(,%[1]v)*$", filterPattern)
	rexp := regexp.MustCompile(pattern)
	matched := rexp.MatchString(selector)
	return matched
}
