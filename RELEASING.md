# Releasing

This document describes the release process for herbst-mud.

## Versioning

All components (SSH server, backend, API) are versioned together as **one version**.

## Release Process

### 1. Prepare for Release

Ensure your local main branch is up to date:
```bash
git checkout main
git pull origin main
```

### 2. Create a Release Tag

Create and push a version tag following semantic versioning (vMAJOR.MINOR.PATCH):

```bash
# Example: v0.2.0
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

### 3. What Happens Automatically

When you push a version tag, the GitHub Actions workflow will:

1. **Run tests** - Execute all Go tests
2. **Build binaries** - Compile for Linux, macOS, and Windows (amd64 + arm64)
3. **Generate release notes** - Parse conventional commits
4. **Create GitHub Release** - Upload binaries and release notes
5. **Update CHANGELOG.md** - Automatically add the new release entry

## Commit Message Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, semicolons, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Test changes
- `build`: Build system or dependency changes
- `ci`: CI configuration changes
- `chore`: Other changes that don't modify src or test files
- `BREAKING CHANGE`: Incompatible changes

### Examples

```bash
# Feature
git commit -m "feat(auth): add OAuth2 login support"

# Bug fix
git commit -m "fix(room): resolve navigation issue in cross-shaped rooms"

# Breaking change
git commit -m "feat(api)!: change user endpoint response format"
```

## CHANGELOG Format

The CHANGELOG is automatically updated with entries for:

- **Added**: New features
- **Changed**: Changes to existing functionality
- **Fixed**: Bug fixes
- **Security**: Security-related changes

## Tag Protection

It is recommended to protect version tags in GitHub settings to prevent accidental deletion.

## Manual Release (if needed)

If you need to create a release manually:

1. Go to https://github.com/samrocksc/herbst-mud/releases
2. Click "Draft a new release"
3. Select a tag (or create new)
4. Fill in release notes
5. Upload binaries
6. Publish release