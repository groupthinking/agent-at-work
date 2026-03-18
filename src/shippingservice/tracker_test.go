// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"regexp"
	"testing"
)

func TestCreateTrackingId(t *testing.T) {
	tests := []struct {
		name string
		salt string
	}{
		{"empty salt", ""},
		{"short salt", "abc"},
		{"long salt", "this is a very long salt that should still work fine and produce a valid tracking id"},
		{"special characters", "!@#$%^&*()_+"},
		{"unicode characters", "こんにちは世界"},
	}

	re := regexp.MustCompile(`^[A-Z]{2}-\d+-\d+$`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateTrackingId(tt.salt)
			if !re.MatchString(got) {
				t.Errorf("CreateTrackingId(%q) = %q, want format like 'AA-123-456'", tt.salt, got)
			}
		})
	}
}

func TestGetRandomLetterCode(t *testing.T) {
	foundZ := false
	for i := 0; i < 1000; i++ {
		got := getRandomLetterCode()
		if got < 65 || got > 90 {
			t.Errorf("getRandomLetterCode() = %d, want between 65 (A) and 90 (Z)", got)
		}
		if got == 90 {
			foundZ = true
		}
	}
	if !foundZ {
		t.Error("getRandomLetterCode() did not generate 'Z' in 1000 iterations")
	}
}

func TestGetRandomNumber(t *testing.T) {
	tests := []struct {
		digits int
	}{
		{0},
		{1},
		{3},
		{7},
		{10},
	}

	for _, tt := range tests {
		got := getRandomNumber(tt.digits)
		if len(got) != tt.digits {
			t.Errorf("getRandomNumber(%d) = %q, length = %d, want %d", tt.digits, got, len(got), tt.digits)
		}
		for _, r := range got {
			if r < '0' || r > '9' {
				t.Errorf("getRandomNumber(%d) = %q, contains non-digit %q", tt.digits, got, r)
			}
		}
	}
}
