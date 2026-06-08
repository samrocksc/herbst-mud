Feature: Dashboard Active NPCs Count
  As a game administrator
  I want to see the correct count of active NPC instances on the dashboard
  So that I can monitor my world's NPC population at a glance

  Background:
    Given I am authenticated as admin "sma" with password "sma"

  Scenario: Dashboard shows active NPC instances for a world
    Given world "Ooze Surfers" (id=2) has NPC template "Gizmo"
    And   NPC template "Gizmo" has at least 1 active instance (hitpoints > 0) in world 2
    When  I navigate to the admin dashboard
    And   I select world "2" from the world dropdown
    Then  the "Active NPCs" stat card should display "1"

  Scenario: Dashboard NPC count responds to world switching
    Given world "herbst-mud" (id=1) has no NPC instances
    And   world "Ooze Surfers" (id=2) has 1 active NPC instance
    When  I navigate to the dashboard with world "herbst-mud" selected
    Then  the "Active NPCs" stat card should display "0"
    When  I switch to world "Ooze Surfers"
    Then  the "Active NPCs" stat card should display "1"
    When  I switch back to world "herbst-mud"
    Then  the "Active NPCs" stat card should display "0"

  Scenario: NPC list page correctly counts instances (regression guard)
    Given world "Ooze Surfers" is selected
    When  I navigate to "/npcs"
    Then  the "Instances" column should show "1" for template "Gizmo"
    And   the "Instances" column should show "0" for template "Teamo"
