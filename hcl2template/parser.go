package hcl2template

import (
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

type Parser struct {
	*hclparse.Parser

	// List of possible provisioners names.
	ProvisionersSchema *hcl.BodySchema

	// List of possible post-provisioners names.
	PostProvisionersSchema *hcl.BodySchema
}

func (p *Parser) ParseFile(filename string) (*hcl.File, hcl.Diagnostics) {
	if strings.HasSuffix(filename, ".json") {
		return p.ParseJSONFile(filename)
	}
	return p.ParseHCLFile(filename)
}
