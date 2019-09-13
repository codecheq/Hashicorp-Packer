package hcl2template

import (
	"github.com/hashicorp/hcl2/hcl"
)

const (
	outputLabel        = "output"
	provisionnersLabel = "provision"
)

var buildSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: outputLabel, LabelNames: []string{"type", "name"}},
		{Type: provisionnersLabel},
	},
}

type Build struct {
	// Ordered list of provisioner groups
	ProvisionerGroups ProvisionerGroups

	// Ordered list of output stanzas
	Outputs Outputs

	HCL2Ref HCL2Ref
}

type Builds []Build

func (builds *Builds) decodeConfig(block *hcl.Block) hcl.Diagnostics {
	build := Build{}

	content, diags := block.Body.Content(buildSchema)
	for _, block := range content.Blocks {
		switch block.Type {
		case outputLabel:
			output := Output{}
			moreDiags := output.decodeConfig(block)
			diags = append(diags, moreDiags...)
			build.Outputs = append(build.Outputs, output)
		case provisionnersLabel:
			pg := ProvisionerGroup{}
			moreDiags := pg.decodeConfig(block)
			diags = append(diags, moreDiags...)
			build.ProvisionerGroups = append(build.ProvisionerGroups, pg)
		}
	}

	*builds = append((*builds), build)
	return diags
}
