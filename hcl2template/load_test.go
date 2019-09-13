package hcl2template

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl2/hcl"
)

func TestPackerConfig_Load(t *testing.T) {
	type fields struct {
		Sources map[SourceRef]*Source

		Variables PackerV1Variables

		Builds Builds
	}
	tests := []struct {
		name             string
		fields           fields
		filename         string
		wantPackerConfig *PackerConfig
		wantDiags        bool
	}{
		{
			"valid " + sourceLabel + " load",
			fields{}, "testdata/sources/vb-iso.tf", &PackerConfig{
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
			"duplicate " + sourceLabel,
			fields{
				Sources: map[SourceRef]*Source{
					SourceRef{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1204",
					}: &Source{
						Type: "virtualbox-iso",
						Name: "vb-ubuntu-1204",
					},
				},
			}, "testdata/sources/vb-iso.tf", &PackerConfig{
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

		{"valid variables load",
			fields{}, "testdata/variables/basic.tf", &PackerConfig{
				Variables: PackerV1Variables{
					"image_name": "foo-image-{{user `my_secret`}}",
					"key":        "value",
					"my_secret":  "foo",
				},
			},
			false,
		},

		{"valid " + buildLabel + " load",
			fields{}, "testdata/build/basic.tf", &PackerConfig{
				Builds: Builds{
					Build{
						Outputs: Outputs{
							Output{
								Type: "aws_ami",
								Name: "{{user `image_name`}}-aws-ubuntu-16.04",
							},
							Output{
								Type: "aws_ami",
								Name: "{{user `image_name`}}-vb-ubuntu-12.04",
							},
							Output{
								Type: "aws_ami",
								Name: "{{user `image_name`}}-vmw-ubuntu-16.04",
							},
						},
					},
					Build{
						Outputs: Outputs{
							Output{
								Type: "aws_ami",
								Name: "fooooobaaaar",
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
			cfg := &PackerConfig{
				Sources:   tt.fields.Sources,
				Variables: tt.fields.Variables,
				Builds:    tt.fields.Builds,
			}
			diags := cfg.Load(tt.filename)
			if tt.wantDiags == (diags == nil) {
				t.Errorf("PackerConfig.Load() unexpected diagnostics. %s", diags)
			}
			if diff := cmp.Diff(cfg, tt.wantPackerConfig, cmpopts.IgnoreTypes(HCL2Ref{}), cmpopts.IgnoreInterfaces(struct{ hcl.Expression }{})); diff != "" {
				t.Errorf("PackerConfig.Load() wrong packer config. %s", diff)
			}
			if t.Failed() {
				t.Fatal()
			}
		})
	}
}
