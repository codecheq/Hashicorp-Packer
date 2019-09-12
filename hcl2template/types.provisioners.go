package hcl2template

import (
	"github.com/hashicorp/hcl2/hcl"
)

type ProvisionerGroups []ProvisionerGroup

type ProvisionerGroup []Provisioner

type Provisioner struct {
	Type string

	HCL2Ref HCL2Ref
}

func (provisionerGroup *ProvisionerGroup) decodeConfig(block *hcl.Block) hcl.Diagnostics {
	provisionerGroup := ProvisionerGroup{}

	content, diags := block.Body.Content(buildSchema)
	for _, block := range content.Blocks {
		switch block.Type {
		case outputLabel:
			output := Output{}
			moreDiags := output.decodeConfig(block)
			diags = append(diags, moreDiags...)
			build.Outputs = append(build.Outputs, output)
		}
	}

	*builds = append((*builds), build)
	return diags
}
