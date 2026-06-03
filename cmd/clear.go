package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gonzolino/omapctl/internal/ceph"
	"github.com/spf13/cobra"
)

func newClearCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "clear <pool> <oid>",
		Short: "Remove all omap entries from an object",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, oid := args[0], args[1]

			if !force {
				fmt.Fprintf(os.Stderr, "This will delete ALL omap entries for object %q in pool %q.\nType the object ID to confirm: ", oid, pool)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				input := strings.TrimSpace(scanner.Text())
				if input != oid {
					return fmt.Errorf("confirmation did not match object ID %q — aborted", oid)
				}
			}

			conn := connFromContext(cmd.Context())
			ioctx, err := ceph.OpenIOContext(conn, pool)
			if err != nil {
				return err
			}
			defer ioctx.Destroy()

			if err := ioctx.CleanOmap(oid); err != nil {
				return fmt.Errorf("clear omap on %q in pool %q: %w", oid, pool, err)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")

	return cmd
}
