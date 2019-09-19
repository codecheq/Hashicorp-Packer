package hcl2template

import (
	// "github.com/hashicorp/packer/common/template"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/packer/template"
	"github.com/mitchellh/mapstructure"
)

// PackerConfig represents a loaded packer config
type PackerConfig struct {
	Sources map[SourceRef]*Source

	Variables PackerV1Variables

	Builds Builds

	Communicators []*Communicator
}

func (pkrCfg *PackerConfig) ToTemplate() (*template.Template, error) {
	var result template.Template
	// var errs error

	result.Comments = nil                           // HAHA !
	result.Variables = pkrCfg.Variables.Variables() // TODO(azr): make pkrCfg.Variables the right type
	// TODO(azr): add sensitive variables

	result.Builders()
	for i, rawB := range pkrCfg.Builders {
		var b Builder
		if err := mapstructure.WeakDecode(rawB, &b); err != nil {
			errs = multierror.Append(errs, fmt.Errorf(
				"builder %d: %s", i+1, err))
			continue
		}

		// Set the raw configuration and delete any special keys
		b.Config = rawB.(map[string]interface{})

		delete(b.Config, "name")
		delete(b.Config, "type")

		if len(b.Config) == 0 {
			b.Config = nil
		}

		// If there is no type set, it is an error
		if b.Type == "" {
			errs = multierror.Append(errs, fmt.Errorf(
				"builder %d: missing 'type'", i+1))
			continue
		}

		// The name defaults to the type if it isn't set
		if b.Name == "" {
			b.Name = b.Type
		}

		// If this builder already exists, it is an error
		if _, ok := result.Builders[b.Name]; ok {
			errs = multierror.Append(errs, fmt.Errorf(
				"builder %d: builder with name '%s' already exists",
				i+1, b.Name))
			continue
		}

		// Append the builders
		result.Builders[b.Name] = &b
	}

	// // Gather all the post-processors
	// if len(r.PostProcessors) > 0 {
	// 	result.PostProcessors = make([][]*PostProcessor, 0, len(r.PostProcessors))
	// }
	// for i, v := range r.PostProcessors {
	// 	// Parse the configurations. We need to do this because post-processors
	// 	// can take three different formats.
	// 	configs, err := r.parsePostProcessor(i, v)
	// 	if err != nil {
	// 		errs = multierror.Append(errs, err)
	// 		continue
	// 	}

	// 	// Parse the PostProcessors out of the configs
	// 	pps := make([]*PostProcessor, 0, len(configs))
	// 	for j, c := range configs {
	// 		var pp PostProcessor
	// 		if err := r.decoder(&pp, nil).Decode(c); err != nil {
	// 			errs = multierror.Append(errs, fmt.Errorf(
	// 				"post-processor %d.%d: %s", i+1, j+1, err))
	// 			continue
	// 		}

	// 		// Type is required
	// 		if pp.Type == "" {
	// 			errs = multierror.Append(errs, fmt.Errorf(
	// 				"post-processor %d.%d: type is required", i+1, j+1))
	// 			continue
	// 		}

	// 		// Set the raw configuration and delete any special keys
	// 		pp.Config = c

	// 		// The name defaults to the type if it isn't set
	// 		if pp.Name == "" {
	// 			pp.Name = pp.Type
	// 		}

	// 		delete(pp.Config, "except")
	// 		delete(pp.Config, "only")
	// 		delete(pp.Config, "keep_input_artifact")
	// 		delete(pp.Config, "type")
	// 		delete(pp.Config, "name")

	// 		if len(pp.Config) == 0 {
	// 			pp.Config = nil
	// 		}

	// 		pps = append(pps, &pp)
	// 	}

	// 	result.PostProcessors = append(result.PostProcessors, pps)
	// }

	// // Gather all the provisioners
	// if len(r.Provisioners) > 0 {
	// 	result.Provisioners = make([]*Provisioner, 0, len(r.Provisioners))
	// }
	// for i, v := range r.Provisioners {
	// 	var p Provisioner
	// 	if err := r.decoder(&p, nil).Decode(v); err != nil {
	// 		errs = multierror.Append(errs, fmt.Errorf(
	// 			"provisioner %d: %s", i+1, err))
	// 		continue
	// 	}

	// 	// Type is required before any richer validation
	// 	if p.Type == "" {
	// 		errs = multierror.Append(errs, fmt.Errorf(
	// 			"provisioner %d: missing 'type'", i+1))
	// 		continue
	// 	}

	// 	// Set the raw configuration and delete any special keys
	// 	p.Config = v.(map[string]interface{})

	// 	delete(p.Config, "except")
	// 	delete(p.Config, "only")
	// 	delete(p.Config, "override")
	// 	delete(p.Config, "pause_before")
	// 	delete(p.Config, "type")
	// 	delete(p.Config, "timeout")

	// 	if len(p.Config) == 0 {
	// 		p.Config = nil
	// 	}

	// 	result.Provisioners = append(result.Provisioners, &p)
	// }

	// // If we have errors, return those with a nil result
	// if errs != nil {
	// 	return nil, errs
	// }

	return &result, nil
}
