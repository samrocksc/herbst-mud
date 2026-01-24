---
requires_bdd: false
---
# SSH Server Behavior (IMPLEMENTED WITH WISH LIBRARY)

As a player
I want to connect to the MUD via SSH without requiring authentication
So that I can easily access the game

## Implementation Summary

The SSH server has been refactored to use the Wish library which is specifically designed for building terminal applications like MUDs:

- Runs on port 4444 using Wish SSH server
- Allows passwordless connections with no authentication required
- Uses Bubble Tea for the terminal UI rendering
- Displays a welcome message using Bubble Tea on connection
- Supports multiple concurrent connections naturally through Wish

## Test Scenarios

1. Connect to SSH server without authentication
2. See welcome screen after connecting
3. Handle multiple concurrent connections

## Usage

Players can now connect with:

```bash
ssh localhost -p 4444
```

No password is required. The SSH authentication is disabled for ease of access.

## Benefits of Using Wish

- Purpose-built for terminal applications like MUDs
- Secure communication without HTTPS certificates
- Automatic PTY handling and window resizing
- Built-in middlewares for logging and access control
- No risk of accidentally sharing a shell
- Uses modern Go SSH implementation