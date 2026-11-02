package testutil

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringDiff(t *testing.T) {
	tests := []struct {
		name       string
		s1         string
		s2         string
		wantStart  int
		wantEnd    int
		wantS1Diff string
		wantS2Diff string
	}{
		{
			name:       "identical strings",
			s1:         "hello world",
			s2:         "hello world",
			wantStart:  -1,
			wantEnd:    -1,
			wantS1Diff: "",
			wantS2Diff: "",
		},
		{
			name:       "empty strings",
			s1:         "",
			s2:         "",
			wantStart:  -1,
			wantEnd:    -1,
			wantS1Diff: "",
			wantS2Diff: "",
		},
		{
			name:       "one empty string",
			s1:         "hello",
			s2:         "",
			wantStart:  0,
			wantEnd:    5,
			wantS1Diff: "hello",
			wantS2Diff: "",
		},
		{
			name:       "different at start",
			s1:         "hello world",
			s2:         "goodbye world",
			wantStart:  0,
			wantEnd:    5,
			wantS1Diff: "hello ",
			wantS2Diff: "goodbye ",
		},
		{
			name:       "different at end",
			s1:         "hello world",
			s2:         "hello mars",
			wantStart:  6,
			wantEnd:    11,
			wantS1Diff: " world",
			wantS2Diff: " mars",
		},
		{
			name:       "different in middle",
			s1:         "hello beautiful world",
			s2:         "hello cruel world",
			wantStart:  6,
			wantEnd:    14,
			wantS1Diff: " beautiful",
			wantS2Diff: " cruel",
		},
		{
			name:       "s1 longer than s2",
			s1:         "hello world extra",
			s2:         "hello world",
			wantStart:  11,
			wantEnd:    17,
			wantS1Diff: "d extra",
			wantS2Diff: "d",
		},
		{
			name:       "s2 longer than s1",
			s1:         "hello world",
			s2:         "hello world extra",
			wantStart:  11,
			wantEnd:    11,
			wantS1Diff: "d",
			wantS2Diff: "d extra",
		},
		{
			name:       "completely different",
			s1:         "abc",
			s2:         "xyz",
			wantStart:  0,
			wantEnd:    3,
			wantS1Diff: "abc",
			wantS2Diff: "xyz",
		},
		{
			name:       "single character difference",
			s1:         "cat",
			s2:         "bat",
			wantStart:  0,
			wantEnd:    1,
			wantS1Diff: "ca",
			wantS2Diff: "ba",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, gotS1Diff, gotS2Diff := StringDiff(tt.s1, tt.s2)

			if gotStart != tt.wantStart {
				t.Errorf("StringDiff() test \"%v\" start = %v, want %v", tt.name, gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("StringDiff() test \"%v\" end = %v, want %v", tt.name, gotEnd, tt.wantEnd)
			}
			if gotS1Diff != tt.wantS1Diff {
				t.Errorf("StringDiff() test \"%v\" s1Diff = %q, want %q", tt.name, gotS1Diff, tt.wantS1Diff)
			}
			if gotS2Diff != tt.wantS2Diff {
				t.Errorf("StringDiff() test \"%v\" s2Diff = %q, want %q", tt.name, gotS2Diff, tt.wantS2Diff)
			}
		})
	}
}

