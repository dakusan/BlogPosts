#!/bin/bash
#This uses Window's native ping since the Cygwin ping is sorely lacking in options
#"-n 1"=Only runs once, "-w 3000"=Timeout after 3 seconds
#The grep strings are also directly tailored for Window's native ping
PING_COMMAND=$(
	echo 'C:/Windows/System32/PING.EXE -n 1 -w 3000 $DOMAIN |';
	echo 'grep -iP "^(Request timed out|Reply from|Ping request could not find)"';
)

i=0 #The subdomain counter
STARTTIME=`date +%s.%N` #This holds the timestamp of the end of the previous loop

#Infinite loop
while true
do
	#Get the domain to run. This requires a domain that has a wildcard as a primary subdomain
	DOMAIN="$i.castledragmire.com"

	#Output the time, domain name, and ping output
	echo `date +%s` "$DOMAIN" $(eval $PING_COMMAND)

	#If less than a second has passed, sleep up to 1 second
	ENDTIME=`date +%s.%N`
	SLEEPTIME=$(echo "1 - ($ENDTIME - $STARTTIME)" | bc)
	STARTTIME=$ENDTIME
	if [ $(echo "$SLEEPTIME>0" | bc) -eq 1 ]; then
		sleep $SLEEPTIME
		STARTTIME=$(echo "$STARTTIME + $SLEEPTIME" | bc)
	fi

	#Increment the subdomain counter
	let i+=1
done