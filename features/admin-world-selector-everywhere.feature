Feature: Admin — World selector available on all content pages

  Background:
    Given I am authenticated as an admin

  Scenario: Switch world from the races page
    Given I navigate to "/races"
    When I select world "Ooze Surfers" from the top navigation
    Then the active world indicator should show "Ooze Surfers"
    And the races list should refresh for world "Ooze Surfers"
