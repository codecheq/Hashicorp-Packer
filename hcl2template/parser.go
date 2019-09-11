package hcl2template

import (
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

var parser = &Parser{hclparse.NewParser()}

type Parser struct {
	*hclparse.Parser
}

func (p *Parser) ParseFile(filename string) (*hcl.File, hcl.Diagnostics) {
	if strings.HasSuffix(filename, ".json") {
		return p.ParseJSONFile(filename)
	}
	return p.ParseHCLFile(filename)
}
