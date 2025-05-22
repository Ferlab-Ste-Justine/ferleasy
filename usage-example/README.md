# About

This is a repo to tryout ferleasy.

It uses a git store to manage both the releases and the state and requires some github configuration to setup.

# Requirements

To run this example, you need to have a github account and golang 1.23 or newer setup on your computer.

# Setup

Under your user, create a private repo called **ferleasy-playground**. This repo will be used by ferleasy to manage both its releases and its state.

Also, fork the following repo under your user: https://github.com/Ferlab-Ste-Justine/ferlease-playground

From there, edit the **config.yml** file and change `<YourGithubUser>` for your github user.

Also, if you prefer not to use your personal ssh key to run this example, feel to create a separate deploy key for both repos change change all ssh key references on the **config.yml** file to point to this key.

Here are instructions to setup a deploy key for repositories on github: https://docs.github.com/en/authentication/connecting-to-github-with-ssh/managing-deploy-keys

# Usage

Run the **refresh.sh** to compile ferleasy and place its binary in the example directory.

Run the following commands to simulate a end user adding some releases: 

```
./ferleasy add --release=test --params="Project=MyProject"
./ferleasy add --release=test2 --params="Project=MyProject"
```

Run the following command to simulate the end user listing the releases they have added:
```
./ferleasy list
```

Run the following command to simulate an ops backend recurring job applying the user added ferlease and syncing it into its state:
```
./ferleasy sync
```

Run the following command to simulate the end user removing one of their release:
```
./ferleasy remove --release=test
```

Run the following command to simulate the ops backend recurring job applying the user removed ferlease and syncing it into its state:
```
./ferleasy sync
```