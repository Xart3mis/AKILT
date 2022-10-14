CLIENT_DIR=.\\Client\\

SERVER_DIR=.\\Server\\

CLIENT_BIN=$(CLIENT_DIR)bin\\

SERVER_BIN=$(SERVER_DIR)bin\\

COMPILE="go build"

COMPILE_SERVER = $(COMPILE) -ldflags \"-w -s\"
COMPILE_CLIENT=$(COMPILE) -ldflags \"-w -s -H=windowsgui\"


all: program

program: server client

server: ./Server/*
	cd $(SERVER_DIR)
	set GOOS=windows; set GOARCH=amd64; go build -ldflags "-w -s" -o $(SERVER_BIN)server-windows-amd64.exe .
	set GOOS=windows; set GOARCH=386; go build -ldflags "-w -s" -o $(SERVER_BIN)server-windows-386.exe .
	set GOOS=windows; set GOARCH=arm64; go build -ldflags "-w -s" -o $(SERVER_BIN)server-windows-arm64.exe .
	set GOOS=windows; set GOARCH=arm; go build -ldflags "-w -s" -o $(SERVER_BIN)server-windows-arm.exe .
	
	set GOOS=darwin; set GOARCH=arm64; go build -ldflags "-w -s" -o $(SERVER_BIN)server-darwin-arm64.exe .
	set GOOS=darwin; set GOARCH=amd64; go build -ldflags "-w -s" -o $(SERVER_BIN)server-darwin-amd64.exe .
	
	set GOOS=linux; set GOARCH=amd64; go build -ldflags "-w -s" -o $(SERVER_BIN)server-linux-amd64.exe .
	set GOOS=linux; set GOARCH=386; go build -ldflags "-w -s" -o $(SERVER_BIN)server-linux-386.exe .
	set GOOS=linux; set GOARCH=arm; go build -ldflags "-w -s" -o $(SERVER_BIN)server-linux-arm.exe .
	set GOOS=linux; set GOARCH=arm64; go build -ldflags "-w -s" -o $(SERVER_BIN)server-linux-arm64.exe .

	cd ..\\

client:
	cd $(CLIENT_DIR)
	set GOOS=windows; set GOARCH=amd64; go build -ldflags "-w -s" -o $(CLIENT_BIN)client-windows-amd64.exe .
	set GOOS=windows; set GOARCH=386; go build -ldflags "-w -s" -o $(CLIENT_BIN)client-windows-386.exe .
	set GOOS=windows; set GOARCH=arm; go build -ldflags "-w -s" -o $(CLIENT_BIN)client-windows-arm.exe .
	set GOOS=windows; set GOARCH=arm64; go build -ldflags "-w -s" -o $(CLIENT_BIN)client-windows-arm64.exe .

clean:
	del /Q $(CLIENT_BIN)*
	del /Q $(SERVER_BIN)*

.PHONY: clean
