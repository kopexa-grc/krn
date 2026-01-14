// Copyright (c) Kopexa GRC
// SPDX-License-Identifier: Apache-2.0

// Package krn implements Kopexa Resource Names (KRN) following Google's Resource Name Design.
//
// KRN Format:
//
//	//kopexa.com/{service}/{resource-type}/{resource-id}[/{sub-resource-type}/{sub-resource-id}][@{version}]
//
// Examples:
//
//	//kopexa.com/catalog/frameworks/iso27001
//	//kopexa.com/catalog/frameworks/iso27001/controls/a-5-1
//	//kopexa.com/catalog/frameworks/iso27001/controls/a-5-1@v2
//	//kopexa.com/isms/tenants/acme-corp/control-implementations/ci-123
package krn

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Domain is the base domain for all KRNs.
const Domain = "kopexa.com"

// Service constants define the available Kopexa services.
const (
	ServiceCatalog = "catalog"
	ServiceISMS    = "isms"
	ServiceOrg     = "org"
	ServiceAudit   = "audit"
	ServicePolicy  = "policy"
)

// validServices contains all valid service names.
var validServices = map[string]bool{
	ServiceCatalog: true,
	ServiceISMS:    true,
	ServiceOrg:     true,
	ServiceAudit:   true,
	ServicePolicy:  true,
}

// Error types for KRN parsing and validation.
var (
	ErrEmptyKRN          = errors.New("krn: empty KRN string")
	ErrInvalidKRN        = errors.New("krn: invalid KRN format")
	ErrInvalidDomain     = errors.New("krn: invalid domain")
	ErrInvalidService    = errors.New("krn: invalid service")
	ErrInvalidResourceID = errors.New("krn: invalid resource ID")
	ErrInvalidVersion    = errors.New("krn: invalid version format")
	ErrResourceNotFound  = errors.New("krn: resource not found")
)

