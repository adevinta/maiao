# Getting Started

## Installation

### Unix users

The simplest way to install maiao is to download the binary from github.
To do so, visit the [releases](https://github.com/adevinta/maiao/releases) page and download the version you want to install.
Then run

```
mv <downloadsDir>/git-review-`uname -s`-`uname -m` /usr/local/bin/git-review
chmod +x /usr/local/bin/git-review
```

!!! warning "MacOS users"
    To make sure the binary is not quarantined run: `xattr -d com.apple.quarantine /usr/local/bin/git-review`

To build a standalone binary, you will need to run `go build -o /usr/local/bin/git-review ./cmd/maiao` 

### Windows users

The relevant Windows binary can be downloaded in the [releases](https://github.com/adevinta/maiao/releases) page

Download the `git-review-windows-<arch>` where `<arch>` is the architecture of your machine.
If you are unsure, usually, amd64 is the default windows architecture.


## Configuration

If you have ever installed and configured `adv` command line, you should already be good to go.
If not, you will need to configure a github token in your netrc file. To do so, follow these steps:

Create your personal [access token](https://github.company.example.com/settings/tokens) with `repo` privileges
and create a file `~/.netrc` having the following content:

```
machine github.company.example.com
  login <firstname>.<lastname>@company.example.com
  password <your token>
```

Your environment is then ready to run maiao.

There will be some other configuration required for each repository you need to run maiao on.
Don't worry, maiao will prompt and suggest you to perform the setup if the configuration is missing.

