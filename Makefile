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
	$(GO_BUILD) -o $(NAME) -v $(DIR)
	/bin/bash -c 'if grep -q "$(NAME)" /etc/passwd; then echo "ghostdbservice user already exists!"; else useradd $(NAME) -s /sbin/nologin -M && echo "ghostdbservice user created"; fi'
	/bin/bash -c 'if grep -q "ghostdb" /etc/group; then echo "Ghostdb group already exists!"; else groupadd ghostdb && echo "Ghostdb group created!"; fi'
	sudo cp $(NAME) /bin/
	sudo cp init/ghostdb.service /lib/systemd/system
	sudo chmod 755 /lib/systemd/system/ghostdb.service
	sudo chown -R ghostdbservice:ghostdbservice /home/ghostdbservice
	systemctl daemon-reload
	systemctl start ghostdb