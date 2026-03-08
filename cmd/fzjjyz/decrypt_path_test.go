package main

import "testing"

func TestSafeDefaultOutputFromHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "normal", input: "report.txt", want: "report.txt"},
		{name: "traversal", input: "../secret.txt", want: "secret.txt"},
		{name: "absolute unix", input: "/tmp/data.bin", want: "data.bin"},
		{name: "absolute windows", input: `C:\Windows\temp\foo.log`, want: "foo.log"},
		{name: "empty", input: "", wantErr: true},
		{name: "dot", input: ".", wantErr: true},
		{name: "dotdot", input: "..", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := safeDefaultOutputFromHeader(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("input %q, got %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
