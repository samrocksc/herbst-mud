# SSH Client UI Overhaul Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Polish the SSH MUD client with a fantasy color theme, clearer screen layouts, better onboarding, and a restructured playing screen.

**Architecture:** The herbst/ package uses bubbletea with a model/view/update pattern. Auth screens render via standalone functions in ui_screens.go. The playing screen renders inline in game_model.go View(). The style system in style.go provides shared color/style vars. All screens use split-panel layout (output top, input bottom).

**Tech Stack:** Go, bubbletea, lipgloss, charmbracelet SSH/wish

**Spec:** `docs/superpowers/specs/2026-05-16-ssh-client-ui-overhaul-design.md`

---

### Task 1: Update Color Palette to Fantasy Theme

**Files:**
- Modify: `herbst/style.go:1-353`

**Overview:** Replace the post-apocalyptic neon palette (pink/teal/orange) with a fantasy theme (deep purples, rich golds, bright blues). Keep the same style struct/layout — just change color values and add new shared styles for the header bar.

- [ ] **Step 1: Replace color variables**

Change the color palette section:

```go
var (
    // Fantasy primary colors
    PrimaryPurple = lipgloss.Color("129") // #AF00FF - Deep purple
    PrimaryGold   = lipgloss.Color("220") // #FFD700 - Rich gold
    AccentBlue    = lipgloss.Color("75")  // #5F87FF - Bright blue

    // Status colors (keep same semantics, adjust shades)
    StatusRed     = lipgloss.Color("196") // Red - danger/low HP
    StatusGreen   = lipgloss.Color("46")  // Green - healthy/good
    StatusYellow  = lipgloss.Color("226") // Yellow - warning
    StatusBlue    = lipgloss.Color("75")  // Blue - mana/magic

    // Neutral colors
    TextWhite    = lipgloss.Color("15")  // Primary text
    TextGray     = lipgloss.Color("8")   // Secondary text
    TextDarkGray = lipgloss.Color("236") // Dark backgrounds
    BorderColor  = lipgloss.Color("240") // Borders

    // Exit colors (for navigation)
    ExitVisited = lipgloss.Color("46")  // Green - visited
    ExitKnown   = lipgloss.Color("220") // Gold - known
    ExitNew     = lipgloss.Color("15")  // White - new
)

// Header bar colors
var (
    HeaderBg    = lipgloss.Color("54")  // Dark purple background
    HeaderFg    = lipgloss.Color("220") // Gold text
    StatusBarBg = lipgloss.Color("236") // Dark bg for status
)
```

- [ ] **Step 2: Update typography styles to use new colors**

```go
var (
    TitleStyle = lipgloss.NewStyle().
            Foreground(PrimaryGold).
            Background(TextDarkGray).
            Bold(true).
            Padding(0, 1).
            Align(lipgloss.Center)

    SubtitleStyle = lipgloss.NewStyle().
            Foreground(PrimaryPurple).
            Padding(0, 1).
            Align(lipgloss.Center)

    HeaderStyle = lipgloss.NewStyle().
            Foreground(AccentBlue).
            Bold(true).
            Padding(0, 1)
)
```

- [ ] **Step 3: Add a HeaderBarStyle for the playing screen top bar**

```go
var (
    HeaderBarStyle = lipgloss.NewStyle().
            Background(HeaderBg).
            Foreground(HeaderFg).
            Bold(true).
            Padding(0, 1).
            Width(80).
            Align(lipgloss.Center)

    HeaderBarRoomStyle = lipgloss.NewStyle().
            Foreground(HeaderFg).
            Bold(true)

    HeaderBarInfoStyle = lipgloss.NewStyle().
            Foreground(AccentBlue)
)
```

- [ ] **Step 4: Update border styles to use purple/gold**

