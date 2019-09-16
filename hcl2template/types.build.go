package hcl2template

import (
	"github.com/hashicorp/hcl2/hcl"
)

const (
	buildFromLabel          = "from"
	buildProvisionnersLabel = "provision"
)

var buildSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: buildFromLabel, LabelNames: []string{"src"}},
		{Type: buildProvisionnersLabel},
	},
}

type Build struct {
	// Ordered list of provisioner groups
	ProvisionerGroups ProvisionerGroups

	// Ordered list of output stanzas
	Froms BuildFromList

	HCL2Ref HCL2Ref
}

type Builds []*Build

func (p *Parser) decodeBuildConfig(block *hcl.Block) (*Build, hcl.Diagnostics) {
	build := &Build{}

	content, diags := block.Body.Content(buildSchema)
	for _, block := range content.Blocks {
		switch block.Type {
		case buildFromLabel:
			bf := BuildFrom{}
			moreDiags := bf.decodeConfig(block)
			diags = append(diags, moreDiags...)
			build.Froms = append(build.Froms, bf)
		case buildProvisionnersLabel:
			pg, moreDiags := p.decodeProvisionerGroup(block)
			diags = append(diags, moreDiags...)
			build.ProvisionerGroups = append(build.ProvisionerGroups, pg)
		}
	}

	return build, diags
}
