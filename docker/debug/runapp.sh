#!/bin/sh
service ssh start

# git the latest code
cd /root/m4
git pull
cd /root

while :
do
  go build m4/samplewebhook

  ./samplewebhook -logpath=/home/LogFiles/
  sleep .1
  rm -f samplewebhook
done
