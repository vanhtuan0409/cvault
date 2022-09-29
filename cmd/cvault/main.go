package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cobra.OnInitialize(initConfig)
	rootCmd := &cobra.Command{
		Use:          "cvault",
		Short:        "A cold vault encryption tool",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	rootCmd.PersistentFlags().StringP("key-id", "k", "", "KMS key id")
	rootCmd.PersistentFlags().StringP("store", "s", "local://.", "Location of storage")

	viper.BindPFlag("keyId", rootCmd.PersistentFlags().Lookup("key-id"))
	viper.BindPFlag("store", rootCmd.PersistentFlags().Lookup("store"))

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	kmsClient := kms.NewFromConfig(cfg)
	AddEncryptCommand(kmsClient, rootCmd)
	AddDecryptCommand(kmsClient, rootCmd)
	AddListCommand(rootCmd)
	AddRemoveCommand(rootCmd)
	AddPeekCommand(kmsClient, rootCmd)

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
