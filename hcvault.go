package cvault

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type inferFn = func() string

func InferVaultToken() string {
	fn := combineInfers(
		inferFromEnv,
		inferFromHelper,
		inferFromHomeToken,
	)
	return fn()
}

func combineInfers(fns ...inferFn) inferFn {
	return func() string {
		for _, fn := range fns {
			if token := fn(); token != "" {
				return token
			}
		}
		return ""
	}
}

func inferFromEnv() string {
	return os.Getenv("VAULT_TOKEN")
}

func inferFromHomeToken() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	tokenPath := filepath.Join(home, ".vault-token")
	fileToken, err := os.ReadFile(tokenPath)
	if err != nil {
		return ""
	}
	return string(fileToken)
}

func inferFromHelper() string {
	cfgFile := func() string {
		if f := os.Getenv("VAULT_CONFIG_FILE"); f != "" {
			return f
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".vault")
	}()
	cfgData, err := os.ReadFile(cfgFile)
	if err != nil {
		return ""
	}

	config := struct {
		TokenHelper string `hcl:"token_helper"`
	}{}
	if err := hclsimple.Decode("vault.hcl", cfgData, nil, &config); err != nil {
		return ""
	}
	cmd := exec.Command(config.TokenHelper, "get")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func InferVaultTlsConfig() (config *tls.Config) {
	config = &tls.Config{}
	caPath := os.Getenv("VAULT_CAPATH")
	if caPath == "" {
		return
	}
	pool, err := x509.SystemCertPool()
	if err != nil {
		return
	}
	certPem, err := ioutil.ReadFile(caPath)
	if err != nil {
		return
	}
	pool.AppendCertsFromPEM(certPem)
	config.RootCAs = pool
	return
}
