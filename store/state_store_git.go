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
		store.Config.Auth.SshKeyPath,
		store.Config.Auth.KnownHostsPath,
		store.Config.Auth.User,
	)
	if sshCredsErr != nil {
		return sshCredsErr
	}

	store.sshCreds = sshCreds

	return nil
}

func (store *StoreGit[T]) ReadContent() (T, error) {
	var content T

	_, mem, cloneErr := git.MemCloneGitRepo(store.Config.Url, store.Config.Ref, 1, store.sshCreds)
	if cloneErr != nil {
		return content, cloneErr
	}

	exists, existsErr := mem.FileExists(path.Join(store.Config.Path, store.Key))
	if existsErr != nil {
		return content, existsErr
	}

	if !exists {
		return content, nil
	}

	rawContent, getErr := mem.GetFileContent(path.Join(store.Config.Path, store.Key))
	if getErr != nil {
		return content, getErr
	}

	desErr := yaml.Unmarshal([]byte(rawContent), &content)
	if desErr != nil {
		return content, errors.New(fmt.Sprintf("Error deserializing content info: %s", desErr.Error()))
	}

	return content, nil
}

func (store *StoreGit[T]) WriteContent(content T) error {
	pushErr := git.PushChanges(func() (*git.GitRepository, error) {
		output, err := yaml.Marshal(&content)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error serializing the content info: %s", err.Error()))
		}
		
		repo, mem, cloneErr := git.MemCloneGitRepo(store.Config.Url, store.Config.Ref, 1, store.sshCreds)
		if cloneErr != nil {
			return nil, cloneErr
		}
	
		setErr := mem.SetFileContent(path.Join(store.Config.Path, store.Key), string(output))
		if setErr != nil {
			return nil, setErr
		}

		var signature *git.CommitSignatureKey
		var signatureErr error
		if store.Config.CommitSignature.Key != "" {
			signature, signatureErr = git.GetSignatureKey(store.Config.CommitSignature.Key, store.Config.CommitSignature.Passphrase)
			if signatureErr != nil {
				return nil, signatureErr
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
			return nil, comErr
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