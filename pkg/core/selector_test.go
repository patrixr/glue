package core_test

import (
	"fmt"
	"testing"

	. "github.com/patrixr/glue/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestSelector(t *testing.T) {
	testCases := []struct {
		selector   string
		levels     []string
		shouldPass bool
	}{
		// Basic group selection
		{"other", []string{"group1", "subgroup"}, false},
		{"group1", []string{"root"}, false},
		{"group1", []string{"root", "group1"}, false},
		{"root.group1", []string{"root", "group1"}, true},
		{"group1,group2", []string{"group1"}, true},
		{"group1,group2", []string{"group2"}, true},
		{"group1", []string{"group3"}, false},
		{"group1,group2", []string{"group3"}, false},

		// Wildcards
		{"*", []string{"group1"}, true},
		{"*", []string{"group1", "subgroup"}, true},
		{"**", []string{"group1", "subgroup"}, true},
		{"**", []string{"group1", "subgroup", "subsubgroup"}, true},
		{"*.*", []string{"group1", "subgroup", "subsubgroup"}, true},
		{"other,**", []string{"group1", "subgroup", "subsubgroup"}, true},
		{"group*", []string{"group1"}, true},
		{"group*", []string{"group2"}, true},
		{"group*", []string{"group2", "subgroup"}, true},
		{"*gr*", []string{"thegroup"}, true},
		{"group*", []string{"other"}, false},
		{"*.internal", []string{"group1", "internal"}, true},
		{"*.internal", []string{"group2", "other"}, false},
		{"*.internal", []string{"internal"}, false}, // a nested internal group
		{"*.internal", []string{"group1"}, false},
		{"group2.*", []string{"group2", "child"}, true},

		// Exclusion
		{"~group2", []string{"group1"}, true},
		{"~group2", []string{"group2"}, false},
		{"~group2", []string{"group2", "subgroup"}, false},
		{"~root.group2", []string{"root", "group2"}, false},
		{"~root.group2", []string{"root"}, true},
		{"root,~root.group2", []string{"root", "group2"}, false},

		// Nested group control
		{"group2.internal", []string{"group2"}, false},
		{"group2.internal", []string{"group2", "internal"}, true},
		{"group2.internal", []string{"group2", "other"}, false},
		{"~group2.internal", []string{"group2"}, true},
		{"group2,~group2.internal", []string{"group2"}, true},
		{"root.group2,~root.group2.internal", []string{"root"}, false},

		// Root level exclusion with wildcard
		{"root,~group2.*", []string{"root"}, true},
		{"root,~root.group2.*", []string{"root", "group2"}, true},
		{"~group2.*", []string{"group2", "child"}, false},
		{"something_else,~group2.*", []string{"group2", "child"}, false},
		{"root,~root.group2.*", []string{"root", "group2"}, true},

		{"~root.group2.internal", []string{"root", "group2", "internal"}, false},
	}

	for _, tc := range testCases {
		selector := NewSelector(tc.selector)
		result, err := selector.Test(tc.levels)

		assert.NoError(t, err)

		if tc.shouldPass {
			assert.True(t, result, "Selector(%q).test(%v) should be true", tc.selector, tc.levels)
		} else {
			assert.False(t, result, "Selector(%q).test(%v) should be false", tc.selector, tc.levels)
		}
	}
}

func TestValidateFilter(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		filter string
		valid  bool
	}{
		{"group1", true},
		{"group1.subgroup", true},
		{"group1.subgroup.item", true},
		{"group1_with_underscore", true},
		{"group1-with-hyphen", true},
		{"group1:with:colon", true},
		{"group1/with/slash", false},
		{"group1$with$dollar", true},
		{"group1#with#hash", true},
		{"group1*with*star", true},
		{"group1 with spaces", true},
		{"group1.*", true},
		{"group1.**", true},
		{"", false},
		{".", false},
		{"group1.", false},
		{".subgroup", false},
		{"group1..subgroup", false},
		{"group1/subgroup.item", false},
		{"group1:subgroup.item", true},
		{"group1#subgroup.item", true},
		{"group1$subgroup.item", true},
		{"group1*subgroup.item", true},
		{"group1_subgroup.item", true},
		{"group1-subgroup.item", true},
		{"invalid~", false},
		{"group1.invalid~", false},
		{"~group1", true},
	}

	for _, tc := range testCases {
		t.Run(tc.filter, func(t *testing.T) {
			isValid := ValidSelectorString(tc.filter)
			if tc.valid {
				assert.True(isValid, fmt.Sprintf("Filter '%s' should be valid", tc.filter))
			} else {
				assert.False(isValid, fmt.Sprintf("Filter '%s' should be invalid", tc.filter))
			}
		})
	}
}
