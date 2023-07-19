package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/awskms"
	"github.com/google/tink/go/integration/gcpkms"
	"github.com/google/tink/go/integration/hcvault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault"
	"github.com/vanhtuan0409/cvault/aead/aesgcm"
)

func main() {
	cobra.OnInitialize(initConfig)
	rootCmd := &cobra.Command{
		Use:          "cvault",
		Short:        "A cold vault encryption tool",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			keyId := viper.GetString("keyId")
			if keyId == "" {
				return errors.New("invalid key id")
			}

			storeUrl := viper.GetString("store")
			if storeUrl == "" {
				return errors.New("invalid store url")
			}

			if err := registerTinkKey(keyId); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	rootCmd.PersistentFlags().StringP("key-id", "k", "", "KMS key id")
	rootCmd.PersistentFlags().StringP("store", "s", "local://.", "Location of storage")
	rootCmd.PersistentFlags().String("vault-token", "", "HC vault token")
	rootCmd.PersistentFlags().String("pass-prompt", "", "Prompt for AES passphrase. Only when use key aesgcm://")

	viper.BindPFlag("vaultToken", rootCmd.PersistentFlags().Lookup("vault-token"))
	viper.BindPFlag("keyId", rootCmd.PersistentFlags().Lookup("key-id"))
	viper.BindPFlag("store", rootCmd.PersistentFlags().Lookup("store"))
	viper.BindPFlag("passPrompt", rootCmd.PersistentFlags().Lookup("pass-prompt"))

	AddEncryptCommand(rootCmd)
	AddDecryptCommand(rootCmd)
	AddListCommand(rootCmd)
	AddRemoveCommand(rootCmd)
	AddPeekCommand(rootCmd)
	AddEditCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	cvaultCfgDir := filepath.Join(cfgDir, "cvault")
	viper.AddConfigPath(".")
	viper.AddConfigPath(cvaultCfgDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
}

func registerTinkKey(keyId string) (err error) {
	var client registry.KMSClient

	switch {
	case strings.HasPrefix(keyId, cvault.TinkAwsKms):
		client, err = awskms.NewClient(keyId)
	case strings.HasPrefix(keyId, cvault.TinkGcpKms):
		client, err = gcpkms.NewClient(keyId)
	case strings.HasPrefix(keyId, cvault.TinkHcVault):
		token := viper.GetString("vaultToken")
		if token == "" {
			token = cvault.InferVaultToken()
		}
		client, err = hcvault.NewClient(keyId, nil, token)
	case strings.HasPrefix(keyId, cvault.AesGcm):
		promptScript := viper.GetString("passPrompt")
		client, err = aesgcm.NewClient(keyId, promptScript)
	default:
		return
	}
	if err != nil {
		return
	}

	registry.RegisterKMSClient(client)
	return nil
}
