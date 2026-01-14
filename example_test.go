// Copyright (c) Kopexa GRC
// SPDX-License-Identifier: Apache-2.0

package krn_test

import (
	"fmt"

	"github.com/kopexa-grc/krn"
)

func ExampleParse() {
	k, err := krn.Parse("//kopexa.com/catalog/frameworks/iso27001")
	if err != nil {
		panic(err)
	}
	fmt.Println(k.Service())
	fmt.Println(k.Basename())
	// Output:
	// catalog
	// iso27001
}

func ExampleParse_withVersion() {
	k, err := krn.Parse("//kopexa.com/catalog/frameworks/iso27001@v1.2.3")
	if err != nil {
		panic(err)
	}
	fmt.Println(k.Version())
	fmt.Println(k.HasVersion())
	// Output:
	// v1.2.3
	// true
}

func ExampleMustParse() {
	k := krn.MustParse("//kopexa.com/catalog/frameworks/iso27001/controls/a-5-1")
	fmt.Println(k.Depth())
	fmt.Println(k.Basename())
	// Output:
	// 2
	// a-5-1
}

func ExampleIsValid() {
	fmt.Println(krn.IsValid("//kopexa.com/catalog/frameworks/iso27001"))
	fmt.Println(krn.IsValid("invalid"))
	// Output:
	// true
	// false
}

func ExampleNew() {
	k, err := krn.New(krn.ServiceCatalog).
		Resource("frameworks", "iso27001").
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(k.String())
	// Output:
	// //kopexa.com/catalog/frameworks/iso27001
}

func ExampleNew_nested() {
	k, err := krn.New(krn.ServiceCatalog).
		Resource("frameworks", "iso27001").
		Resource("controls", "a-5-1").
		Version("v2").
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(k.String())
	// Output:
	// //kopexa.com/catalog/frameworks/iso27001/controls/a-5-1@v2
}

func ExampleNewChild() {
	parent := krn.MustParse("//kopexa.com/catalog/frameworks/iso27001")
	child, err := krn.NewChild(parent, "controls", "a-5-1")
	if err != nil {
		panic(err)
	}
	fmt.Println(child.String())
	// Output:
	// //kopexa.com/catalog/frameworks/iso27001/controls/a-5-1
}

func ExampleNewChildFromString() {
	child, err := krn.NewChildFromString(
		"//kopexa.com/isms/tenants/acme-corp",
		"workspaces",
		"main",
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(child.String())
	// Output:
	// //kopexa.com/isms/tenants/acme-corp/workspaces/main
}

func ExampleKRN_ResourceID() {
	k := krn.MustParse("//kopexa.com/catalog/frameworks/iso27001/controls/a-5-1")
	frameworkID, _ := k.ResourceID("frameworks")
	controlID, _ := k.ResourceID("controls")
	fmt.Println(frameworkID)
	fmt.Println(controlID)
	// Output:
	// iso27001
	// a-5-1
}

func ExampleKRN_Parent() {
	k := krn.MustParse("//kopexa.com/catalog/frameworks/iso27001/controls/a-5-1")
	parent := k.Parent()
	fmt.Println(parent.String())
	// Output:
	// //kopexa.com/catalog/frameworks/iso27001
}

func ExampleKRN_WithVersion() {
	k := krn.MustParse("//kopexa.com/catalog/frameworks/iso27001")
	versioned, err := k.WithVersion("v1")
	if err != nil {
		panic(err)
	}
	fmt.Println(versioned.String())
	// Output:
	// //kopexa.com/catalog/frameworks/iso27001@v1
}

func ExampleKRN_WithoutVersion() {
	k := krn.MustParse("//kopexa.com/catalog/frameworks/iso27001@v1")
	unversioned := k.WithoutVersion()
	fmt.Println(unversioned.String())
	// Output:
	// //kopexa.com/catalog/frameworks/iso27001
}

func ExampleKRN_Segments() {
	k := krn.MustParse("//kopexa.com/isms/tenants/acme/workspaces/main")
	for _, seg := range k.Segments() {
		fmt.Printf("%s: %s\n", seg.Collection, seg.ResourceID)
	}
	// Output:
	// tenants: acme
	// workspaces: main
}

func ExampleGetResource() {
	id, err := krn.GetResource("//kopexa.com/catalog/frameworks/iso27001", "frameworks")
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
	// Output:
	// iso27001
}

func ExampleIsValidResourceID() {
	fmt.Println(krn.IsValidResourceID("valid-id"))
	fmt.Println(krn.IsValidResourceID("-invalid"))
	// Output:
	// true
	// false
}

func ExampleIsValidVersion() {
	fmt.Println(krn.IsValidVersion("v1.2.3"))
	fmt.Println(krn.IsValidVersion("latest"))
	fmt.Println(krn.IsValidVersion("invalid"))
	// Output:
	// true
	// true
	// false
}

func ExampleSafeResourceID() {
	fmt.Println(krn.SafeResourceID("Hello World!"))
	fmt.Println(krn.SafeResourceID("-leading-dash"))
	// Output:
	// Hello-World
	// leading-dash
}
