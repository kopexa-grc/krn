# KRN - Kopexa Resource Names

[![Go Reference](https://pkg.go.dev/badge/github.com/kopexa-grc/krn.svg)](https://pkg.go.dev/github.com/kopexa-grc/krn)
[![Go Report Card](https://goreportcard.com/badge/github.com/kopexa-grc/krn)](https://goreportcard.com/report/github.com/kopexa-grc/krn)
[![CI](https://github.com/kopexa-grc/krn/actions/workflows/ci.yml/badge.svg)](https://github.com/kopexa-grc/krn/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kopexa-grc/krn/branch/main/graph/badge.svg)](https://codecov.io/gh/kopexa-grc/krn)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A Go package for working with Kopexa Resource Names (KRN), following [Google's Resource Name Design](https://cloud.google.com/apis/design/resource_names).

## Installation

```bash
go get github.com/kopexa-grc/krn
```

## KRN Specification

### Format

```
//kopexa.com/{collection}/{resource-id}[/{collection}/{resource-id}][@{version}]
```

### Components

| Component | Description | Example |
|-----------|-------------|---------|
| Domain | Always `kopexa.com` | `kopexa.com` |
| Collection | Resource type (plural) | `frameworks`, `controls`, `tenants` |
| Resource ID | Unique identifier | `iso27001`, `a-5-1`, `acme-corp` |
| Version | Optional version tag | `@v1`, `@v1.2.3`, `@latest`, `@draft` |

### Examples

```
//kopexa.com/controls/ctrl-123
//kopexa.com/frameworks/iso27001
//kopexa.com/frameworks/iso27001/controls/a-5-1
//kopexa.com/frameworks/iso27001/controls/a-5-1@v2
//kopexa.com/tenants/acme-corp/workspaces/main
```

## Usage

### Parsing KRNs

```go
import "github.com/kopexa-grc/krn"

// Parse a KRN string
k, err := krn.Parse("//kopexa.com/frameworks/iso27001")
if err != nil {
    log.Fatal(err)
}

// Use MustParse when you're confident the string is valid
k := krn.MustParse("//kopexa.com/frameworks/iso27001")

// Check if a string is a valid KRN
if krn.IsValid(s) {
    // ...
}
```

### Building KRNs

```go
// Build a simple KRN
k, err := krn.New().
    Resource("frameworks", "iso27001").
    Build()
// Result: //kopexa.com/frameworks/iso27001

// Build a nested KRN with version
k, err := krn.New().
    Resource("frameworks", "iso27001").
    Resource("controls", "a-5-1").
    Version("v2").
    Build()
// Result: //kopexa.com/frameworks/iso27001/controls/a-5-1@v2
```

### Creating Child KRNs

```go
parent := krn.MustParse("//kopexa.com/frameworks/iso27001")

// Create a child KRN
child, err := krn.NewChild(parent, "controls", "a-5-1")
// Result: //kopexa.com/frameworks/iso27001/controls/a-5-1

// Or from a string
child, err := krn.NewChildFromString(
    "//kopexa.com/tenants/acme-corp",
    "workspaces",
    "main",
)
```

### Extracting Information

```go
k := krn.MustParse("//kopexa.com/frameworks/iso27001/controls/a-5-1@v1")

k.Path()              // "frameworks/iso27001/controls/a-5-1"
k.Version()           // "v1"
k.HasVersion()        // true
k.Basename()          // "a-5-1"
k.BasenameCollection() // "controls"
k.Depth()             // 2

// Get resource ID by collection
frameworkID, err := k.ResourceID("frameworks") // "iso27001"
controlID := k.MustResourceID("controls")      // "a-5-1"

// Check if a collection exists
k.HasResource("frameworks") // true
k.HasResource("policies")   // false

// Get all segments
for _, seg := range k.Segments() {
    fmt.Printf("%s: %s\n", seg.Collection, seg.ResourceID)
}
```

### Working with Parents

```go
k := krn.MustParse("//kopexa.com/frameworks/iso27001/controls/a-5-1")

parent := k.Parent()
// Result: //kopexa.com/frameworks/iso27001

// Root resources return nil
root := krn.MustParse("//kopexa.com/frameworks/iso27001")
root.Parent() // nil
```

### Version Manipulation

```go
k := krn.MustParse("//kopexa.com/frameworks/iso27001")

// Add version
versioned, err := k.WithVersion("v1.2.3")
// Result: //kopexa.com/frameworks/iso27001@v1.2.3

// Remove version
k2 := krn.MustParse("//kopexa.com/frameworks/iso27001@v1")
unversioned := k2.WithoutVersion()
// Result: //kopexa.com/frameworks/iso27001
```

### Comparison

```go
k1 := krn.MustParse("//kopexa.com/frameworks/iso27001")
k2 := krn.MustParse("//kopexa.com/frameworks/iso27001")

k1.Equals(k2)                                        // true
k1.EqualsString("//kopexa.com/frameworks/iso27001")  // true
```

### Control Mapping Example

A common use case is mapping controls between frameworks:

```go
// Framework A control
controlA := krn.MustParse("//kopexa.com/frameworks/iso27001/controls/a-5-1")

// Framework B control that maps to it
controlB := krn.MustParse("//kopexa.com/frameworks/nist-csf/controls/pr-ac-1")

// In a mapping file:
// - krn: //kopexa.com/frameworks/iso27001/controls/a-5-1
//   maps_to:
//     - krn: //kopexa.com/frameworks/nist-csf/controls/pr-ac-1
```

### Utility Functions

```go
// Quick resource extraction from string
id, err := krn.GetResource("//kopexa.com/frameworks/iso27001", "frameworks")
// Result: "iso27001"

// Validate resource IDs
krn.IsValidResourceID("valid-id")    // true
krn.IsValidResourceID("-invalid")    // false

// Validate versions
krn.IsValidVersion("v1.2.3")  // true
krn.IsValidVersion("latest")  // true
krn.IsValidVersion("invalid") // false

// Convert strings to valid resource IDs
krn.SafeResourceID("Hello World!") // "Hello-World"
```

## Resource ID Rules

Resource IDs must follow these rules:

- Length: 1-200 characters
- Allowed characters: `a-z`, `A-Z`, `0-9`, `-`, `_`, `.`
- Cannot start or end with `-` or `.`

## Version Formats

Supported version formats:

- Semantic: `v1`, `v1.2`, `v1.2.3`
- Keywords: `latest`, `draft`

## Error Handling

The package exports typed errors for precise error handling:

```go
import "errors"

k, err := krn.Parse(input)
if err != nil {
    switch {
    case errors.Is(err, krn.ErrEmptyKRN):
        // Handle empty input
    case errors.Is(err, krn.ErrInvalidKRN):
        // Handle invalid format
    case errors.Is(err, krn.ErrInvalidDomain):
        // Handle wrong domain
    case errors.Is(err, krn.ErrInvalidResourceID):
        // Handle invalid resource ID
    case errors.Is(err, krn.ErrInvalidVersion):
        // Handle invalid version format
    case errors.Is(err, krn.ErrResourceNotFound):
        // Handle missing resource
    }
}
```

## License

Copyright (c) Kopexa GRC

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