```go
var (
    BoxStyle = lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(PrimaryPurple).
            Padding(1, 2).
            Align(lipgloss.Left)

    DoubleBoxStyle = lipgloss.NewStyle().
            Border(lipgloss.DoubleBorder()).
            BorderForeground(PrimaryGold).
            Padding(1, 2)

    ThickBoxStyle = lipgloss.NewStyle().
            Border(lipgloss.ThickBorder()).
            BorderForeground(PrimaryGold).
            Padding(1, 2)

    HighlightBoxStyle = lipgloss.NewStyle().
            Border(lipgloss.NormalBorder(), false, false, false, true).
            BorderForeground(PrimaryPurple).
            Padding(0, 0, 0, 1)
)
```

- [ ] **Step 5: Update box style Width values to be dynamic**

The box styles currently hardcode Width(58). Remove the Width from BoxStyle, DoubleBoxStyle, ThickBoxStyle so they adapt to their container. The width will be set at render time based on terminal width.

Remove `.Width(58)` from BoxStyle, DoubleBoxStyle, and ThickBoxStyle.

- [ ] **Step 6: Update backward-compatibility aliases**

Update the aliases at the bottom of style.go to match the new primary colors:

```go
var (
    pink  = PrimaryPurple  // was AccentPink, now purple
    cyan  = AccentBlue     // was PrimaryTeal, now blue
    // ... rest stay the same
    purple = PrimaryPurple
)
```

- [ ] **Step 7: Build and verify**

```bash
cd /home/sam/GitHub/herbst-mud && make build-all
```

Expect: clean compile. No functional changes yet.

- [ ] **Step 8: Commit**

```bash
git add herbst/style.go
git commit -m "feat: update SSH client color palette to fantasy theme (purple/gold/blue)"
```

---

### Task 2: Rewrite Auth Screen Renderers

**Files:**
- Modify: `herbst/ui_screens.go:1-232`

**Overview:** Rewrite the screen layout functions (welcomeScreen, loginScreen, registerScreen, worldSelectScreen, characterSelectScreen) to use the new fantasy palette. Keep the split-panel approach but improve spacing, add gold/purple headers, and styled instruction text.

- [ ] **Step 1: Rewrite welcomeScreen**

```go
func welcomeScreen(width, height int, inputView string) string {
    if width < 40 {
        width = 80
    }

    // Calculate split
    outputHeight := height * 60 / 100
    if outputHeight < 8 {
        outputHeight = 8
    }

    var outputContent strings.Builder

    // Title banner
    outputContent.WriteString("\n")
    titleLine := lipgloss.NewStyle().
        Bold(true).
        Foreground(PrimaryGold).
        Render("✦ HERBST MUD ✦")
    outputContent.WriteString(lipgloss.NewStyle().
        Width(width).
        Align(lipgloss.Center).
        Render(titleLine))
    outputContent.WriteString("\n\n")

    // Subtitle
    subtitleLine := lipgloss.NewStyle().
        Foreground(PrimaryPurple).
        Italic(true).
        Render("A World of Adventure Awaits")
    outputContent.WriteString(lipgloss.NewStyle().
        Width(width).
        Align(lipgloss.Center).
        Render(subtitleLine))
    outputContent.WriteString("\n\n")

    // Menu
    menuItems := []struct {
        key  string
        desc string
    }{
        {"1", "Login"},
        {"2", "Register"},
        {"3", "Quit"},
    }
    for _, item := range menuItems {
        keyStyle := lipgloss.NewStyle().Foreground(AccentBlue).Bold(true).Render(item.key + ".")
        nameStyle := lipgloss.NewStyle().Foreground(TextWhite).Render(item.desc)
        outputContent.WriteString(fmt.Sprintf("  %s  %s\n", keyStyle, nameStyle))
    }
    outputContent.WriteString("\n")
    outputContent.WriteString(lipgloss.NewStyle().
        Foreground(TextGray).
        Width(width).
        Align(lipgloss.Center).
        Render("Type a number or command"))

    outputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryPurple).
        Padding(0, 1).
        Width(width).
        Height(outputHeight)
    inputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryGold).
        Padding(0, 1).
        Width(width).
        Height(height - outputHeight - 1)

    var sb strings.Builder
    sb.WriteString(outputStyle.Render(outputContent.String()))
    sb.WriteString("\n")
    sb.WriteString(inputStyle.Render(inputView))
    return sb.String()
}
```

