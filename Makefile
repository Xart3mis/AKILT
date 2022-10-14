CLIENT_DIR=.\\Client\\

SERVER_DIR=.\\Server\\

CLIENT_BIN=$(CLIENT_DIR)bin\\

SERVER_BIN=$(SERVER_DIR)bin\\

COMPILE=go build

COMPILE_SERVER = $(COMPILE) -ldflags "-w -s"
COMPILE_CLIENT=$(COMPILE) -ldflags "-w -s -H=windowsgui"

.PHONY: all
all: server client

server: .FORCE
	set GOOS=windows
	
	set GOARCH=amd64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-windows-amd64.exe $(SERVER_DIR)
	set GOARCH=386
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-windows-386.exe $(SERVER_DIR)
	set GOARCH=arm64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-windows-arm64.exe $(SERVER_DIR)	
	set GOARCH=arm
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-windows-arm.exe $(SERVER_DIR)
	
	
	set GOOS=darwin
	
	set GOARCH=arm64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-darwin-arm64 $(SERVER_DIR)
	set GOARCH=amd64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-darwin-amd64 $(SERVER_DIR)


	set GOOS=linux
	
	set GOARCH=amd64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-linux-amd64 $(SERVER_DIR)
	set GOARCH=386
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-linux-386 $(SERVER_DIR)
	set GOARCH=arm
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-linux-arm $(SERVER_DIR)
	set GOARCH=arm64
	$(COMPILE_SERVER) -o $(SERVER_BIN)server-linux-arm64 $(SERVER_DIR)
	
	upx $(SERVER_BIN)*.exe

client: .FORCE
	set GOOS=windows

	set GOARCH=amd64
	$(COMPILE_CLIENT) -o $(CLIENT_BIN)client-windows-amd64.exe $(CLIENT_DIR)
	set GOARCH=386
	$(COMPILE_CLIENT) -o $(CLIENT_BIN)client-windows-386.exe $(CLIENT_DIR)
	set GOARCH=arm
	$(COMPILE_CLIENT) -o $(CLIENT_BIN)client-windows-arm.exe $(CLIENT_DIR)
	set GOARCH=arm64
	$(COMPILE_CLIENT) -o $(CLIENT_BIN)client-windows-arm64.exe $(CLIENT_DIR)

	upx $(CLIENT_BIN)*.exe

clean:
	del /Q $(CLIENT_BIN)*
	del /Q $(SERVER_BIN)*

.PHONY: clean
.PHONY: .FORCE