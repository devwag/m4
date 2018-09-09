#!/bin/sh
cd /root

service ssh start

mkdir -p /home/LogFiles

file=restart.txt
echo "restart" > $file

echo "\nStarting server ($WEBSITE_ROLE_INSTANCE_ID) ..."

# git the latest code
cd /go/src/m4
git pull
cd /root

while [ -f "$file" ]
do
#  rm $file
  go build m4/sampleapp
  ./app
  sleep .5
  rm -f app
done

