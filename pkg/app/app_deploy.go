package app

import (
	"fmt"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	ctldep "github.com/k14s/kapp-controller/pkg/deploy"
	"github.com/k14s/kapp-controller/pkg/exec"
)

func (a *App) deploy(tplOutput string, changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {
	err := a.blockDeletion()
	if err != nil {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Blocking for deploy: %s", err))
	}

	if len(a.app.Spec.Deploy) != 1 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected exactly one deploy option"))
	}

	var result exec.CmdRunResult

	for _, dep := range a.app.Spec.Deploy {
		switch {
		case dep.Kapp != nil:
			cancelCh, closeCancelCh := a.newCancelCh()
			defer closeCancelCh()

			kapp, err := a.newKapp(*dep.Kapp, cancelCh)
			if err != nil {
				return exec.NewCmdRunResultWithErr(fmt.Errorf("Preparing kapp: %s", err))
			}

			result = kapp.Deploy(tplOutput, a.startFlushingAllStatusUpdates, changedFunc)

		default:
			result.AttachErrorf("%s", fmt.Errorf("Unsupported way to deploy"))
		}

		if result.Error != nil {
			break
		}
	}

	return result
}

func (a *App) delete(changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {
	if len(a.app.Spec.Deploy) != 1 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected exactly one deploy option"))
	}

	var result exec.CmdRunResult

	for _, dep := range a.app.Spec.Deploy {
		switch {
		case dep.Kapp != nil:
			cancelCh, closeCancelCh := a.newCancelCh()
			defer closeCancelCh()

			kapp, err := a.newKapp(*dep.Kapp, cancelCh)
			if err != nil {
				return exec.NewCmdRunResultWithErr(fmt.Errorf("Preparing kapp: %s", err))
			}

			result = kapp.Delete(a.startFlushingAllStatusUpdates, changedFunc)

		default:
			result.AttachErrorf("%s", fmt.Errorf("Unsupported way to delete"))
		}

		if result.Error != nil {
			break
		}
	}

	if result.Error == nil {
		err := a.unblockDeletion()
		if err != nil {
			return exec.NewCmdRunResultWithErr(fmt.Errorf("Unblocking for deploy: %s", err))
		}
	}

	return result
}

func (a *App) inspect() exec.CmdRunResult {
	if len(a.app.Spec.Deploy) != 1 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected exactly one deploy option"))
	}

	var result exec.CmdRunResult

	for _, dep := range a.app.Spec.Deploy {
		switch {
		case dep.Kapp != nil:
			cancelCh, closeCancelCh := a.newCancelCh()
			defer closeCancelCh()

			kapp, err := a.newKapp(*dep.Kapp, cancelCh)
			if err != nil {
				return exec.NewCmdRunResultWithErr(fmt.Errorf("Preparing kapp: %s", err))
			}

			result = kapp.Inspect()

		default:
			result.AttachErrorf("%s", fmt.Errorf("Unsupported way to inspect"))
		}

		if result.Error != nil {
			break
		}
	}

	return result
}

func (a *App) newKapp(kapp v1alpha1.AppDeployKapp, cancelCh chan struct{}) (*ctldep.Kapp, error) {
	genericOpts := ctldep.GenericOpts{Name: a.app.Name, Namespace: a.app.Namespace}

	return a.deployFactory.NewKapp(kapp, a.app.Spec.ServiceAccountName,
		a.app.Spec.Cluster, genericOpts, cancelCh)
}
