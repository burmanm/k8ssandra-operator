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

- [ENHANCEMENT] [#1115](https://github.com/k8ssandra/k8ssandra-operator/issues/1115) Add a validation check for the projected pod names length
* [CHANGE] [#1050](https://github.com/k8ssandra/k8ssandra-operator/issues/1050) Remove unnecessary requeues in the Medusa controllers
* [ENHANCEMENT] [#1161](https://github.com/k8ssandra/k8ssandra-operator/issues/1161) Update cass-operator Helm chart to 0.46.1. Adds containerPort for cass-operator metrics and changes cass-config-builder base from UBI7 to UBI8