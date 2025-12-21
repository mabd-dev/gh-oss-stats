# GitHub Token Setup Guide

This guide shows you how to create and configure a GitHub Personal Access Token (PAT) for `gh-oss-stats`.

## Why Do I Need a Token?

GitHub API rate limits:
- **Without token**: 60 requests/hour (very limited)
- **With token**: 5,000 requests/hour (recommended)

For any non-trivial usage, you **need a token**.

## Quick Setup (Recommended)

### Step 1: Create a GitHub Token

1. Go to https://github.com/settings/tokens
2. Click **"Generate new token"** → **"Generate new token (classic)"**
3. Give it a name: `gh-oss-stats`
4. **No scopes needed** - leave all checkboxes unchecked (read-only access is enough)
5. Click **"Generate token"**
6. **Copy the token** (starts with `ghp_...`)

### Step 2: Set Environment Variable

Add to your shell config file:

**For Bash** (`~/.bashrc` or `~/.bash_profile`):
```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

**For Zsh** (`~/.zshrc`):
```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

**For Fish** (`~/.config/fish/config.fish`):
```fish
set -x GITHUB_TOKEN ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Step 3: Reload Your Shell

```bash
# Reload config
source ~/.bashrc  # or ~/.zshrc

# Verify it's set
echo $GITHUB_TOKEN
```

### Step 4: Done!

Now just run:
```bash
gh-oss-stats --user YOUR_USERNAME
```

The tool will automatically use `$GITHUB_TOKEN`.

---

## Security Best Practices

### ✅ DO

- Use **classic tokens** with **no scopes** (read-only is enough)
- Store in environment variable or CI/CD secrets
- Rotate tokens periodically
- Use different tokens for different tools
- Set expiration dates on tokens

### ❌ DON'T

- Commit tokens to git repositories
- Share tokens publicly
- Use tokens with unnecessary scopes
- Store tokens in plain text files with wrong permissions
- Use the same token everywhere

---

## Token Scopes Explained

For `gh-oss-stats`, you need **NO scopes** (public read-only access).

If you accidentally added scopes:
- `repo` - Full repo access (not needed, too permissive)
- `public_repo` - Public repo write access (not needed)
- `read:user` - Read user data (not needed)

