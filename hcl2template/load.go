package hcl2template

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

const (
	sourceLabel       = "source"
	variablesLabel    = "variables"
	buildLabel        = "build"
	communicatorLabel = "communicator"
)

var configSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: sourceLabel, LabelNames: []string{"type", "name"}},
		{Type: variablesLabel},
		{Type: buildLabel},
		{Type: communicatorLabel, LabelNames: []string{"type", "name"}},
	},
}

// Parse filename content into cfg.
//
// Parse may be called multiple times with the same cfg on a different file.
//
// Parse returns as complete a config as we can manage, even if there are
// errors, since a partial result can be useful for careful analysis by
// development tools such as text editor extensions.
func (p *Parser) Parse(filename string, cfg *PackerConfig) hcl.Diagnostics {
	if cfg == nil {
		cfg = &PackerConfig{}
	}

	f, diags := p.ParseFile(filename)
	if diags.HasErrors() {
		return diags
	}

	content, moreDiags := f.Body.Content(configSchema)
	diags = append(diags, moreDiags...)

	for _, block := range content.Blocks {
		switch block.Type {
		case sourceLabel:
			if cfg.Sources == nil {
				cfg.Sources = map[SourceRef]*Source{}
			}
			source := &Source{}
			moreDiags := source.decodeConfig(block)
			diags = append(diags, moreDiags...)

			ref := source.Ref()
			if existing := cfg.Sources[ref]; existing != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate " + sourceLabel + " block",
					Detail: fmt.Sprintf("This "+sourceLabel+" block has the "+
						"same builder type and name as a previous block declared "+
						"at %s. Each source must have a unique name per builder type.",
						existing.HCL2Ref.DeclRange),
					Subject: &source.HCL2Ref.DeclRange,
				})
				continue
			}
			cfg.Sources[ref] = source

		case variablesLabel:
			if cfg.Variables == nil {
				cfg.Variables = PackerV1Variables{}
			}

			moreDiags := cfg.Variables.decodeConfig(block)
			diags = append(diags, moreDiags...)

		case buildLabel:
			build, moreDiags := p.decodeBuildConfig(block)
			diags = append(diags, moreDiags...)
			cfg.Builds = append(cfg.Builds, build)

		default:
			panic(fmt.Sprintf("unexpected block type %q", block.Type)) // TODO(azr): err
		}
	}

	return diags
}
