package stringutil

import (
	"sort"
	"strings"
)

// Compact will modify a and remove empty string items
func Compact(a []string) []string {
	b := a[:0]
	for _, x := range a {
		if x != "" {
			b = append(b, x)
		}
	}
	return b
}

// Contain s if include in list
func Contain(list []string, s string) bool {
	for _, l := range list {
		if l == s {
			return true
		}
	}

	return false
}

// Unique In-place deduplicate (comparable)
func Unique(list []string) []string {
	if len(list) == 0 {
		return list
	}

	sort.Strings(list)
	j := 0
	for i := 1; i < len(list); i++ {
		if list[j] == list[i] {
			continue
		}
		j++
		// preserve the original data
		list[i], list[j] = list[j], list[i]
		// only set what is required
		// list[j] = list[i]
	}

	return list[:j+1]
}

// EqualIgnoreCase 会对字符串进行比较, 忽略字符串的大小写
func EqualIgnoreCase(lhs, rhs string) bool {
	return strings.ToLower(lhs) == strings.ToLower(rhs)
}
