package controller

import (
	"context"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	networkingv1alpha2 "github.com/adyanth/cloudflare-operator/api/v1alpha2"
)

func TestGetAPIDetailsUsesDefaultSecretKeys(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("add corev1 scheme: %v", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cloudflare-secrets",
			Namespace: "default",
		},
		Data: map[string][]byte{
			networkingv1alpha2.DefaultCloudflareAPITokenSecretKey: []byte("token-value"),
		},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()
	spec := networkingv1alpha2.TunnelSpec{
		Cloudflare: networkingv1alpha2.CloudflareDetails{
			Secret: "cloudflare-secrets",
			Domain: "example.com",
		},
	}

	api, gotSecret, err := getAPIDetails(context.Background(), c, logr.Discard(), spec, networkingv1alpha2.TunnelStatus{}, "default")
	if err != nil {
		t.Fatalf("get api details: %v", err)
	}
	if api == nil {
		t.Fatal("expected api details")
	}
	if gotSecret == nil {
		t.Fatal("expected secret")
	}
	if gotSecret.Name != secret.Name {
		t.Fatalf("expected secret %q, got %q", secret.Name, gotSecret.Name)
	}
}

func TestGetAPIDetailsUsesCustomAPITokenKey(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("add corev1 scheme: %v", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cloudflare-secrets",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"token": []byte("token-value"),
		},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()
	spec := networkingv1alpha2.TunnelSpec{
		Cloudflare: networkingv1alpha2.CloudflareDetails{
			Secret:      "cloudflare-secrets",
			Domain:      "example.com",
			ApiTokenKey: "token",
		},
	}

	if _, _, err := getAPIDetails(context.Background(), c, logr.Discard(), spec, networkingv1alpha2.TunnelStatus{}, "default"); err != nil {
		t.Fatalf("get api details with custom token key: %v", err)
	}
}

func TestGetAPIDetailsErrorsWhenConfiguredKeyMissing(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("add corev1 scheme: %v", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cloudflare-secrets",
			Namespace: "default",
		},
		Data: map[string][]byte{},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()
	spec := networkingv1alpha2.TunnelSpec{
		Cloudflare: networkingv1alpha2.CloudflareDetails{
			Secret:      "cloudflare-secrets",
			Domain:      "example.com",
			ApiTokenKey: "token",
		},
	}

	_, _, err := getAPIDetails(context.Background(), c, logr.Discard(), spec, networkingv1alpha2.TunnelStatus{}, "default")
	if err == nil {
		t.Fatal("expected error when configured key is missing")
	}
	if !strings.Contains(err.Error(), "token") {
		t.Fatalf("expected error to mention missing configured key, got %q", err)
	}
}

func TestGetAPIDetailsUsesLegacyCustomTokenKey(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("add corev1 scheme: %v", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cloudflare-secrets",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"legacy-token": []byte("token-value"),
		},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()
	spec := networkingv1alpha2.TunnelSpec{
		Cloudflare: networkingv1alpha2.CloudflareDetails{
			Secret:               "cloudflare-secrets",
			Domain:               "example.com",
			CLOUDFLARE_API_TOKEN: "legacy-token",
		},
	}

	if _, _, err := getAPIDetails(context.Background(), c, logr.Discard(), spec, networkingv1alpha2.TunnelStatus{}, "default"); err != nil {
		t.Fatalf("get api details with legacy token key: %v", err)
	}
}
