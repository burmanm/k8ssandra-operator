# Monitoring using Prometheus

## Installing and configuring Prometheus for monitoring

`k8ssandra-operator` has integrations with Prometheus which allow for the simple rollout of Prometheus ServiceMonitors for both Stargate and Cassandra Datacenters.

### Prerequisites

To use Prometheus for monitoring, you need to have the Prometheus operator installed on your Kubernetes (k8s) cluster. A simple way to install Prometheus operator is via [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus). 

```
kubectl -k https://github.com/prometheus-operator/kube-prometheus
```

The Prometheus operator installs the ServiceMonitor CRD, which is the integration point we use to tell Prometheus how to find the Stargate and Cassandra pods and what endpoints on those pods to scrape.

### Permissions

When Prometheus operator is installed via `kube-prometheus`, it does not provide itself with sufficient permissions to watch pods in other namespaces.

To address this, you should add the following RBAC permissions to the `ClusterRole` `prometheus-k8s`.

```
  - verbs:
    - get
    - list
    - watch
    apiGroups:
    - ''
    resources:
    - pods
    - endpoints
    - services
```

### Enabling Prometheus for Cassandra Datacenters

There are two ways to enable Prometheus monitoring for CassandraDatacenters. As with most Cassandra related configurations, it can be enabled cluster-wide as per the below:

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
spec:
  cassandra:
    telemetry: 
      prometheus:
        enabled: true
```

It can also be enabled (or disabled) for individual DCs within the cluster via individual `cassandraTelemetry` fields which reside within the `DCs` array of the CRD:

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
spec:
  cassandra:
    datacenters:
    - metadata: 
        name: dc1
      telemetry: 
        prometheus:
          enabled: true
```

The DC level setting will take precedence in the event of a conflict between cluster-level and DC level settings.

### Enabling Prometheus for Stargate Deployments

Like the Cassandra Datacenter monitoring, Stargate monitoring can be enabled for every DC within the cluster using the cluster template:

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
spec:
  stargate:
    telemetry: 
      prometheus:
        enabled: true
```

It can also be enabled and disabled for individual DCs.

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
spec:
cassandra:
    datacenters:
    - metadata: 
        name: dc1
      stargate:
        telemetry: 
          prometheus:
            enabled: true
```

### Filtering metrics

Cassandra provides a lot of metrics which can create some overload, especially when there are many tables in a cluster. [Filtering rules for MCAC](https://github.com/datastax/metric-collector-for-apache-cassandra/blob/master/config/metric-collector.yaml#L9-L72) can be defined in the telemetry spec:

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
spec:
  cassandra:
    telemetry: 
      prometheus:
        enabled: true
        mcacMetricFilters:
          - "deny:org.apache.cassandra.metrics.Table"
          - "allow:org.apache.cassandra.metrics.Table.LiveSSTableCount"
```

When no filter is explicitly defined in the spec, default K8ssandra v1.x filters will be applied.

## Using Prometheus

Prometheus will generally be running in the `monitoring` namespace. To connect to the Prometheus GUI for an initial inspection of what metrics it is scraping, you can run the following `kubectl` command to forward a port.

```
kubectl port-forward -n monitoring services/prometheus-k8s 10000:9090
```

By navigating to `http://localhost:10000` in your web browser, you can then view the Prometheus GUI and examine the metrics flowing into your Prometheus instance from the k8ssandra-operator components. Searching for `collectd`, `mcac` or `stargate` in the Prometheus metrics search box will help you to find the metrics you're after.

It is worth noting that Prometheus typically comes up well ahead of Cassandra and Stargate, so these metrics may not appear until the cluster is fully bootstrapped.