# Copyright (c) 2020, Jake Grogan
# All rights reserved.
# 
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are met:
# 
#  * Redistributions of source code must retain the above copyright notice, this
#    list of conditions and the following disclaimer.
# 
#  * Redistributions in binary form must reproduce the above copyright notice,
#    this list of conditions and the following disclaimer in the documentation
#    and/or other materials provided with the distribution.
# 
#  * Neither the name of the copyright holder nor the names of its
#    contributors may be used to endorse or promote products derived from
#    this software without specific prior written permission.
# 
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
# AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
# IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
# DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
# FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
# DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
# SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
# CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
# OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
# OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
NAME=ghostdb
NAME_WIN=ghostdb.exe
DIR=./cmd
CONF_DIR=/etc/ghostdb

all: build install group

build:
	$(GO_BUILD) -o $(NAME) -v $(DIR)

build-win:
	$(GO_BUILD) -o $(NAME_WIN) -v $(DIR)

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
	sudo rm $(NAME)
