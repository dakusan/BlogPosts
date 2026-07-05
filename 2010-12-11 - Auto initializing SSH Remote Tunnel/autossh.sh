#!/bin/bash
#Variables
USERNAME="user"
PROXYSERVER="192.168.192.168"
PROXYPORT=6822

for ((;;)) #Infinite loop to keep the tunnel open
do
	ssh "$USERNAME@$PROXYSERVER" -o ExitOnForwardFailure=yes -o TCPKeepAlive=yes -o ServerAliveCountMax=2 -o ServerAliveInterval=10 -N -R $PROXYPORT:localhost:22 &>> TUNNEL_LOG #Create the SSH tunnel
	echo "Restarting: " `date` >> TUNNEL_LOG #Write to the log file "TUNNEL_LOG" whenever a restart is attempted
	sleep 1 #Wait 1 second in between connection attempts
done