- [ ] **Step 2: Rewrite loginScreen**

```go
func loginScreen(width, height int, message, messageType string, inputView string) string {
    if width < 40 {
        width = 80
    }

    outputHeight := height * 55 / 100
    if outputHeight < 8 {
        outputHeight = 8
    }

    var outputContent strings.Builder
    outputContent.WriteString("\n")
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(PrimaryGold).
        Width(width).
        Align(lipgloss.Center).
        Render("✦ ACCOUNT LOGIN ✦")
    outputContent.WriteString(title)
    outputContent.WriteString("\n\n")

    if message != "" {
        styled := styleMessage(message, messageType)
        outputContent.WriteString(styled)
        outputContent.WriteString("\n\n")
    }

    instructions := lipgloss.NewStyle().
        Foreground(TextGray).
        Render("Type 'register' to create a new account")
    outputContent.WriteString(instructions)
    outputContent.WriteString("\n")
    quitText := lipgloss.NewStyle().
        Foreground(TextGray).
        Render("Type 'quit' or press Ctrl+C to exit")
    outputContent.WriteString(quitText)

    outputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryPurple).
        Padding(0, 1).
        Width(width).
        Height(outputHeight)
    inputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryGold).
        Padding(0, 1).
        Width(width).
        Height(height - outputHeight - 1)

    var sb strings.Builder
    sb.WriteString(outputStyle.Render(outputContent.String()))
    sb.WriteString("\n")
    sb.WriteString(inputStyle.Render(inputView))
    return sb.String()
}
```

- [ ] **Step 3: Rewrite registerScreen**

```go
func registerScreen(width, height int, message, messageType string, inputView string) string {
    if width < 40 {
        width = 80
    }

    outputHeight := height * 55 / 100
    if outputHeight < 8 {
        outputHeight = 8
    }

    var outputContent strings.Builder
    outputContent.WriteString("\n")
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(PrimaryGold).
        Width(width).
        Align(lipgloss.Center).
        Render("✦ CREATE ACCOUNT ✦")
    outputContent.WriteString(title)
    outputContent.WriteString("\n\n")

    if message != "" {
        styled := styleMessage(message, messageType)
        outputContent.WriteString(styled)
        outputContent.WriteString("\n\n")
    }

    instructions := lipgloss.NewStyle().
        Foreground(TextGray).
        Render("Press Esc to go back")
    outputContent.WriteString(instructions)

    outputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryPurple).
        Padding(0, 1).
        Width(width).
        Height(outputHeight)
    inputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryGold).
        Padding(0, 1).
        Width(width).
        Height(height - outputHeight - 1)

    var sb strings.Builder
    sb.WriteString(outputStyle.Render(outputContent.String()))
    sb.WriteString("\n")
    sb.WriteString(inputStyle.Render(inputView))
    return sb.String()
}
```

- [ ] **Step 4: Rewrite worldSelectScreen**

```go
func worldSelectScreen(width, height int, displayContent string, inputView string) string {
    if width < 40 {
        width = 80
    }

    outputHeight := height * 55 / 100
    if outputHeight < 8 {
        outputHeight = 8
    }

    var outputContent strings.Builder
    outputContent.WriteString("\n")
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(PrimaryGold).
        Width(width).
        Align(lipgloss.Center).
        Render("✦ SELECT WORLD ✦")
    outputContent.WriteString(title)
    outputContent.WriteString("\n\n")
    outputContent.WriteString(displayContent)
    outputContent.WriteString("\n")
    hint := lipgloss.NewStyle().
        Foreground(TextGray).
        Render("Type the number or name of a world. 'b' to go back.")
    outputContent.WriteString(hint)

    outputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryPurple).
        Padding(0, 1).
        Width(width).
        Height(outputHeight)
    inputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryGold).
        Padding(0, 1).
        Width(width).
        Height(height - outputHeight - 1)

    var sb strings.Builder
    sb.WriteString(outputStyle.Render(outputContent.String()))
    sb.WriteString("\n")
    sb.WriteString(inputStyle.Render(inputView))
    return sb.String()
}
```

