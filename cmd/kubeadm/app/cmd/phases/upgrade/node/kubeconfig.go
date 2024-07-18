/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package node holds phases of "kubeadm upgrade node".
package node

import (
	"fmt"

	"github.com/pkg/errors"

	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
	"k8s.io/kubernetes/cmd/kubeadm/app/features"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/upgrade"
)

// NewKubeconfigPhase creates a kubeadm workflow phase that implements handling of kubeconfig upgrade.
func NewKubeconfigPhase() workflow.Phase {
	phase := workflow.Phase{
		Name:   "kubeconfig",
		Short:  "Upgrade kubeconfig files for this node",
		Run:    runKubeconfig(),
		Hidden: true,
		InheritFlags: []string{
			options.DryRun,
			options.KubeconfigPath,
		},
	}
	return phase
}

func runKubeconfig() func(c workflow.RunData) error {
	return func(c workflow.RunData) error {
		data, ok := c.(Data)
		if !ok {
			return errors.New("kubeconfig phase invoked with an invalid data struct")
		}

		// if this is not a control-plane node, this phase should not be executed
		if !data.IsControlPlaneNode() {
			fmt.Println("[upgrade] Skipping phase. Not a control plane node.")
			return nil
		}

		// otherwise, retrieve all the info required for kubeconfig upgrade
		cfg := data.InitCfg()
		dryRun := data.DryRun()

		if features.Enabled(cfg.FeatureGates, features.ControlPlaneKubeletLocalMode) {
			if err := upgrade.UpdateKubeletLocalMode(cfg, dryRun); err != nil {
				return errors.Wrap(err, "failed to update kubelet local mode")
			}
		}

		fmt.Println("[upgrade] The kubeconfig for this node was successfully updated!")

		return nil
	}
}
