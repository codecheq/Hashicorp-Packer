package hcl2template

import (
	"github.com/hashicorp/packer/template"
)

// PackerConfig represents a loaded packer config
type PackerConfig struct {
	Sources map[SourceRef]*Source

	Variables PackerV1Variables

	Builds Builds

	Communicators []*Communicator
}

func (pkrCfg *PackerConfig) ToTemplate() *template.Template {
	return nil
}
