#!/bin/bash
#Fill in variables: USERNAME, PROXYSERVER, PROXYPORT, TUNNEL_SCRIPT_PATH, TUNNEL_SCRIPT_DIRECTORY, TUNNEL_SCRIPT_EXECUTABLE, TUNNEL_SCRIPT_SERVICE_NAME
echo -e '#!/bin/bash\nfor ((;;))\ndo\n  ssh USERNAME@PROXYSERVER -o TCPKeepAlive=yes -o ExitOnForwardFailure=yes -o ServerAliveCountMax=2 -o ServerAliveInterval=10 -N -R PROXYPORT:localhost:22 &>> TUNNEL_LOG\n  echo "Restarting: " `date` >> TUNNEL_LOG\n  sleep 1\ndone' > TUNNEL_SCRIPT_PATH #This creates the above script
echo -e '#!/bin/bash\ncd TUNNEL_SCRIPT_DIRECTORY\n./TUNNEL_SCRIPT_EXECUTABLE &' > /etc/init.d/TUNNEL_SCRIPT_SERVICE_NAME #This creates the init.d daemon script. I have the script set the working path before running the executable so the log file stays in the same directory
chmod u+x TUNNEL_SCRIPT_PATH /etc/init.d/TUNNEL_SCRIPT_SERVICE_NAME #Set all scripts as executable
update-rc.d TUNNEL_SCRIPT_SERVICE_NAME defaults #Set the run levels at which the script runs