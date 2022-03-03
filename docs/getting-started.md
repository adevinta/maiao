# Getting Started

## The Maiao workflow

Maiao aims to take advantage of this classic situation:

![](img/code_reviews_tweet.jpeg)

Maiao encourages you to create small, self-contained pull requests. You commit each change sequentially, and run
the `git review` command. It will create a PR for each commit, each PR depending on the previous PR.

This will make PR reviews easier, quicker, and more meaningful.

### Example

Let's say you're working on a new feature. This feature requires changes to 2 files, and 2 corresponding test files.
With Maiao, what you would do is:

1. Identify which of the files should be modified first.
2. Apply the change.
3. Write the corresponding tests, fix broken tests if needed.
4. Add the changed files to git with `git add`.
5. Commit these changes with a meaningful message with `git commit`.
6. Apply the changes to the second file and corresponding tests.
4. Add the changed files to git with `git add`.
7. Commit these changes with a meaningful message with `git commit`.
8. Run `git review`.

Now you'll have 2 PRs published in GitHub, and your colleagues will be delighted as to how elegantly you've structured
them. Each PR is a self-contained change, and all tests pass on each step.

### Advanced usage

Let's say you receive an insightful comment on the **second** PR, pointing how confusing the name of the `asdf` variable
is. You then:

1. Change the name of the variable
2. Add the changes with `git add`
3. Commit those changes as a fixup to the last commit with `git commit --fixup HEAD`.
4. Run `git review`.

Now you see that change in your pull request, and your colleague is happy.

A second colleague points out that the function you defined in the **first commit** is inefficient. Then you:

1. Fix the function
2. Add the changes with `git add`
3. Run `git log`
4. Copy the hash of the corresponding **first commit**. (For example, this corresponds to the hexadecimal value
    `cf5fd9a0baf2fd899ed0ef45629dc1f1c1b7af87`).
5. Commit those changes as a fixup to that specific commit with `git commit --fixup <commit-hash>`.
6. Run `git review`.

Your change is now published to the corresponding PR. You and your colleagues are satisfied.

## Installation

### Unix users

The simplest way to install maiao is to download the binary from github. To do so, visit
the [releases](https://github.com/adevinta/maiao/releases) page and download the version you want to install. Then run

```
mv <downloadsDir>/git-review-`uname -s`-`uname -m` /usr/local/bin/git-review
chmod +x /usr/local/bin/git-review
```

!!! warning "MacOS users"
To make sure the binary is not quarantined run: `xattr -d com.apple.quarantine /usr/local/bin/git-review`

To build a standalone binary, you will need to run `go build -o /usr/local/bin/git-review ./cmd/maiao`

### Windows users

The relevant Windows binary can be downloaded in the [releases](https://github.com/adevinta/maiao/releases) page

Download the `git-review-windows-<arch>` where `<arch>` is the architecture of your machine. If you are unsure, usually,
amd64 is the default windows architecture.

## Configuration

If you have ever installed and configured `adv` command line, you should already be good to go. If not, you will need to
configure a github token in your netrc file. To do so, follow these steps:

Create your personal [access token](https://github.company.example.com/settings/tokens) with `repo` privileges and
create a file `~/.netrc` having the following content:

```
machine github.company.example.com
  login <firstname>.<lastname>@company.example.com
  password <your token>
```

Your environment is then ready to run maiao.

There will be some other configuration required for each repository you need to run maiao on. Don't worry, maiao will
prompt and suggest you to perform the setup if the configuration is missing.

