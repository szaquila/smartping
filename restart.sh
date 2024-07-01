#!/bin/bash

# ps aux|grep bin/smartping|grep -v grep|cut -c 9-15|xargs kill -s 9
pkill -9 smartping;sleep 1s
/var/www/html/smartping/start.sh

