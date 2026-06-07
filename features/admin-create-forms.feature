Feature: Admin Create Forms Manual QA
  As a game administrator
  I want to manually test the admin panel create forms
  So that I can verify all admin create forms work correctly

  Background:
    Given I am logged into the admin panel at http://100.67.206.65:5173
      * with username: sma
      * with password: sma

  Scenario: Create an NPC Template
    When I navigate to Content > NPCs > + Add Template
    Then the NPC create form renders without errors
    When I fill "Name" with "QA-TEST NPC"
    And I fill "Level" with "1"
    And I fill "XP Value" with "10"
    And I click "Create Template"
    Then I am redirected to the NPCs list
    And "QA-TEST NPC" appears in the list
    When I delete "QA-TEST NPC"
    Then "QA-TEST NPC" should no longer exist in the list

  Scenario: Create an Item Template
    When I navigate to Content > Items > + Add Item
    Then the Item create form renders without errors
    When I fill "Name" with "QA-TEST Item"
    And I fill "Description" with "Test item for QA purposes"
    And I click "Create Item"
    Then I am redirected to the Items list
    And "QA-TEST Item" appears in the list
    When I delete "QA-TEST Item"
    Then "QA-TEST Item" should no longer exist in the list

  Scenario: Create an Ability
    When I navigate to Content > Abilities > + Add Ability
    Then the Ability create form renders without errors
    When I fill "Name" with "QA-TEST Ability"
    And I fill "Description" with "Test ability for QA purposes"
    And I click "Create Ability"
    Then I am redirected to the Abilities list
    And "QA-TEST Ability" appears in the list
    When I delete "QA-TEST Ability"
    Then "QA-TEST Ability" should no longer exist in the list

  Scenario: Create a Quest
    When I navigate to Content > Quests > + Add Quest
    Then the Quest create form renders without errors
    When I fill "Title" with "QA-TEST Quest"
    And I fill "Description" with "Test quest for QA purposes"
    And I click "Create Quest"
    Then I am redirected to the Quests list
    And "QA-TEST Quest" appears in the list
    When I delete "QA-TEST Quest"
    Then "QA-TEST Quest" should no longer exist in the list

  Scenario: Create a Trigger
    When I navigate to Content > Triggers > + Add Trigger
    Then the Trigger create form renders without errors
    When I fill "Name" with "QA-TEST Trigger"
    And I fill "Description" with "Test trigger for QA purposes"
    And I click "Create Trigger"
    Then I am redirected to the Triggers list
    And "QA-TEST Trigger" appears in the list
    When I delete "QA-TEST Trigger"
    Then "QA-TEST Trigger" should no longer exist in the list

  Scenario: Create a Social
    When I navigate to Social > + Add Social
    Then the Social create form renders without errors
    When I fill "Name" with "QA-TEST Social"
    And I fill "Message" with "Test social message for QA purposes"
    And I click "Create Social"
    Then I am redirected to the Socials list
    And "QA-TEST Social" appears in the list
    When I delete "QA-TEST Social"
    Then "QA-TEST Social" should no longer exist in the list

  Scenario: Create a Skill / Competency
    When I navigate to Content > Skills > + Add Skill
    Then the Skill create form renders without errors
    When I fill "Name" with "QA-TEST Skill"
    And I fill "Description" with "Test skill for QA purposes"
    And I click "Create Skill"
    Then I am redirected to the Skills list
    And "QA-TEST Skill" appears in the list
    When I delete "QA-TEST Skill"
    Then "QA-TEST Skill" should no longer exist in the list

  Scenario: Cleanup - Delete all QA test data
    Given I am on the admin panel
    When I delete all items matching "QA-TEST NPC"
    And I delete all items matching "QA-TEST Item"
    And I delete all items matching "QA-TEST Ability"
    And I delete all items matching "QA-TEST Quest"
    And I delete all items matching "QA-TEST Trigger"
    And I delete all items matching "QA-TEST Social"
    And I delete all items matching "QA-TEST Skill"
    Then all QA-TEST items should be removed from the system
