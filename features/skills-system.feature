Feature: Skills System
  As a player
  I want weapon and magic proficiency skills that affect combat performance
  So that my character's expertise matters in battle

  Background:
    Given I have created a character
    And my character has skills initialized

  Scenario: Character starts with Warrior skills
    Given I create a Warrior class character
    When I view my starting skills
    Then I should have "blades" at level 1
    And I should have "brawling" at level 1
    And my skills should be stored in the CharacterSkill table

  Scenario: Skill levels affect damage
    Given I have blades skill at level 25
    When I attack with a blade weapon
    Then my damage should receive +0% bonus (novice tier)

  Scenario: Skill level 26-50 grants bonus
    Given I have blades skill at level 50
    When I attack with a blade weapon
    Then my damage should receive +10% bonus

  Scenario: Skill level 51-75 grants bonus
    Given I have blades skill at level 75
    When I attack with a blade weapon
    Then my damage should receive +25% bonus

  Scenario: Skill level 76-90 grants bonus
    Given I have blades skill at level 90
    When I attack with a blade weapon
    Then my damage should receive +50% bonus

  Scenario: Skill level 91-99 grants bonus
    Given I have blades skill at level 99
    When I attack with a blade weapon
    Then my damage should receive +75% bonus

  Scenario: Skill level 100 mastery bonus
    Given I have blades skill at level 100
    When I attack with a blade weapon
    Then my damage should receive +100% bonus
    And I should be recognized as a Master

  Scenario: Skills display in profile
    Given I have various skills at different levels
    When I view my character profile
    Then I should see all my weapon/magic skills
    And each skill should show its current level
    And each skill should show its bonus tier

  Scenario: Skills API - GET character skills
    Given I have a character with skills
    When I call GET /characters/:id/skills
    Then I should receive a list of all character skills
    And each skill should include id, type, and level

  Scenario: Skills API - PUT character skills
    Given I have a character
    When I call PUT /characters/:id/skills
    Then I should be able to update skill levels
    And the response should confirm the updated skills

  Scenario: Skill improves through combat use
    Given I have blades skill at level 10
    And I am in combat with a blade weapon
    When I successfully hit enemies
    Then my blades skill should gain XP
    And my skill level should eventually increase with practice

  Scenario: All weapon categories available
    Given I view available weapon skills
    When I check the skill list
    Then I should see: blades, staves, knives, martial, brawling, tech
    And each skill should show its primary stat

  Scenario: All magic categories available
    Given I view available magic skills
    When I check the skill list
    Then I should see: fire_magic, water_magic, wind_magic
    And each magic skill should show its primary stat (INT or WIS)

  Scenario: Blades skill - weapon types
    Given I have blades skill
    When I use different blade weapons
    Then machetes should use blades skill
    And swords should use blades skill
    And cleavers should use blades skill
    And scrap blades should use blades skill

  Scenario: Staves skill - weapon types
    Given I have staves skill
    When I use different staff weapons
    Then hockey sticks should use staves skill
    And bows should use staves skill
    And spears should use staves skill
    And pool cues should use staves skill

  Scenario: Brawling skill - improvised weapons
    Given I have brawling skill
    When I use improvised weapons
    Then table legs should use brawling skill
    And bottles should use brawling skill
    And chains should use brawling skill
    And fists should use brawling skill

  Scenario: Tech skill - rare weapons
    Given I have tech skill
    When I use tech weapons
    Then laser pistols should use tech skill
    And beam rifles should use tech skill
    And EMP devices should use tech skill

  Scenario: Skill unlocks talents at certain levels
    Given I have blades skill
    When my skill reaches required level
    Then new blade-related talents should become available
    And I should be notified of unlocked abilities