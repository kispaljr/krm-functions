package main

import (
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Run is the entry point of the KRM function
func Run(rl *fn.ResourceList) (success bool, err error) {
	success = true // success return value is not useful for us
	for _, ko := range rl.Items.Where(fn.IsGroupKind(schema.ParseGroupKind("NFTopology.topology.nephio.org"))) {
		err = ko.RemoveAnnotationsIfEmpty()
		if err != nil {
			rl.Results.ErrorE(err)
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
