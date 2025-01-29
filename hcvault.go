package cvault

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/vault/api"
	"github.com/tink-crypto/tink-go-hcvault/v2/integration/hcvault"
	"github.com/tink-crypto/tink-go/v2/core/registry"
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

func VaultClient(keyId string, token string) (registry.KMSClient, error) {
	cfg := api.DefaultConfig()
	vurl, err := url.Parse(keyId)
	if err != nil {
		return nil, err
	}

	schema := "https://"
	insecure, _ := strconv.ParseBool(vurl.Query().Get("insecure"))
	if insecure {
		schema = "http://"
	}
	cfg.Address = schema + vurl.Host

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if token == "" {
		token = InferVaultToken()
	}
	client.SetToken(token)

	return hcvault.NewClientWithAEADOptions(keyId, client.Logical(), hcvault.WithLegacyContextParamater())
}
