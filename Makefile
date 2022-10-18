CLIENT_DIR=.\\Client\\

SERVER_DIR=.\\Server\\

CLIENT_BIN=$(CLIENT_DIR)bin\\

SERVER_BIN=$(SERVER_DIR)bin\\

COMPILE=go build

COMPILE_SERVER = $(COMPILE) -ldflags "-w -s"
COMPILE_CLIENT=$(COMPILE) -ldflags "-w -s -H=windowsgui"

.PHONY: all
all: clean server client

server: .FORCE
	go env -w GOOS=windows
	
	go env -w GOARCH=amd64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-windows-amd64.exe $(SERVER_DIR)
	go env -w GOARCH=386
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-windows-386.exe $(SERVER_DIR)
	go env -w GOARCH=arm64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-windows-arm64.exe $(SERVER_DIR)	
	go env -w GOARCH=arm
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-windows-arm.exe $(SERVER_DIR)
	
	
	go env -w GOARCH=amd64
	go env -w GOOS=darwin
	
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-darwin-amd64 $(SERVER_DIR)

	go env -w GOARCH=arm64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-darwin-arm64 $(SERVER_DIR)


	go env -w GOOS=linux
	
	go env -w GOARCH=amd64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-linux-amd64 $(SERVER_DIR)
	go env -w GOARCH=386
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-linux-386 $(SERVER_DIR)
	go env -w GOARCH=arm
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-linux-arm $(SERVER_DIR)
	go env -w GOARCH=arm64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-linux-arm64 $(SERVER_DIR)
	
	upx $(SERVER_BIN)server-windows-amd64.exe
	upx $(SERVER_BIN)server-windows-386.exe

client: .FORCE
	go env -w GOOS=windows

	go env -w GOARCH=amd64
	$(COMPILE_CLIENT) -o $(CLIENT_BIN)client-windows-amd64.exe $(CLIENT_DIR)
	
	upx $(CLIENT_BIN)*.exe

clean:
	del /Q $(CLIENT_BIN)*
	del /Q $(SERVER_BIN)*

.PHONY: clean
.PHONY: .FORCE