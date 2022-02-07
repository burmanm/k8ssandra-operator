/*
Copyright 2021.

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

package v1alpha1

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/bombsimon/logrusr"
	"github.com/k8ssandra/k8ssandra-operator/pkg/clientcache"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	admissionv1 "k8s.io/api/admission/v1"

	//+kubebuilder:scaffold:imports
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	reaperapi "github.com/k8ssandra/k8ssandra-operator/apis/reaper/v1alpha1"
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx context.Context
var cancel context.CancelFunc

func TestWebhook(t *testing.T) {
	require := require.New(t)
	ctx, cancel = context.WithCancel(context.TODO())

	log := logrusr.NewLogger(logrus.New())
	logf.SetLogger(log)

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: false,
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{filepath.Join("..", "..", "..", "config", "webhook")},
		},
	}

	cfg, err := testEnv.Start()
	require.NoError(err)
	require.NotNil(cfg)

	defer cancel()
	defer testEnv.Stop()

	scheme := runtime.NewScheme()
	err = AddToScheme(scheme)
	require.NoError(err)

	err = corev1.AddToScheme(scheme)
	require.NoError(err)

	err = admissionv1.AddToScheme(scheme)
	require.NoError(err)

	err = AddToScheme(scheme)
	require.NoError(err)

	err = reaperapi.AddToScheme(scheme)
	require.NoError(err)

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	require.NoError(err)
	require.NotNil(k8sClient)

	// start webhook server using Manager
	webhookInstallOptions := &testEnv.WebhookInstallOptions
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme,
		Host:               webhookInstallOptions.LocalServingHost,
		Port:               webhookInstallOptions.LocalServingPort,
		CertDir:            webhookInstallOptions.LocalServingCertDir,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	require.NoError(err)

	clientCache := clientcache.New(k8sClient, k8sClient, scheme)
	clientCache.AddClient("envtest", k8sClient)
	err = (&K8ssandraCluster{}).SetupWebhookWithManager(mgr, clientCache)
	require.NoError(err)

	//+kubebuilder:scaffold:webhook

	go func() {
		err = mgr.Start(ctx)
		require.NoError(err)
	}()

	// wait for the webhook server to get ready
	dialer := &net.Dialer{Timeout: time.Second}
	addrPort := fmt.Sprintf("%s:%d", webhookInstallOptions.LocalServingHost, webhookInstallOptions.LocalServingPort)
	require.Eventually(func() bool {
		conn, err := tls.DialWithDialer(dialer, "tcp", addrPort, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}, 2*time.Second, 300*time.Millisecond)

	t.Run("TestCreateValidations", testValidationCreate)
	t.Run("TestUpdateValidations", testValidationUpdate)
}

func testValidationCreate(t *testing.T) {
	require := require.New(t)
	createNamespace(require, "create-namespace")
	cluster := createMinimalClusterObj("create-test", "create-namespace")

	err := k8sClient.Create(ctx, cluster)
	require.NoError(err)

	cluster.Spec.Cassandra.Datacenters[0].K8sContext = "wrong"
	err = k8sClient.Create(ctx, cluster)
	require.Error(err)
}

func testValidationUpdate(t *testing.T) {
	require := require.New(t)
	createNamespace(require, "update-namespace")
	cluster := createMinimalClusterObj("update-test", "update-namespace")

	cluster.Spec.Reaper = &reaperapi.ReaperClusterTemplate{
		ReaperTemplate: reaperapi.ReaperTemplate{
			Keyspace: "original",
		},
	}

	err := k8sClient.Create(ctx, cluster)
	require.NoError(err)

	cluster.Spec.Reaper.ReaperTemplate.Keyspace = "modified"
	err = k8sClient.Update(ctx, cluster)
	require.Error(err)
}

func createNamespace(require *require.Assertions, namespace string) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	err := k8sClient.Create(ctx, ns)
	require.NoError(err)
}

func createMinimalClusterObj(name, namespace string) *K8ssandraCluster {
	return &K8ssandraCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: K8ssandraClusterSpec{
			Cassandra: &CassandraClusterTemplate{
				Datacenters: []CassandraDatacenterTemplate{
					{
						K8sContext: "envtest",
						Size:       1,
					},
				},
			},
		},
	}
}
