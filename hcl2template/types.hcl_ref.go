package hcl2template

import "github.com/hashicorp/hcl2/hcl"

// reference to the source definition in configuration text file
type HCL2Ref struct {
	// reference to the source definition in configuration text file
	DeclRange hcl.Range

	// remainder of unparsed body
	Remain hcl.Body
}
