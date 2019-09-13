package hcl2template

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
)

type ProvisionerGroups []ProvisionerGroup

type ProvisionerGroup struct {
	Communicator hcl.Expression

	HCL2Ref HCL2Ref
}

var provisionerGroupSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{"communicator", false},
	},
}

func (provisionerGroup *ProvisionerGroup) decodeConfig(block *hcl.Block) hcl.Diagnostics {
	provisionerGroup.HCL2Ref.DeclRange = block.DefRange

	var b struct {
		Communicator hcl.Expression `hcl:"communicator"`
		Remain       hcl.Body       `hcl:",remain"`
	}

	diags := gohcl.DecodeBody(block.Body, nil, &b)

	provisionerGroup.Communicator = b.Communicator
	provisionerGroup.HCL2Ref.Remain = b.Remain
	return diags
}
