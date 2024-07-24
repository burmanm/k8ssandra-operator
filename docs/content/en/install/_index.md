---
title: "Install K8ssandra Operator"
linkTitle: "Install K8ssandra Operator"
no_list: true
weight: 3
description: >
  How to install and manage an Apache Cassandra® cluster (or multiple clusters) in Kubernetes using K8ssandra Operator.
---

These quickstart topics step through the K8ssandra Operator install steps for single- and multi-cluster deployments of a `K8ssandraCluster` custom resource, Cassandra, and related services in a local Kubernetes environment.

* [Prerequisites]({{< relref "install/local/" >}}) for K8ssandra Operator installs.

## Quickstarts
* [Quickstart **single-cluster** install]({{< relref "install/local/single-cluster-helm/" >}}) of K8ssandra Operator with `helm`.
* [Quickstart **multi-cluster** install]({{< relref "install/local/multi-cluster-helm/" >}}) of K8ssandra Operator with `helm`. 
* [Quickstart **single-cluster** install]({{< relref "install/local/single-cluster-kustomize/" >}}) of K8ssandra Operator with `kustomize`.
* [Quickstart **multi-cluster** install]({{< relref "install/local/multi-cluster-kustomize/" >}}) of K8ssandra Operator with `kustomize`.

## Cloud Provider Guides
If you are using a cloud provider, explore the following topics for cloud-specific guidance on installing K8ssandra Operator with Helm or Kustomize.

* [Amazon Elastic Kubernetes Service]({{< relref "install/eks/" >}}) (EKS)
* [Azure Kubernetes Service]({{< relref "install/aks/" >}}) (AKS)
* [DigitalOcean Kubernetes]({{< relref "install/doks/" >}}) (DOKS)
* [Google Kubernetes Engine]({{< relref "install/gke/" >}}) (GKE)

## Upgrade

* [Upgrade notes]({{< relref "install/upgrade/" >}}) (Upgrade notes)

**Tip:** For an architectural overview of K8ssandra Operator and its new `K8ssandraCluster` custom resource, see the [K8ssandra Operator]({{< relref "components/k8ssandra-operator/" >}}) component page.
