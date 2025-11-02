#!/bin/bash

# Test script for MUD server movement

echo "Testing MUD server movement commands..."

# Start the server in the background
cd /home/sam/GitHub/challenges/makeathing
./mudserver &
SERVER_PID=$!

# Wait a moment for the server to start
sleep 2

# Test connecting and sending commands
echo "Testing connection and basic commands..."
timeout 5s expect << 'EOF'
spawn ssh -o StrictHostKeyChecking=no -p 2222 localhost
expect ">"
send "help\r"
expect "Available commands"
expect ">"
send "look\r"
expect "You are in:"
expect ">"
send "north\r"
expect "You cannot go north"
expect ">"
send "up\r"
expect "You move up"
expect "You are now in:"
expect ">"
send "down\r"
expect "You move down"
expect "You are now in:"
expect ">"
send "quit\r"
expect "Goodbye!"
EOF

# Kill the server
kill $SERVER_PID

echo "Test completed."