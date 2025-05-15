package cmd

import (
	"github.com/Ferlab-Ste-Justine/ferleasy/config"
	"github.com/Ferlab-Ste-Justine/ferleasy/state"
	"github.com/Ferlab-Ste-Justine/ferleasy/store"

	"github.com/spf13/cobra"
)

func generateAddCmd(confPath *string) *cobra.Command {
	var environment string
	var service string
	var release string
	var customParams map[string]string

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "End-user client command to add a release",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(*confPath)
			AbortOnErr(err)

			relStore, relStoreErr := store.GetStore[state.Entries](&conf.ReleasesStore, "releases.yml", "releases.lock")
			AbortOnErr(relStoreErr)

			newEntry := state.Entry{
				Environment: environment,
				Service: service,
				Release: release,
				CustomParams: customParams,
			}

			policyErr := newEntry.ApplyPolicy(&conf.EntryPolicy)
			AbortOnErr(policyErr)

			opErr := store.ProcessStoreContent[state.Entries](func(entries state.Entries) (state.Entries, error) {
				entries.Add(newEntry)
				return entries, nil
			}, relStore)
			AbortOnErr(opErr)
		},
	}

	addCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environment of the ferlease to add")
	addCmd.Flags().StringVarP(&service, "service", "s", "", "Service of the ferlease to add")
	addCmd.Flags().StringVarP(&release, "release", "r", "", "Release of the ferlease to add")
	addCmd.Flags().StringToStringVarP(&customParams, "params", "p", map[string]string{}, "Custom parameters to add to the ferlease")

	return addCmd
}