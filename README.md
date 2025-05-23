# About

Ferleasy is a tool to integrate less technical users in a gitops workflow using ferlease: https://github.com/Ferlab-Ste-Justine/ferlease

It provides 3 commands for end-users to manage fearlease releases in a end-user configuration state:
- **add**: Command to add a release
- **list**: Command to list the releases
- **remove**: Command to remove a release

From there, backend devs will create the fearlease templates ferleasy will use to create releases in the git code and create a cron job that will run the **sync** command at a reasonable frequency.

The **sync** command will look at its internal state where it stores the ferleases that were applied, compare it with the user configuration and remove or add ferlease releases in the git code as needed.

# Fields

The client allow users to specify 3 fixed fields for each release: **environment**, **service** and **release**. Those 3 fields together form the unique id of the release. Not all 3 fields are needed if you won't use them for templating, but at least one of those fields is needed to separate release.

Additionally, **custom_params** (just **params** in the client) allows you to specify other relevant key/value pairs for templating purposes (but are otherwise not used to identify releases in the tool's state).

Custom parameters can be added when creating a release in the command line using the **params** parameter like so:

```
ferleasy add ... --params="<name>=value" ...
```

Note that a end-user can modify the same release in-place by changing the custom parameters which will cause ferleasy to redo the ferlease release. Be careful with this behavior if you modify the template files or use custom parameters in the **filesystem-conventions.yml** part of the template (we recommend you don't do the later with custom parameters that are not fixed), see limitations below.

# Stores

ferleasy supports 3 diffent store types currently, but for the end-user specified releases and for the synchronization state in the backend: **s3**, **etcd** and **git**. You may use different storage solutions for the end-user specified releases and the synchronization state.

**etcd** has the benefit of supporting locking and is the best solution to use for the backend synchronization state if there is any expected concurrency. **git** could technically also support a lock (at the cost of some commit noise), but we opted not to do so due for now due to time constraints.

Note that a **git** store would be an elegant solution to use for the end-user part of the storage (which should be pretty low concurrency anyways) if you want to use gitops end-to-end. If the releases specification is in a git repo, it also makes it easier for backend developers to provide support by making edits as needed directly in the repo.

# Policies

Ferleasy supports both default values (for convenience) and fixed values (for security if the same templates are reused) for all the fields.

For custom parameters, the default and fixed values will apply only for the keys that are specified.

Default values are for the client only and will be used for fields the user omit.

Fixed values are both for the client and the backend. The client commands will abort in error if a field with a fixed value has a different value. And should the user try to bypass this by changing their configuration file, then assuming you specified the same fixed values in the configurations for the sync command in the backend job, then the job will abort in error without applying the release.

# Configuration

All ferleasy commands expect a configuration file in yaml which can be specified with the **config** command line argument and defaults to **config.yml** in the current directory.

Below are all the entries in the configuration:

- **releases_store**: Store configuration for user releases. See below for the entries.
- **entry_policy**: Configuration for the entry policies.
  - **default**: Default values if the end-user omits those fields when creating a release. Only the entries with default values need to be specified here.
    - **environment**: Default value for the environment.
    - **service**: Default value for the service.
    - **release**: Default value for the release.
    - **custom_parameters**: Default value for some custom parameters. They can be enumerated as a key/value pairs nested under this field. See usage example if still unsure.
  - **fixed**: Fixed values that will be enforced both for end-user commands and the sync commands for the fields that have them.
    - **environment**: Fixed value for the environment.
    - **service**: Fixed value for the service.
    - **release**: Fixed value for the release.
    - **custom_parameters**: Fixed value for some custom parameters. They can be enumerated as a key/value pairs nested under this field. See usage example if still unsure.
- **sync**: Configuration for the backend sync command only. Not needed for the client-side commands.
  - **author**: Git author information for the commits. Can optionally be omitted if that is setup on the system where the tool runs and match the ssh credentials provided.
    - **name**: Name of the author
    - **email**: Email of the author
  - **commit_message**: Default commit message for all orchestrations. This property can be templatized.
  - **push_retries**: In the unlikely even that a barrage of upstream commits keep blocking a gitops operation, how many times to retry before giving up.
  - **push_retry_interval**: : If a gitop operation is blocked by an upstream commit, how long to wait before re-cloning, re-commiting and re-attempting the push. Should be a string in golang duration format.
  - **orchestrations**: Specification for ferlease templates. We invite you to look at the **orchestrations** part of ferlease configuration here for details: https://github.com/Ferlab-Ste-Justine/ferlease#configuration-file
  - **state_store**: Store for the backend state when using the sync command. See below for the entries.

The store configuration entries are:
- **etcd**: Configuration for an etcd store
  - **prefix**: Etcd key prefix to store under
  - **endpoints**: List of endpoints for the etcd cluster. Each entry should have an `<ip>:<port>` format.
  - **connection_timeout**: Connection timeout in a golang duration string format.
  - **request_timeout**: Request timeout in a golang duration string format.
  - **retry_interval**: Retry interval for failed requests in a golang duration string format.
  - **retries**: Number of time to retry a failed request
  - **lock_ttl**: Lock ttl before it expires, in second. Useful to free lock if ferleasy crashes or is forcefully stopped after obtaining it without freeing it.
  - **auth**: Authentification to the etcd cluster.
    - **ca_cert**: Path to a CA cert file to authentify the etcd servers.
    - **client_cert**: Path to the client's cert file if certificate authentication is used.
    - **client_key**: Path to the client's private key if certificate authentication is used.
    - **password_auth**: Path to a yaml file containing the user credentials if password authentication is used. It should have two keys: **username** and **password**.
- **git**: Configuration for a git store
  - **url**: Url of the git repo
  - **ref**: Branch to commit to
  - **Path**: Path to store under in the git repo
  - **accepted_signatures**: Optional directory containing the public part of trusted gpg signatures. If the top commit in the desired repo's branch is not signed by one of those keys, ferleasy will abort in error.
  - **commit_signature**: Optional signature key to sign commits with
    - **key**: Path to a file containing the signing key
    - **Passphrase**: Path to a file containing the passphrase to decrypt the signing key if it is encrypted with a passphrase.
  - **commit_message**: Commit message to use when commiting content changes.
  - **push_retries**: In the unlikely even that a barrage of upstream commits keep blocking a push, how many times to retry before giving up.
  - **push_retry_interval**: : If a push is blocked by an upstream commit, how long to wait before re-cloning, re-commiting and re-attempting the push. Should be a string in golang duration format.
  - **auth**: Authentication to the git server.
    - **ssh_key**: Ssh key to use to authentify with the git server. It should be the path to a file containing the key, not the key itself.
    - **known_key**: Path to a file containing the git server's ssh fingerprint. Used to authentify the server.
    - **user**: User to user to identify as with the git server. Can be left empty for many git providers, but some like Gitea require it.
- **s3**: Configuration for an s3 store
  - **endpoint**: Endpoint of the s3 store. It should have an `<ip|domain>:<port>` format.
  - **bucket**: s3 bucket to store in
  - **path**: Path prefix in the s3 bucket to store under
  - **region**: Region to use if your s3 store is region-specific
  - **connection_timeout**: Timeout for connection attempts in goland duration string format.
  - **request_timeout**: Timeout for request attempts in golang duration string format.
  - **auth**: Authentication to the s3 store.
    - **ca_cert**: Optional path to a CA cert file to authentify the s3 servers.
    - **key_auth**: Path to the client's key credentials for the s3 store. It should be in yaml and contain 2 keys: **access_key** and **secret_key**.
# Limitations

## Field names

Ferlease was initially used to deploy services, hence the name of the fields (ie, environment, service and release).

Ferlease can be used to potentially orchestrate anything by thinking of the term **service** as a descriptor for what you are templating (ex: **user** if you want to orchestrate users) and **release** as a unique id for what you are trying to add/remove (ex: id of a specific user for users).

For example, if we follow the users example, if you wanted to create a user with id **marc01** and email (additional templating parameter) **marc@email.com** for the production environment, you might type the following to add a release:

```
ferleasy add --environment=prod --service=user --release=marc01 --params="email=marc@email.com"

```

We'll try it out as it is and if the terminology is too confusing for our end users, we'll make the field names presented by the command line client configurable to more intuitively describe the terminology of specific application cases for end users.

## Template Modifications

Currently, support for changing backend templates for existing releases is poor, mostly due to time constraints. 

If you change the backend template, existing releases won't be impacted and new releases will. You could technically for the new templates by adding/modifying a bogus custom parameter on each existing release.

However, you should be very careful about that, because when ferlease re-applies templates, it works well to modify existing files in place or add new files, but it is liable to leave orphan files behind if template files are removed or the path where the template is applied changes. There is some work that can be done to improve this in the repo directories that ferlease creates and manages (ie, terraform modules or fluxcd app directories ferlease creates), but changes to the **filesystem-conventions.yml** part of the template (ie, where the template is applied in the repo) is a non-trivial problem as ferlease (the tool ferleasy manages essentially) is a stateless tool.

The medium-term solution we envision is to improve ferlease support for cleaning up the directories it explicitly manages (ie: terraform modules and fluxcd app directories), have ferleasy store information about the templates it applies, have ferleasy re-apply the release when those template changes, but also abort with an error if the **filesystem-conventions.yml** part of the template is changed.

## Logs

We noticed that storage solutions, go-git especially, may be a little noisy with their logging on the prompt which may be a little annoying for end users. Depending on the feedback we get, we may investigate what we can do to silence that logging.