// Validation patterns.
var (
	// resourceIDPattern validates resource IDs: 1-200 chars, alphanumeric plus - _ .
	// Cannot start or end with - or .
	resourceIDPattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9._-]{0,198}[a-zA-Z0-9])?$|^[a-zA-Z0-9]$`)

	// versionPattern validates version strings: v1, v1.2, v1.2.3, latest, draft
	versionPattern = regexp.MustCompile(`^(v\d+(\.\d+){0,2}|latest|draft)$`)
)

// Segment represents a collection/resource-id pair in a KRN path.
type Segment struct {
	Collection string
	ResourceID string
}

// KRN represents a Kopexa Resource Name.
type KRN struct {
	service  string
	segments []Segment
	version  string
}

// Parse parses a KRN string and returns a KRN struct.
func Parse(s string) (*KRN, error) {
	if s == "" {
		return nil, ErrEmptyKRN
	}

	// Must start with //
	if !strings.HasPrefix(s, "//") {
		return nil, fmt.Errorf("%w: must start with //", ErrInvalidKRN)
	}

	// Remove // prefix
	s = s[2:]

	// Extract version if present
	var version string
	if idx := strings.LastIndex(s, "@"); idx != -1 {
		version = s[idx+1:]
		s = s[:idx]
		if !IsValidVersion(version) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidVersion, version)
		}
	}

	// Split by /
	parts := strings.Split(s, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("%w: must have at least domain/service/collection/id", ErrInvalidKRN)
	}

	// Validate domain
	if parts[0] != Domain {
		return nil, fmt.Errorf("%w: expected %s, got %s", ErrInvalidDomain, Domain, parts[0])
	}

	// Validate service
	service := parts[1]
	if !validServices[service] {
		return nil, fmt.Errorf("%w: %s", ErrInvalidService, service)
	}

	// Parse resource path (must be pairs of collection/id)
	resourcePath := parts[2:]
	if len(resourcePath)%2 != 0 {
		return nil, fmt.Errorf("%w: resource path must be pairs of collection/id", ErrInvalidKRN)
	}

	segments := make([]Segment, 0, len(resourcePath)/2)
	for i := 0; i < len(resourcePath); i += 2 {
		collection := resourcePath[i]
		resourceID := resourcePath[i+1]

		if collection == "" {
			return nil, fmt.Errorf("%w: empty collection name", ErrInvalidKRN)
		}
		if !IsValidResourceID(resourceID) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidResourceID, resourceID)
		}

		segments = append(segments, Segment{
			Collection: collection,
			ResourceID: resourceID,
		})
	}

	return &KRN{
		service:  service,
		segments: segments,
		version:  version,
	}, nil
}

// MustParse parses a KRN string and panics on error.
func MustParse(s string) *KRN {
	krn, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return krn
}

// IsValid checks if a string is a valid KRN.
func IsValid(s string) bool {
	_, err := Parse(s)
	return err == nil
}

// IsValidResourceID checks if a string is a valid resource ID.
func IsValidResourceID(id string) bool {
	if id == "" || len(id) > 200 {
		return false
	}
	return resourceIDPattern.MatchString(id)
}

// IsValidVersion checks if a string is a valid version.
func IsValidVersion(v string) bool {
	return versionPattern.MatchString(v)
}

// SafeResourceID converts a string to a valid resource ID by replacing invalid characters.
func SafeResourceID(s string) string {
	if s == "" {
		return ""
	}

	// Replace invalid characters with -
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			result.WriteRune(r)
		} else {
			result.WriteRune('-')
		}
	}

	res := result.String()

	// Trim leading/trailing - and .
	res = strings.Trim(res, "-.")

	// Truncate to 200 characters
	if len(res) > 200 {
		res = res[:200]
		// Make sure we don't end with - or .
		res = strings.TrimRight(res, "-.")
	}

	return res
}

// GetResource extracts a resource ID from a KRN string by collection name.
func GetResource(krnString, collection string) (string, error) {
	k, err := Parse(krnString)
	if err != nil {
		return "", err
	}
	return k.ResourceID(collection)
}

// String returns the string representation of the KRN.
func (k *KRN) String() string {
	var sb strings.Builder
	sb.WriteString("//")
	sb.WriteString(Domain)
	sb.WriteString("/")
	sb.WriteString(k.service)

	for _, seg := range k.segments {
		sb.WriteString("/")
		sb.WriteString(seg.Collection)
		sb.WriteString("/")
		sb.WriteString(seg.ResourceID)
	}

	if k.version != "" {
		sb.WriteString("@")
		sb.WriteString(k.version)
	}

	return sb.String()
}

// Service returns the service name.
func (k *KRN) Service() string {
	return k.service
}

// Path returns the resource path without domain and service.
func (k *KRN) Path() string {
	var parts []string
	for _, seg := range k.segments {
		parts = append(parts, seg.Collection, seg.ResourceID)
	}
	return strings.Join(parts, "/")
}

// Version returns the version string, or empty if no version.
func (k *KRN) Version() string {
	return k.version
}

// HasVersion returns true if the KRN has a version.
func (k *KRN) HasVersion() bool {
	return k.version != ""
}

// ResourceID returns the resource ID for a given collection.
func (k *KRN) ResourceID(collection string) (string, error) {
	for _, seg := range k.segments {
		if seg.Collection == collection {
			return seg.ResourceID, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrResourceNotFound, collection)
}

// MustResourceID returns the resource ID for a given collection, panics if not found.
func (k *KRN) MustResourceID(collection string) string {
	id, err := k.ResourceID(collection)
	if err != nil {
		panic(err)
	}
	return id
}

// HasResource returns true if the KRN has a resource with the given collection.
func (k *KRN) HasResource(collection string) bool {
	for _, seg := range k.segments {
		if seg.Collection == collection {
			return true
		}
	}
	return false
}

// Basename returns the last resource ID in the path.
func (k *KRN) Basename() string {
	if len(k.segments) == 0 {
		return ""
	}
	return k.segments[len(k.segments)-1].ResourceID
}

// BasenameCollection returns the last collection name in the path.
func (k *KRN) BasenameCollection() string {
	if len(k.segments) == 0 {
		return ""
	}
	return k.segments[len(k.segments)-1].Collection
}

// Parent returns a new KRN without the last segment, or nil if this is a root resource.
func (k *KRN) Parent() *KRN {
	if len(k.segments) <= 1 {
		return nil
	}

	newSegments := make([]Segment, len(k.segments)-1)
	copy(newSegments, k.segments[:len(k.segments)-1])

	return &KRN{
		service:  k.service,
		segments: newSegments,
		version:  "", // Parent doesn't inherit version
	}
}

// WithVersion returns a new KRN with the specified version.
func (k *KRN) WithVersion(version string) (*KRN, error) {
	if !IsValidVersion(version) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidVersion, version)
	}

	newSegments := make([]Segment, len(k.segments))
	copy(newSegments, k.segments)

	return &KRN{
		service:  k.service,
		segments: newSegments,
		version:  version,
	}, nil
}

// WithoutVersion returns a new KRN without the version.
func (k *KRN) WithoutVersion() *KRN {
	newSegments := make([]Segment, len(k.segments))
	copy(newSegments, k.segments)

	return &KRN{
		service:  k.service,
		segments: newSegments,
		version:  "",
	}
}

// Equals checks if two KRNs are equal.
func (k *KRN) Equals(other *KRN) bool {
	if other == nil {
		return false
	}
	return k.String() == other.String()
}

// EqualsString checks if the KRN equals another KRN string.
func (k *KRN) EqualsString(other string) bool {
	otherKRN, err := Parse(other)
	if err != nil {
		return false
	}
	return k.Equals(otherKRN)
}

// Segments returns a copy of all segments in the KRN.
func (k *KRN) Segments() []Segment {
	result := make([]Segment, len(k.segments))
	copy(result, k.segments)
	return result
}

// Depth returns the number of resource levels in the KRN.
func (k *KRN) Depth() int {
	return len(k.segments)
}

// NewChild creates a new KRN as a child of the given parent.
func NewChild(parent *KRN, collection, resourceID string) (*KRN, error) {
	if parent == nil {
		return nil, fmt.Errorf("%w: parent cannot be nil", ErrInvalidKRN)
	}
	if collection == "" {
		return nil, fmt.Errorf("%w: collection cannot be empty", ErrInvalidKRN)
	}
	if !IsValidResourceID(resourceID) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidResourceID, resourceID)
	}

	newSegments := make([]Segment, len(parent.segments)+1)
	copy(newSegments, parent.segments)
	newSegments[len(parent.segments)] = Segment{
		Collection: collection,
		ResourceID: resourceID,
	}

	return &KRN{
		service:  parent.service,
		segments: newSegments,
		version:  "", // Child doesn't inherit version
	}, nil
}

// NewChildFromString creates a new KRN as a child of the given parent KRN string.
func NewChildFromString(parentKRN, collection, resourceID string) (*KRN, error) {
	parent, err := Parse(parentKRN)
	if err != nil {
		return nil, err
	}
	return NewChild(parent, collection, resourceID)
}

// Builder provides a fluent API for building KRNs.
type Builder struct {
	service  string
	segments []Segment
	version  string
	err      error
}

// New creates a new KRN builder with the specified service.
func New(service string) *Builder {
	b := &Builder{
		service:  service,
		segments: make([]Segment, 0),
	}

	if !validServices[service] {
		b.err = fmt.Errorf("%w: %s", ErrInvalidService, service)
	}

	return b
}

// Resource adds a resource segment to the builder.
func (b *Builder) Resource(collection, resourceID string) *Builder {
	if b.err != nil {
		return b
	}

	if collection == "" {
		b.err = fmt.Errorf("%w: collection cannot be empty", ErrInvalidKRN)
		return b
	}

	if !IsValidResourceID(resourceID) {
		b.err = fmt.Errorf("%w: %s", ErrInvalidResourceID, resourceID)
		return b
	}

	b.segments = append(b.segments, Segment{
		Collection: collection,
		ResourceID: resourceID,
	})
	return b
}

// Version sets the version for the KRN.
func (b *Builder) Version(version string) *Builder {
	if b.err != nil {
		return b
	}

	if !IsValidVersion(version) {
		b.err = fmt.Errorf("%w: %s", ErrInvalidVersion, version)
		return b
	}

	b.version = version
	return b
}

// Build creates the KRN. Returns nil and error if any error occurred during building.
func (b *Builder) Build() (*KRN, error) {
	if b.err != nil {
		return nil, b.err
	}

	if len(b.segments) == 0 {
		return nil, fmt.Errorf("%w: must have at least one resource", ErrInvalidKRN)
	}

	return &KRN{
		service:  b.service,
		segments: b.segments,
		version:  b.version,
	}, nil
}

// MustBuild creates the KRN and panics on error.
func (b *Builder) MustBuild() *KRN {
	k, err := b.Build()
	if err != nil {
		panic(err)
	}
	return k
}
