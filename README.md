# About

Ferleasy is a tool to integrate less technical users in a gitops workflow using ferlease: https://github.com/Ferlab-Ste-Justine/ferlease

It provides 3 commands for end-users to manage fearlease releases in a end-user configuration state:
- **add**: Command to add a release
- **list**: Command to list the releases
- **remove**: Command to remove a release

From there, backend devs will create the fearlease templates ferleasy will use to create releases in the git code and create a cron job that will run the **sync** command at a reasonable frequency.

The **sync** command will look at its internal state where it stores the ferleases that were applied, compare it with the user configuration and remove or add ferlease releases in the git code as needed.

# Fields

TODO

# Stores

TODO

# Config

TODO

# Limitations

TODO