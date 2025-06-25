package testutil

import (
	"math"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		x        interface{}
		y        interface{}
		opts     []cmp.Option
		expected bool
	}{
		{
			name:     "identical integers",
			x:        42,
			y:        42,
			expected: true,
		},
		{
			name:     "different integers",
			x:        42,
			y:        43,
			expected: false,
		},
		{
			name:     "identical strings",
			x:        "hello",
			y:        "hello",
			expected: true,
		},
		{
			name:     "different strings",
			x:        "hello",
			y:        "world",
			expected: false,
		},
		{
			name:     "identical slices",
			x:        []int{1, 2, 3},
			y:        []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "different slices",
			x:        []int{1, 2, 3},
			y:        []int{1, 2, 4},
			expected: false,
		},
		{
			name:     "identical maps",
			x:        map[string]int{"a": 1, "b": 2},
			y:        map[string]int{"a": 1, "b": 2},
			expected: true,
		},
		{
			name:     "different maps",
			x:        map[string]int{"a": 1, "b": 2},
			y:        map[string]int{"a": 1, "b": 3},
			expected: false,
		},
		{
			name:     "identical structs",
			x:        struct{ Name string }{"test"},
			y:        struct{ Name string }{"test"},
			expected: true,
		},
		{
			name:     "different structs",
			x:        struct{ Name string }{"test1"},
			y:        struct{ Name string }{"test2"},
			expected: false,
		},
		{
			name:     "nil values",
			x:        nil,
			y:        nil,
			expected: true,
		},
		{
			name:     "one nil value",
			x:        nil,
			y:        42,
			expected: false,
		},
		{
			name:     "NaN float64 values",
			x:        math.NaN(),
			y:        math.NaN(),
			expected: true, // NaNs should compare equal with our options
		},
		{
			name:     "NaN float32 values",
			x:        float32(math.NaN()),
			y:        float32(math.NaN()),
			expected: true, // NaNs should compare equal with our options
		},
		{
			name:     "normal float64 values",
			x:        3.14,
			y:        3.14,
			expected: true,
		},
		{
			name:     "different float64 values",
			x:        3.14,
			y:        2.71,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equal(tt.x, tt.y, tt.opts...)
			if result != tt.expected {
				t.Errorf("Equal(%v, %v) = %v, want %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestEqualWithBigRat(t *testing.T) {
	tests := []struct {
		name     string
		x        *big.Rat
		y        *big.Rat
		expected bool
	}{
		{
			name:     "identical big.Rat values",
			x:        big.NewRat(1, 2),
			y:        big.NewRat(1, 2),
			expected: true,
		},
		{
			name:     "equivalent big.Rat values",
			x:        big.NewRat(1, 2),
			y:        big.NewRat(2, 4),
			expected: true,
		},
		{
			name:     "different big.Rat values",
			x:        big.NewRat(1, 2),
			y:        big.NewRat(1, 3),
			expected: false,
		},
		{
			name:     "both nil big.Rat values",
			x:        nil,
			y:        nil,
			expected: true,
		},
		{
			name:     "one nil big.Rat value",
			x:        big.NewRat(1, 2),
			y:        nil,
			expected: false,
		},
		{
			name:     "other nil big.Rat value",
			x:        nil,
			y:        big.NewRat(1, 2),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equal(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("Equal(%v, %v) = %v, want %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name    string
		x       interface{}
		y       interface{}
		opts    []cmp.Option
		isEmpty bool // whether diff should be empty
	}{
		{
			name:    "identical values",
			x:       42,
			y:       42,
			isEmpty: true,
		},
		{
			name:    "different integers",
			x:       42,
			y:       43,
			isEmpty: false,
		},
		{
			name:    "different strings",
			x:       "hello",
			y:       "world",
			isEmpty: false,
		},
		{
			name:    "different slices",
			x:       []int{1, 2, 3},
			y:       []int{1, 2, 4},
			isEmpty: false,
		},
		{
			name:    "different maps",
			x:       map[string]int{"a": 1, "b": 2},
			y:       map[string]int{"a": 1, "b": 3},
			isEmpty: false,
		},
		{
			name:    "NaN values should be equal",
			x:       math.NaN(),
			y:       math.NaN(),
			isEmpty: true,
		},
		{
			name:    "big.Rat equivalent values",
			x:       big.NewRat(1, 2),
			y:       big.NewRat(2, 4),
			isEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := Diff(tt.x, tt.y, tt.opts...)

			if tt.isEmpty && diff != "" {
				t.Errorf("Diff(%v, %v) = %q, want empty string", tt.x, tt.y, diff)
			}
			if !tt.isEmpty && diff == "" {
				t.Errorf("Diff(%v, %v) = empty string, want non-empty", tt.x, tt.y)
			}
		})
	}
}

func TestIsDeepEqual(t *testing.T) {
	tests := []struct {
		name            string
		funcName        string
		testSectionName string
		result          interface{}
		expected        interface{}
		shouldFail      bool
	}{
		{
			name:            "identical structs",
			funcName:        "TestFunc",
			testSectionName: "struct test",
			result:          struct{ Name string }{"test"},
			expected:        struct{ Name string }{"test"},
			shouldFail:      false,
		},
		{
			name:            "different structs",
			funcName:        "TestFunc",
			testSectionName: "struct difference",
			result:          struct{ Name string }{"test1"},
			expected:        struct{ Name string }{"test2"},
			shouldFail:      true,
		},
		{
			name:            "identical maps",
			funcName:        "TestFunc",
			testSectionName: "map test",
			result:          map[string]int{"a": 1},
			expected:        map[string]int{"a": 1},
			shouldFail:      false,
		},
		{
			name:            "different maps",
			funcName:        "TestFunc",
			testSectionName: "map difference",
			result:          map[string]int{"a": 1},
			expected:        map[string]int{"a": 2},
			shouldFail:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			testIsDeepEqual(mockT, tt.funcName, tt.testSectionName, tt.result, tt.expected)

			if tt.shouldFail && !mockT.errorCalled {
				t.Errorf("Expected test to fail but it passed")
			}
			if !tt.shouldFail && mockT.errorCalled {
				t.Errorf("Expected test to pass but it failed with: %s", mockT.errorMessage)
			}
		})
	}
}

func TestIsMapVersionDeepEqual(t *testing.T) {
	type testStruct struct {
		Name  string
		Value int
	}

	tests := []struct {
		name            string
		funcName        string
		testSectionName string
		result          interface{}
		expected        interface{}
		shouldFail      bool
	}{
		{
			name:            "identical maps",
			funcName:        "TestFunc",
			testSectionName: "map test",
			result:          map[string]int{"a": 1, "b": 2},
			expected:        map[string]int{"a": 1, "b": 2},
			shouldFail:      false,
		},
		{
			name:            "different maps",
			funcName:        "TestFunc",
			testSectionName: "map difference",
			result:          map[string]int{"a": 1, "b": 2},
			expected:        map[string]int{"a": 1, "b": 3},
			shouldFail:      true,
		},
		{
			name:            "identical structs converted to maps",
			funcName:        "TestFunc",
			testSectionName: "struct to map test",
			result:          testStruct{"test", 42},
			expected:        testStruct{"test", 42},
			shouldFail:      false,
		},
		{
			name:            "different structs converted to maps",
			funcName:        "TestFunc",
			testSectionName: "struct to map difference",
			result:          testStruct{"test1", 42},
			expected:        testStruct{"test2", 42},
			shouldFail:      true,
		},
		{
			name:            "struct vs map comparison",
			funcName:        "TestFunc",
			testSectionName: "struct vs map",
			result:          testStruct{"test", 42},
			expected:        map[string]interface{}{"Name": "test", "Value": 42},
			shouldFail:      false,
		},
		{
			name:            "empty maps",
			funcName:        "TestFunc",
			testSectionName: "empty maps",
			result:          map[string]int{},
			expected:        map[string]int{},
			shouldFail:      false,
		},
		{
			name:            "nil comparison",
			funcName:        "TestFunc",
			testSectionName: "nil test",
			result:          nil,
			expected:        nil,
			shouldFail:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			testIsMapVersionDeepEqual(mockT, tt.funcName, tt.testSectionName, tt.result, tt.expected)

			if tt.shouldFail && !mockT.errorCalled {
				t.Errorf("Expected test to fail but it passed")
			}
			if !tt.shouldFail && mockT.errorCalled {
				t.Errorf("Expected test to pass but it failed with: %s", mockT.errorMessage)
			}
		})
	}
}

func TestAsMap(t *testing.T) {
	type testStruct struct {
		Name  string
		Value int
		Flag  bool
	}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "already a map",
			input:    map[string]int{"a": 1, "b": 2},
			expected: map[string]int{"a": 1, "b": 2},
		},
		{
			name:  "struct to map conversion",
			input: testStruct{"test", 42, true},
			expected: map[string]interface{}{
				"Name":  "test",
				"Value": 42,
				"Flag":  true,
			},
		},
		{
			name:     "empty struct",
			input:    struct{}{},
			expected: map[string]interface{}{},
		},
		{
			name:     "nil value",
			input:    nil,
			expected: nil,
		},
		{
			name:  "pointer to struct",
			input: &testStruct{"test", 42, false},
			expected: map[string]interface{}{
				"Name":  "test",
				"Value": 42,
				"Flag":  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := asMap(tt.input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("asMap(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test the default comparison options work correctly
func TestDefaultCmpOptions(t *testing.T) {
	t.Run("NaN float64 comparison", func(t *testing.T) {
		x := math.NaN()
		y := math.NaN()
		if !Equal(x, y) {
			t.Error("NaN float64 values should compare equal")
		}
	})

	t.Run("NaN float32 comparison", func(t *testing.T) {
		x := float32(math.NaN())
		y := float32(math.NaN())
		if !Equal(x, y) {
			t.Error("NaN float32 values should compare equal")
		}
	})

	t.Run("big.Rat comparison", func(t *testing.T) {
		x := big.NewRat(1, 2)
		y := big.NewRat(2, 4) // equivalent to 1/2
		if !Equal(x, y) {
			t.Error("Equivalent big.Rat values should compare equal")
		}
	})

	t.Run("big.Rat nil comparison", func(t *testing.T) {
		var x, y *big.Rat
		if !Equal(x, y) {
			t.Error("Nil big.Rat values should compare equal")
		}
	})
}

// Test with custom options
func TestEqualWithCustomOptions(t *testing.T) {
	// Custom option that ignores case for strings
	ignoreCase := cmp.Transformer("ignoreCase", func(s string) string {
		return strings.ToLower(s)
	})

	tests := []struct {
		name     string
		x        interface{}
		y        interface{}
		opts     []cmp.Option
		expected bool
	}{
		{
			name:     "case sensitive without option",
			x:        "Hello",
			y:        "hello",
			expected: false,
		},
		{
			name:     "case insensitive with option",
			x:        "Hello",
			y:        "hello",
			opts:     []cmp.Option{ignoreCase},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equal(tt.x, tt.y, tt.opts...)
			if result != tt.expected {
				t.Errorf("Equal(%v, %v, %v) = %v, want %v", tt.x, tt.y, tt.opts, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkEqual(b *testing.B) {
	x := map[string]interface{}{
		"name":    "test",
		"value":   42,
		"enabled": true,
		"items":   []int{1, 2, 3, 4, 5},
	}
	y := map[string]interface{}{
		"name":    "test",
		"value":   42,
		"enabled": true,
		"items":   []int{1, 2, 3, 4, 5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(x, y)
	}
}

func BenchmarkDiff(b *testing.B) {
	x := map[string]interface{}{
		"name":    "test1",
		"value":   42,
		"enabled": true,
		"items":   []int{1, 2, 3, 4, 5},
	}
	y := map[string]interface{}{
		"name":    "test2",
		"value":   43,
		"enabled": false,
		"items":   []int{1, 2, 3, 4, 6},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(x, y)
	}
}

func BenchmarkAsMap(b *testing.B) {
	type testStruct struct {
		Name    string
		Value   int
		Enabled bool
		Items   []int
	}

	s := testStruct{
		Name:    "test",
		Value:   42,
		Enabled: true,
		Items:   []int{1, 2, 3, 4, 5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		asMap(s)
	}
}

// testIsDeepEqual is a wrapper that accepts our interface
func testIsDeepEqual(t testingTInterface, funcName string, testSectionName string, result interface{}, expected interface{}) {
	testIsMapVersionDeepEqual(t, funcName, testSectionName, result, expected)
	testIsStringVersionDeepEqual(t, funcName, testSectionName, result, expected)
}

// testIsMapVersionDeepEqual is a wrapper that accepts our interface
func testIsMapVersionDeepEqual(t testingTInterface, funcName string, testSectionName string, result interface{}, expected interface{}) {
	// Convert to maps if not already maps
	expectedMap := asMap(expected)
	resultMap := asMap(result)

	if !reflect.DeepEqual(resultMap, expectedMap) {
		t.Errorf("%v() test \"%v\" Result: %v Want: %v", funcName, testSectionName, result, expected)
		return
	}
}
