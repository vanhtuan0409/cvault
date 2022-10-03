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

			if err := registerTinkKey(keyId, cmd); err != nil {
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

	viper.BindPFlag("keyId", rootCmd.PersistentFlags().Lookup("key-id"))
	viper.BindPFlag("store", rootCmd.PersistentFlags().Lookup("store"))

	AddEncryptCommand(rootCmd)
	AddDecryptCommand(rootCmd)
	AddListCommand(rootCmd)
	AddRemoveCommand(rootCmd)
	AddPeekCommand(rootCmd)

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

func registerTinkKey(keyId string, cmd *cobra.Command) (err error) {
	var client registry.KMSClient

	switch {
	case strings.HasPrefix(keyId, cvault.TinkAwsKms):
		client, err = awskms.NewClient(keyId)
	case strings.HasPrefix(keyId, cvault.TinkGcpKms):
		client, err = gcpkms.NewClient(keyId)
	case strings.HasPrefix(keyId, cvault.TinkHcVault):
		token := cmd.Flag("vault-token").Value.String()
		client, err = hcvault.NewClient(keyId, nil, token)
	default:
		return
	}
	if err != nil {
		return
	}

	registry.RegisterKMSClient(client)
	return nil
}
