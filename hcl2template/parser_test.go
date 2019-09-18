package hcl2template

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

func TestParser_Parse(t *testing.T) {
	defaultParser := &Parser{
		hclparse.NewParser(),
		&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{Type: "shell"},
				{Type: "upload", LabelNames: []string{"source", "destination"}},
			}},
		&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{Type: "amazon-import"},
			}},
	}

	type args struct {
		filename string
	}
	tests := []struct {
		name      string
		fields    *Parser
		args      args
		wantCfg   *PackerConfig
		wantDiags bool
	}{
		{"complete",
			defaultParser,
			args{"testdata/complete"},
			&PackerConfig{
				Sources: map[SourceRef]*Source{
					SourceRef{
						Type: "virtualbox-iso",
						Name: "ubuntu-1204",
					}: {
						Type: "virtualbox-iso",
						Name: "ubuntu-1204",
					},
					SourceRef{
						Type: "amazon-ebs",
						Name: "ubuntu-1604",
					}: {
						Type: "amazon-ebs",
						Name: "ubuntu-1604",
					},
					SourceRef{
						Type: "amazon-ebs",
						Name: "{{user `image_name`}}-ubuntu-1.0",
					}: {
						Type: "amazon-ebs",
						Name: "{{user `image_name`}}-ubuntu-1.0",
					},
				},
				Communicators: []*Communicator{
					{Type: "ssh", Name: "vagrant"},
				},
				Variables: PackerV1Variables{
					"image_name": "foo-image-{{user `my_secret`}}",
					"key":        "value",
					"my_secret":  "foo",
				},
				Builds: Builds{
					{
						Froms: BuildFromList{
							{
								Src: "src.amazon-ebs.ubuntu-1604",
							},
							{
								Src: "src.virtualbox-iso.ubuntu-1204",
							},
						},
						ProvisionerGroups: ProvisionerGroups{
							&ProvisionerGroup{
								Provisioners: []Provisioner{
									{
										&hcl.Block{
											Type: "shell",
										},
									},
									{
										&hcl.Block{
											Type: "shell",
										},
									},
									{
										&hcl.Block{
											Type:   "upload",
											Labels: []string{"log.go", "/tmp"},
										},
									},
								},
							},
						},
						PostProvisionerGroups: ProvisionerGroups{
							&ProvisionerGroup{
								Provisioners: []Provisioner{
									{
										&hcl.Block{
											Type: "amazon-import",
										},
									},
								},
							},
						},
					},
					&Build{
						Froms: BuildFromList{
							{
								Src: "src.amazon.{{user `image_name`}-ubuntu-1.0",
							},
						},
						ProvisionerGroups: ProvisionerGroups{
							{
								Provisioners: []Provisioner{
									{
										&hcl.Block{
											Type: "shell",
										},
									},
								},
							},
						},
					},
				},
			}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Parser:                 tt.fields.Parser,
				ProvisionersSchema:     tt.fields.ProvisionersSchema,
				PostProvisionersSchema: tt.fields.PostProvisionersSchema,
			}
			gotCfg, gotDiags := p.Parse(tt.args.filename)
			if tt.wantDiags == (gotDiags == nil) {
				t.Errorf("Parser.Parse() unexpected diagnostics. %s", gotDiags)
			}
			if diff := cmp.Diff(tt.wantCfg, gotCfg,
				cmpopts.IgnoreTypes(HCL2Ref{}),
				cmpopts.IgnoreTypes([]hcl.Range{}),
				cmpopts.IgnoreTypes(hcl.Range{}),
				cmpopts.IgnoreInterfaces(struct{ hcl.Expression }{}),
				cmpopts.IgnoreInterfaces(struct{ hcl.Body }{}),
			); diff != "" {
				t.Errorf("Parser.Parse() wrong packer config. %s", diff)
			}

		})
	}
}
