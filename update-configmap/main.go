package main

import (
	"fmt"
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
)

// Run is the entry point of the KRM function
func Run(rl *fn.ResourceList) (success bool, err error) {
	success = true // success return value is not useful for us

	// parse FunctionConfig
	newData, found, err := rl.FunctionConfig.NestedStringMap("data")
	if !found {
		err = fmt.Errorf("mandatory field '.data' is missing from FunctionConfig. It should be a ConfigMap.")
	}
	if err != nil {
		rl.LogResult(err)
		return
	}
	targetName, found := newData["targetName"]
	if !found || targetName == "" {
		err = fmt.Errorf("mandatory field '.data.targetName' is missing from FunctionConfig")
		rl.LogResult(err)
		return
	}
	delete(newData, "targetName")

	// find target
	var target *fn.KubeObject
	for _, item := range rl.Items {
		if item.GetName() == targetName {
			if item.GetKind() == "ConfigMap" {
				target = item
				break
			}
			_, found, dataErr := item.NestedStringMap("data")
			if found && dataErr == nil {
				target = item
				// do not 'break', but continue to check if there is also a ConfigMap with the same name
			}
		}
	}
	if target == nil {
		target = fn.NewEmptyKubeObject()
		target.SetAPIVersion("v1")
		target.SetKind("ConfigMap")
		target.SetName(targetName)
		target.SetAnnotation(kioutil.PathAnnotation, targetName+".yaml")
		rl.Items = append(rl.Items, target)
	}
	targetData := target.UpsertMap("data")
	for key, value := range newData {
		err = targetData.SetNestedString(value, key)
		if err != nil {
			rl.LogResult(err)
			return
		}
	}
	return
}

func main() {
	if err := fn.AsMain(fn.ResourceListProcessorFunc(Run)); err != nil {
		os.Exit(1)
	}
}
