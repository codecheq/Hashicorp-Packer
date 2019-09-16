package hcl2template

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
)

type BuildFromList []BuildFrom

type BuildFrom struct {
	// source to take config from
	Src string

	HCL2Ref HCL2Ref
}

func (bf *BuildFrom) decodeConfig(block *hcl.Block) hcl.Diagnostics {

	bf.Src = block.Labels[0]
	bf.HCL2Ref.DeclRange = block.DefRange

	var b struct {
		From   string
		Config hcl.Body `hcl:",remain"`
	}
	diags := gohcl.DecodeBody(block.Body, nil, &b)

	// if !hclsyntax.ValidIdentifier(bf.Src) {
	// 	diags = append(diags, &hcl.Diagnostic{
	// 		Severity: hcl.DiagError,
	// 		Summary:  "Invalid builder type",
	// 		Detail: "A " + sourceLabel + " type must start with a letter and " +
	// 			"may contain only letters, digits, underscores, and dashes.",
	// 		Subject: &block.DefRange,
	// 	})
	// }

	block.Body.Content(buildSchema)

	return diags
}
