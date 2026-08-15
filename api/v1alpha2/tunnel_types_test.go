package v1alpha2

import "testing"

func TestAPITokenSecretKeyPrefersLegacyOverrideWhenCamelCaseFieldHasDefaultValue(t *testing.T) {
	details := CloudflareDetails{
		ApiTokenKey:          DefaultCloudflareAPITokenSecretKey,
		CLOUDFLARE_API_TOKEN: "legacy-token",
	}

	if got := details.APITokenSecretKey(); got != "legacy-token" {
		t.Fatalf("expected legacy override to win when apiTokenKey is just the default, got %q", got)
	}
}

func TestAPITokenSecretKeyPrefersExplicitCamelCaseOverride(t *testing.T) {
	details := CloudflareDetails{
		ApiTokenKey:          "token",
		CLOUDFLARE_API_TOKEN: "legacy-token",
	}

	if got := details.APITokenSecretKey(); got != "token" {
		t.Fatalf("expected explicit camelCase override, got %q", got)
	}
}
