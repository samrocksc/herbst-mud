# BDD Integration Testing Framework

This directory contains Behavior-Driven Development (BDD) style integration tests for the MUD server.

## Test Structure

The tests follow the BDD Given-When-Then format:

```
GIVEN [initial context]
WHEN [event occurs]
THEN [ensure some outcomes]
```

## Implemented BDD Scenarios

### Authentication Flow

**Scenario 1: User authentication prompt**
```gherkin
GIVEN the server is started
WHEN I ssh into the server
THEN I am presented with the username and password challenge
```

**Scenario 2: Successful user authentication**
```gherkin
GIVEN the server is started
WHEN I ssh into the server
WHEN I enter the username `nelly` and password `password`
THEN I am in the starting room
```

## Test Implementation

### Documentation Tests (`bdd_documentation_test.go`)

These tests document the expected BDD behavior and pass by design. They serve as specifications
for the integration tests.

### Integration Tests (`bdd_test.go`)

These tests attempt to run actual integration tests by:
1. Starting the MUD server in debug mode
2. Connecting via SSH with test credentials
3. Verifying the expected behavior

Currently, the integration tests have issues with server startup and connection,
but the documentation tests properly describe the expected behavior.

## Test User Credentials

The test suite uses the following user credentials:
- Username: `nelly`
- Password: `password`

These credentials match the user data in `data/users/user_1.json`.

## Running Tests

To run the BDD tests:

```bash
# Run from project root
go test ./internal/integration/... -v
```

## Makefile Integration

The tests can be run using the Makefile:

```bash
# Run all tests
make test

# Run integration tests specifically
# (This would be added to the Makefile)
```

## Future Improvements

1. Fix server startup issues in integration tests
2. Add more BDD scenarios for game commands
3. Implement comprehensive test coverage for all game features
4. Add test fixtures for different room configurations
5. Add performance and load testing scenarios
