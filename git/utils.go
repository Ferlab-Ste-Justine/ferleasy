package git

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

    git "github.com/Ferlab-Ste-Justine/git-sdk"
)

func ExpanHomeDir(fpath string) string {
	if strings.HasPrefix(fpath, "~/") {
		homeDir, homeDirErr := os.UserHomeDir()
		if homeDirErr == nil {
			return filepath.Join(homeDir, fpath[2:])
		}
	}

	return fpath
}

func VerifyRepoSignatures(repo *git.GitRepository, signaturesPath string) error {
	keys := []string{}
	err := filepath.Walk(signaturesPath, func(fPath string, fInfo fs.FileInfo, fErr error) error {
		if fErr != nil {
			return fErr
		}

		if fInfo.IsDir() {
			return nil
		}

		key, keyErr := os.ReadFile(fPath)
		if keyErr != nil {
			return errors.New(fmt.Sprintf("Error reading accepted signature: %s", keyErr.Error()))
		}

		keys = append(keys, string(key))

		return nil
	})

	if err != nil {
		return err
	}

	return git.VerifyTopCommit(repo, keys)
}