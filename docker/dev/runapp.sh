#!/bin/sh
cd /root

service ssh start

mkdir -p /home/LogFiles

file=restart.txt
echo "restart" > $file

while [ -f "$file" ]
do
  echo "\nStarting server ($WEBSITE_ROLE_INSTANCE_ID) ...\n"
  rm $file
  cd /go/src/m4
  git pull
  cd /root
  go build m4/app
  ./app
done

