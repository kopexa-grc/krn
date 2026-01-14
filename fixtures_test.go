// Copyright (c) Kopexa GRC
// SPDX-License-Identifier: Apache-2.0

package krn

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

// Fixtures represents the structure of testcases.json
type Fixtures struct {
	Version        string               `json:"version"`
	Parse          ParseFixtures        `json:"parse"`
	RoundTrip      []string             `json:"roundTrip"`
	Validation     ValidationFixtures   `json:"validation"`
	SafeResourceID []SafeResourceIDCase `json:"safeResourceId"`
	Operations     OperationsFixtures   `json:"operations"`
	ErrorCodes     []string             `json:"errorCodes"`
}

type ParseFixtures struct {
	Valid   []ValidParseCase   `json:"valid"`
	Invalid []InvalidParseCase `json:"invalid"`
}

type ValidParseCase struct {
	Name     string        `json:"name"`
	Input    string        `json:"input"`
	Expected ParseExpected `json:"expected"`
}

type ParseExpected struct {
	Service            *string           `json:"service"`
	Version            *string           `json:"version"`
	Segments           []SegmentExpected `json:"segments"`
	Depth              *int              `json:"depth"`
	Basename           *string           `json:"basename"`
	BasenameCollection *string           `json:"basenameCollection"`
	FullDomain         *string           `json:"fullDomain"`
	Path               *string           `json:"path"`
}

type SegmentExpected struct {
	Collection string `json:"collection"`
	ResourceID string `json:"resourceId"`
}

type InvalidParseCase struct {
	Name          string `json:"name"`
	Input         string `json:"input"`
	ExpectedError string `json:"expectedError"`
}

type ValidationFixtures struct {
	ResourceID ResourceIDValidation `json:"resourceId"`
	Version    VersionValidation    `json:"version"`
	Service    ServiceValidation    `json:"service"`
}

type ResourceIDValidation struct {
	Valid     []string `json:"valid"`
	Invalid   []string `json:"invalid"`
	MaxLength int      `json:"maxLength"`
}

type VersionValidation struct {
	Valid   []string `json:"valid"`
	Invalid []string `json:"invalid"`
}

type ServiceValidation struct {
	Valid   []string `json:"valid"`
	Invalid []string `json:"invalid"`
}

type SafeResourceIDCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type OperationsFixtures struct {
	Parent         []ParentCase         `json:"parent"`
	WithVersion    []WithVersionCase    `json:"withVersion"`
	WithoutVersion []WithoutVersionCase `json:"withoutVersion"`
	WithService    []WithServiceCase    `json:"withService"`
	WithoutService []WithoutServiceCase `json:"withoutService"`
	Child          []ChildCase          `json:"child"`
	ResourceID     []ResourceIDCase     `json:"resourceId"`
}

type ParentCase struct {
	Input    string  `json:"input"`
	Expected *string `json:"expected"`
}

type WithVersionCase struct {
	Input    string `json:"input"`
	Version  string `json:"version"`
	Expected string `json:"expected"`
}

type WithoutVersionCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type WithServiceCase struct {
	Input    string `json:"input"`
	Service  string `json:"service"`
	Expected string `json:"expected"`
}

type WithoutServiceCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type ChildCase struct {
	Input      string `json:"input"`
	Collection string `json:"collection"`
	ResourceID string `json:"resourceId"`
	Expected   string `json:"expected"`
}

type ResourceIDCase struct {
	Input         string  `json:"input"`
	Collection    string  `json:"collection"`
	Expected      *string `json:"expected"`
	ExpectedError *string `json:"expectedError"`
}

func loadFixtures(t *testing.T) *Fixtures {
	t.Helper()
	data, err := os.ReadFile("fixtures/testcases.json")
	if err != nil {
		t.Fatalf("failed to read fixtures: %v", err)
	}
	var fixtures Fixtures
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatalf("failed to parse fixtures: %v", err)
	}
	return &fixtures
}

func errorCodeToError(code string) error {
	switch code {
	case "EMPTY_KRN":
		return ErrEmptyKRN
	case "INVALID_KRN":
		return ErrInvalidKRN
	case "INVALID_DOMAIN":
		return ErrInvalidDomain
	case "INVALID_RESOURCE_ID":
		return ErrInvalidResourceID
	case "INVALID_VERSION":
		return ErrInvalidVersion
	case "RESOURCE_NOT_FOUND":
		return ErrResourceNotFound
	default:
		return nil
	}
}

