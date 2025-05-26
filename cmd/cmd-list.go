package cmd

import (
	"fmt"

	"github.com/Ferlab-Ste-Justine/ferleasy/config"
	"github.com/Ferlab-Ste-Justine/ferleasy/state"
	"github.com/Ferlab-Ste-Justine/ferleasy/store"

	"github.com/spf13/cobra"
)

func generateListCmd(confPath *string) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "End-user client command to list all releases",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(*confPath)
			AbortOnErr(err)

			relStore, relStoreErr := store.GetStore[state.Entries](&conf.ReleasesStore, "releases.yml", "releases.lock")
			AbortOnErr(relStoreErr)

			releases, releasesErr := relStore.ReadContent()
			AbortOnErr(releasesErr)

			idx := 0
			for _, release := range releases {
				if idx == 0 {
					fmt.Println("--------------------")
				}
				fmt.Printf("Environment: %s\nService: %s\nRelease: %s\n", release.Environment, release.Service, release.Release)
				for key, val := range release.CustomParams {
					fmt.Printf("\tParam %s: %s\n", key, val)
				}
				fmt.Println("--------------------")
				idx += 1
			}
		},
	}

	return listCmd
}