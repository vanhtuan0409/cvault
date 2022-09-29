package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanhtuan0409/cvault/storage"
)

func AddListCommand(root *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List encrypted file in store",
		RunE: func(cmd *cobra.Command, args []string) error {
			storeUrl := viper.GetString("store")
			if storeUrl == "" {
				return errors.New("invalid store url")
			}
			ctx := cmd.Context()

			s, err := storage.GetStorage(storeUrl)
			if err != nil {
				return err
			}

			files, err := s.List(ctx)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 3, 4, 1, '\t', tabwriter.AlignRight)
			defer w.Flush()
			for _, entry := range files {
				timeStr := entry.LastModified.In(time.Local).Format(time.RFC3339)
				fmt.Fprintf(w, "%s\t%s\n", timeStr, entry.Key)
			}
			return nil
		},
	}

	root.AddCommand(listCmd)
}
