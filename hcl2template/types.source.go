package hcl2template

import (
	"fmt"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

// A source field in an HCL file will load into the Source type.
//
type Source struct {
	// Type of source; ex: virtualbox-iso
	Type string
	// Given name; if any
	Name string

	HCL2Ref HCL2Ref
}

func (source *Source) decodeConfig(block *hcl.Block) hcl.Diagnostics {

	source.Type = block.Labels[0]
	source.Name = block.Labels[1]
	source.HCL2Ref.DeclRange = block.DefRange

	var b struct {
		Remain hcl.Body `hcl:",remain"`
	}
	diags := gohcl.DecodeBody(block.Body, nil, &b)

	if !hclsyntax.ValidIdentifier(source.Type) {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + sourceLabel + " type",
			Detail: "A " + sourceLabel + " type must start with a letter and " +
				"may contain only letters, digits, underscores, and dashes.",
			Subject: &block.LabelRanges[0],
		})
	}

	//
	// TODO(azr): validate this later on; after interpolating pkr v1 variables
	// ? or simply use HCL type variables ?
	//
	// if !hclsyntax.ValidIdentifier(source.Name) {
	//  diags = append(diags, &hcl.Diagnostic{
	//      Severity: hcl.DiagError,
	//      Summary:  "Invalid source name",
	//      Detail:   "A " + sourceLabel + " name must start with a letter and may contain only letters, digits, underscores, and dashes.",
	//      Subject:  &block.LabelRanges[1],
	//  })
	// }

	source.HCL2Ref.Remain = b.Remain

	return diags
}

func (source *Source) Ref() SourceRef {
	return SourceRef{
		Type: source.Type,
		Name: source.Name,
	}
}

type SourceRef struct {
	Type string
	Name string
}

// NoSource is the zero value of sourceRef, representing the absense of an
// source.
var NoSource SourceRef

func sourceRefFromAbsTraversal(t hcl.Traversal) (SourceRef, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	if len(t) != 3 {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + sourceLabel + " reference",
			Detail:   "A " + sourceLabel + " reference must have three parts separated by periods: the keyword \"" + sourceLabel + "\", the builder type name, and the source name.",
			Subject:  t.SourceRange().Ptr(),
		})
		return NoSource, diags
	}

	if t.RootName() != sourceLabel {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + sourceLabel + " reference",
			Detail:   "The first part of an source reference must be the keyword \"" + sourceLabel + "\".",
			Subject:  t[0].SourceRange().Ptr(),
		})
		return NoSource, diags
	}
	btStep, ok := t[1].(hcl.TraverseAttr)
	if !ok {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + sourceLabel + " reference",
			Detail:   "The second part of an " + sourceLabel + " reference must be an identifier giving the builder type of the " + sourceLabel + ".",
			Subject:  t[1].SourceRange().Ptr(),
		})
		return NoSource, diags
	}
	nameStep, ok := t[2].(hcl.TraverseAttr)
	if !ok {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + sourceLabel + " reference",
			Detail:   "The third part of an " + sourceLabel + " reference must be an identifier giving the name of the " + sourceLabel + ".",
			Subject:  t[2].SourceRange().Ptr(),
		})
		return NoSource, diags
	}

	return SourceRef{
		Type: btStep.Name,
		Name: nameStep.Name,
	}, diags
}

func (r SourceRef) String() string {
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}
