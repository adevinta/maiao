# Maiao

![Main branch build](https://github.com/github/adevinta/maiao/workflows/go.yml/badge.svg)
![License](https://img.shields.io/github/license/adevinta/maiao)
![GitHub all releases downloads](https://img.shields.io/github/downloads/adevinta/maiao/total)

**Gerrit-style code review workflow for GitHub**

Maiao brings the power of **stacked pull requests** to GitHub, enabling you to break large features into small, reviewable commits where each commit becomes its own PR.

## 🎯 What is Maiao?

Maiao provides the `git review` command that:

- **Creates one PR per commit** in your branch
- **Stacks PRs automatically** with proper parent-child dependencies
- **Manages fixups elegantly** using `git commit --fixup`
- **Tracks commits via Change-IDs** (using the Gerrit commit-msg hook)
- **Auto-rebases your stack** when PRs get merged

## 🚀 Quick Example

```bash
# Make multiple commits
git commit -m "Add user authentication"
git commit -m "Add authorization middleware"
git commit -m "Add admin endpoints"

# Create stacked PRs for all commits
git review
```

**Result**: Three GitHub PRs created:

- PR #1: `Add user authentication` → `main`
- PR #2: `Add authorization middleware` → PR #1
- PR #3: `Add admin endpoints` → PR #2

## ✨ Key Benefits

- **Granular Reviews**: Each commit reviewed independently for faster, focused feedback
- **Clear History**: One logical change per PR maintains clean git history
- **Easy Fixups**: Address review feedback with `git commit --fixup <sha>`
- **Automatic Stacking**: Tool manages PR dependencies automatically
- **Merge Detection**: Stack updates automatically when PRs merge
- **Rebase Integration**: Handles upstream changes gracefully

## 📚 Documentation

- **[Getting Started](getting-started.md)** - Installation and workflow guide
- **[How Does It Work](how-does-it-work.md)** - Technical details and architecture
- **[Pricing](pricing.md)** - Free and open source

## 🎓 Learn About Stacked Diffs

Maiao implements the **stacked diffs** methodology. Learn more:

- **[Stacked Diffs Versus Pull Requests](https://jg.gg/2018/09/29/stacked-diffs-versus-pull-requests/)** - Jackson Gabbard's foundational article on the philosophy
- **[Graphite's Guide to Stacked Diffs](https://graphite.dev/guides/stacked-diffs)** - Comprehensive guide to the workflow and best practices
- **[The Pragmatic Engineer: Stacked Diffs](https://newsletter.pragmaticengineer.com/p/stacked-diffs)** - Industry analysis and adoption patterns (paywalled)

## 🎪 Why "Maiao"?

As Maiao encourages users to create smaller and nicer commits in their pull requests, it has been given the name of a tiny island:

![](img/inspired.jpg)

## 🤝 Contributing

Contributions are welcome! See [CONTRIBUTING.md](../CONTRIBUTING.md) for details.

## 📜 License

MIT License - see [DISCLAIMER.md](../DISCLAIMER.md) for details.
