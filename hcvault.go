package cvault

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/cliconfig"
	"github.com/tink-crypto/tink-go-hcvault/v2/integration/hcvault"
	"github.com/tink-crypto/tink-go/v2/core/registry"
)

func loadVaultToken() string {
	helper, err := cliconfig.DefaultTokenHelper()
	if err != nil {
		return ""
	}
	token, _ := helper.Get()
	return token
}

func setVaultAddress(cfg *api.Config, keyId string) error {
	vurl, err := url.Parse(keyId)
	if err != nil {
		return err
	}

	schema := "https://"
	insecure, _ := strconv.ParseBool(vurl.Query().Get("insecure"))
	if insecure {
		schema = "http://"
	}

	cfg.Address = fmt.Sprintf("%s%s", schema, vurl.Host)
	return nil
}

func VaultClient(keyId string) (registry.KMSClient, error) {
	cfg := api.DefaultConfig()
	if err := setVaultAddress(cfg, keyId); err != nil {
		return nil, err
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	client.SetToken(loadVaultToken())
	return hcvault.NewClientWithAEADOptions(keyId, client.Logical(), hcvault.WithLegacyContextParamater())
}
