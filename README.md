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
//{service}.kopexa.com/{collection}/{resource-id}[/{collection}/{resource-id}][@{version}]
```

### Components

| Component | Description | Example |
|-----------|-------------|---------|
| Service | Optional service subdomain | `catalog`, `isms`, `policy` |
| Domain | Always `kopexa.com` | `kopexa.com` |
| Collection | Resource type (plural) | `frameworks`, `controls`, `tenants` |
| Resource ID | Unique identifier | `iso27001`, `5.1.1`, `acme-corp` |
| Version | Optional version tag | `@v1`, `@v1.2.3`, `@latest`, `@draft` |

### Examples

```
//kopexa.com/frameworks/iso27001
//kopexa.com/frameworks/iso27001/controls/5.1.1
//catalog.kopexa.com/frameworks/iso27001
//catalog.kopexa.com/frameworks/iso27001/controls/5.1.1@v2
//isms.kopexa.com/tenants/acme-corp/workspaces/main
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

// Build a KRN with service
k, err := krn.New().
    Service("catalog").
    Resource("frameworks", "iso27001").
    Build()
// Result: //catalog.kopexa.com/frameworks/iso27001
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
k := krn.MustParse("//catalog.kopexa.com/frameworks/iso27001/controls/5.1.1@v1")

k.Service()           // "catalog"
k.HasService()        // true
k.FullDomain()        // "catalog.kopexa.com"
k.Path()              // "frameworks/iso27001/controls/5.1.1"
k.Version()           // "v1"
k.HasVersion()        // true
k.Basename()          // "5.1.1"
k.BasenameCollection() // "controls"
k.Depth()             // 2

// Get resource ID by collection
frameworkID, err := k.ResourceID("frameworks") // "iso27001"
controlID := k.MustResourceID("controls")      // "5.1.1"

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

### Service Manipulation

```go
k := krn.MustParse("//kopexa.com/frameworks/iso27001")

// Add service
withService, err := k.WithService("catalog")
// Result: //catalog.kopexa.com/frameworks/iso27001

// Remove service
k2 := krn.MustParse("//catalog.kopexa.com/frameworks/iso27001")
withoutService := k2.WithoutService()
// Result: //kopexa.com/frameworks/iso27001
```

### Comparison

```go
k1 := krn.MustParse("//kopexa.com/frameworks/iso27001")
k2 := krn.MustParse("//kopexa.com/frameworks/iso27001")

k1.Equals(k2)                                        // true
k1.EqualsString("//kopexa.com/frameworks/iso27001")  // true
```

### Framework Versioning

Compliance frameworks often have different editions (e.g., ISO 27001:2013 vs ISO 27001:2022).
Include the edition year or version number in the framework resource ID:

```go
// ISO 27001:2022
k := krn.MustParse("//catalog.kopexa.com/frameworks/iso27001-2022/controls/5.1.1")

// ISO 27001:2013 (different control numbering)
k := krn.MustParse("//catalog.kopexa.com/frameworks/iso27001-2013/controls/A.5.1.1")

// NIST CSF 2.0
k := krn.MustParse("//catalog.kopexa.com/frameworks/nist-csf-2.0/controls/GV.OC-01")

// CIS AWS Benchmark v1.4.0
k := krn.MustParse("//catalog.kopexa.com/frameworks/cis-aws-1.4.0/controls/1.1.1")

// SOC 2 Type II 2017
k := krn.MustParse("//catalog.kopexa.com/frameworks/soc2-2017/controls/CC1.1")
```

**Naming conventions:**
- Use lowercase with hyphens: `iso27001-2022`, `nist-csf-2.0`
- Include the year or version: `iso27001-2022`, `cis-aws-1.4.0`
- Keep it readable: `pci-dss-4.0` not `pcidssv4.0`

**Note:** The `@version` suffix (e.g., `@v1`, `@draft`) is for versioning the KRN content itself
(like draft vs published mappings), not for framework editions.

### Control Mapping Example

A common use case is mapping controls between frameworks:

```go
// ISO 27001:2022 control
controlA := krn.MustParse("//catalog.kopexa.com/frameworks/iso27001-2022/controls/5.1.1")

// CIS AWS v1.5.0 control that maps to it
controlB := krn.MustParse("//catalog.kopexa.com/frameworks/cis-aws-1.5.0/controls/1.1.1")

// NIST CSF 2.0 control
controlC := krn.MustParse("//catalog.kopexa.com/frameworks/nist-csf-2.0/controls/GV.OC-01")

// In a mapping file:
// - krn: //catalog.kopexa.com/frameworks/iso27001-2022/controls/5.1.1
//   maps_to:
//     - krn: //catalog.kopexa.com/frameworks/cis-aws-1.5.0/controls/1.1.1
//     - krn: //catalog.kopexa.com/frameworks/nist-csf-2.0/controls/GV.OC-01
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

// Validate service names
krn.IsValidService("catalog") // true
krn.IsValidService("Catalog") // false (must be lowercase)
krn.IsValidService("1svc")    // false (can't start with number)

// Convert strings to valid resource IDs
krn.SafeResourceID("Hello World!") // "Hello-World"
```

## Service Name Rules

Service names must follow DNS label rules:

- Length: 1-63 characters
- Allowed characters: `a-z`, `0-9`, `-`
- Must start with a letter
- Cannot end with `-`

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