func TestStringDiffStart(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{
			name: "identical strings",
			s1:   "hello",
			s2:   "hello",
			want: -1,
		},
		{
			name: "empty strings",
			s1:   "",
			s2:   "",
			want: -1,
		},
		{
			name: "different at start",
			s1:   "hello",
			s2:   "world",
			want: 0,
		},
		{
			name: "same prefix",
			s1:   "hello world",
			s2:   "hello mars",
			want: 6,
		},
		{
			name: "one string is prefix of another",
			s1:   "hello",
			s2:   "hello world",
			want: 5,
		},
		{
			name: "reverse prefix case",
			s1:   "hello world",
			s2:   "hello",
			want: 5,
		},
		{
			name: "one empty string",
			s1:   "hello",
			s2:   "",
			want: 0,
		},
		{
			name: "other empty string",
			s1:   "",
			s2:   "hello",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringDiffStart(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("stringDiffStart(%q, %q) = %v, want %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestStringDiffEnd(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{
			name: "identical strings",
			s1:   "hello",
			s2:   "hello",
			want: -1,
		},
		{
			name: "empty strings",
			s1:   "",
			s2:   "",
			want: -1,
		},
		{
			name: "different at end",
			s1:   "hello world",
			s2:   "hello mars",
			want: 11,
		},
		{
			name: "same suffix",
			s1:   "good morning",
			s2:   "bad morning",
			want: 3,
		},
		{
			name: "completely different",
			s1:   "abc",
			s2:   "xyz",
			want: 3,
		},
		{
			name: "one string longer",
			s1:   "hello world",
			s2:   "hello",
			want: 11,
		},
		{
			name: "other string longer",
			s1:   "hello",
			s2:   "hello world",
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringDiffEnd(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("stringDiffEnd(%q, %q) = %v, want %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "single character",
			input: "a",
			want:  "a",
		},
		{
			name:  "simple string",
			input: "hello",
			want:  "olleh",
		},
		{
			name:  "palindrome",
			input: "racecar",
			want:  "racecar",
		},
		{
			name:  "string with spaces",
			input: "hello world",
			want:  "dlrow olleh",
		},
		{
			name:  "unicode characters",
			input: "café",
			want:  "éfac",
		},
		{
			name:  "mixed case",
			input: "Hello World",
			want:  "dlroW olleH",
		},
		{
			name:  "numbers and symbols",
			input: "123!@#",
			want:  "#@!321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reverseString(tt.input)
			if got != tt.want {
				t.Errorf("reverseString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsStringVersionDeepEqual(t *testing.T) {
	tests := []struct {
		name       string
		result     interface{}
		expected   interface{}
		shouldFail bool
	}{
		{
			name:       "identical integers",
			result:     42,
			expected:   42,
			shouldFail: false,
		},
		{
			name:       "identical strings",
			result:     "hello",
			expected:   "hello",
			shouldFail: false,
		},
		{
			name:       "different integers",
			result:     42,
			expected:   43,
			shouldFail: true,
		},
		{
			name:       "different strings",
			result:     "hello",
			expected:   "world",
			shouldFail: true,
		},
		{
			name:       "different types same string representation",
			result:     42,
			expected:   "42",
			shouldFail: false,
		},
		{
			name:       "slice comparison",
			result:     []int{1, 2, 3},
			expected:   []int{1, 2, 3},
			shouldFail: false,
		},
		{
			name:       "different slice comparison",
			result:     []int{1, 2, 3},
			expected:   []int{1, 2, 4},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal, err := IsStringVersionDeepEqual(tt.result, tt.expected)
			if failed := err != nil || !equal; failed != tt.shouldFail {
				t.Errorf("IsStringVersionDeepEqual(%v, %v) = %v, %v; want failure: %v", tt.result, tt.expected, equal, err, tt.shouldFail)
			}
		})
	}
}

// mockTestingT is a mock implementation of testing.T for testing the test utility
type mockTestingT struct {
	errorCalled  bool
	errorMessage string
}

func (m *mockTestingT) Errorf(format string, args ...interface{}) {
	m.errorCalled = true
	m.errorMessage = strings.TrimSpace(fmt.Sprintf(format, args...))
}

func (m *mockTestingT) Helper() {}

// TestIsDeepEqualTest checks that IsDeepEqualTest reports failures with the category and description
func TestIsDeepEqualTest(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name       string
		result     interface{}
		expected   interface{}
		shouldPass bool
	}{
		{
			name:       "identical structs",
			result:     person{"John", 30},
			expected:   person{"John", 30},
			shouldPass: true,
		},
		{
			name:       "different structs",
			result:     person{"John", 30},
			expected:   person{"Jane", 25},
			shouldPass: false,
		},
		{
			name:       "different integers",
			result:     42,
			expected:   100,
			shouldPass: false,
		},
		{
			name:       "same string representation, different map values",
			result:     map[string]interface{}{"a": 1},
			expected:   map[string]interface{}{"a": "1"},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			got := IsDeepEqualTest(mockT, tt.result, tt.expected, "TestFunc", tt.name)

			if got != tt.shouldPass {
				t.Errorf("IsDeepEqualTest() = %v, want %v", got, tt.shouldPass)
			}
			if tt.shouldPass && mockT.errorCalled {
				t.Errorf("Expected to pass but Errorf was called with: %s", mockT.errorMessage)
			}
			if !tt.shouldPass && !mockT.errorCalled {
				t.Errorf("Expected to fail but Errorf was not called")
			}
			if mockT.errorCalled {
				if !strings.Contains(mockT.errorMessage, "TestFunc") || !strings.Contains(mockT.errorMessage, tt.name) {
					t.Errorf("Error message should contain category and description, got: %s", mockT.errorMessage)
				}
			}
		})
	}
}

// TestGetDiffRange tests the getDiffRange helper function
func TestGetDiffRange(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		start    int
		end      int
		expected string
	}{
		{
			name:     "normal range in middle",
			s:        "hello world",
			start:    2,
			end:      7,
			expected: "ello wo", // start-1=1 to end+1=8: "ello wo"
		},
		{
			name:     "range at start",
			s:        "hello world",
			start:    0,
			end:      3,
			expected: "hell", // start=0 to end+1=4: "hell"
		},
		{
			name:     "range at end",
			s:        "hello world",
			start:    8,
			end:      11,
			expected: "orld", // start-1=7 to end=11: "orld"
		},
		{
			name:     "single character range",
			s:        "hello",
			start:    2,
			end:      3,
			expected: "ell", // start-1=1 to end+1=4: "ell"
		},
		{
			name:     "empty string",
			s:        "",
			start:    0,
			end:      0,
			expected: "",
		},
		{
			name:     "invalid negative start",
			s:        "hello",
			start:    -1,
			end:      3,
			expected: "",
		},
		{
			name:     "invalid negative end",
			s:        "hello",
			start:    1,
			end:      -1,
			expected: "",
		},
		{
			name:     "start beyond string length",
			s:        "hello",
			start:    10,
			end:      15,
			expected: "",
		},
		{
			name:     "end beyond string length",
			s:        "hello",
			start:    2,
			end:      10,
			expected: "",
		},
		{
			name:     "start equals string length",
			s:        "hello",
			start:    5,
			end:      5,
			expected: "o", // start-1=4 to end=5: "o"
		},
		{
			name:     "full string with context",
			s:        "abc",
			start:    1,
			end:      2,
			expected: "ab", // start-1=0 to end+1=3, but end+1 > len-1, so end+1=2: "ab"
		},
		{
			name:     "range covers whole string",
			s:        "hello",
			start:    0,
			end:      5,
			expected: "hello", // start=0 to end=5: "hello"
		},
		{
			name:     "middle character with context",
			s:        "abcde",
			start:    2,
			end:      3,
			expected: "bcd", // start-1=1 to end+1=4: "bcd"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDiffRange(tt.s, tt.start, tt.end)
			if got != tt.expected {
				t.Errorf("getDiffRange(%q, %d, %d) = %q, want %q",
					tt.s, tt.start, tt.end, got, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkStringDiff(b *testing.B) {
	s1 := "This is a long string that will be used for benchmarking the StringDiff function"
	s2 := "This is a long string that will be used for benchmarking the StringDiff method"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringDiff(s1, s2)
	}
}

func BenchmarkReverseString(b *testing.B) {
	s := "This is a moderately long string for benchmarking the reverse function"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reverseString(s)
	}
}

func BenchmarkStringDiffStart(b *testing.B) {
	s1 := "This is a long string that will be used for benchmarking"
	s2 := "This is a long string that will be used for testing"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stringDiffStart(s1, s2)
	}
}

// Test edge cases and error conditions
func TestStringDiffEdgeCases(t *testing.T) {
	// Test with very long strings
	longStr1 := strings.Repeat("a", 10000) + "different" + strings.Repeat("b", 10000)
	longStr2 := strings.Repeat("a", 10000) + "changed" + strings.Repeat("b", 10000)

	start, end, diff1, diff2 := StringDiff(longStr1, longStr2)

	if start == -1 {
		t.Error("Expected difference to be found in long strings")
	}

	if !strings.Contains(diff1, "different") {
		t.Error("Expected diff1 to contain 'different'")
	}

	if !strings.Contains(diff2, "changed") {
		t.Error("Expected diff2 to contain 'changed'")
	}

	// Verify the indices are reasonable
	if start < 0 || start > len(longStr1) {
		t.Errorf("Start index %d is out of bounds for string length %d", start, len(longStr1))
	}

	if end < 0 || end > len(longStr1) {
		t.Errorf("End index %d is out of bounds for string length %d", end, len(longStr1))
	}
}

func TestReverseStringUnicodeEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "emoji",
			input: "🚀🌟💫",
			want:  "💫🌟🚀",
		},
		{
			name:  "mixed ascii and unicode",
			input: "hello🌍world",
			want:  "dlrow🌍olleh",
		},
		{
			name:  "accented characters",
			input: "naïve café",
			want:  "éfac evïan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reverseString(tt.input)
			if got != tt.want {
				t.Errorf("reverseString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
