# Changesets Management

This project uses [Changesets](https://github.com/changesets/changesets) to manage changelogs and versioning, even though it's a Go project rather than a JavaScript project.

## How to Update CHANGELOG.md

1. **Add a new changeset**:
   ```bash
   npx changeset
   ```
   This will prompt you to select the type of change (patch, minor, or major) and ask you to write a summary of the changes.

2. **Update versions and changelog**:
   ```bash
   npx changeset version
   ```
   This will update the `package.json` version and generate/update the `CHANGELOG.md` file.

3. **Update the VERSION file**:
   Since this is a Go project that uses a VERSION file for tracking, you'll also need to manually update the VERSION file to match the version in package.json:
   ```bash
   # Check the version in package.json
   grep '"version"' package.json
   
   # Update VERSION file to match
   echo "0.1.1" > VERSION  # Replace with actual version
   ```

4. **Commit the changes**:
   ```bash
   git add .
   git commit -m "chore: update changelog and version"
   ```

## How It Works

- Changesets are stored in the `.changeset` directory
- Each changeset is a markdown file that describes changes for a specific version bump
- The `npx changeset version` command consumes these changesets to update the changelog and version
- For this Go project, we maintain both `package.json` and `VERSION` files for version tracking

## Project Versioning

The project version is tracked in:
- `package.json` - For changesets compatibility
- `VERSION` - For Go project compatibility
- `CHANGELOG.md` - For documenting changes

These should always be kept in sync.

## Release Process

To create a new release:

1. Create a new changeset with `npx changeset`
2. Update versions with `npx changeset version`
3. Update the VERSION file to match
4. Commit and tag the release:
   ```bash
   git add .
   git commit -m "chore: prepare release vX.Y.Z"
   git tag vX.Y.Z  # Use the version from package.json/VERSION
   git push origin main --tags
   ```
5. GitHub Actions will automatically create a release when a new tag is pushed.