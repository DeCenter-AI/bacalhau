package node

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bacalhau-project/bacalhau/cmd/util"
	"github.com/bacalhau-project/bacalhau/pkg/bacerrors"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/apimodels"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/client/v2"
)

type NodeActionCmd struct {
	action  string
	message string
}

func NewActionCmd(action apimodels.NodeAction) *cobra.Command {
	actionCmd := &NodeActionCmd{
		action:  string(action),
		message: "",
	}

	cmd := &cobra.Command{
		Use:           fmt.Sprintf("%s [id]", action),
		Short:         action.Description(),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// initialize a new or open an existing repo merging any config file(s) it contains into cfg.
			cfg, err := util.SetupRepoConfig(cmd)
			if err != nil {
				return fmt.Errorf("failed to setup repo: %w", err)
			}
			// create an api client
			api, err := util.NewAPIClientManager(cmd, cfg).GetAuthenticatedAPIClient()
			if err != nil {
				return fmt.Errorf("failed to create api client: %w", err)
			}
			return actionCmd.run(cmd, args, api)
		},
	}

	cmd.Flags().StringVarP(&actionCmd.message, "message", "m", "", "Message to include with the action")
	return cmd
}

func (n *NodeActionCmd) run(cmd *cobra.Command, args []string, api client.API) error {
	ctx := cmd.Context()
	nodeID := args[0]

	response, err := api.Nodes().Put(ctx, &apimodels.PutNodeRequest{
		NodeID:  nodeID,
		Action:  n.action,
		Message: n.message,
	})
	if err != nil {
		return bacerrors.Wrapf(err, "failed to %s node %s", n.action, nodeID)
	}

	if response.Success {
		cmd.Println("Ok")
	} else {
		cmd.PrintErrf("Failed to %s node %s: %s\n", n.action, nodeID, response.Error)
	}

	return nil
}
