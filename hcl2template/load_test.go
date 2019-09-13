package hcl2template

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

func TestParser_Parse(t *testing.T) {
	defaultParser := &Parser{hclparse.NewParser(), &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "shell"},
			{Type: "upload", LabelNames: []string{"source", "destination"}},
		},
	}}

	type fields struct {
		Parser *hclparse.Parser
	}
	type args struct {
		filename string
		cfg      *PackerConfig
	}
	tests := []struct {
		name             string
		parser           *Parser
		args             args
		wantPackerConfig *PackerConfig
		wantDiags        bool
	}{
		{
			"valid " + sourceLabel + " load",
			defaultParser,
			args{"testdata/sources/vb-iso.tf", new(PackerConfig)},
			&PackerConfig{
				Sources: map[SourceRef]*Source{
					SourceRef{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1204",
					}: &Source{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1204",
					},
					SourceRef{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1604",
					}: &Source{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1604",
					},
				},
			},
			false,
		},

		{
			"valid " + communicatorLabel + " load",
			defaultParser,
			args{"testdata/communicator/basic.tf", new(PackerConfig)},
			&PackerConfig{
				Communicators: []*Communicator{
					{Type: "ssh", Name: "vagrant"},
				},
			},
			false,
		},

		{
			"duplicate " + sourceLabel, defaultParser,
			args{"testdata/sources/vb-iso.tf", &PackerConfig{
				Sources: map[SourceRef]*Source{
					SourceRef{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1204",
					}: {
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1204",
					},
				},
			},
			},
			&PackerConfig{
				Sources: map[SourceRef]*Source{
					SourceRef{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1204",
					}: &Source{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1204",
					},
					SourceRef{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1604",
					}: &Source{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1604",
					},
				},
			},
			true,
		},

		{"valid variables load", defaultParser,
			args{"testdata/variables/basic.tf", new(PackerConfig)},
			&PackerConfig{
				Variables: PackerV1Variables{
					"image_name": "foo-image-{{user `my_secret`}}",
					"key":        "value",
					"my_secret":  "foo",
				},
			},
			false,
		},

		{"valid " + buildLabel + " load", defaultParser,
			args{"testdata/build/basic.tf", new(PackerConfig)},
			&PackerConfig{
				Builds: Builds{
					&Build{
						Outputs: Outputs{
							Output{
								Type: "aws_ami",
								Name: "{{user `image_name`}}-aws-ubuntu-16.04",
							},
							Output{
								Type: "aws_ami",
								Name: "{{user `image_name`}}-vb-ubuntu-12.04",
							},
						},
						ProvisionerGroups: ProvisionerGroups{
							&ProvisionerGroup{
								Provisioners: []Provisioner{
									Provisioner{
										&hcl.Block{
											Type: "shell",
										},
									},
									Provisioner{
										&hcl.Block{
											Type: "shell",
										},
									},
									Provisioner{
										&hcl.Block{
											Type:   "upload",
											Labels: []string{"log.go", "/tmp"},
										},
									},
								},
							},
						},
					},
					&Build{
						Outputs: Outputs{
							Output{
								Type: "aws_ami",
								Name: "fooooobaaaar",
							},
						},
						ProvisionerGroups: ProvisionerGroups{
							&ProvisionerGroup{
								Provisioners: []Provisioner{
									Provisioner{
										&hcl.Block{
											Type: "shell",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.parser
			diags := p.Parse(tt.args.filename, tt.args.cfg)
			if tt.wantDiags == (diags == nil) {
				t.Errorf("PackerConfig.Load() unexpected diagnostics. %s", diags)
			}
			if diff := cmp.Diff(tt.wantPackerConfig, tt.args.cfg,
				cmpopts.IgnoreTypes(HCL2Ref{}),
				cmpopts.IgnoreTypes([]hcl.Range{}),
				cmpopts.IgnoreTypes(hcl.Range{}),
				cmpopts.IgnoreInterfaces(struct{ hcl.Expression }{}),
				cmpopts.IgnoreInterfaces(struct{ hcl.Body }{}),
			); diff != "" {
				t.Errorf("PackerConfig.Load() wrong packer config. %s", diff)
			}
			if t.Failed() {
				t.Fatal()
			}
		})
	}
}
