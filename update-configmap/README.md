update-configmap KRM function
============================

## Overview

Updates key-value pairs in any ConfigMap-like KRM resource.

## Usage

Finds the ConfigMap whose name matches with the `targetName` FunctionConfig parameter, and updates the key-value pairs in its `data` field with the key-value pairs given in FunctionConfig.

The KRM function actually works for not just ConfigMaps, but for any KRM resource that has a `data` field with a type of `map[string]string`.

Can be used both declaratively and imperatively.

## FunctionConfig

The KRM function supposed to be configured via a ConfigMap. 
The ConfigMap must contain the `targetName` key that contains the name of the ConfigMap-like resource to update.
All other key-value pairs in FunctionConfig are used to update the target.

## Example

This will ensure that the ConfigMap named `test-data` in the package has a `.data.newKey` field and contains the value `newValue`:

```yaml
apiVersion: kpt.dev/v1
kind: Kptfile
metadata:
  annotations:
    config.kubernetes.io/local-config: "true"
  name: testpkg
pipeline:
  mutators:
    - image: ghcr.io/kispaljr/update-configmap:v1
      configMap:
        targetName: test-data
        newKey: newValue
```
