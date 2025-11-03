# Changesets Management for Makeathing MUD

This document explains how to use [Changesets](https://github.com/changesets/changesets) to manage the changelog and versioning for the Makeathing MUD project.

## Overview

Changesets provides a way to manage versioning and changelogs in a systematic way. Even though this is a Go project, we've integrated Changesets by adding a minimal `package.json` file that allows us to use the Changesets CLI tool.

## Prerequisites

Ensure you have Node.js and npm installed:

```bash
node --version
npm --version
```

## Workflow

### 1. Adding a New Changeset

When you make changes that should be included in the next release, create a changeset:

```bash
npx changeset
```

You'll be prompted to:
1. Select which package(s) are affected (in our case, just "makeathing-mud")
2. Choose the type of version bump:
   - `patch` (0.0.1) - Bug fixes, small changes
   - `minor` (0.1.0) - New features, backwards compatible
   - `major` (1.0.0) - Breaking changes
3. Write a summary of the changes

This creates a new markdown file in the `.changeset` directory with a unique name.

### 2. Updating the Changelog and Versions

When you're ready to release a new version:

```bash
npx changeset version
```

This command will:
- Consume all changesets in the `.changeset` directory
- Update the version in `package.json`
- Generate/update `CHANGELOG.md` with the changes
- Remove the consumed changeset files

### 3. Update the VERSION File

Since this is a Go project that uses a VERSION file, you need to manually update it to match the version in `package.json`:

```bash
# Check the new version
grep '"version"' package.json

# Update VERSION file (replace 0.1.2 with the actual version)
echo "0.1.2" > VERSION
```

### 4. Commit and Tag the Release

```bash
git add .
git commit -m "chore: prepare release v0.1.2"
git tag v0.1.2  # Use the version from package.json/VERSION
```

### 5. Push the Changes

```bash
git push origin main
git push origin v0.1.2
```

The GitHub Actions workflow will automatically create a release when a new tag is pushed.

## Directory Structure

- `.changeset/` - Contains changeset files and configuration
- `CHANGELOG.md` - Generated changelog file
- `package.json` - Minimal file for Changesets compatibility
- `VERSION` - Go project version file (must be manually updated)

## Best Practices

1. **Create changesets early**: Create a changeset right after making significant changes
2. **Be descriptive**: Write clear, concise summaries of changes
3. **Categorize correctly**: Choose the appropriate version bump type
4. **Keep files in sync**: Remember to update VERSION after running `npx changeset version`
5. **Review before committing**: Check CHANGELOG.md and package.json versions before committing

## Example Workflow

```bash
# After making changes
npx changeset
# Select "minor" for new feature
# Write "Add new combat system with advanced stats"

# When ready to release
npx changeset version
grep '"version"' package.json  # Check new version (e.g., 0.2.0)
echo "0.2.0" > VERSION  # Update VERSION file

# Commit and tag
git add .
git commit -m "chore: prepare release v0.2.0"
git tag v0.2.0
git push origin main --tags
```

This workflow ensures that all changes are properly documented and versioned, making it easier for users and contributors to understand the evolution of the project.