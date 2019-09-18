package hcl2template

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

type Communicator struct {
	// Type of communicator; ex: ssh
	Type string
	// Given name
	Name string

	HCL2Ref HCL2Ref
}

func (p *Parser) decodeCommunicatorConfig(block *hcl.Block) (*Communicator, hcl.Diagnostics) {

	output := &Communicator{}
	output.Type = block.Labels[0]
	output.Name = block.Labels[1]
	output.HCL2Ref.DeclRange = block.DefRange

	diags := hcl.Diagnostics{}

	if !hclsyntax.ValidIdentifier(output.Name) {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + communicatorLabel + " type",
			Detail: "A " + communicatorLabel + " type must start with a letter and " +
				"may contain only letters, digits, underscores, and dashes.",
			Subject: &block.DefRange,
		})
	}

	return output, diags
}
