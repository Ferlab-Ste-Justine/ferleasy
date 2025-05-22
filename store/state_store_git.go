package store

import (
	"errors"
	"fmt"
	"path"

	gitutils "github.com/Ferlab-Ste-Justine/ferleasy/git"
	git "github.com/Ferlab-Ste-Justine/git-sdk"
	yaml "gopkg.in/yaml.v2"
)

type StoreGit[T any] struct {
	Key      string
	Config   gitutils.GitConfig
	sshCreds *git.SshCredentials
}

func (store *StoreGit[T]) Initialize() error {
	sshCreds, sshCredsErr := git.GetSshCredentials(
		gitutils.ExpanHomeDir(store.Config.Auth.SshKey),
		gitutils.ExpanHomeDir(store.Config.Auth.KnownKey),
		store.Config.Auth.User,
	)
	if sshCredsErr != nil {
		return errors.New(fmt.Sprintf("Error getting ssh credentials to authentify with git store: %s", sshCredsErr.Error()))
	}

	store.sshCreds = sshCreds

	return nil
}

func (store *StoreGit[T]) ReadContent() (T, error) {
	var content T

	repo, mem, cloneErr := git.MemCloneGitRepo(store.Config.Url, store.Config.Ref, 1, store.sshCreds)
	if cloneErr != nil {
		return content, errors.New(fmt.Sprintf("Error cloning git store repo: %s", cloneErr.Error()))
	}

	if store.Config.AcceptedSignatures != "" {
		verifyErr := gitutils.VerifyRepoSignatures(repo, store.Config.AcceptedSignatures)
		if verifyErr != nil {
			return content, errors.New(fmt.Sprintf("Error verifying git store repo signatures: %s", verifyErr.Error()))
		}
	}

	exists, existsErr := mem.FileExists(path.Join(store.Config.Path, store.Key))

	if existsErr != nil {
		return content, errors.New(fmt.Sprintf("Error checking git store content existence: %s", existsErr.Error()))
	}

	if !exists {
		return content, nil
	}

	rawContent, getErr := mem.GetFileContent(path.Join(store.Config.Path, store.Key))
	if getErr != nil {
		return content, errors.New(fmt.Sprintf("Error reading git store content: %s", getErr.Error()))
	}

	desErr := yaml.Unmarshal([]byte(rawContent), &content)
	if desErr != nil {
		return content, errors.New(fmt.Sprintf("Error deserializing git store content: %s", desErr.Error()))
	}

	return content, nil
}

func (store *StoreGit[T]) WriteContent(content T) error {
	pushErr := git.PushChanges(func() (*git.GitRepository, error) {
		output, serErr := yaml.Marshal(&content)
		if serErr != nil {
			return nil, errors.New(fmt.Sprintf("Error serializing git store content: %s", serErr.Error()))
		}
		
		repo, mem, cloneErr := git.MemCloneGitRepo(store.Config.Url, store.Config.Ref, 1, store.sshCreds)
		if cloneErr != nil {
			return nil, errors.New(fmt.Sprintf("Error cloning git store repo: %s", cloneErr.Error()))
		}

		if store.Config.AcceptedSignatures != "" {
			verifyErr := gitutils.VerifyRepoSignatures(repo, store.Config.AcceptedSignatures)
			if verifyErr != nil {
				return nil, errors.New(fmt.Sprintf("Error verifying git store repo signatures: %s", verifyErr.Error()))
			}
		}

		setErr := mem.SetFileContent(path.Join(store.Config.Path, store.Key), string(output))
		if setErr != nil {
			return nil, errors.New(fmt.Sprintf("Error writing git store content: %s", setErr.Error()))
		}

		var signature *git.CommitSignatureKey
		var signatureErr error
		if store.Config.CommitSignature.Key != "" {
			signature, signatureErr = git.GetSignatureKey(store.Config.CommitSignature.Key, store.Config.CommitSignature.Passphrase)
			if signatureErr != nil {
				return nil, errors.New(fmt.Sprintf("Error getting signature key to commit to git store: %s", signatureErr.Error()))
			}
		}

		changes, comErr := git.CommitFiles(
			repo, 
			[]string{path.Join(store.Config.Path, store.Key)}, 
			store.Config.CommitMessage,
			git.CommitOptions{
				Name: store.Config.Author.Name,
				Email: store.Config.Author.Email,
				SignatureKey: signature,
			},
		)
		if comErr != nil {
			return nil, errors.New(fmt.Sprintf("Error commiting change to git store: %s", comErr.Error()))
		}

		if !changes {
			return nil, nil
		}

		return repo, nil
	}, store.Config.Ref, store.sshCreds, store.Config.PushRetries, store.Config.PushRetryInterval)

	return pushErr
}

func (st *StoreGit[T]) Lock() error {
	return nil
}

func (st *StoreGit[T]) ReleaseLock() error {
	return nil
}

func (st *StoreGit[T]) Cleanup() error {
	return nil
}