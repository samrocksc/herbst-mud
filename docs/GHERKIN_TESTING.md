# Gherkin Testing Guide

## Overview

This project uses Gherkin syntax for BDD (Behavior-Driven Development) testing. Feature files are stored in `features/` and follow the Given-When-Then format.

## Feature File Format

```gherkin
Feature: Room Navigation
  As a player
  I want to move between rooms
  So that I can explore the MUD world

  Scenario: Move north from starting room
    Given I am in "The Hole" starting room
    When I type "north"
    Then I should be in the "North Room"
    And I should see the room description
```

## Running Tests

```bash
# Run all tests
make test

# Run Go tests
cd server && go test ./...
cd herbst && go test ./...

# Run BDD tests (if Cucumber is set up)
cucumber features/
```

## Writing Tests

1. **Feature files** go in `features/` directory
2. **Step definitions** go in `features/step_definitions/` 
3. **Support files** go in `features/support/`

## Test Structure

- `features/*.md` - Feature specifications
- `features/step_definitions/*.go` - Go step implementations
- `features/support/*.go` - Test hooks and configuration

## Notes

- See `docs/TESTING.md` for more details
- Each feature in `features/` should have corresponding tests
- Mark features as `requires_bdd: true` in the frontmatter if BDD tests are needed