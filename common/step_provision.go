package common

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

// StepProvision runs the provisioners.
//
// Uses:
//   communicator packer.Communicator
//   hook         packer.Hook
//   ui           packer.Ui
//
// Produces:
//   <nothing>
type StepProvision struct {
	Comm packer.Communicator
}

func PopulateProvisionHookData(state multistep.StateBag) packer.ProvisionHookData {
	hookData := packer.NewProvisionHookData()

	// Add WinRMPassword to runtime data
	WinRMPassword, ok := state.GetOk("winrm_password")
	if ok {
		hookData.WinRMPassword = WinRMPassword.(string)
	}
	return hookData
}

func (s *StepProvision) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	comm := s.Comm
	if comm == nil {
		raw, ok := state.Get("communicator").(packer.Communicator)
		if ok {
			comm = raw.(packer.Communicator)
		}
	}

	hook := state.Get("hook").(packer.Hook)
	ui := state.Get("ui").(packer.Ui)

	hookData := PopulateProvisionHookData(state)

	// Run the provisioner in a goroutine so we can continually check
	// for cancellations...
	log.Println("Running the provision hook")
	errCh := make(chan error, 1)
	go func() {
		errCh <- hook.Run(ctx, packer.HookProvision, ui, comm, &hookData)
	}()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				state.Put("error", err)
				return multistep.ActionHalt
			}

			return multistep.ActionContinue
		case <-ctx.Done():
			log.Printf("Cancelling provisioning due to context cancellation: %s", ctx.Err())
			return multistep.ActionHalt
		case <-time.After(1 * time.Second):
			if _, ok := state.GetOk(multistep.StateCancelled); ok {
				log.Println("Cancelling provisioning due to interrupt...")
				return multistep.ActionHalt
			}
		}
	}
}

func (*StepProvision) Cleanup(multistep.StateBag) {}
