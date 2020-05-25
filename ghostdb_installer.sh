#!/bin/bash

wget https://ghostdb.github.io/linux/latest
wget https://ghostdb.github.io/linux/ghostdb.service

if [ grep -q "ghostdbservice" /etc/passwd ]; then 
    echo "ghostdbservice user already exists!";
else 
    useradd ghostdbservice -s /sbin/nologin -M && echo "ghostdbservice user created"; 
fi

if [ grep -q "ghostdbservice" /etc/group ]; then
    echo "ghostdbservice group already exists!";
else 
    groupadd ghostdbservice && echo "Ghostdb group created!";
fi

cp ghostdb /bin/
cp ghostdb.service /lib/systemd/system
chmod 755 /lib/systemd/system/ghostdb.service
chown -R ghostdbservice:ghostdbservice /home/ghostdbservice
systemctl daemon-reload
systemctl start ghostdb