package cmd

import (
	"github.com/Ferlab-Ste-Justine/ferleasy/config"
	"github.com/Ferlab-Ste-Justine/ferleasy/state"
	"github.com/Ferlab-Ste-Justine/ferleasy/store"

	fercmd "github.com/Ferlab-Ste-Justine/ferlease/cmd"
	ferconf "github.com/Ferlab-Ste-Justine/ferlease/config"
	"github.com/spf13/cobra"
)

func generateSyncCmd(confPath *string) *cobra.Command {
	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Backend synchronization command to synchronize the ferleases with the desired client state",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.GetConfig(*confPath)
			AbortOnErr(err)

			relStore, relStoreErr := store.GetStore[state.Entries](&conf.ReleasesStore, "releases.yml", "releases.lock")
			AbortOnErr(relStoreErr)

			stateStore, stateStoreErr := store.GetStore[state.State](&conf.Sync.StateStore, "state.yml", "state.lock")
			AbortOnErr(stateStoreErr)

			entries, entriesErr := relStore.ReadContent()
			AbortOnErr(entriesErr)

			for _, entry := range entries {
				policyErr := entry.CheckPolicy(&conf.EntryPolicy)
				AbortOnErr(policyErr)
			}

			opErr := store.ProcessStoreContent[state.State](func(state state.State) (state.State, error) {
				diff := state.Entries.Diff(&entries)

				for _, entry := range diff.Remove {
					ferConf := &ferconf.Config{
						Operation: "teardown",
						Environment: entry.Environment,
						Service: entry.Service,
						Release: entry.Release,
						CustomParams: entry.CustomParams,
						Author: conf.Sync.Author,
						CommitMessage: conf.Sync.CommitMessage,
						PushRetries: conf.Sync.PushRetries,
						PushRetryInterval: conf.Sync.PushRetryInterval,
						Orchestrations: conf.Sync.Orchestrations,
					}
					processErr := ferConf.Process()
					AbortOnErr(processErr)

					teardownErr := fercmd.Teardown(ferConf)
					AbortOnErr(teardownErr)

					state.Entries.Remove(entry)
				}

				for _, entry := range diff.Update {
					ferConf := &ferconf.Config{
						Operation: "release",
						Environment: entry.Environment,
						Service: entry.Service,
						Release: entry.Release,
						CustomParams: entry.CustomParams,
						Author: conf.Sync.Author,
						CommitMessage: conf.Sync.CommitMessage,
						PushRetries: conf.Sync.PushRetries,
						PushRetryInterval: conf.Sync.PushRetryInterval,
						Orchestrations: conf.Sync.Orchestrations,
					}
					processErr := ferConf.Process()
					AbortOnErr(processErr)

					releaseErr := fercmd.Release(ferConf)
					AbortOnErr(releaseErr)

					state.Entries.Add(entry)
				}

				for _, entry := range diff.Add {
					ferConf := &ferconf.Config{
						Operation: "release",
						Environment: entry.Environment,
						Service: entry.Service,
						Release: entry.Release,
						CustomParams: entry.CustomParams,
						Author: conf.Sync.Author,
						CommitMessage: conf.Sync.CommitMessage,
						PushRetries: conf.Sync.PushRetries,
						PushRetryInterval: conf.Sync.PushRetryInterval,
						Orchestrations: conf.Sync.Orchestrations,
					}
					processErr := ferConf.Process()
					AbortOnErr(processErr)

					releaseErr := fercmd.Release(ferConf)
					AbortOnErr(releaseErr)

					state.Entries.Add(entry)
				}

				return state, nil
			}, stateStore)
			AbortOnErr(opErr)
		},
	}

	return syncCmd
}