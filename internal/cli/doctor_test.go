// Copyright 2026 sidhartha1s. Licensed under Apache-2.0. See LICENSE.

package cli

import "testing"

// Test that the interstitial detection function recognizes known bot walls.
func TestLooksLikeDoctorInterstitial_Cloudflare(t *testing.T) {
	body := []byte(`<html><head><title>Just a moment...</title></head><body>Checking your browser</body></html>`)
	vendor := looksLikeDoctorInterstitial(body)
	if vendor != "Cloudflare" {
		t.Errorf("expected Cloudflare, got %q", vendor)
	}
}

func TestLooksLikeDoctorInterstitial_Akamai(t *testing.T) {
	body := []byte(`<html><head><title>Access Denied</title></head><body>Akamai - Request Unsuccessful</body></html>`)
	vendor := looksLikeDoctorInterstitial(body)
	if vendor != "Akamai" {
		t.Errorf("expected Akamai, got %q", vendor)
	}
}

func TestLooksLikeDoctorInterstitial_AWSWAF(t *testing.T) {
	body := []byte(`<html><head><title>Request Blocked</title></head><body>AWS WAF has blocked this request</body></html>`)
	vendor := looksLikeDoctorInterstitial(body)
	if vendor != "AWS WAF" {
		t.Errorf("expected AWS WAF, got %q", vendor)
	}
}

func TestLooksLikeDoctorInterstitial_NoMatch(t *testing.T) {
	body := []byte(`{"status": "ok", "message": "API is healthy"}`)
	vendor := looksLikeDoctorInterstitial(body)
	if vendor != "" {
		t.Errorf("expected empty string for JSON response, got %q", vendor)
	}
}

func TestLooksLikeDoctorInterstitial_EmptyBody(t *testing.T) {
	vendor := looksLikeDoctorInterstitial([]byte{})
	if vendor != "" {
		t.Errorf("expected empty string for empty body, got %q", vendor)
	}
}

// Test that a normal HTML page with a title doesn't false-positive.
func TestLooksLikeDoctorInterstitial_NormalHTML(t *testing.T) {
	body := []byte(`<html><head><title>Welcome to our API</title></head><body>Documentation here</body></html>`)
	vendor := looksLikeDoctorInterstitial(body)
	if vendor != "" {
		t.Errorf("expected empty string for normal HTML, got %q", vendor)
	}
}
