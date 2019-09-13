package hcl2template

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
)

type ProvisionerGroup struct {
	Communicator hcl.Expression

	Provisioners []Provisioner
	HCL2Ref      HCL2Ref
}

type Provisioner struct {
	*hcl.Block
}

var provisionerGroupSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{"communicator", false},
	},
}

type ProvisionerGroups []*ProvisionerGroup

func (p *Parser) decodeProvisionerGroup(block *hcl.Block) (*ProvisionerGroup, hcl.Diagnostics) {

	var b struct {
		Communicator hcl.Expression `hcl:"communicator"`
		Remain       hcl.Body       `hcl:",remain"`
	}

	diags := gohcl.DecodeBody(block.Body, nil, &b)

	pg := &ProvisionerGroup{}
	pg.Communicator = b.Communicator
	pg.HCL2Ref.DeclRange = block.DefRange
	pg.HCL2Ref.Remain = b.Remain

	s := provisionerGroupSchema
	s.Attributes = append(s.Attributes, p.ProvisionersSchema.Attributes...)
	s.Blocks = append(s.Blocks, p.ProvisionersSchema.Blocks...)

	content, moreDiags := pg.HCL2Ref.Remain.Content(&s)
	diags = append(diags, moreDiags...)

	for _, block := range content.Blocks {
		p := Provisioner{block}
		pg.Provisioners = append(pg.Provisioners, p)
	}

	return pg, diags
}