- [ ] **Step 5: Rewrite characterSelectScreen**

```go
func characterSelectScreen(width, height int, message, messageType string, inputView string) string {
    if width < 40 {
        width = 80
    }

    outputHeight := height * 55 / 100
    if outputHeight < 8 {
        outputHeight = 8
    }

    var outputContent strings.Builder
    outputContent.WriteString("\n")
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(PrimaryGold).
        Width(width).
        Align(lipgloss.Center).
        Render("✦ SELECT CHARACTER ✦")
    outputContent.WriteString(title)
    outputContent.WriteString("\n\n")

    if message != "" {
        styled := styleMessage(message, messageType)
        outputContent.WriteString(styled)
        outputContent.WriteString("\n\n")
    }

    hint := lipgloss.NewStyle().
        Foreground(TextGray).
        Render("Type number to select, 'n' for new character, 'b' for back")
    outputContent.WriteString(hint)

    outputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryPurple).
        Padding(0, 1).
        Width(width).
        Height(outputHeight)
    inputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryGold).
        Padding(0, 1).
        Width(width).
        Height(height - outputHeight - 1)

    var sb strings.Builder
    sb.WriteString(outputStyle.Render(outputContent.String()))
    sb.WriteString("\n")
    sb.WriteString(inputStyle.Render(inputView))
    return sb.String()
}
```

- [ ] **Step 6: Build and verify**

```bash
cd /home/sam/GitHub/herbst-mud && make build-all
```

Expect: clean compile.

- [ ] **Step 7: Commit**

```bash
git add herbst/ui_screens.go
git commit -m "feat: rewrite auth screen renderers with fantasy theme"
```

---

### Task 3: Restructure Playing Screen with Header Bar

**Files:**
- Modify: `herbst/game_model.go:491-562` (ScreenPlaying section of View)
- Modify: `herbst/model.go` (add new fields if needed)

**Overview:** Restructure the playing screen view into a clear three-zone layout: header bar at top, scrollable content (room + messages), status bar, and input area.

- [ ] **Step 1: Add header bar rendering function to game_model.go**

```go
// renderHeaderBar builds the top header bar with room/area name and compact stats
func (m *model) renderHeaderBar(width int) string {
    headerContent := fmt.Sprintf("%s  |  Lvl %d %s",
        lipgloss.NewStyle().Foreground(PrimaryGold).Bold(true).Render(m.roomName),
        m.characterLevel,
        m.characterRace,
    )
    return HeaderBarStyle.Width(width).Render(headerContent)
}
```

- [ ] **Step 2: Restructure the ScreenPlaying section of View()**

Replace lines 491-561 of game_model.go with:

