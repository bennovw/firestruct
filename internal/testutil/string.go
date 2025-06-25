package testutil

import (
	"fmt"
	"testing"
)

// IsStringVersionDeepEqual checks if the string representation of the result is equal to the expected result
func IsStringVersionDeepEqual(t *testing.T, funcName string, testSectionName string, result interface{}, expected interface{}) {
	rStr := fmt.Sprint(result)
	eStr := fmt.Sprint(expected)

	if rStr != eStr {
		start, end, rDiff, eDiff := StringDiff(rStr, eStr)
		if start == -1 {
			// Strings are equal, no need to report
			t.Errorf("%v() test \"%v\""+
				"unexpected difference in string representation",
				funcName, testSectionName)
			return
		}
		t.Errorf("%v() test \"%v\" Result: %q Want: %q\n"+
			"Difference starts at index %d and ends at index %d\n"+
			"Result: %q\n"+
			"Expected: %q",
			funcName, testSectionName, rStr, eStr,
			start, end,
			rDiff, eDiff)
		return
	}
}

// StringDiff finds the difference between two strings and returns the start and end indices of the differing sections
func StringDiff(s1, s2 string) (int, int, string, string) {
	if s1 == s2 {
		return -1, -1, "", ""
	}

	start := stringDiffStart(s1, s2)
	if start == -1 {
		return -1, -1, "", ""
	}

	// If the strings are different, perform a backwards comparison to find the end of the differing section
	end := stringDiffEnd(s1, s2)
	end2 := stringDiffEnd(s2, s1)

	start2 := start
	if start > len(s2) {
		start2 = len(s2)
	}

	// Find the differing substrings including the start and end indices
	s1Diff := getDiffRange(s1, start, end)
	s2Diff := getDiffRange(s2, start2, end2)

	return start, end, s1Diff, s2Diff
}

// getDiffRange returns a substring that includes one character before start and one character after end (when possible)
func getDiffRange(s string, start, end int) string {
	if start < 0 || end < 0 || start > len(s) || end > len(s) {
		return ""
	}

	rangeStart := start
	if start > 0 {
		rangeStart = start - 1
	}

	rangeEnd := end
	if end < len(s)-1 {
		rangeEnd = end + 1
	}

	return s[rangeStart:rangeEnd]
}

// stringDiffStart finds string divergence points in two strings, returning the start indices of the differing sections
func stringDiffStart(s1, s2 string) int {
	if s1 == s2 {
		return -1
	}

	smallest := len(s1)
	if len(s2) < smallest {
		smallest = len(s2)
	}

	start := 0
	for start < smallest && s1[start] == s2[start] {
		start++
	}

	return start
}

// stringDiffEnd finds the difference between two strings starting in reverse order.
func stringDiffEnd(s1, s2 string) int {
	// reverse the strings
	reversedS1 := reverseString(s1)
	reversedS2 := reverseString(s2)

	// find the difference in the reversed strings
	start := stringDiffStart(reversedS1, reversedS2)
	if start == -1 {
		return -1 // strings are equal
	}

	// Calculate the end index in the original string
	end := len(s1) - start

	return end
}

// reverseString reverses a string and returns the reversed version.
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
