# Contributing to DWBDD-chat

First of all, thanks for your interest in contributing to DWBDD-chat! This is a college project but I hope that it'll be a big project one day! This guide will help you get started.

## AI Usage

It's totally fine to use AI (GPT, Gemini, Claude, etc.) to write code. However, any AI-generated code that appears unverified or poorly understood by the contributor will be rejected, even if it is functional.

## Getting started

1. Fork the repository and clone it:

```sh
git clone https://github.com/<your-username>/DWBDD-chat.git
cd DWBDD-chat
```

2. Create a branch for your work

```sh
git checkout -b <branch_name>
```

## Making changes

### Before you start

- Check if an [issue](https://github.com/ItsMeViipeR/DWBDD-chat/issues) or a [pull request](https://github.com/ItsMeViipeR/DWBDD-chat/pulls) is opened to avoid duplicating work.
- For larger changes, open an issue first to discuss the approach.

## Code style

- **Go**: Run `gofmt` or install `gopls` in your favorite editor before committing. The project uses Go v1.25.8.
- **Svelte**: Keep up the file tree as is. If a function can be reused in another page, add it to `src/lib`.

## Commit messages

Write clear, concise commit messages. Use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) prefixes:
- `feat:` --- for a new feature
- `fix:` --- bug fix
- `refactor` --- code restructuring without any behavior change

Example: `feat: add 2FA for user login`

If you have a long text for your commit use this Git syntax:

```sh
git commit -m "concise commit message" -m "Any additions to the commit message"`
```

## Submitting a Pull Request

1. Make sure your code compiles
2. Push your branch and open a Pull Request against `main`
3. Add a description to your Pull Request:
- Summarize what you did and why.
- Include screenshots or recordings for UI changes.
- Link related issues (i.e "Closes #14").
4. Keep PRs focused — avoid unrelated drive-by changes. If you spot something else that needs fixing, open a separate PR.
5. Don't force-push after review has started. Add new commits instead — they get squashed on merge.

### License

By contributing, you agree that your contributions will be licensed under the [MIT License](https://github.com/ItsMeViipeR/DWBDD-chat/blob/main/LICENSE)