```go
case ScreenPlaying:
    width := m.width
    height := m.height
    if width < 40 {
        width = 80
    }
    if height < 10 {
        height = 24
    }

    // Layout proportions
    headerHeight := 1
    inputHeight := height * 20 / 100
    if inputHeight < 3 {
        inputHeight = 3
    }
    statusHeight := height * 10 / 100
    if statusHeight < 3 {
        statusHeight = 3
    }
    contentHeight := height - headerHeight - inputHeight - statusHeight - 3
    if contentHeight < 5 {
        contentHeight = 5
    }

    // Header bar
    s.WriteString(m.renderHeaderBar(width))
    s.WriteString("\n")

    // Main content (room + messages)
    roomInfo := fmt.Sprintf("%s\n\nExits: %s",
        m.roomDesc,
        m.formatExitsWithColor())

    if len(m.messageHistory) > 0 {
        msgContent := m.buildOutputContent()
        if msgContent != "" {
            roomInfo += "\n\n" + msgContent
        }
    }

    if m.viewport.Width != width {
        m.viewport.Width = width
    }
    m.viewport.Height = contentHeight
    m.viewport.SetContent(roomInfo)

    contentStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryPurple).
        Width(width)
    s.WriteString(contentStyle.Render(m.viewport.View()))
    s.WriteString("\n")

    // Status bar
    regenActive := !m.inCombat && m.characterHP < m.characterMaxHP && m.characterHP > 0
    statsLine := MiniStatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana, regenActive)
    debugInfo := ""
    if m.debugMode {
        debugInfo = " " + lipgloss.NewStyle().Foreground(StatusYellow).Bold(true).Render(fmt.Sprintf("[Room: %d]", m.currentRoom))
    }
    statusStyle := lipgloss.NewStyle().
        Foreground(AccentBlue).
        Background(StatusBarBg).
        Bold(true).
        Width(width).
        Padding(0, 1)
    s.WriteString(statusStyle.Render(statsLine + debugInfo))
    s.WriteString("\n")

    // Input area
    inputStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PrimaryGold).
        Padding(0, 1).
        Width(width).
        Height(inputHeight - 2)
    s.WriteString(inputStyle.Render(promptStyle.Render("> ") + m.textInput.View()))
```

- [ ] **Step 3: Build and verify**

```bash
cd /home/sam/GitHub/herbst-mud && make build-all
```

Expect: clean compile.

- [ ] **Step 4: Commit**

```bash
git add herbst/game_model.go
git commit -m "feat: restructure playing screen with header bar layout"
```

---

### Task 4: Add Onboarding Messages and Polish Auth Messages

**Files:**
- Modify: `herbst/auth.go:1-835`

**Overview:** Add onboarding welcome messages for new characters. Polish the formatting of auth-related messages (login success, world selection, character selection).

- [ ] **Step 1: Create onboarding content function**

Add to auth.go:

```go
// getWelcomeMessage returns an onboarding message for new characters
func (m *model) getWelcomeMessage() string {
    msg := lipgloss.NewStyle().Bold(true).Foreground(PrimaryGold).Render("Welcome to Herbst MUD!")
    msg += "\n\n"
    msg += lipgloss.NewStyle().Foreground(AccentBlue).Render("Essential Commands:")
    msg += "\n"
    commands := []struct{ cmd, desc string }{
        {"look", "Examine your surroundings"},
        {"north / south / east / west", "Move between rooms"},
        {"say <text>", "Speak to others in the room"},
        {"help", "Show all available commands"},
        {"who", "See who's online"},
    }
    for _, c := range commands {
        cmdStyle := lipgloss.NewStyle().Foreground(PrimaryGold).Bold(true).Render(c.cmd)
        msg += fmt.Sprintf("  %s — %s\n", cmdStyle, c.desc)
    }
    msg += "\n"
    msg += lipgloss.NewStyle().Foreground(TextGray).Italic(true).Render("Tip: try 'look' to see where you are!")
    return msg
}
```

- [ ] **Step 2: Add onboarding trigger to loadCharacter**

In `loadCharacter()` (auth.go around line 635), after the character is loaded and screen is set to ScreenPlaying, add the onboarding message. Add a field `hasSeenOnboarding bool` to the model if needed, or check if the character has no visited rooms:

```go
// Add after setting visitedRooms in loadCharacter:
if !m.visitedRooms[m.currentRoom] {
    m.AppendMessage(m.getWelcomeMessage(), "info")
}
```

The visitedRooms map is initialized empty in main.go, and `m.visitedRooms[m.currentRoom] = true` is set in loadCharacter — but it's set AFTER this check would run. So we need to check before the assignment:

In loadCharacter, change this block:

