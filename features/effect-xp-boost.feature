Feature: XP Boost Effect
  As a game administrator
  I want to create XP multiplier effects for consumable items
  So that players can temporarily increase their XP gain rate

  Scenario: Create XP Boost effect definition
    Given I am authenticated as an admin
    When I POST to "/api/effects" with:
      | name        | XP Boost           |
      | effect_type | xp_multiplier      |
      | parameters  | {"multiplier": 2.0, "duration_secs": 300} |
      | stack_mode  | replace            |
      | description | Doubles XP gain for 5 minutes |
    Then the response status should be 201
    And the response should contain:
      | name         | XP Boost      |
      | effect_type  | xp_multiplier |
      | parameters.multiplier | 2   |
      | duration_secs | 300        |

  Scenario: Apply XP Boost effect to character
    Given an effect "XP Boost" exists
    And a character "smack" exists
    When I POST to "/api/characters/{char_id}/effects/apply" with:
      | effect_id  | {effect_id} |
      | applied_by_id | 3       |
      | duration_secs | 300     |
    Then the response status should be 201
    And the active effect should have:
      | effect_type | xp_multiplier |
      | multiplier  | 2             |
      | is_active   | true          |

  Scenario: XP multiplier increases XP gain
    Given a character "smack" with XP 1000
    And an active XP Boost effect on character "smack" with multiplier 2
    When I POST to "/api/events" with:
      | type        | quest_complete |
      | payload.character_id | 5        |
      | payload.xp_reward    | 500      |
    Then the response status should be 200
    And character "smack" should have XP 2000 (500 * 2 multiplier)

  Scenario: XP multiplier expires after duration
    Given an XP Boost effect on character "smack" with duration 300 seconds
    When the effect expires (after 300 seconds)
    Then character "smack" should have:
      | xp_multiplier | 1.0         |
      | xp_multiplier_end | null    |

  Scenario: Consume XP Potion to apply effect
    Given an item template "XP Potion" with effect_type: xp_multiplier
    And an instance of "XP Potion" owned by character "smack"
    When character "smack" uses the XP Potion
    Then character "smack" should receive an active XP Boost effect
    And the effect should have multiplier: 2 and duration: 300 seconds
