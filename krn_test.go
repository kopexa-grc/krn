// Copyright (c) Kopexa GRC
// SPDX-License-Identifier: Apache-2.0

package krn

import (
	"errors"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   error
		checkFunc func(t *testing.T, k *KRN)
	}{
		{
			name:  "simple KRN",
			input: "//kopexa.com/frameworks/iso27001",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Depth() != 1 {
					t.Errorf("expected depth 1, got %d", k.Depth())
				}
				if k.Basename() != "iso27001" {
					t.Errorf("expected basename iso27001, got %s", k.Basename())
				}
			},
		},
		{
			name:  "nested KRN",
			input: "//kopexa.com/frameworks/iso27001/controls/a-5-1",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Depth() != 2 {
					t.Errorf("expected depth 2, got %d", k.Depth())
				}
				if k.Basename() != "a-5-1" {
					t.Errorf("expected basename a-5-1, got %s", k.Basename())
				}
				if k.BasenameCollection() != "controls" {
					t.Errorf("expected basename collection controls, got %s", k.BasenameCollection())
				}
			},
		},
		{
			name:  "KRN with version",
			input: "//kopexa.com/frameworks/iso27001@v2",
			checkFunc: func(t *testing.T, k *KRN) {
				if !k.HasVersion() {
					t.Error("expected HasVersion to be true")
				}
				if k.Version() != "v2" {
					t.Errorf("expected version v2, got %s", k.Version())
				}
			},
		},
		{
			name:  "KRN with semantic version",
			input: "//kopexa.com/frameworks/iso27001@v1.2.3",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Version() != "v1.2.3" {
					t.Errorf("expected version v1.2.3, got %s", k.Version())
				}
			},
		},
		{
			name:  "KRN with latest version",
			input: "//kopexa.com/frameworks/iso27001@latest",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Version() != "latest" {
					t.Errorf("expected version latest, got %s", k.Version())
				}
			},
		},
		{
			name:  "KRN with draft version",
			input: "//kopexa.com/frameworks/iso27001@draft",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Version() != "draft" {
					t.Errorf("expected version draft, got %s", k.Version())
				}
			},
		},
		{
			name:  "deep nested KRN",
			input: "//kopexa.com/tenants/acme-corp/control-implementations/ci-123/evidences/ev-456",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Depth() != 3 {
					t.Errorf("expected depth 3, got %d", k.Depth())
				}
			},
		},
		{
			name:  "controls KRN",
			input: "//kopexa.com/controls/ctrl-123",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Basename() != "ctrl-123" {
					t.Errorf("expected basename ctrl-123, got %s", k.Basename())
				}
			},
		},
		{
			name:  "resource ID with dots",
			input: "//kopexa.com/frameworks/iso.27001.2022",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Basename() != "iso.27001.2022" {
					t.Errorf("expected basename iso.27001.2022, got %s", k.Basename())
				}
			},
		},
		{
			name:  "resource ID with underscores",
			input: "//kopexa.com/frameworks/iso_27001_2022",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Basename() != "iso_27001_2022" {
					t.Errorf("expected basename iso_27001_2022, got %s", k.Basename())
				}
			},
		},
		{
			name:  "resource ID with mixed characters",
			input: "//kopexa.com/frameworks/ISO-27001_v2.0",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Basename() != "ISO-27001_v2.0" {
					t.Errorf("expected basename ISO-27001_v2.0, got %s", k.Basename())
				}
			},
		},
		{
			name:  "single character resource ID",
			input: "//kopexa.com/frameworks/x",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Basename() != "x" {
					t.Errorf("expected basename x, got %s", k.Basename())
				}
			},
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: ErrEmptyKRN,
		},
		{
			name:    "missing prefix",
			input:   "kopexa.com/frameworks/iso27001",
			wantErr: ErrInvalidKRN,
		},
		{
			name:    "wrong domain",
			input:   "//google.com/frameworks/iso27001",
			wantErr: ErrInvalidDomain,
		},
		{
			name:    "missing resource",
			input:   "//kopexa.com",
			wantErr: ErrInvalidKRN,
		},
		{
			name:    "odd number of path segments",
			input:   "//kopexa.com/frameworks",
			wantErr: ErrInvalidKRN,
		},
		{
			name:    "empty collection",
			input:   "//kopexa.com//iso27001",
			wantErr: ErrInvalidKRN,
		},
		{
			name:    "nested empty collection",
			input:   "//kopexa.com/frameworks/iso27001//a-5-1",
			wantErr: ErrInvalidKRN,
		},
		{
			name:    "invalid resource ID - starts with dash",
			input:   "//kopexa.com/frameworks/-iso27001",
			wantErr: ErrInvalidResourceID,
		},
		{
			name:    "invalid resource ID - ends with dash",
			input:   "//kopexa.com/frameworks/iso27001-",
			wantErr: ErrInvalidResourceID,
		},
		{
			name:    "invalid resource ID - starts with dot",
			input:   "//kopexa.com/frameworks/.iso27001",
			wantErr: ErrInvalidResourceID,
		},
		{
			name:    "invalid resource ID - ends with dot",
			input:   "//kopexa.com/frameworks/iso27001.",
			wantErr: ErrInvalidResourceID,
		},
		{
			name:    "invalid resource ID - contains space",
			input:   "//kopexa.com/frameworks/iso 27001",
			wantErr: ErrInvalidResourceID,
		},
		{
			name:    "invalid resource ID - contains special char",
			input:   "//kopexa.com/frameworks/iso!27001",
			wantErr: ErrInvalidResourceID,
		},
		{
			name:    "nested invalid resource ID",
			input:   "//kopexa.com/frameworks/iso27001/controls/-invalid",
			wantErr: ErrInvalidResourceID,
		},
		{
			name:    "invalid version format",
			input:   "//kopexa.com/frameworks/iso27001@invalid",
			wantErr: ErrInvalidVersion,
		},
		{
			name:    "invalid version - missing v prefix",
			input:   "//kopexa.com/frameworks/iso27001@1.0",
			wantErr: ErrInvalidVersion,
		},
		{
			name:    "invalid version - too many parts",
			input:   "//kopexa.com/frameworks/iso27001@v1.2.3.4",
			wantErr: ErrInvalidVersion,
		},
		// Service-based KRNs
		{
			name:  "KRN with service",
			input: "//catalog.kopexa.com/frameworks/iso27001",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Service() != "catalog" {
					t.Errorf("expected service catalog, got %s", k.Service())
				}
				if !k.HasService() {
					t.Error("expected HasService to be true")
				}
				if k.FullDomain() != "catalog.kopexa.com" {
					t.Errorf("expected full domain catalog.kopexa.com, got %s", k.FullDomain())
				}
				if k.String() != "//catalog.kopexa.com/frameworks/iso27001" {
					t.Errorf("expected string //catalog.kopexa.com/frameworks/iso27001, got %s", k.String())
				}
			},
		},
		{
			name:  "KRN with service and version",
			input: "//isms.kopexa.com/tenants/acme-corp@v1",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Service() != "isms" {
					t.Errorf("expected service isms, got %s", k.Service())
				}
				if k.Version() != "v1" {
					t.Errorf("expected version v1, got %s", k.Version())
				}
			},
		},
		{
			name:  "KRN with service - nested",
			input: "//policy.kopexa.com/frameworks/iso27001/controls/a-5-1",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Service() != "policy" {
					t.Errorf("expected service policy, got %s", k.Service())
				}
				if k.Depth() != 2 {
					t.Errorf("expected depth 2, got %d", k.Depth())
				}
			},
		},
		{
			name:  "KRN without service has empty service",
			input: "//kopexa.com/frameworks/iso27001",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Service() != "" {
					t.Errorf("expected empty service, got %s", k.Service())
				}
				if k.HasService() {
					t.Error("expected HasService to be false")
				}
				if k.FullDomain() != "kopexa.com" {
					t.Errorf("expected full domain kopexa.com, got %s", k.FullDomain())
				}
			},
		},
		{
			name:    "invalid service name - uppercase",
			input:   "//CATALOG.kopexa.com/frameworks/iso27001",
			wantErr: ErrInvalidDomain,
		},
		{
			name:    "invalid service name - starts with number",
			input:   "//1catalog.kopexa.com/frameworks/iso27001",
			wantErr: ErrInvalidDomain,
		},
		{
			name:    "invalid service name - starts with dash",
			input:   "//-catalog.kopexa.com/frameworks/iso27001",
			wantErr: ErrInvalidDomain,
		},
		// Control number-style resource IDs (like "5.1.1")
		{
			name:  "control number with dots",
			input: "//kopexa.com/frameworks/iso27001/controls/5.1.1",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Basename() != "5.1.1" {
					t.Errorf("expected basename 5.1.1, got %s", k.Basename())
				}
				id, err := k.ResourceID("controls")
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if id != "5.1.1" {
					t.Errorf("expected control ID 5.1.1, got %s", id)
				}
			},
		},
		{
			name:  "control number with dots and service",
			input: "//catalog.kopexa.com/frameworks/iso27001/controls/5.1.1",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Service() != "catalog" {
					t.Errorf("expected service catalog, got %s", k.Service())
				}
				if k.Basename() != "5.1.1" {
					t.Errorf("expected basename 5.1.1, got %s", k.Basename())
				}
			},
		},
		{
			name:  "control number dash style",
			input: "//kopexa.com/frameworks/iso27001/controls/5-1-1",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Basename() != "5-1-1" {
					t.Errorf("expected basename 5-1-1, got %s", k.Basename())
				}
			},
		},
		{
			name:  "CIS control number style",
			input: "//kopexa.com/frameworks/cis-aws/controls/1.1.1",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Basename() != "1.1.1" {
					t.Errorf("expected basename 1.1.1, got %s", k.Basename())
				}
			},
		},
		{
			name:  "NIST control number style",
			input: "//kopexa.com/frameworks/nist-csf/controls/PR.AC-1",
			checkFunc: func(t *testing.T, k *KRN) {
				if k.Basename() != "PR.AC-1" {
					t.Errorf("expected basename PR.AC-1, got %s", k.Basename())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := Parse(tt.input)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, k)
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	t.Run("valid KRN", func(t *testing.T) {
		k := MustParse("//kopexa.com/frameworks/iso27001")
		if k.Basename() != "iso27001" {
			t.Errorf("expected basename iso27001, got %s", k.Basename())
		}
	})

	t.Run("invalid KRN panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()
		MustParse("invalid")
	})
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"//kopexa.com/frameworks/iso27001", true},
		{"//kopexa.com/frameworks/iso27001@v1", true},
		{"", false},
		{"invalid", false},
		{"//google.com/frameworks/iso27001", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValid(tt.input); got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidResourceID(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"valid", true},
		{"Valid123", true},
		{"with-dash", true},
		{"with_underscore", true},
		{"with.dot", true},
		{"a", true},
		{"ab", true},
		{"a1", true},
		{"1a", true},
		{"ABC123", true},
		{"", false},
		{"-starts-with-dash", false},
		{"ends-with-dash-", false},
		{".starts-with-dot", false},
		{"ends-with-dot.", false},
		{"has space", false},
		{"has@symbol", false},
		{"has/slash", false},
		{strings.Repeat("a", 200), true},
		{strings.Repeat("a", 201), false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidResourceID(tt.input); got != tt.want {
				t.Errorf("IsValidResourceID(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidVersion(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"v1", true},
		{"v12", true},
		{"v123", true},
		{"v1.2", true},
		{"v1.2.3", true},
		{"v10.20.30", true},
		{"latest", true},
		{"draft", true},
		{"", false},
		{"1", false},
		{"1.0", false},
		{"v", false},
		{"v1.2.3.4", false},
		{"v1.", false},
		{"v.1", false},
		{"version1", false},
		{"release", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidVersion(tt.input); got != tt.want {
				t.Errorf("IsValidVersion(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidService(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"catalog", true},
		{"isms", true},
		{"policy", true},
		{"audit", true},
		{"org", true},
		{"a", true},
		{"ab", true},
		{"a1", true},
		{"service-name", true},
		{"my-service-123", true},
		{"", false},
		{"Catalog", false},      // uppercase not allowed
		{"ISMS", false},         // all uppercase not allowed
		{"1service", false},     // can't start with number
		{"-service", false},     // can't start with dash
		{"service-", false},     // can't end with dash (DNS label rules)
		{"Service", false},      // mixed case not allowed
		{"service_name", false}, // underscores not allowed
		{"service.name", false}, // dots not allowed
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidService(tt.input); got != tt.want {
				t.Errorf("IsValidService(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSafeResourceID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"valid", "valid"},
		{"with space", "with-space"},
		{"with@symbol", "with-symbol"},
		{"with/slash", "with-slash"},
		{"-leading-dash", "leading-dash"},
		{"trailing-dash-", "trailing-dash"},
		{".leading-dot", "leading-dot"},
		{"trailing-dot.", "trailing-dot"},
		{"multiple---dashes", "multiple---dashes"},
		{"", ""},
		{strings.Repeat("a", 250), strings.Repeat("a", 200)},
		{strings.Repeat("a", 199) + "-", strings.Repeat("a", 199)},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := SafeResourceID(tt.input); got != tt.want {
				t.Errorf("SafeResourceID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetResource(t *testing.T) {
	tests := []struct {
		name       string
		krnString  string
		collection string
		want       string
		wantErr    error
	}{
		{
			name:       "find framework",
			krnString:  "//kopexa.com/frameworks/iso27001",
			collection: "frameworks",
			want:       "iso27001",
		},
		{
			name:       "find nested resource",
			krnString:  "//kopexa.com/frameworks/iso27001/controls/a-5-1",
			collection: "controls",
			want:       "a-5-1",
		},
		{
			name:       "resource not found",
			krnString:  "//kopexa.com/frameworks/iso27001",
			collection: "controls",
			wantErr:    ErrResourceNotFound,
		},
		{
			name:       "invalid KRN",
			krnString:  "invalid",
			collection: "frameworks",
			wantErr:    ErrInvalidKRN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResource(tt.krnString, tt.collection)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestKRN_String(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "//kopexa.com/frameworks/iso27001",
			want:  "//kopexa.com/frameworks/iso27001",
		},
		{
			input: "//kopexa.com/frameworks/iso27001/controls/a-5-1",
			want:  "//kopexa.com/frameworks/iso27001/controls/a-5-1",
		},
		{
			input: "//kopexa.com/frameworks/iso27001@v1",
			want:  "//kopexa.com/frameworks/iso27001@v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			k := MustParse(tt.input)
			if got := k.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestKRN_Path(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "//kopexa.com/frameworks/iso27001",
			want:  "frameworks/iso27001",
		},
		{
			input: "//kopexa.com/frameworks/iso27001/controls/a-5-1",
			want:  "frameworks/iso27001/controls/a-5-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			k := MustParse(tt.input)
			if got := k.Path(); got != tt.want {
				t.Errorf("Path() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestKRN_RelativeResourceName(t *testing.T) {
	k := MustParse("//kopexa.com/frameworks/iso27001")
	if k.RelativeResourceName() != k.Path() {
		t.Errorf("RelativeResourceName() should equal Path()")
	}
}

func TestKRN_ResourceID(t *testing.T) {
	k := MustParse("//kopexa.com/frameworks/iso27001/controls/a-5-1")

	t.Run("find framework", func(t *testing.T) {
		got, err := k.ResourceID("frameworks")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got != "iso27001" {
			t.Errorf("got %q, want %q", got, "iso27001")
		}
	})

	t.Run("find control", func(t *testing.T) {
		got, err := k.ResourceID("controls")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got != "a-5-1" {
			t.Errorf("got %q, want %q", got, "a-5-1")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := k.ResourceID("nonexistent")
		if !errors.Is(err, ErrResourceNotFound) {
			t.Errorf("expected ErrResourceNotFound, got %v", err)
		}
	})
}

func TestKRN_MustResourceID(t *testing.T) {
	k := MustParse("//kopexa.com/frameworks/iso27001")

	t.Run("found", func(t *testing.T) {
		got := k.MustResourceID("frameworks")
		if got != "iso27001" {
			t.Errorf("got %q, want %q", got, "iso27001")
		}
	})

	t.Run("not found panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()
		k.MustResourceID("nonexistent")
	})
}

func TestKRN_HasResource(t *testing.T) {
	k := MustParse("//kopexa.com/frameworks/iso27001/controls/a-5-1")

	if !k.HasResource("frameworks") {
		t.Error("expected HasResource(frameworks) to be true")
	}
	if !k.HasResource("controls") {
		t.Error("expected HasResource(controls) to be true")
	}
	if k.HasResource("nonexistent") {
		t.Error("expected HasResource(nonexistent) to be false")
	}
}

func TestKRN_Parent(t *testing.T) {
	t.Run("has parent", func(t *testing.T) {
		k := MustParse("//kopexa.com/frameworks/iso27001/controls/a-5-1")
		parent := k.Parent()
		if parent == nil {
			t.Fatal("expected parent, got nil")
		}
		if parent.String() != "//kopexa.com/frameworks/iso27001" {
			t.Errorf("got %q, want %q", parent.String(), "//kopexa.com/frameworks/iso27001")
		}
	})

	t.Run("no parent for root", func(t *testing.T) {
		k := MustParse("//kopexa.com/frameworks/iso27001")
		parent := k.Parent()
		if parent != nil {
			t.Errorf("expected nil parent, got %v", parent)
		}
	})

	t.Run("parent does not inherit version", func(t *testing.T) {
		k := MustParse("//kopexa.com/frameworks/iso27001/controls/a-5-1@v1")
		parent := k.Parent()
		if parent.HasVersion() {
			t.Error("parent should not inherit version")
		}
	})
}

func TestKRN_WithVersion(t *testing.T) {
	k := MustParse("//kopexa.com/frameworks/iso27001")

	t.Run("add version", func(t *testing.T) {
		versioned, err := k.WithVersion("v1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if versioned.Version() != "v1" {
			t.Errorf("got %q, want %q", versioned.Version(), "v1")
		}
		// Original should be unchanged
		if k.HasVersion() {
			t.Error("original should not have version")
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		_, err := k.WithVersion("invalid")
		if !errors.Is(err, ErrInvalidVersion) {
			t.Errorf("expected ErrInvalidVersion, got %v", err)
		}
	})
}

func TestKRN_WithoutVersion(t *testing.T) {
	k := MustParse("//kopexa.com/frameworks/iso27001@v1")

	unversioned := k.WithoutVersion()
	if unversioned.HasVersion() {
		t.Error("expected no version")
	}
	if unversioned.String() != "//kopexa.com/frameworks/iso27001" {
		t.Errorf("got %q, want %q", unversioned.String(), "//kopexa.com/frameworks/iso27001")
	}
	// Original should be unchanged
	if !k.HasVersion() {
		t.Error("original should still have version")
	}
}

func TestKRN_WithService(t *testing.T) {
	k := MustParse("//kopexa.com/frameworks/iso27001")

	t.Run("add service", func(t *testing.T) {
		withService, err := k.WithService("catalog")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if withService.Service() != "catalog" {
			t.Errorf("got %q, want %q", withService.Service(), "catalog")
		}
		if withService.String() != "//catalog.kopexa.com/frameworks/iso27001" {
			t.Errorf("got %q, want %q", withService.String(), "//catalog.kopexa.com/frameworks/iso27001")
		}
		// Original should be unchanged
		if k.HasService() {
			t.Error("original should not have service")
		}
	})

	t.Run("invalid service", func(t *testing.T) {
		_, err := k.WithService("Invalid")
		if !errors.Is(err, ErrInvalidDomain) {
			t.Errorf("expected ErrInvalidDomain, got %v", err)
		}
	})

	t.Run("preserves version", func(t *testing.T) {
		versioned := MustParse("//kopexa.com/frameworks/iso27001@v1")
		withService, err := versioned.WithService("catalog")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if withService.Version() != "v1" {
			t.Errorf("expected version v1, got %s", withService.Version())
		}
		if withService.String() != "//catalog.kopexa.com/frameworks/iso27001@v1" {
			t.Errorf("got %q, want %q", withService.String(), "//catalog.kopexa.com/frameworks/iso27001@v1")
		}
	})
}

func TestKRN_WithoutService(t *testing.T) {
	k := MustParse("//catalog.kopexa.com/frameworks/iso27001")

	withoutService := k.WithoutService()
	if withoutService.HasService() {
		t.Error("expected no service")
	}
	if withoutService.String() != "//kopexa.com/frameworks/iso27001" {
		t.Errorf("got %q, want %q", withoutService.String(), "//kopexa.com/frameworks/iso27001")
	}
	// Original should be unchanged
	if !k.HasService() {
		t.Error("original should still have service")
	}
}

func TestKRN_ServicePreservation(t *testing.T) {
	t.Run("Parent preserves service", func(t *testing.T) {
		k := MustParse("//catalog.kopexa.com/frameworks/iso27001/controls/a-5-1")
		parent := k.Parent()
		if parent.Service() != "catalog" {
			t.Errorf("expected service catalog, got %s", parent.Service())
		}
		if parent.String() != "//catalog.kopexa.com/frameworks/iso27001" {
			t.Errorf("got %q, want %q", parent.String(), "//catalog.kopexa.com/frameworks/iso27001")
		}
	})

	t.Run("WithVersion preserves service", func(t *testing.T) {
		k := MustParse("//catalog.kopexa.com/frameworks/iso27001")
		versioned, err := k.WithVersion("v1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if versioned.Service() != "catalog" {
			t.Errorf("expected service catalog, got %s", versioned.Service())
		}
	})

	t.Run("WithoutVersion preserves service", func(t *testing.T) {
		k := MustParse("//catalog.kopexa.com/frameworks/iso27001@v1")
		unversioned := k.WithoutVersion()
		if unversioned.Service() != "catalog" {
			t.Errorf("expected service catalog, got %s", unversioned.Service())
		}
	})

	t.Run("NewChild preserves service", func(t *testing.T) {
		parent := MustParse("//catalog.kopexa.com/frameworks/iso27001")
		child, err := NewChild(parent, "controls", "a-5-1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if child.Service() != "catalog" {
			t.Errorf("expected service catalog, got %s", child.Service())
		}
		if child.String() != "//catalog.kopexa.com/frameworks/iso27001/controls/a-5-1" {
			t.Errorf("got %q, want %q", child.String(), "//catalog.kopexa.com/frameworks/iso27001/controls/a-5-1")
		}
	})
}

func TestKRN_Equals(t *testing.T) {
	k1 := MustParse("//kopexa.com/frameworks/iso27001")
	k2 := MustParse("//kopexa.com/frameworks/iso27001")
	k3 := MustParse("//kopexa.com/frameworks/iso27002")

	if !k1.Equals(k2) {
		t.Error("expected k1 to equal k2")
	}
	if k1.Equals(k3) {
		t.Error("expected k1 to not equal k3")
	}
	if k1.Equals(nil) {
		t.Error("expected k1 to not equal nil")
	}
}

func TestKRN_EqualsString(t *testing.T) {
	k := MustParse("//kopexa.com/frameworks/iso27001")

	if !k.EqualsString("//kopexa.com/frameworks/iso27001") {
		t.Error("expected equality")
	}
	if k.EqualsString("//kopexa.com/frameworks/iso27002") {
		t.Error("expected inequality")
	}
	if k.EqualsString("invalid") {
		t.Error("expected inequality for invalid string")
	}
}

func TestKRN_Segments(t *testing.T) {
	k := MustParse("//kopexa.com/frameworks/iso27001/controls/a-5-1")
	segments := k.Segments()

	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}

	if segments[0].Collection != "frameworks" || segments[0].ResourceID != "iso27001" {
		t.Errorf("unexpected first segment: %+v", segments[0])
	}
	if segments[1].Collection != "controls" || segments[1].ResourceID != "a-5-1" {
		t.Errorf("unexpected second segment: %+v", segments[1])
	}

	// Verify it's a copy
	segments[0].Collection = "modified"
	original := k.Segments()
	if original[0].Collection == "modified" {
		t.Error("Segments should return a copy")
	}
}

func TestNewChild(t *testing.T) {
	parent := MustParse("//kopexa.com/frameworks/iso27001")

	t.Run("valid child", func(t *testing.T) {
		child, err := NewChild(parent, "controls", "a-5-1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if child.String() != "//kopexa.com/frameworks/iso27001/controls/a-5-1" {
			t.Errorf("got %q, want %q", child.String(), "//kopexa.com/frameworks/iso27001/controls/a-5-1")
		}
	})

	t.Run("nil parent", func(t *testing.T) {
		_, err := NewChild(nil, "controls", "a-5-1")
		if !errors.Is(err, ErrInvalidKRN) {
			t.Errorf("expected ErrInvalidKRN, got %v", err)
		}
	})

	t.Run("empty collection", func(t *testing.T) {
		_, err := NewChild(parent, "", "a-5-1")
		if !errors.Is(err, ErrInvalidKRN) {
			t.Errorf("expected ErrInvalidKRN, got %v", err)
		}
	})

	t.Run("invalid resource ID", func(t *testing.T) {
		_, err := NewChild(parent, "controls", "-invalid")
		if !errors.Is(err, ErrInvalidResourceID) {
			t.Errorf("expected ErrInvalidResourceID, got %v", err)
		}
	})

	t.Run("child does not inherit version", func(t *testing.T) {
		versionedParent := MustParse("//kopexa.com/frameworks/iso27001@v1")
		child, err := NewChild(versionedParent, "controls", "a-5-1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if child.HasVersion() {
			t.Error("child should not inherit version")
		}
	})
}

func TestNewChildFromString(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		child, err := NewChildFromString("//kopexa.com/frameworks/iso27001", "controls", "a-5-1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if child.String() != "//kopexa.com/frameworks/iso27001/controls/a-5-1" {
			t.Errorf("got %q", child.String())
		}
	})

	t.Run("invalid parent", func(t *testing.T) {
		_, err := NewChildFromString("invalid", "controls", "a-5-1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestBuilder(t *testing.T) {
	t.Run("simple build", func(t *testing.T) {
		k, err := New().
			Resource("frameworks", "iso27001").
			Build()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if k.String() != "//kopexa.com/frameworks/iso27001" {
			t.Errorf("got %q", k.String())
		}
	})

	t.Run("nested build", func(t *testing.T) {
		k, err := New().
			Resource("frameworks", "iso27001").
			Resource("controls", "a-5-1").
			Build()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if k.String() != "//kopexa.com/frameworks/iso27001/controls/a-5-1" {
			t.Errorf("got %q", k.String())
		}
	})

	t.Run("with version", func(t *testing.T) {
		k, err := New().
			Resource("frameworks", "iso27001").
			Version("v1").
			Build()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if k.String() != "//kopexa.com/frameworks/iso27001@v1" {
			t.Errorf("got %q", k.String())
		}
	})

	t.Run("with service", func(t *testing.T) {
		k, err := New().
			Service("catalog").
			Resource("frameworks", "iso27001").
			Build()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if k.String() != "//catalog.kopexa.com/frameworks/iso27001" {
			t.Errorf("got %q", k.String())
		}
		if k.Service() != "catalog" {
			t.Errorf("expected service catalog, got %s", k.Service())
		}
	})

	t.Run("with service and version", func(t *testing.T) {
		k, err := New().
			Service("isms").
			Resource("tenants", "acme-corp").
			Version("v1").
			Build()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if k.String() != "//isms.kopexa.com/tenants/acme-corp@v1" {
			t.Errorf("got %q", k.String())
		}
	})

	t.Run("invalid service", func(t *testing.T) {
		_, err := New().
			Service("Invalid").
			Resource("frameworks", "iso27001").
			Build()
		if !errors.Is(err, ErrInvalidDomain) {
			t.Errorf("expected ErrInvalidDomain, got %v", err)
		}
	})

	t.Run("empty collection", func(t *testing.T) {
		_, err := New().
			Resource("", "iso27001").
			Build()
		if !errors.Is(err, ErrInvalidKRN) {
			t.Errorf("expected ErrInvalidKRN, got %v", err)
		}
	})

	t.Run("invalid resource ID", func(t *testing.T) {
		_, err := New().
			Resource("frameworks", "-invalid").
			Build()
		if !errors.Is(err, ErrInvalidResourceID) {
			t.Errorf("expected ErrInvalidResourceID, got %v", err)
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		_, err := New().
			Resource("frameworks", "iso27001").
			Version("invalid").
			Build()
		if !errors.Is(err, ErrInvalidVersion) {
			t.Errorf("expected ErrInvalidVersion, got %v", err)
		}
	})

	t.Run("no resources", func(t *testing.T) {
		_, err := New().Build()
		if !errors.Is(err, ErrInvalidKRN) {
			t.Errorf("expected ErrInvalidKRN, got %v", err)
		}
	})

	t.Run("error propagation stops further operations", func(t *testing.T) {
		_, err := New().
			Resource("", "iso27001").  // Error here
			Resource("controls", "a"). // Should not panic
			Version("v1").             // Should not panic
			Build()
		if !errors.Is(err, ErrInvalidKRN) {
			t.Errorf("expected ErrInvalidKRN, got %v", err)
		}
	})
}

func TestBuilder_MustBuild(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		k := New().
			Resource("frameworks", "iso27001").
			MustBuild()
		if k.String() != "//kopexa.com/frameworks/iso27001" {
			t.Errorf("got %q", k.String())
		}
	})

	t.Run("invalid panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()
		New().Resource("", "iso27001").MustBuild()
	})
}

func TestKRN_EmptySegments(t *testing.T) {
	// Test the edge case where a KRN is created with zero segments
	// This shouldn't happen via Parse, but test the methods handle it gracefully
	k := &KRN{
		segments: []Segment{},
	}

	if k.Basename() != "" {
		t.Errorf("expected empty basename, got %q", k.Basename())
	}
	if k.BasenameCollection() != "" {
		t.Errorf("expected empty basename collection, got %q", k.BasenameCollection())
	}
}

// Benchmarks

func BenchmarkParse(b *testing.B) {
	inputs := []string{
		"//kopexa.com/frameworks/iso27001",
		"//kopexa.com/frameworks/iso27001/controls/a-5-1",
		"//kopexa.com/tenants/acme-corp/control-implementations/ci-123/evidences/ev-456@v1.2.3",
	}

	for _, input := range inputs {
		b.Run(input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = Parse(input)
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	krns := []*KRN{
		MustParse("//kopexa.com/frameworks/iso27001"),
		MustParse("//kopexa.com/frameworks/iso27001/controls/a-5-1"),
		MustParse("//kopexa.com/tenants/acme-corp/control-implementations/ci-123/evidences/ev-456@v1.2.3"),
	}

	for _, k := range krns {
		b.Run(k.String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = k.String()
			}
		})
	}
}

func BenchmarkBuilder(b *testing.B) {
	b.Run("simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = New().
				Resource("frameworks", "iso27001").
				Build()
		}
	})

	b.Run("nested", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = New().
				Resource("frameworks", "iso27001").
				Resource("controls", "a-5-1").
				Build()
		}
	})

	b.Run("with-version", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = New().
				Resource("frameworks", "iso27001").
				Version("v1.2.3").
				Build()
		}
	})
}

func BenchmarkIsValidResourceID(b *testing.B) {
	ids := []string{
		"simple",
		"with-dashes-and_underscores",
		strings.Repeat("a", 200),
	}

	for _, id := range ids {
		b.Run(id[:min(len(id), 20)], func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = IsValidResourceID(id)
			}
		})
	}
}
