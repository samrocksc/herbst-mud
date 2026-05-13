Every feature must have a test.

- Use standard go testing libraries for unit tests.
- For integration tests use recommended framework from cucumber for BDD testing.
- Only write BDD features for SSH events that can be authenticated with SSH
- Each of the features should have a corresponding feature file in the features directory.
- Check for "Requires BDD Test" in the feature file's frontmatter