func TestFixtures_Parse_Valid(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.Parse.Valid {
		t.Run(tc.Name, func(t *testing.T) {
			k, err := Parse(tc.Input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.Expected.Service != nil && k.Service() != *tc.Expected.Service {
				t.Errorf("service: got %q, want %q", k.Service(), *tc.Expected.Service)
			}
			if tc.Expected.Version != nil && k.Version() != *tc.Expected.Version {
				t.Errorf("version: got %q, want %q", k.Version(), *tc.Expected.Version)
			}
			if tc.Expected.Depth != nil && k.Depth() != *tc.Expected.Depth {
				t.Errorf("depth: got %d, want %d", k.Depth(), *tc.Expected.Depth)
			}
			if tc.Expected.Basename != nil && k.Basename() != *tc.Expected.Basename {
				t.Errorf("basename: got %q, want %q", k.Basename(), *tc.Expected.Basename)
			}
			if tc.Expected.BasenameCollection != nil && k.BasenameCollection() != *tc.Expected.BasenameCollection {
				t.Errorf("basenameCollection: got %q, want %q", k.BasenameCollection(), *tc.Expected.BasenameCollection)
			}
			if tc.Expected.FullDomain != nil && k.FullDomain() != *tc.Expected.FullDomain {
				t.Errorf("fullDomain: got %q, want %q", k.FullDomain(), *tc.Expected.FullDomain)
			}
			if tc.Expected.Path != nil && k.Path() != *tc.Expected.Path {
				t.Errorf("path: got %q, want %q", k.Path(), *tc.Expected.Path)
			}
			if tc.Expected.Segments != nil {
				segments := k.Segments()
				if len(segments) != len(tc.Expected.Segments) {
					t.Errorf("segments length: got %d, want %d", len(segments), len(tc.Expected.Segments))
				} else {
					for i, seg := range tc.Expected.Segments {
						if segments[i].Collection != seg.Collection {
							t.Errorf("segment[%d].collection: got %q, want %q", i, segments[i].Collection, seg.Collection)
						}
						if segments[i].ResourceID != seg.ResourceID {
							t.Errorf("segment[%d].resourceId: got %q, want %q", i, segments[i].ResourceID, seg.ResourceID)
						}
					}
				}
			}
		})
	}
}

func TestFixtures_Parse_Invalid(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.Parse.Invalid {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := Parse(tc.Input)
			if err == nil {
				t.Fatalf("expected error %s, got nil", tc.ExpectedError)
			}
			expectedErr := errorCodeToError(tc.ExpectedError)
			if !errors.Is(err, expectedErr) {
				t.Errorf("expected error %v, got %v", expectedErr, err)
			}
		})
	}
}

