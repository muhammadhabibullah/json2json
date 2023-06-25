package json2json

import "testing"

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		jsonInput map[string]any
		want      any
		wantErr   bool
	}{
		{
			name:      "simple string",
			input:     "STRING('hello world')",
			jsonInput: map[string]any{},
			want:      "hello world",
			wantErr:   false,
		},
		{
			name:      "empty string",
			input:     "STRING()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple int",
			input:     "INT(123)",
			jsonInput: map[string]any{},
			want:      int64(123),
			wantErr:   false,
		},
		{
			name:      "empty int",
			input:     "INT()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple float",
			input:     "FLOAT(123.45)",
			jsonInput: map[string]any{},
			want:      123.45,
			wantErr:   false,
		},
		{
			name:      "simple float with precision",
			input:     "FLOAT(123.45, 2)",
			jsonInput: map[string]any{},
			want:      123.45,
			wantErr:   false,
		},
		{
			name:      "error float",
			input:     "FLOAT('two point five')",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty float",
			input:     "FLOAT()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "error float with precision",
			input:     "FLOAT(123.45, 'two')",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple bool",
			input:     "BOOL(TRUE)",
			jsonInput: map[string]any{},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "empty bool",
			input:     "BOOL()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple set",
			input:     "SET(STRING('hello world'))",
			jsonInput: map[string]any{},
			want:      "hello world",
			wantErr:   false,
		},
		{
			name:      "simple var",
			input:     `VAR('a', 'b')`,
			jsonInput: map[string]any{},
			want:      "a",
			wantErr:   false,
		},
		{
			name:      "default var",
			input:     `VAR([a], 'b')`,
			jsonInput: map[string]any{},
			want:      "b",
			wantErr:   false,
		},
		{
			name:      "error var",
			input:     `VAR([a])`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty var",
			input:     `VAR()`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple len from string",
			input:     `LEN('hello world')`,
			jsonInput: map[string]any{},
			want:      11,
			wantErr:   false,
		},
		{
			name:  "simple len from input array",
			input: `LEN([a])`,
			jsonInput: map[string]any{
				"a": []any{1, 2, 3},
			},
			want:    3,
			wantErr: false,
		},
		{
			name:  "simple len from input map",
			input: `LEN([a])`,
			jsonInput: map[string]any{
				"a": map[string]any{
					"b": 1,
					"c": 2,
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name:      "error len",
			input:     `LEN(123)`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty len",
			input:     `LEN()`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple slice_str",
			input:     `SLICE_STR('hello world', 0, 5)`,
			jsonInput: map[string]any{},
			want:      "hello",
			wantErr:   false,
		},
		{
			name:      "error slice_str invalid start",
			input:     `SLICE_STR('hello world', 12, 5)`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "error slice_str second arg",
			input:     `SLICE_STR('hello world', 'a', 5)`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "error slice_str third arg",
			input:     `SLICE_STR('hello world', 0, 'a')`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty slice_str",
			input:     `SLICE_STR()`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty set",
			input:     "SET()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple if",
			input:     "IF(TRUE, STRING('hello world'), STRING('goodbye world'))",
			jsonInput: map[string]any{},
			want:      "hello world",
			wantErr:   false,
		},
		{
			name:      "simple if else",
			input:     "IF(FALSE, STRING('hello world'), STRING('goodbye world'))",
			jsonInput: map[string]any{},
			want:      "goodbye world",
			wantErr:   false,
		},
		{
			name:      "error if expr",
			input:     "IF('a', 'b', 'c')",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty if",
			input:     "IF()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:  "simple switch",
			input: "SWITCH([a], INT(1), 'one', 'zero')",
			jsonInput: map[string]any{
				"a": int64(1),
			},
			want:    "one",
			wantErr: false,
		},
		{
			name:      "simple switch default",
			input:     "SWITCH([a], INT(1), 'one', 'zero')",
			jsonInput: map[string]any{},
			want:      "zero",
			wantErr:   false,
		},
		{
			name:      "simple switch default only",
			input:     "SWITCH([a], 'zero')",
			jsonInput: map[string]any{},
			want:      "zero",
			wantErr:   false,
		},
		{
			name:      "error odd switch",
			input:     "SWITCH([a], INT(1), 'zero')",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "error empty switch",
			input:     "SWITCH()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple and",
			input:     "AND(TRUE, TRUE)",
			jsonInput: map[string]any{},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "simple and false",
			input:     "AND(TRUE, FALSE)",
			jsonInput: map[string]any{},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "simple and error",
			input:     "AND(TRUE, 'a')",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty and",
			input:     `AND()`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple or",
			input:     "OR(TRUE, FALSE)",
			jsonInput: map[string]any{},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "simple or false",
			input:     "OR(FALSE, FALSE)",
			jsonInput: map[string]any{},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "simple or error",
			input:     "OR('a', FALSE)",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty or",
			input:     `OR()`,
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple gte",
			input:     "GTE(INT(1), INT(1))",
			jsonInput: map[string]any{},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "simple gte false",
			input:     "GTE(INT(1), INT(2))",
			jsonInput: map[string]any{},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "simple gte first error",
			input:     "GTE('a', INT(1))",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple gte second error",
			input:     "GTE(INT(1), 'a')",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty gte",
			input:     "GTE()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple gt",
			input:     "GT(INT(2), INT(1))",
			jsonInput: map[string]any{},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "simple gt false",
			input:     "GT(INT(1), INT(2))",
			jsonInput: map[string]any{},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "simple gt first error",
			input:     "GT('a', INT(1))",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple gt second error",
			input:     "GT(INT(1), 'a')",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty gt",
			input:     "GT()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple lte",
			input:     "LTE(INT(1), INT(1))",
			jsonInput: map[string]any{},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "simple lte false",
			input:     "LTE(INT(2), INT(1))",
			jsonInput: map[string]any{},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "simple lte first error",
			input:     "LTE('a', INT(1))",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple lte second error",
			input:     "LTE(INT(1), 'a')",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty lte",
			input:     "LTE()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple lt",
			input:     "LT(INT(1), INT(2))",
			jsonInput: map[string]any{},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "simple lt false",
			input:     "LT(INT(2), INT(1))",
			jsonInput: map[string]any{},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "simple lt first error",
			input:     "LT('a', INT(1))",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "simple lt second error",
			input:     "LT(INT(1), 'a')",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty lt",
			input:     "LT()",
			jsonInput: map[string]any{},
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := NewParser(tt.jsonInput)
			got, err := p.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveWhitespace(t *testing.T) {
	input := "STRING('a b c') "
	p := NewParser(map[string]any{})
	trimmed := p.removeWhitespace(input)
	if trimmed != "STRING('a b c')" {
		t.Errorf("error trimming whitespace: %s", trimmed)
	} else {
		t.Logf("success trimming whitespace: %s", trimmed)
	}
}

func TestSplitByComma(t *testing.T) {
	input := "VAR([abc],[def]),STRING('a,b')"
	p := NewParser(map[string]any{})
	split := p.splitArgs(input)
	if splitLen := len(split); splitLen != 2 {
		t.Errorf("error splitting by comma got %d results: %v", splitLen, split)
	} else {
		t.Logf("success splitting by comma got %d results: %v", splitLen, split)
	}
}
