package hcl2template

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
)

type PackerV1Variables map[string]string

// decodeConfig decodes a "variables" section the way packer 1 used to
func (variables *PackerV1Variables) decodeConfig(block *hcl.Block) hcl.Diagnostics {
	return gohcl.DecodeBody(block.Body, nil, variables)
}
