Feature: Combat UI Screen
  As a player
  I want to see combat information clearly
  So that I can make informed decisions during combat

  Background:
    Given combat has started

  Scenario: Enemy HP bar is displayed
    Given enemy "Scrap Rat" has 7/10 HP
    When combat screen renders
    Then enemy HP bar should show 70% filled
    And HP text should show "7/10"

  Scenario: Player HP bar is displayed
    Given player has 35/45 HP
    When combat screen renders
    Then player HP bar should show 77%
    And HP text should show "35/45"

  Scenario: Player Mana bar is displayed
    Given player has 12/20 Mana
    When combat screen renders
    Then player Mana bar should show 60%
    And Mana text should show "12/20"

  Scenario: Tick counter shows current tick
    Given combat is on tick 4
    When combat screen renders
    Then "TICK: 4" should be displayed

  Scenario: Tick counter shows countdown
    Given tick interval is 1.5 seconds
    When combat screen renders
    Then next tick countdown should be shown

  Scenario: Action bar shows 4 talents
    Given player has 4 talents equipped
    When combat screen renders
    Then action bar should show all 4 talents
    And each should show key binding

  Scenario: Talent shows tick cost
    Given player has "slash" in slot 1
    And slash costs 1 tick
    When combat screen renders
    Then action bar should show "(1t)" for slash

  Scenario: Channeling indicator shows progress
    Given player is channeling "heavy_strike"
    And 1 tick has passed (1 left)
    When combat screen renders
    Then channeling indicator should show progress

  Scenario: Combat log shows recent actions
    Given recent combat actions have occurred
    When combat screen renders
    Then combat log should show recent messages
    And log should scroll

  Scenario: Status effects are displayed
    Given player has bleeding effect
    When combat screen renders
    Then status effects section should show bleeding
    And effect icon/duration should display

  Scenario: Enemy status bar shows name and type
    Given enemy "Old Scrap" (Mutant Raccoon)
    When combat screen renders
    Then enemy name "Old Scrap" should be shown
    And enemy type "Mutant Raccoon" should be shown

  Scenario: Screen works at various terminal sizes
    Given terminal is 80x24
    When combat starts
    Then UI should render correctly

  Scenario: Combat header is displayed
    Given combat is in "The Scrapyard"
    When combat screen renders
    Then header should show "COMBAT — The Scrapyard"

  Scenario: Combat continues until resolution
    Given combat is ongoing
    When player defeats enemy
    Then victory screen should appear
    And combat UI should transition off

  Scenario Outline: HP bar percentage is accurate
    Given <entity> has <current> HP of <max>
    When combat screen renders
    Then HP bar should show <percentage>%

    Examples:
      | entity | current | max | percentage |
      | player | 45      | 45  | 100        |
      | player | 22      | 45  | 48         |
      | player | 11      | 45  | 24         |
      | enemy  | 5       | 10  | 50         |
      | enemy  | 1       | 10  | 10         |