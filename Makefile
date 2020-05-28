GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
NAME=ghostdb
DIR=./cmd/ghostdb
CONF_DIR=/etc/ghostdb

all: build install group

build:
	$(GO_BUILD) -o $(NAME) -v $(DIR)

install-dev:
	$(GO_BUILD) -o $(NAME) -v $(DIR)
	/bin/bash -c 'if grep -q "$(NAME)" /etc/passwd; then echo "ghostdbservice user already exists!"; else useradd $(NAME) -s /sbin/nologin -M && echo "ghostdbservice user created"; fi'
	
	/bin/bash -c 'if grep -q "$(NAME)" /etc/group; then echo "Ghostdb group already exists!"; else groupadd ghostdb && echo "Ghostdb group created!"; fi'
	
	/bin/bash -c 'if [ -d "$(CONF_DIR)" ]; then echo "ghostdb config directory already exists!"; else mkdir /etc/ghostdb && echo "ghostdb config direcotry created"; chown -R ghostdbservice:ghostdbservice /etc/ghostdb; fi'
	
	sudo cp $(NAME) /bin/
	sudo cp init/ghostdb.service /lib/systemd/system
	sudo cp config/ghostdbConf.json $(CONF_DIR)
	
	sudo chmod 755 /bin/$(NAME)
	sudo chmod 755 /lib/systemd/system/ghostdb.service
	sudo chmod 755 $(CONF_DIR)/ghostdbConf.json
	sudo chown -R ghostdbservice:ghostdbservice /home/ghostdbservice
	systemctl daemon-reload
	systemctl start ghostdb
