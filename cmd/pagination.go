package cmd

import (
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/spf13/cobra"
)

func addPaginationFlags(cmd *cobra.Command, defaultLimit int) {
	cmd.Flags().IntP(api.PaginationLimit, "", defaultLimit, "Number of items to fetch")
	cmd.Flags().IntP(api.PaginationOffset, "", 0, "Offset for fetch")
}
