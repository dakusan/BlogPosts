<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Auto_initializing_SSH_Remote_Tunnel">Auto initializing SSH Remote Tunnel</a><br>
Date: 2010-12-11<br>
SubTitle: SSH is so incredibly useful<br>
Tags: Examples, Networking, OpenSSH, Port Forwarding, Remote Control, Servers, SSH, SSH Tunnels, Tunneling<br><br>

<p>I often find SSH tunnels absolutely indispensable in my line of work for multiple reasons including secure proxies (tunneling) over insecure connections and connecting computers and programs together over difficult network setups involving NATs and firewalls.</p>
<p>One such example of this I ran into recently is that I have a server machine (hereby called “the client”) that I wanted to make sure I have access to no matter where it is. For this I created an auto initializing SSH remote port tunnel to a server with a static IP Address (hereby called “the proxy server”) which attempts to keep itself open when there is problems.</p>
<p>The first step of this was to create the following bash script on the client that utilizes the OpenSSH’s client to connect to an OpenSSH server on the proxy server for tunneling:</p>
</div>

```bash
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
```

<div>

The parts of the command that create the SSH tunnel are as follows:
<table class=SSHTunnelTable>
	<tr><th>Part of Command										</th><th>Description</th></tr>
	<tr><td>ssh													</td><td>The OpenSSH client application</td></tr>
	<tr><td><kbd>USERNAME</kbd>@<kbd>PROXYSERVER</kbd>			</td><td>The proxy server and username on said server to connect to</td></tr>
	<tr><td>-o ExitOnForwardFailure=yes							</td><td>Automatically terminate the SSH session if the remote port forward fails</td></tr>
	<tr><td>-o TCPKeepAlive=yes<br>-o ServerAliveCountMax=2<br>-o ServerAliveInterval=10		</td><td>Make sure the SSH connection is working, and if not, reinitialize it. The connection fails if server keepalive packets that are sent every 10 seconds are not received twice in a row, or if TCP protocol keepalive fails</td></tr>
	<tr><td>-N													</td><td>“Do not execute a remote command.  This is useful for just forwarding ports” (This means no interactive shell is run)</td></tr>
	<tr><td>-R <kbd>PROXYPORT</kbd>:localhost:22				</td><td>Establish a port of <kbd>PROXYPORT</kbd> on the proxy server that sends all data to port 22 (ssh) on the client (localhost)</td></tr>
	<tr><td>&>> TUNNEL_LOG										</td><td>Send all output from both stdin and stderr to the log file “TUNNEL_LOG”</td></tr>
</table><br>

<p>For this to work, you will need to set up <a href="https://www.castledragmire.com/Posts/OpenSSH_RSA_Authentication_public_key_file_format">public key</a> authentication between the client and utilized user on the proxy server. To do this, you have to run “ssh-keygen” on the client. When it has finished, you must copy the contents of your public key file (most likely at “~/.ssh/id_rsa.pub”) to the “~/.ssh/authorized_keys” file on the user account you are logging into on the proxy server. You will also have to log into this account at least once from the client so the proxy server’s information is in the client’s “known_hosts” file.</p>
<p>For security reasons you may want the user on the proxy server to have a <i>nologin</i> shell set since it is only being used for tunneling, and the above mentioned public key will allow access without a password.</p>
<p>For network access reasons, it might also be worth it to have the proxy port and the ssh port on the proxy server set to commonly accessible ports on all network setups (that a firewall wouldn’t block). Or you may <b>NOT</b> want to have it on common ports for other security reasons 🙂.</p>
<p>If you want the proxy port on the proxy server accessible from other computers (not only the proxy server), you will have to turn on “GatewayPorts” (set to “yes”) in the proxy server’s sshd config, most likely located at “/etc/ssh/sshd_config”.</p><br>

<p>The next step is to create a daemon that calls this script. The normal method for this in Linux would be to use <i>inittab</i>. This can be a bit cumbersome though with Linux distributions that use upstart, like Ubuntu, so I just cheated in it and created the following script to initialize a daemon through rc.d:</p>
</div>

Fill in variables: <kbd>USERNAME</kbd>, <kbd>PROXYSERVER</kbd>, <kbd>PROXYPORT</kbd>, <kbd>TUNNEL_SCRIPT_PATH</kbd>, <kbd>TUNNEL_SCRIPT_DIRECTORY</kbd>, <kbd>TUNNEL_SCRIPT_EXECUTABLE</kbd>, <kbd>TUNNEL_SCRIPT_SERVICE_NAME</kbd>
```bash
#!/bin/bash
echo -e '#!/bin/bash\nfor ((;;))\ndo\n  ssh USERNAME@PROXYSERVER -o TCPKeepAlive=yes -o ExitOnForwardFailure=yes -o ServerAliveCountMax=2 -o ServerAliveInterval=10 -N -R PROXYPORT:localhost:22 &>> TUNNEL_LOG\n  echo "Restarting: " `date` >> TUNNEL_LOG\n  sleep 1\ndone' > TUNNEL_SCRIPT_PATH #This creates the above script
echo -e '#!/bin/bash\ncd TUNNEL_SCRIPT_DIRECTORY\n./TUNNEL_SCRIPT_EXECUTABLE &' > /etc/init.d/TUNNEL_SCRIPT_SERVICE_NAME #This creates the init.d daemon script. I have the script set the working path before running the executable so the log file stays in the same directory
chmod u+x TUNNEL_SCRIPT_PATH /etc/init.d/TUNNEL_SCRIPT_SERVICE_NAME #Set all scripts as executable
update-rc.d TUNNEL_SCRIPT_SERVICE_NAME defaults #Set the run levels at which the script runs
```

<div>
</div>