func TestFixtures_RoundTrip(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, input := range fixtures.RoundTrip {
		t.Run(input, func(t *testing.T) {
			k, err := Parse(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if k.String() != input {
				t.Errorf("round-trip failed: got %q, want %q", k.String(), input)
			}
		})
	}
}

func TestFixtures_Validation_ResourceID(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, id := range fixtures.Validation.ResourceID.Valid {
		t.Run("valid_"+id, func(t *testing.T) {
			if !IsValidResourceID(id) {
				t.Errorf("expected %q to be valid", id)
			}
		})
	}

	for _, id := range fixtures.Validation.ResourceID.Invalid {
		t.Run("invalid_"+id, func(t *testing.T) {
			if IsValidResourceID(id) {
				t.Errorf("expected %q to be invalid", id)
			}
		})
	}

	maxLen := fixtures.Validation.ResourceID.MaxLength
	t.Run("maxLength", func(t *testing.T) {
		if !IsValidResourceID(strings.Repeat("a", maxLen)) {
			t.Errorf("expected %d chars to be valid", maxLen)
		}
		if IsValidResourceID(strings.Repeat("a", maxLen+1)) {
			t.Errorf("expected %d chars to be invalid", maxLen+1)
		}
	})
}

func TestFixtures_Validation_Version(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, v := range fixtures.Validation.Version.Valid {
		t.Run("valid_"+v, func(t *testing.T) {
			if !IsValidVersion(v) {
				t.Errorf("expected %q to be valid", v)
			}
		})
	}

	for _, v := range fixtures.Validation.Version.Invalid {
		t.Run("invalid_"+v, func(t *testing.T) {
			if IsValidVersion(v) {
				t.Errorf("expected %q to be invalid", v)
			}
		})
	}
}

func TestFixtures_Validation_Service(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, s := range fixtures.Validation.Service.Valid {
		t.Run("valid_"+s, func(t *testing.T) {
			if !IsValidService(s) {
				t.Errorf("expected %q to be valid", s)
			}
		})
	}

	for _, s := range fixtures.Validation.Service.Invalid {
		t.Run("invalid_"+s, func(t *testing.T) {
			if IsValidService(s) {
				t.Errorf("expected %q to be invalid", s)
			}
		})
	}
}

func TestFixtures_SafeResourceID(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.SafeResourceID {
		t.Run(tc.Input, func(t *testing.T) {
			got := SafeResourceID(tc.Input)
			if got != tc.Expected {
				t.Errorf("SafeResourceID(%q) = %q, want %q", tc.Input, got, tc.Expected)
			}
		})
	}

	t.Run("truncates_to_200", func(t *testing.T) {
		input := strings.Repeat("a", 250)
		expected := strings.Repeat("a", 200)
		if got := SafeResourceID(input); got != expected {
			t.Errorf("expected truncation to 200, got %d chars", len(got))
		}
	})
}

func TestFixtures_Operations_Parent(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.Operations.Parent {
		t.Run(tc.Input, func(t *testing.T) {
			k := MustParse(tc.Input)
			parent := k.Parent()
			if tc.Expected == nil {
				if parent != nil {
					t.Errorf("expected nil parent, got %v", parent)
				}
			} else {
				if parent == nil {
					t.Fatalf("expected parent, got nil")
				}
				if parent.String() != *tc.Expected {
					t.Errorf("parent: got %q, want %q", parent.String(), *tc.Expected)
				}
			}
		})
	}
}

func TestFixtures_Operations_WithVersion(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.Operations.WithVersion {
		t.Run(tc.Input+"@"+tc.Version, func(t *testing.T) {
			k := MustParse(tc.Input)
			versioned, err := k.WithVersion(tc.Version)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if versioned.String() != tc.Expected {
				t.Errorf("got %q, want %q", versioned.String(), tc.Expected)
			}
		})
	}
}

func TestFixtures_Operations_WithoutVersion(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.Operations.WithoutVersion {
		t.Run(tc.Input, func(t *testing.T) {
			k := MustParse(tc.Input)
			unversioned := k.WithoutVersion()
			if unversioned.String() != tc.Expected {
				t.Errorf("got %q, want %q", unversioned.String(), tc.Expected)
			}
		})
	}
}

func TestFixtures_Operations_WithService(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.Operations.WithService {
		t.Run(tc.Input+"+"+tc.Service, func(t *testing.T) {
			k := MustParse(tc.Input)
			withService, err := k.WithService(tc.Service)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if withService.String() != tc.Expected {
				t.Errorf("got %q, want %q", withService.String(), tc.Expected)
			}
		})
	}
}

func TestFixtures_Operations_WithoutService(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.Operations.WithoutService {
		t.Run(tc.Input, func(t *testing.T) {
			k := MustParse(tc.Input)
			withoutService := k.WithoutService()
			if withoutService.String() != tc.Expected {
				t.Errorf("got %q, want %q", withoutService.String(), tc.Expected)
			}
		})
	}
}

func TestFixtures_Operations_Child(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.Operations.Child {
		t.Run(tc.Input+"/"+tc.Collection+"/"+tc.ResourceID, func(t *testing.T) {
			k := MustParse(tc.Input)
			child, err := NewChild(k, tc.Collection, tc.ResourceID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if child.String() != tc.Expected {
				t.Errorf("got %q, want %q", child.String(), tc.Expected)
			}
		})
	}
}

func TestFixtures_Operations_ResourceID(t *testing.T) {
	fixtures := loadFixtures(t)

	for _, tc := range fixtures.Operations.ResourceID {
		t.Run(tc.Input+"/"+tc.Collection, func(t *testing.T) {
			k := MustParse(tc.Input)
			got, err := k.ResourceID(tc.Collection)
			if tc.ExpectedError != nil {
				if err == nil {
					t.Fatalf("expected error %s, got nil", *tc.ExpectedError)
				}
				expectedErr := errorCodeToError(*tc.ExpectedError)
				if !errors.Is(err, expectedErr) {
					t.Errorf("expected error %v, got %v", expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tc.Expected != nil && got != *tc.Expected {
					t.Errorf("got %q, want %q", got, *tc.Expected)
				}
			}
		})
	}
}

func TestFixtures_ErrorCodes(t *testing.T) {
	fixtures := loadFixtures(t)

	expectedErrors := map[string]error{
		"EMPTY_KRN":           ErrEmptyKRN,
		"INVALID_KRN":         ErrInvalidKRN,
		"INVALID_DOMAIN":      ErrInvalidDomain,
		"INVALID_RESOURCE_ID": ErrInvalidResourceID,
		"INVALID_VERSION":     ErrInvalidVersion,
		"RESOURCE_NOT_FOUND":  ErrResourceNotFound,
	}

	if len(fixtures.ErrorCodes) != len(expectedErrors) {
		t.Errorf("expected %d error codes, got %d", len(expectedErrors), len(fixtures.ErrorCodes))
	}

	for _, code := range fixtures.ErrorCodes {
		if _, ok := expectedErrors[code]; !ok {
			t.Errorf("unexpected error code: %s", code)
		}
	}
}
