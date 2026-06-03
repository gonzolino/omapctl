package cmd

import (
	"fmt"

	"github.com/gonzolino/omapctl/internal/ceph"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <pool> <oid> <key>...",
		Short: "Delete one or more omap keys from an object",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, oid := args[0], args[1]
			keys := args[2:]

			conn := connFromContext(cmd.Context())
			ioctx, err := ceph.OpenIOContext(conn, pool)
			if err != nil {
				return err
			}
			defer ioctx.Destroy()

			if err := ioctx.RmOmapKeys(oid, keys); err != nil {
				return fmt.Errorf("delete omap keys on %q in pool %q: %w", oid, pool, err)
			}
			return nil
		},
	}
}