```go
m.visitedRooms[m.currentRoom] = true

// Determine reconnect room
targetRoomID := m.determineReconnectRoom()
m.loadRoom(targetRoomID)
```

To:

```go
// Check if this is a new character (first room visit)
isNewCharacter := !m.visitedRooms[m.currentRoom]
m.visitedRooms[m.currentRoom] = true

// Determine reconnect room
targetRoomID := m.determineReconnectRoom()
m.loadRoom(targetRoomID)

// Show onboarding for new characters
if isNewCharacter {
    m.AppendMessage(m.getWelcomeMessage(), "info")
}
```

- [ ] **Step 3: Polish world selection display**

Update `displayWorlds()` in auth.go to use the new theme colors:

```go
func (m *model) displayWorlds() string {
    var buf bytes.Buffer

    if len(availableWorlds) == 0 {
        buf.WriteString(lipgloss.NewStyle().Foreground(TextGray).Render("Fetching available worlds..."))
        buf.WriteString("\n\n")
    } else {
        for idx, world := range availableWorlds {
            numStyle := lipgloss.NewStyle().Foreground(AccentBlue).Bold(true).Render(fmt.Sprintf("%d.", idx+1))
            nameStyle := lipgloss.NewStyle().Foreground(TextWhite).Render(world)
            line := fmt.Sprintf("  %s  %s", numStyle, nameStyle)
            if world == m.currentWorld {
                line += lipgloss.NewStyle().Foreground(PrimaryGold).Render("  [ACTIVE]")
            }
            buf.WriteString(line)
            buf.WriteString("\n")
        }
        buf.WriteString("\n")
    }

    return buf.String()
}
```

- [ ] **Step 4: Polish character selection display**

Update `displayCharacters()` in auth.go:

```go
func (m *model) displayCharacters() string {
    var buf bytes.Buffer

    worldLabel := lipgloss.NewStyle().Foreground(AccentBlue).Bold(true).Render("World:")
    worldVal := lipgloss.NewStyle().Foreground(TextWhite).Render(m.currentWorld)
    buf.WriteString(fmt.Sprintf("  %s  %s\n\n", worldLabel, worldVal))

    if len(m.selectedWorldCharacters) == 0 {
        buf.WriteString(lipgloss.NewStyle().Foreground(TextGray).Render("  No characters in this world."))
        buf.WriteString("\n")
        buf.WriteString(lipgloss.NewStyle().Foreground(TextGray).Render("  Type 'n' to create a new character."))
        buf.WriteString("\n\n")
    } else {
        for idx, char := range m.selectedWorldCharacters {
            numStyle := lipgloss.NewStyle().Foreground(AccentBlue).Bold(true).Render(fmt.Sprintf("%d.", idx+1))
            nameStyle := lipgloss.NewStyle().Foreground(PrimaryGold).Bold(true).Render(char.Name)
            details := fmt.Sprintf("Lvl %d %s %s", char.Level, char.Race, char.Class)
            detailsStyle := lipgloss.NewStyle().Foreground(TextGray).Render(details)
            hpLabel := lipgloss.NewStyle().Foreground(StatusRed).Render("HP")
            hpStyle := lipgloss.NewStyle().Foreground(TextWhite).Render(fmt.Sprintf("%d/%d", char.Hitpoints, char.MaxHitpoints))
            buf.WriteString(fmt.Sprintf("  %s  %s\n", numStyle, nameStyle))
            buf.WriteString(fmt.Sprintf("       %s  %s  %s\n", detailsStyle, hpLabel+":", hpStyle))
        }
        buf.WriteString("\n")
    }

    return buf.String()
}
```

- [ ] **Step 5: Fix import for lipgloss in auth.go**

auth.go currently only imports `textinput` from bubbles. After the changes above, `displayWorlds()` and `displayCharacters()` reference `lipgloss`. The file doesn't currently import lipgloss directly — it uses style.go's global vars. Wait — actually `displayCharacters()` now uses `lipgloss.NewStyle()` directly. Check if lipgloss is already imported:

