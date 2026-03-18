Feature: User Entity
  As a developer
  I want a User entity model
  So that I can manage user data in the system

  Scenario: User entity fields
    Given I need to model a User
    When I create the User entity
    Then it should have the following fields:
      | field       | type   |
      | ID          | int    |
      | Username    | string |
      | Email       | string |
      | Password    | string |
      | Role        | string |
      | CreatedAt   | time   |
      | UpdatedAt   | time   |

  Scenario: User relationships
    Given a User entity exists
    Then it should have a one-to-many relationship with Characters
    And it should have a one-to-many relationship with Sessions