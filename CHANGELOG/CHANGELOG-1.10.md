# Changelog

Changelog for the K8ssandra Operator, new PRs should update the `unreleased` section below with entries describing the changes like:

```markdown
* [CHANGE]
* [FEATURE]
* [ENHANCEMENT]
* [BUGFIX]
* [DOCS]
* [TESTING]
```

When cutting a new release, update the `unreleased` heading to the tag being generated and date, like `## vX.Y.Z - YYYY-MM-DD` and create a new placeholder section for  `unreleased` entries.

## unreleased

* [CHANGE] [#1088](https://github.com/k8ssandra/k8ssandra-operator/issues/1088) Use the Scarf proxy for image coordinates
* [ENHANCEMENT] [#1073](https://github.com/k8ssandra/k8ssandra-operator/issues/1073) Add a namespace label to the Cassandra metrics 
* [BUGFIX] [#1060](https://github.com/k8ssandra/k8ssandra-operator/issues/1060) Fix restore mapping shuffling nodes when restoring in place
* [BUGFIX] [#1061](https://github.com/k8ssandra/k8ssandra-operator/issues/1061) Point to cass-config-builder 1.0.7 for arm64 compatibility
* [ENHANCEMENT] [#956](https://github.com/k8ssandra/k8ssandra-operator/issues/956) Enable linting in the project
* [BUGFIX] [#1102](https://github.com/k8ssandra/k8ssandra-operator/issues/1102) Update gRPC maximum receive size to 512MB. Note, the operator might need more max memory than the default to take advantage of this.