Feature: Quest Backend API
  Quest management system handles quest definitions, character quest state, and quest operations

  Background:
    Given quest service is running
    And quest database is initialized

  # ─── Quest Definition & Storage ────────────────────────────────

  Scenario: Quest definitions are stored in database
    Given quest definition "Prove Yourself" exists:
      | field              | value                    |
      | id                 | quest_prove_yourself     |
      | name               | Prove Yourself           |
      | description        | Show your worth to Guard Marco |
      | level_requirement  | 1                        |
      | xp_reward          | 150                      |
      | item_rewards       | repair_kit               |
      | currency_reward    | 50                       |
    When quest is registered in system
    Then quest can be retrieved by ID

  Scenario: Quest objectives are linked to quest
    Given quest "Prove Yourself" has objectives:
      | objective_type | target        | count |
      | kill            | Scrap Rat     | 3     |
      | return          | Guard Marco   | 1     |
    When retrieving quest definition
    Then all objectives are included

  Scenario: Quest difficulty affects rewards
    Given quest has difficulty tier <tier>
    When calculating rewards
    Then XP multiplier is <multiplier>
    And currency multiplier is <multiplier>

    Examples:
      | tier | multiplier |
      | 1    | 1.0         |
      | 3    | 1.5         |
      | 5    | 2.0         |
      | 10   | 3.0         |

  # ─── Character Quest State ─────────────────────────────────────

  Scenario: Character quest state tracks progress
    Given character "Player1" has quest "Prove Yourself" active
    When checking quest state
    Then state includes:
      | field             | value                    |
      | quest_id          | quest_prove_yourself     |
      | status            | active                   |
      | objectives        | (in progress)            |
      | started_at        | (timestamp)              |

  Scenario: Quest progress persists across sessions
    Given character has quest "Prove Yourself" with 2/3 Scrap Rats killed
    When character reconnects
    Then quest progress is preserved
    And objective count is 2/3

  Scenario: Quest state includes completion timestamp
    Given character completed quest "First Blood"
    When retrieving quest state
    Then completed_at timestamp is set

  # ─── Quest Operations ──────────────────────────────────────────

  Scenario: Accept quest adds to character
    Given NPC offers quest "Prove Yourself"
    When character accepts quest
    Then quest is added to character's active quests
    And objectives are initialized to 0 progress

  Scenario: Decline quest does not add to character
    Given NPC offers quest "Prove Yourself"
    When character declines quest
    Then quest is not in character quest list
    And quest can be re-offered later

  Scenario: Abandon quest removes from active
    Given character has quest "Prove Yourself" active
    When character abandons quest
    Then quest is removed from active quests
    And progress is lost
    And no penalty applied

  # ─── Objective Progress ────────────────────────────────────────

  Scenario: Kill objective increments on enemy death
    Given character has quest with kill objective "Kill 3 Scrap Rats"
    And current progress is 2
    When enemy "Scrap Rat" is killed
    Then objective progress becomes 3
    And objective marked complete when reaching 3

  Scenario: Collect objective increments on item pickup
    Given character has quest with collect objective "Collect 5 Ooze"
    And current progress is 4
    When character picks up "Ooze Sample"
    Then objective progress becomes 5
    And objective marked complete

  Scenario: Explore objective completes on room entry
    Given quest has explore objective "Enter Old Market"
    When character enters room "Old Market"
    Then objective progress becomes 1
    And objective marked complete

  Scenario: Talk objective completes on NPC interaction
    Given quest has talk objective "Talk to Guard Marco"
    When character talks to "Guard Marco"
    Then objective progress becomes 1
    And objective marked complete

  Scenario: Objectives don't overcount
    Given quest requires 3 kills
    And progress is already 3
    When target is killed again
    Then progress stays at 3
    And no overflow

  # ─── Quest Completion ──────────────────────────────────────────

  Scenario: Quest completes when all objectives done
    Given quest "Prove Yourself" has all objectives at 100%
    When checking quest completion
    Then quest is marked completeable

  Scenario: Turn in quest distributes rewards
    Given quest "Prove Yourself" awards:
      | reward_type | amount |
      | xp          | 150    |
      | currency    | 50     |
      | items       | repair_kit |
    When quest is turned in
    Then XP is added to character
    And currency is added to character
    And items added to inventory

  Scenario: Quest completion updates status to completed
    Given quest "Prove Yourself" is turned in
    When quest state is retrieved
    Then status is "completed"
    And completed_at is set

  Scenario: Completed quest cannot be turned in again
    Given quest "First Blood" is already completed
    When attempting to turn in again
    Then error is returned
    And no duplicate rewards

  # ─── Quest Availability ────────────────────────────────────────

  Scenario: Quest available when level requirement met
    Given quest requires level 5
    When character level is 5
    Then quest is available to accept

  Scenario: Quest unavailable when level too low
    Given quest requires level 10
    When character level is 5
    Then quest is not shown as available
    And cannot be accepted

  Scenario: Quest available after prerequisite completed
    Given quest "Chapter 2" has prerequisite "Chapter 1"
    When "Chapter 1" is completed
    Then "Chapter 2" becomes available

  # ─── Quest Expiration ──────────────────────────────────────────

  Scenario: Time-limited quest has deadline
    Given quest has 30 minute time limit
    When quest is accepted
    Then deadline timestamp is set

  Scenario: Expired quest fails automatically
    Given time-limited quest has expired
    When checking quest status
    Then quest status becomes "failed"
    And quest removed from active
    And no rewards given

  # ─── Hidden & Locked Quests ────────────────────────────────────

  Scenario: Hidden quest doesn't appear in available list
    Given quest "Secret Vault" is marked hidden
    When listing available quests
    Then "Secret Vault" is not shown

  Scenario: Hidden quest unlocks on condition
    Given hidden quest "Secret Vault" unlocks on examining "vault_door" with examine level 50
    When character examines "vault_door" with level 50
    Then "Secret Vault" becomes available

  Scenario: Secret quest revealed through action
    Given quest is triggered by examine
    When examine reveals quest
    Then quest appears in available quests
    And notification shown to player

  # ─── Quest Categories ──────────────────────────────────────────

  Scenario: Main story quest gates progression
    Given main quest "Chapter 1" must be completed
    When "Chapter 1" is incomplete
    Then "Chapter 2" is locked

  Scenario: Side quest is optional
    Given side quest "Scrap Collector" is optional
    When side quest completes
    Then main quest progression unaffected

  Scenario: Daily quest resets
    Given daily quest "Daily Scrap" resets at midnight
    When new day begins
    Then quest becomes available again
    And previous progress cleared

  # ─── API Endpoints ─────────────────────────────────────────────

  Scenario: GET /quests returns available quests for character
    When GET /api/quests is called for character
    Then response includes available, active, completed quests

  Scenario: POST /quests/:id/accept accepts quest
    When POST /api/quests/quest_prove_yourself/accept is called
    Then quest added to active
    And 201 returned

  Scenario: POST /quests/:id/abandon abandons quest
    When POST /api/quests/quest_prove_yourself/abandon is called
    Then quest removed from active
    And 200 returned

  Scenario: GET /quests/:id returns quest details
    When GET /api/quests/quest_prove_yourself is called
    Then quest definition is returned
    And current progress for character is included

  Scenario: POST /quests/:id/turnin turns in completed quest
    Given all objectives are complete
    When POST /api/quests/quest_prove_yourself/turnin is called
    Then rewards distributed
    And 200 returned with reward summary

  # ─── Error Handling ─────────────────────────────────────────────

  Scenario: Invalid quest ID returns 404
    When GET /api/quests/invalid_quest is called
    Then 404 returned

  Scenario: Accepting already-active quest returns error
    Given quest is already active
    When POST /api/quests/:id/accept is called
    Then 400 returned with error message

  Scenario: Turning in incomplete quest returns error
    Given quest has incomplete objectives
    When POST /api/quests/:id/turnin is called
    Then 400 returned with error message

  Scenario: Quest database error handled gracefully
    Given database connection is lost
    When quest operation is attempted
    Then error is logged
    And 503 returned to client
