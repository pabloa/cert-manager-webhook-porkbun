package main

import (
	"os"
	"testing"

	dns "github.com/cert-manager/cert-manager/test/acme"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.
	//

	// Uncomment the below fixture when implementing your custom DNS provider
	fixture := dns.NewFixture(&porkbunDNSProviderSolver{},
		dns.SetResolvedZone(zone),
		dns.SetAllowAmbientCredentials(false),
		dns.SetManifestPath("testdata/porkbun"),
		// dns.SetBinariesPath("_test/kubebuilder/bin"),
	)
	//need to uncomment and  RunConformance delete runBasic and runExtended once https://github.com/cert-manager/cert-manager/pull/4835 is merged
	//fixture.RunConformance(t)
	fixture.RunBasic(t)
	fixture.RunExtended(t)
}

func TestMultiLevelSubdomain(t *testing.T) {
	// Test that the webhook correctly handles multi-level subdomains
	// (e.g., myapp.dev.example.com) by finding the authoritative zone.
	// This was the original bug: the webhook would try to use dev.example.com
	// as a zone instead of example.com, causing Porkbun to return ERROR.

	fixture := dns.NewFixture(&porkbunDNSProviderSolver{},
		dns.SetResolvedZone(zone),
		dns.SetAllowAmbientCredentials(false),
		dns.SetManifestPath("testdata/porkbun"),
		dns.SetDNSName("myapp.dev."+zone),
		dns.SetResolvedFQDN("_acme-challenge.myapp.dev."+zone),
	)

	fixture.RunBasic(t)
}
