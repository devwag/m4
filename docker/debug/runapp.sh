#!/bin/sh
cd /root

service ssh start

file=restart.txt
echo "restart" > $file

echo "\nStarting server ($WEBSITE_ROLE_INSTANCE_ID) ..."

# git the latest code
cd /root/m4
git pull
cd /root

while [ -f "$file" ]
do
#  rm $file
  go build m4/samplewebhook

  ./samplewebhook logpath=${logpath}
  sleep .1
  rm -f app
done

