# How does it work

Each commit is being added a Change-Id header by a commit message hook.
This Change-Id is used to track the rebase of a given commit.

For each Change-Id in your branch, a remote tracking branch is created.

When you run git-review, your branch is rebased on the remote tracking branch,
removing merged Change-Id from your local branch
And force-push all Change-Ids to your repository

Consider you are working on a repository:

```
                 A---B---C topic
                /
           D---E---F---G master origin/master
```

By the time you are working, origin/master has been updated and you end up with the following configuration:


```
                 A---B---C topic
                /
           D---E---F---G master
                        \
                         H---I origin/master
```

When you run `git review`, the tool will rebase topic on origin/master, ensuring it is up to date:

```
           D---E---F---G master
                        \
                         H---I origin/master
                              \
                               A---B---C topic
```

And create 3 pull requests for each of those commits:

```
           D---E---F---G master
                        \
                         H---I origin/master
                              \
                               A PR #1
                                \
                                 B PR #2
                                  \
                                   C PR #3 topic
```

After receiving reviews on your pull request, you will update your branch with the necessary fixes,
and origin/master will again be updated.

**NOTE:** to make a fixup on a given commit A, use the following git command: `git commit --fixup=<SHA>`,
where SHA is the sha of the commit you're fixing.
This way, *maiao* will be able to recognise which commit goes to which PR.


```
           D---E---F---G master
                        \
                         H---I---K---L origin/master
                              \
                               A PR #1
                                \
                                 B PR #2
                                  \
                                   C PR #3
                                    \
                                     fixup(A)---fixup(B)---fixup(C) topic
```

The next run of `git review` will re-order commits to ammend PR#1 with the fixups you have done for A,
PR #2 with the fixups you have done for B, and so on.

Then push it to update the pull requests

```
           D---E---F---G master
                        \
                         H---I---K---L origin/master
                                      \
                                       A---fixup(A) PR #1
                                                   \
                                                    B---fixup(B) PR #2
                                                                \
                                                                 C---fixup(C) PR #3 topic
```

When PR #1 is accepted, another `git review` will detect it and rebase your local branch:

```
           D---E---F---G master
                        \
                         H---I---K---L---A  origin/master
                                          \
                                           B---fixup(B) PR #2
                                                       \
                                                        C---fixup(C) PR #3 topic
```