Current auth.go imports:
```go
import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "github.com/charmbracelet/bubbles/textinput"
)
```

No lipgloss import. Add it:

```go
import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/lipgloss"
)
```

- [ ] **Step 6: Build and verify**

```bash
cd /home/sam/GitHub/herbst-mud && make build-all
```

Expect: clean compile.

- [ ] **Step 7: Run tests**

```bash
cd /home/sam/GitHub/herbst-mud && make test && cd server && go test -v ./...
```

Expect: all tests pass.

- [ ] **Step 8: Commit**

```bash
git add herbst/auth.go
git commit -m "feat: add onboarding messages and polish auth display"
```

---

### Task 5: Polish UI Messages with New Theme Colors

**Files:**
- Modify: `herbst/ui_messages.go:1-77`

**Overview:** Update the message styling to use the new fantasy theme colors for info/success/error message icons.

- [ ] **Step 1: Update styleMessage colors**

Replace the info/success/error icon calls to use new theme colors. The styleMessage function already uses shared styles (successStyle, errorStyle, infoStyle) — those are defined in style.go and reference StatusGreen, StatusRed, StatusBlue which haven't changed. But update the icon styling:

```go
func styleMessage(msg string, msgType string) string {
    if msg == "" {
        return ""
    }

    switch msgType {
    case "success":
        return lipgloss.NewStyle().Foreground(StatusGreen).Bold(true).Render("◆ ") + msg
    case "error":
        return lipgloss.NewStyle().Foreground(StatusRed).Bold(true).Render("▲ ") + msg
    case "info":
        return lipgloss.NewStyle().Foreground(AccentBlue).Render("● ") + msg
    case "damage":
        return lipgloss.NewStyle().Foreground(StatusRed).Render("⚔ ") + msg
    case "heal":
        return lipgloss.NewStyle().Foreground(StatusGreen).Render("♥ ") + msg
    default:
        return msg
    }
}
```

Note: If `AccentBlue` is not yet imported in ui_messages.go, add lipgloss import. The current file only uses `styleMessage` which references package-level vars. Replace `infoStyle.Render("ℹ ")` with inline lipgloss style using AccentBlue.

- [ ] **Step 2: Ensure lipgloss is imported**

Check the current imports. Add lipgloss if needed:

```go
import (
    "strings"
    "github.com/charmbracelet/lipgloss"
)
```

- [ ] **Step 3: Build and verify**

```bash
cd /home/sam/GitHub/herbst-mud && make build-all
```

Expect: clean compile.

- [ ] **Step 4: Commit**

```bash
git add herbst/ui_messages.go
git commit -m "feat: update message icons to fantasy theme style"
```

---

### Task 6: Final Build, Test & Verify

**Files:** None — verification only.

- [ ] **Step 1: Full build**

```bash
cd /home/sam/GitHub/herbst-mud && make build-all
```

Expect: clean compile, no errors.

- [ ] **Step 2: Run all tests**

```bash
cd /home/sam/GitHub/herbst-mud && make test && cd server && go test -v ./... && cd /home/sam/GitHub/herbst-mud/herbst && go test -v ./...
```

Expect: all tests pass.

- [ ] **Step 3: Verify acceptance criteria**

1. Login screen shows styled username/password prompts — check by reading herbst/ui_screens.go
2. World select lists worlds with numbered selection — check herbst/auth.go displayWorlds()
3. Character select shows characters with stats — check herbst/auth.go displayCharacters()
4. Playing screen has header bar + content + status + input — check herbst/game_model.go View()
5. Fantasy color palette applied consistently — check herbst/style.go colors
6. New characters see onboarding hints — check herbst/auth.go loadCharacter()
7. All existing screens (combat, skill select, profile) still render — check game_model.go View()

- [ ] **Step 4: Commit any remaining changes**

```bash
git add -A
git commit -m "🟣 feat: complete SSH client UI overhaul with fantasy theme"
```
