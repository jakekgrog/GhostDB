GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
NAME=ghostdb
DIR=./cmd/ghostdb

all: build install group

build:
	$(GO_BUILD) -o $(NAME) -v $(DIR)

install:
	$(GO_BUILD) -o /bin/$(NAME) -v ($DIR)
	cp init/ghostdb.service /lib/systemd/system
	systemctl daemon-reload
	systemctl start ghostdb

group:
	/bin/bash -c 'if grep -q "ghostdb" /etc/group; then echo "Ghostdb group already exists!"; else groupadd ghostdb && echo "Ghostdb group created!"; fi'