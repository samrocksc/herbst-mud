Feature: User Entity
  As a game developer
  I want a User entity
  So that I can manage user accounts in the system

  Scenario: User entity fields
    Given I need a User entity
    When I define the User model
    Then it should have the following fields:
      | field        | type      |
      | id           | UUID      |
      | username     | string    |
      | email        | string    |
      | passwordHash | string    |
      | createdAt    | timestamp |
      | updatedAt    | timestamp |

  Scenario: User validation
    Given I have user data
    When I validate a User
    Then username should be 3-20 characters
    And email should be valid format
    And password should be minimum 8 characters