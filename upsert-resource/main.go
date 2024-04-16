package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Run is the entry point of the KRM function
func Run(rl *fn.ResourceList) (success bool, err error) {
	success = true // success return value is not useful for us

	inYaml, found, err := rl.FunctionConfig.NestedString("data", "resources")
	if !found {
		err = fmt.Errorf("mandatory field '.data.resources' is missing from function config")
	}
	if err != nil {
		rl.LogResult(err)
		return
	}
	inKObjs, err := fn.ParseKubeObjects([]byte(inYaml))
	if err != nil {
		rl.LogResult(err)
		return
	}
	for _, inKobj := range inKObjs {
		// NOTE: rl.UpsertObjectToItems doesn't handle "path annotations" as expected in this KRM function

		idx := -1
		for i, item := range rl.Items {
			if sameObject(inKobj, item) {
				idx = i
				break
			}
		}
		if idx == -1 {
			err = inKobj.SetAnnotation(kioutil.PathAnnotation, "injected_by_upsert-resource-fn.yaml")
			if err != nil {
				rl.LogResult(err)
				return
			}
			rl.Items = append(rl.Items, inKobj)
		} else {
			err = inKobj.SetAnnotation(kioutil.PathAnnotation, rl.Items[idx].PathAnnotation())
			if err != nil {
				rl.LogResult(err)
				return
			}
			err = inKobj.SetAnnotation(kioutil.IndexAnnotation, rl.Items[idx].GetAnnotation(kioutil.IndexAnnotation))
			if err != nil {
				rl.LogResult(err)
				return
			}
			rl.Items[idx] = inKobj
		}
	}
	return
}

func main() {
	if err := fn.AsMain(fn.ResourceListProcessorFunc(Run)); err != nil {
		os.Exit(1)
	}
}

func resourceIdentifier(o *fn.KubeObject) *yaml.ResourceIdentifier {
	apiVersion := o.GetAPIVersion()
	kind := o.GetKind()
	name := o.GetName()
	ns := o.GetNamespace()
	return &yaml.ResourceIdentifier{
		TypeMeta: yaml.TypeMeta{
			APIVersion: apiVersion,
			Kind:       kind,
		},
		NameMeta: yaml.NameMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func sameObject(o1, o2 *fn.KubeObject) bool {
	id1 := resourceIdentifier(o1)
	id2 := resourceIdentifier(o2)
	return reflect.DeepEqual(id1, id2)
}
