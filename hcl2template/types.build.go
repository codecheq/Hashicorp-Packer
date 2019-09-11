package hcl2template

import (
	"github.com/hashicorp/hcl2/hcl"
)

const (
	builderLabel      = "output"
	provisionersLabel = "provisioners"
)

var buildSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: builderLabel, LabelNames: []string{"type", "name"}},
		{Type: provisionersLabel},
	},
}

type Build struct {
	// Ordered list of provisioners
	// Provisioners []Provisioner

	Outputs []Output

	HCL2Ref HCL2Ref
}

type Output struct {
	// Type of output; ex: amazon-ami
	Type string
	// Given name; if any
	Name    string
	HCL2Ref HCL2Ref
}

type Builds []Build

func (builds *Builds) decodeConfig(block *hcl.Block) hcl.Diagnostics {

	block.Body.Content(buildSchema)

	return nil
}
