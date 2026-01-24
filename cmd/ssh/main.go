package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	tea "github.com/charmbracelet/bubbletea"
)

type welcomeModel struct {
	connected   bool
	welcomeMsg  string
	quitting    bool
}

func newWelcomeModel() tea.Model {
	return welcomeModel{
		connected:  true,
		quitting:   false,
		welcomeMsg: `Welcome to MUD - Multi-User Dungeon!

You've successfully connected to the game server.

This is a passwordless SSH connection - your credentials are handled
when you log into the game itself.

Press any key to continue or Ctrl+C to disconnect.

`,
	}
}

func (m welcomeModel) Init() tea.Cmd {
	return nil
}

func (m welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		default:
			// Any key closes welcome and shows status
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m welcomeModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}
	return m.welcomeMsg
}

func handleSession(conn net.Conn, config *ssh.ServerConfig) {
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v", err)
		return
	}
	defer sshConn.Close()
	
	log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				if req.Type == "shell" {
					req.Reply(true, nil)
				}
				if req.Type == "pty-req" {
					req.Reply(true, nil)
				}
			}
		}(requests)

		// Show welcome screen
		p := tea.NewProgram(newWelcomeModel(), 
			tea.WithInput(channel),
			tea.WithOutput(channel),
			tea.WithoutSignalHandler(),
		)
		
		result, err := p.Run()
		if err != nil {
			log.Printf("Error running welcome screen: %v", err)
		}
		
		// Check if user wants to stay connected
		if model, ok := result.(welcomeModel); ok && !model.quitting {
			// User pressed a key other than Ctrl+C, show them they're staying connected
			fmt.Fprintf(channel, "\nThank you for connecting to MUD!\n")
			fmt.Fprintf(channel, "Game login system coming soon...\n") 
			fmt.Fprintf(channel, "Stay connected or press Ctrl+C to disconnect.\n")
			fmt.Fprintf(channel, "Server will keep this session open.\n")
			
			// Keep the session alive for a bit
			for i := 0; i < 3; i++ {
				time.Sleep(10 * time.Second)
				fmt.Fprintf(channel, ".")
			}
		}
		
		fmt.Fprintf(channel, "\nSession ending. Goodbye!\n")
		time.Sleep(1 * time.Second)
		channel.Close()
	}
}

func loadHostKey() (ssh.Signer, error) {
	keyData, err := ioutil.ReadFile("host_key")
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(keyData)
}

func main() {
	signer, err := loadHostKey()
	if err != nil {
		log.Fatalf("Failed to load host key: %v", err)
	}

	// Configure for passwordless authentication
	config := &ssh.ServerConfig{
		NoClientAuthCallback: func(ctx ssh.ConnMetadata) (*ssh.Permissions, error) {
			// Accept all connections without authentication
			return nil, nil
		},
	}
	config.AddHostKey(signer)

	listener, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("Failed to listen on :4444: %v", err)
	}
	defer listener.Close()

	fmt.Println("MUD SSH server started on :4444")
	fmt.Println("Connect with: ssh localhost -p 4444")
	fmt.Println("(No password required - SSH authentication is disabled)")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
			go handleSession(conn, config)
		}
	}()

	<-sigChan
	fmt.Println("\nShutting down server...")
}
