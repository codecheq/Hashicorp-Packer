package hcl2template

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
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

type Parser struct {
	*hclparse.Parser

	// List of possible provisioners names.
	ProvisionersSchema *hcl.BodySchema

	// List of possible post-provisioners names.
	PostProvisionersSchema *hcl.BodySchema
}

func NewParser(provisioners, postProvisioners map[string]string) *Parser {
	p := &Parser{
		Parser: hclparse.NewParser(),
		ProvisionersSchema: &hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{},
		},
		PostProvisionersSchema: &hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{},
		},
	}
	for provisioner := range provisioners {
		p.ProvisionersSchema.Blocks = append(p.ProvisionersSchema.Blocks, hcl.BlockHeaderSchema{Type: provisioner})
	}

	for pp := range postProvisioners {
		p.PostProvisionersSchema.Blocks = append(p.PostProvisionersSchema.Blocks, hcl.BlockHeaderSchema{Type: pp})
	}

	return p
}

const hcl2FileExt = ".pkr.hcl"

func (p *Parser) Parse(filename string) (*PackerConfig, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	hclFiles := []string{}
	jsonFiles := []string{}
	if strings.HasSuffix(filename, hcl2FileExt) {
		hclFiles = append(hclFiles, hcl2FileExt)
	} else if strings.HasSuffix(filename, ".json") {
		jsonFiles = append(jsonFiles, hcl2FileExt)
	} else {
		fileInfos, err := ioutil.ReadDir(filename)
		if err != nil {
			diag := &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Cannot read hcl directory",
				Detail:   err.Error(),
			}
			diags = append(diags, diag)
		}
		for _, fileInfo := range fileInfos {
			if fileInfo.IsDir() {
				continue
			}
			filename := filepath.Join(filename, fileInfo.Name())
			if strings.HasSuffix(filename, hcl2FileExt) {
				hclFiles = append(hclFiles, filename)
			} else if strings.HasSuffix(filename, ".json") {
				jsonFiles = append(jsonFiles, filename)
			}
		}
	}

	var files []*hcl.File
	for _, filename := range hclFiles {
		f, moreDiags := p.ParseHCLFile(filename)
		diags = append(diags, moreDiags...)
		files = append(files, f)
	}
	for _, filename := range jsonFiles {
		f, moreDiags := p.ParseJSONFile(filename)
		diags = append(diags, moreDiags...)
		files = append(files, f)
	}
	if diags.HasErrors() {
		return nil, diags
	}

	cfg := &PackerConfig{}
	for _, file := range files {
		moreDiags := p.ParseFile(file, cfg)
		diags = append(diags, moreDiags...)
	}
	if diags.HasErrors() {
		return cfg, diags
	}

	return cfg, nil
}

// ParseFile filename content into cfg.
//
// ParseFile may be called multiple times with the same cfg on a different file.
//
// ParseFile returns as complete a config as we can manage, even if there are
// errors, since a partial result can be useful for careful analysis by
// development tools such as text editor extensions.
func (p *Parser) ParseFile(f *hcl.File, cfg *PackerConfig) hcl.Diagnostics {
	var diags hcl.Diagnostics

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
						"at %s. Each "+sourceLabel+" must have a unique name per builder type.",
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

		case communicatorLabel:
			communicator, moreDiags := p.decodeCommunicatorConfig(block)
			diags = append(diags, moreDiags...)
			cfg.Communicators = append(cfg.Communicators, communicator)

		default:
			panic(fmt.Sprintf("unexpected block type %q", block.Type)) // TODO(azr): err
		}
	}

	return diags
}
