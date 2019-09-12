package hcl2template

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

type Outputs []Output

type Output struct {
	// Type of output; ex: amazon-ami
	Type string
	// Given name; if any
	Name string
	// Reference to a source; if any
	From string

	HCL2Ref HCL2Ref
}

func (output *Output) decodeConfig(block *hcl.Block) hcl.Diagnostics {

	output.Type = block.Labels[0]
	output.Name = block.Labels[1]
	output.HCL2Ref.DeclRange = block.DefRange

	var b struct {
		From   string
		Config hcl.Body `hcl:",remain"`
	}
	diags := gohcl.DecodeBody(block.Body, nil, &b)

	if !hclsyntax.ValidIdentifier(output.From) {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid builder type",
			Detail: "A " + sourceLabel + " type must start with a letter and " +
				"may contain only letters, digits, underscores, and dashes.",
			Subject: &block.DefRange,
		})
	}

	block.Body.Content(buildSchema)

	return nil
}
