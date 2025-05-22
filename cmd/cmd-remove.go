package cmd

import (
	"github.com/Ferlab-Ste-Justine/ferleasy/config"
	"github.com/Ferlab-Ste-Justine/ferleasy/state"
	"github.com/Ferlab-Ste-Justine/ferleasy/store"

	"github.com/spf13/cobra"
)

func generateRemoveCmd(confPath *string) *cobra.Command {
	var environment string
	var service string
	var release string

	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "End-user client command to remove a release",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(*confPath)
			AbortOnErr(err)

			removedEntry := state.Entry{
				Environment: environment,
				Service: service,
				Release: release,
			}
			removedEntry.CustomParams = map[string]string{}

			removedEntry.ApplyPolicyDefaults(&conf.EntryPolicy)

			relStore, relStoreErr := store.GetStore[state.Entries](&conf.ReleasesStore, "releases.yml", "releases.lock")
			AbortOnErr(relStoreErr)

			opErr := store.ProcessStoreContent[state.Entries](func(entries state.Entries) (state.Entries, error) {
				if entries == nil {
					entries = state.Entries(map[string]state.Entry{})
				}
				entries.Remove(removedEntry)
				return entries, nil
			}, relStore)
			AbortOnErr(opErr)
		},
	}

	removeCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environment of the ferlease to remove")
	removeCmd.Flags().StringVarP(&service, "service", "s", "", "Service of the ferlease to remove")
	removeCmd.Flags().StringVarP(&release, "release", "r", "", "Release of the ferlease to remove")

	return removeCmd
}