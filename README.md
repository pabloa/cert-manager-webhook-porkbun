<p align="center">
  <img src="https://raw.githubusercontent.com/cert-manager/cert-manager/d53c0b9270f8cd90d908460d69502694e1838f5f/logo/logo-small.png" height="256" width="256" alt="cert-manager project logo" />
</p>

# Porkbun Webhook for cert-manager

An implementation of the cert-manager [`webhook.Solver` interface](https://pkg.go.dev/github.com/cert-manager/cert-manager@v1.12.3/pkg/acme/webhook#Solver) for [Porkbun](https://porkbun.com/). This is based on [cert-manager/webhook-example](https://github.com/cert-manager/webhook-example), with inspiration from [baarde/cert-manager-webhook-ovh](https://github.com/baarde/cert-manager-webhook-ovh)

## Subdomain Support

This webhook properly supports issuing certificates for multi-level subdomains (e.g., `*.dev.example.com`, `*.stage.example.com`). The webhook automatically detects the authoritative Porkbun zone by querying the Porkbun API, ensuring that ACME DNS-01 challenges work correctly even when cert-manager passes incorrect zone information for subdomains. This enables environment-specific wildcard certificates for multi-environment Kubernetes deployments without requiring manual DNS delegation or workarounds.

Note: The test suite does work, but I straight up deleted `main_test.go` because the dependency on `github.com/cert-manager/cert-manager/test/acme` was giving me insane, impossible to resolve versioning conflicts. I'm sure these will be resolved by someone more knowledgeable updating the `go.mod` in the example webhook, at which point I'll add the tests back.

## Installation

### Install cert-manager

Install cert-manager using its [installation documentation](https://cert-manager.io/docs/installation/).

### Install webhook

Add helm repo:

```bash
helm repo add cert-manager-webhook-porkbun https://talinx.github.io/cert-manager-webhook-porkbun
```

[Generate a porkbun API key](https://kb.porkbun.com/article/190-getting-started-with-the-porkbun-api) and create a secret with it:

```yaml
apiVersion: v1
stringData:
  PORKBUN_API_KEY: pk1_yourapikeyhere
  PORKBUN_SECRET_API_KEY: sk1_yoursecretkeyhere
kind: Secret
metadata:
  name: porkbun-secret
  namespace: cert-manager
type: Opaque
```

Install helm chart in a namespace of your choice, e. g. `cert-manager`:

```bash
helm install cert-manager-webhook-porkbun cert-manager-webhook-porkbun/cert-manager-webhook-porkbun -n cert-manager
```

Add an issuer (change the email address; the `groupName` has to match the `groupName` value of the helm chart), e. g.:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-issuer
spec:
  acme:
    email: youremailhere@example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-porkbun-tls
    solvers:
    - dns01:
        webhook:
          groupName: porkbun.talinx.dev
          solverName: porkbun
          config:
            apiKey:
              key: PORKBUN_API_KEY
              name: porkbun-secret
            secretApiKey:
              key: PORKBUN_SECRET_API_KEY
              name: porkbun-secret
```

Add a certificate, e. g.:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: prod-cert
spec:
  secretName: prod-cert
  issuerRef:
    name: letsencrypt-issuer
    kind: ClusterIssuer
  dnsNames:
  - 'mysub.example.com'
```

### Running the test suite

All DNS providers **must** run the DNS01 provider conformance testing suite,
else they will have undetermined behaviour when used with cert-manager.

**It is essential that you configure and run the test suite when creating a
DNS01 webhook.**

An example Go test file has been provided in [main_test.go](https://github.com/pabloa/cert-manager-webhook-porkbun/blob/master/main_test.go).

You can run the test suite with:

```bash
$ TEST_ZONE_NAME=example.com. make test
```

The example file has a number of areas you must fill in and replace with your
own options in order for tests to pass.
