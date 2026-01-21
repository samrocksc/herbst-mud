package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices []string
	cursor  int
	selected map[int]struct{}
}

func newModel() tea.Model {
	return model{
		choices: []string{"Welcome to MUD!", "Start Game", "Help", "Exit"},
		cursor:  0,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "MUD Server\n\n"
	
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	
	s += "\nPress q to quit.\n"
	return s
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
			}
		}(requests)

		p := tea.NewProgram(newModel(), 
			tea.WithInput(channel),
			tea.WithOutput(channel),
			tea.WithoutSignalHandler(),
		)
		
		if _, err := p.Run(); err != nil {
			log.Printf("Error running bubbletea: %v", err)
		}
		
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

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == "mud" && string(pass) == "mud" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}
	config.AddHostKey(signer)

	listener, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("Failed to listen on :4444: %v", err)
	}
	defer listener.Close()

	fmt.Println("MUD SSH server started on :4444")
	fmt.Println("Connect with: ssh mud@localhost -p 4444")
	fmt.Println("Password: mud")

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
