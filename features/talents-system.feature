Feature: Talents and Proficiencies System
  As a player
  I want combat talents that I can equip and use in battle
  So that I can customize my combat playstyle

  Background:
    Given I have created a character
    And my character has talents initialized

  Scenario: Warrior starts with 4 talents
    Given I create a Warrior class character
    When I view my starting talents
    Then I should have "slash" talent
    And I should have "parry" talent
    And I should have "shield_bash" talent
    And I should have "heavy_strike" talent
    And my talents should be stored in the CharacterTalent table

  Scenario: View talent bar
    Given I have talents equipped
    When I view my talent bar
    Then I should see 4 talent slots (keys 1-4)
    And each slot should show the equipped talent name
    And I should see my total available talents

  Scenario: Equip talent to slot
    Given I have learned 8 talents
    And I have 4 talents currently equipped
    When I want to change my loadout
    Then I should be able to equip any learned talent
    And I should see the updated talent bar

  Scenario: Use talent in combat with hotkey
    Given I have "slash" equipped in slot 1
    And I am in combat
    When I press key "1"
    Then I should perform the "slash" attack
    And the appropriate resources should be consumed

  Scenario: Swap talents out of combat
    Given I am not in combat
    And I have learned talents: slash, parry, smash, battle_cry, second_wind
    When I open the talent management screen
    Then I should be able to swap any learned talent into my 4 slots
    And changes should persist until I swap again

  Scenario: Talent requires weapon skill
    Given I want to use "slash" talent
    And I have no blades or knives skill
    When I try to equip "slash"
    Then I should receive an error that skill requirement is not met
    And the required skills should be displayed

  Scenario: Parry talent - no skill requirement
    Given I have any character
    When I equip "parry" talent
    Then I should be able to use it regardless of weapon skills
    And it should deflect incoming attacks when active

  Scenario: Slash talent - requires blades or knives
    Given I have blades skill at level 10
    When I equip "slash" talent
    Then I should be able to use slash
    And slash should perform a blade/knife attack

  Scenario: Smash talent - requires staves or martial
    Given I have staves skill at level 10
    When I equip "smash" talent
    Then I should be able to use smash
    And smash should perform a powerful blunt attack

  Scenario: Crash talent - uses weight for damage
    Given I have "crash" talent equipped
    And my character is size large
    When I use crash in combat
    Then damage should be calculated based on my weight
    And my STR should affect the damage

  Scenario: Shield bash - stun effect
    Given I have "shield_bash" equipped
    When I use shield bash on an enemy
    Then the enemy should receive damage
    And the enemy should have a chance to be stunned

  Scenario: Battle cry - debuff enemies
    Given I have "battle_cry" equipped
    When I use battle cry
    Then enemies should receive -accuracy debuff
    And the effect should last for combat duration

  Scenario: Second wind - heal when low
    Given I have "second_wind" equipped
    And my HP is below 30%
    When I use second wind
    Then I should recover HP
    And I should see the healing message

  Scenario: Hail storm - double attacks
    Given I have "hail_storm" equipped
    When I use hail storm
    Then I should attack twice per cycle
    And the effect should last 2 cycles

  Scenario: Iron will - passive talent
    Given I have "iron_will" equipped
    When I am hit by stun or blind effects
    Then I should resist those effects
    And iron_will should not count toward my 4-slot limit

  Scenario: Heavy strike - slow but strong
    Given I have "heavy_strike" equipped
    And I have blades or staves skill
    When I use heavy strike
    Then the attack should deal more damage
    And the attack should be slower than normal

  Scenario: Talents API - GET character talents
    Given I have a character with talents
    When I call GET /characters/:id/talents
    Then I should receive a list of all character talents
    And I should see which are equipped

  Scenario: Talents API - PUT character talents
    Given I have a character
    When I call PUT /characters/:id/talents
    Then I should be able to update equipped talents
    And the response should confirm the updated loadout

  Scenario: Talent bar persists across sessions
    Given I have talents: slash, parry, smash, crash equipped
    When I log out and log back in
    Then my talent bar should show the same 4 talents
    And my hotkey mappings should be preserved