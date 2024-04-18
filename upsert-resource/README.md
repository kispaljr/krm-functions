upsert-resource KRM function
============================

## Overview

Insert a KRM resource, or if the resource already exists, update the existing resource.

## Usage

This is a re-implementation of the [baseline upsert-resource KRM function](https://github.com/GoogleContainerTools/kpt-functions-catalog/tree/master/functions/go/upsert-resource). Unlike the original, this version was designed to be used declaratively, particularly in the mutation pipeline of KPT packages. 

## FunctionConfig

The KRM function supposed to be configured via a ConfigMap with a single mandatory `resources` key in it. The value of that field should contain the YAML representation of the resources to be upserted.

## Example

The upsert-resource KRM function is primarily useful in the context of a `PackageVariant` resource, e.g.:
```yaml
apiVersion: config.porch.kpt.dev/v1alpha1
kind: PackageVariant
metadata:
  name: example
spec:
  upstream:
    repo: catalog
    package: blueprint
    revision: v1
  downstream:
    repo: deployments
  pipeline:
    mutators:
    - image: ghcr.io/kispaljr/upsert-resource:v1
      configMap:
        resources: |-
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: test-data-1
            namespace: default
          data:
            inserted-key: test-value
          ---
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: test-data-2
            namespace: default
          data:
            inserted-key: test-value
```
The above example will insert or update the `test-data-1` and `test-data-2` ConfigMaps into the downstream package.