package hcl2template

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

const (
	sourceLabel = "source"
)

var configSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: sourceLabel, LabelNames: []string{"type", "name"}},
	},
}

// return as complete a config as we can manage, even if there are
// errors, since a partial result can be useful for careful analysis by
// development tools such as text editor extensions.
func (cfg *PackerConfig) Load(filename string) hcl.Diagnostics {
	if cfg.Sources == nil {
		cfg.Sources = map[SourceRef]*Source{}
	}

	f, diags := parser.ParseFile(filename)
	if diags.HasErrors() {
		return diags
	}

	content, moreDiags := f.Body.Content(configSchema)
	diags = append(diags, moreDiags...)

	for _, block := range content.Blocks {
		switch block.Type {
		case sourceLabel:
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

		default:
			// Only "source" is in our schema, so we can never get here
			panic(fmt.Sprintf("unexpected block type %q", block.Type))
		}
	}

	return diags
}
