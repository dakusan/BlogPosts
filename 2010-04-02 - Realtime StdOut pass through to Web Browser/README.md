<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Realtime_StdOut_pass_through_to_Web_Browser">Realtime StdOut pass through to Web Browser</a><br>
Date: 2010-04-02<br>
SubTitle: Tying it all together<br>
Tags: C, Command Line Utility, Data Passing, Examples, Flash, JavaScript, Perl, Realtime, Web browsers<br><br>

<p>I had the need to pass a program’s [standard] output to a web browser in real time. The best solution for this is to use a combination of programs made in different languages. The following are all of these individual components to accomplish this task.</p>
<p>Please note the C components are only compatible with gcc and bash (<a href="http://www.cygwin.com/">cygwin</a> required for Windows), as MSVC and Windows command prompt are missing vital functionality for this to work.</p>

<hr>

<p>The first component is a server made in C that receives stdin (as a pipe, or typed by the user after line breaks) and sends that data out to a connected client (buffering the output until the client connects).</p>
<p><a href="PassThruServer.c">PassThruServer source</a>, <a href="PassThruServer.exe">PassThruServer compiled Windows executable</a>.</p><br>

Compilation notes:
<ul>
	<li>This compiles as C99 under gcc: <code>gcc PassThruServer.c -o PassThruServer</code></li>
	<li>Define “WINDOWS” when compiling in Windows (pass “-DWINDOWS”)</li>
</ul><br>

Source Code:
</div>

```c
#include <stdio.h>
#include <malloc.h>
#include <fcntl.h>
#include <sys/types.h> 
#include <sys/socket.h>
#include <netinet/in.h>
#include <signal.h>

//The server socket and options
int ServerSocket=0;
const int PortNumber=1234; //The port number to listen in on

//If an error occurs, exit cleanly
int error(char *msg)
{
	//Close the socket if it is still open
	if(ServerSocket)
		close(ServerSocket);
	ServerSocket=0;

	//Output the error message, and return the exit status
	fprintf(stderr, "%s\n", msg);
	return 1;
}

//Termination signals
void TerminationSignal(int sig)
{
	error("SIGNAL causing end of process");
	_exit(sig);
}

int main(int argc, char *argv[])
{
	//Listen for termination signals
	signal(SIGINT, TerminationSignal);
	signal(SIGTERM, TerminationSignal);
	signal(SIGHUP, SIG_IGN); //We want the server to continue running if the environment is closed, so SIGHUP is ignored -- This doesn't work in Windows
	
	//Create the server
	struct sockaddr_in ServerAddr={AF_INET, htons(PortNumber), INADDR_ANY, 0}; //Address/port to listen on
	if((ServerSocket=socket(AF_INET, SOCK_STREAM, 0))<0) //Attempt to create the socket
		return error("ERROR on 'socket' call");
	if(bind(ServerSocket, (struct sockaddr*)&ServerAddr, sizeof(ServerAddr))<0) //Bind the socket to the requested address/port
		return error("ERROR on 'bind' call");
	if(listen(ServerSocket,5)<0) //Attempt to listen on the requested address/port
		return error("ERROR on 'listen' call");

	//Accept a connection from a client
	struct sockaddr_in ClientAddr;
	int ClientAddrLen=sizeof(ClientAddr);
	int ClientSocket=accept(ServerSocket, (struct sockaddr*)&ClientAddr, &ClientAddrLen);
	if(ClientSocket<0) 
		return error("ERROR on 'accept' call");

	//Prepare to receive info from STDIN
		//Create the buffer
		const int BufferSize=1024*10;
		char *Buffer=malloc(BufferSize); //Allocate a 10k buffer
		//STDIN only needs to be set to binary mode in windows
		const int STDINno=fileno(stdin);
		#ifdef WINDOWS
			_setmode(STDINno, _O_BINARY);
		#endif
		//Prepare for blocked listening (select function)
		fcntl(STDINno, F_SETFL, fcntl(STDINno, F_GETFL, 0)|O_NONBLOCK); //Set STDIN as blocking
		fd_set WaitForSTDIN;
		FD_ZERO(&WaitForSTDIN);
		FD_SET(STDINno, &WaitForSTDIN);

	//Receive information from STDIN, and pass directly to the client
	int RetVal=0;
	while(1)
	{
		//Get the next block of data from STDIN
		select(STDINno+1, &WaitForSTDIN, NULL, NULL, NULL); //Wait for data
		size_t AmountRead=fread(Buffer, 1, BufferSize, stdin); //Read the data
		if(feof(stdin) || AmountRead==0) //If input is closed, process is complete
			break;
		
		//Send the data to the client
		if(write(ClientSocket,Buffer,AmountRead)<0) //If error in network connection occurred
		{
			RetVal=error("ERROR on 'write' call");
			break;
		}
	}
	
	//Cleanup
	if(ServerSocket)
		close(ServerSocket);
	free(Buffer);
	
	return RetVal;
}
```

<div>

<hr>

<p>The next component is a Flash applet as the client to receive data. Flash is needed as it can keep a socket open for realtime communication. The applet receives the data and then passes it through to JavaScript for final processing.</p>
<p><a href="client.swf">Compiled Flash Client Applet</a></p><br>

ActionScript 3.0 Code (This goes in frame 1)
</div>

```as
import flash.external.ExternalInterface;
import flash.events.Event;
ExternalInterface.addCallback("OpenSocket", OpenSocket);

function OpenSocket(IP:String, Port:Number):void
{
	SendInfoToJS("Trying to connect");
	var TheSocket:Socket = new Socket();
	TheSocket.addEventListener(Event.CONNECT, function(Success) { SendInfoToJS(Success ? "Connected!" : "Could not connect"); });
	TheSocket.addEventListener(Event.CLOSE, function() { SendInfoToJS("Connection Closed"); });
	TheSocket.addEventListener(IOErrorEvent.IO_ERROR, function() {SendInfoToJS("Could not connect");});
	TheSocket.addEventListener(ProgressEvent.SOCKET_DATA, function(event:ProgressEvent):void { ExternalInterface.call("GetPacket", TheSocket.readUTFBytes(TheSocket.bytesAvailable)); });
	TheSocket.connect(IP, Port);
}
function SendInfoToJS(str:String) { ExternalInterface.call("GetInfoFromFlash", str); }
stop();
```

<div>
<br>

Flash sockets can also be implemented in ActionScript 1.0 Code (I did not include hooking up ActionScript 1.0 with JavaScript in this example. “GetPacket” and “SendInfoToJS” need to be implemented separately. “IP” and “Port” need to also be received separately).
</div>

```as
var NewSock=new XMLSocket();
NewSock.onData=function(msg) { GetPacket(msg); }
NewSock.onConnect=function(Success) { SendInfoToJS(Success ? "Connected!" : "Could not connect"); }
SendInfoToJS(NewSock.connect(IP, Port) ? "Trying to Connect" : "Could not start connecting");
```

<div>

<hr>

<p>JavaScript can then receive (and send) information from (and to) the Flash applet through the following functions.</p>
<ul>
	<li><b>FLASH</b>.<i>OpenSocket</i>(<i>String</i> IP, <i>Number</i> Port): Call this from JavaScript to open a connection to a server. Note the IP <b>MIGHT</b> have to be the domain the script is running on for security errors to not be thrown.</li>
	<li><b>JAVASCRIPT</b>.<i>GetInfoFromFlash</i>(<i>String</i>): This is called from Flash whenever connection information is updated. I have it giving arbitrary strings ATM.</li>
	<li><b>JAVASCRIPT</b>.<i>GetPacket</i>(<i>String</i>): This is called from Flash whenever data is received through the connection.</li>
</ul><br>
<p>This example allows the user to input the IP to connect to that is streaming the output. Connection information is shown in the “ConnectionInfo” DOM object. Received data packets are appended to the document in separate DOM objects.</p>

Source Code: (See <a href="client.html">JavaScript+HTML Source file</a> for all code)
</div>

```js
var isIE=navigator.appName.indexOf("Microsoft")!=-1;
function getFlashMovie(movieName) { return (isIE ? window[movieName] : document[movieName]);  }
function $(s) { return document.getElementById(s); }

function Connect()
{
	getFlashMovie("client").OpenSocket($('IP').value, 1234);
}

function GetInfoFromFlash(Str)
{
	$('ConnectionInfo').firstChild.data=Str;
}

function GetPacket(Str)
{
	var NewDiv=document.createElement('DIV');
	NewDiv.appendChild(document.createTextNode(Str));
	$('Info').appendChild(NewDiv);
}
```

<div>

<hr>

<p>Next is an example application that outputs to stdout. It is important that it flushes stdout after every  output or the communication may not be real time.</p>
<p><a href="inc.c">inc source</a>, <a href="inc.exe">inc compiled Windows executable</a>.</p><br>

<p><b>inc</b> counts from 0 to one less than a number (parameter #1 [default=50]) after a certain millisecond interval (parameter #2 [default=500]).</p>
[Bash] Example:
</div>

```bash
./inc 10 #Counts from 0-9 every half a second
```

<div>
<br>

Source Code:
</div>

```c
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[])
{
	int NumLoops=(argc>1 ? atoi(argv[1]) : 50); //Number of loops to run from passed argument 1. Default is 50 if not specified.
	int LoopWait=(argc>2 ? atoi(argv[2]) : 500); //Number of milliseconds to wait in between each loop from passed argument 2. Default is 500ms if not specified.
	LoopWait*=1000; //Convert to microseconds for usleep

	//Output an incremented number every half a second
	int i=0;
	while(i<NumLoops)
	{
		printf("%u\n", i++);
		fflush(stdout); //Force stdout flush
		usleep(LoopWait); //Wait for half a second
	};
	
	return 0;
}
```

<div>

<hr>

<p>This final component is needed so the Flash applet can connect to a server. Unfortunately, new versions of Flash (at least version 10, might have been before that though) started requiring policies for socket connections >:-(. I don’t think this is a problem if you compile your applet to target an older version of Flash with the ActionScript v1.0 code.</p>
<p>This Perl script creates a server on port 843 to respond to Flash policy requests, telling any Flash applet from any domain to allow connections to go through to any port on the computer (IP). It requires Perl, and root privileges on Linux to bind to a port <1024 (<b>su</b> to root or run with <b>sudo</b>).</p>
<p><a href="FlashSocketPolicy.pl">Flash Socket Policy Server</a></p><br>

Source Code:
</div>

```perl
#!/usr/bin/perl
use warnings;
use strict;

#Listen for kill signals
$SIG{'QUIT'}=$SIG{'INT'}=$SIG{__DIE__} = sub
{
	close Server;
	print "Socket Policy Server Ended: $_[0]\n";
	exit;
};

#Start the server:
use Socket;
use IO::Handle;
my $FlashPolicyPort=843;
socket(Server, PF_INET, SOCK_STREAM, getprotobyname('tcp')) or die "'socket' call: $!"; #Open the socket
setsockopt(Server, SOL_SOCKET, SO_REUSEADDR, 1) or die "'setsockopt' call: $!"; #Allow reusing of port/address if in TIME_WAIT state
bind(Server, sockaddr_in($FlashPolicyPort,INADDR_ANY)) or die "'bind' call: $!"; #Listen on port $FlashPolicyPort for connections from any INET adapter
listen(Server,SOMAXCONN) or die "'listen' call: $!"; #Start listening for connections
Server->autoflush(1); #Do not buffer output

#Infinite loop that accepts connections
$/ = "\0"; #Reset terminator from new line to null char
while(my $paddr=accept(Client,Server))
{
	Client->autoflush(1); #Do not buffer IO
	if(<Client> =~ /.*policy\-file.*/i) { #If client requests policy file...
		print Client '<cross-domain-policy><allow-access-from domain="*" to-ports="*" /></cross-domain-policy>'.$/; #Output policy info: Allow any flash applets from any domain to connect
	}
	close Client; #Close the client
}
```

<div>

This could very easily be converted to another better [less resource intensive] language too.

<hr>

<b><u>How to tie all of this together</u></b>
<ol>
	<li>	Start the servers
		<ul>
			<li>
				In your [bash] command shell, execute the following
<pre><code>Server/FlashSocketPolicy.pl & #Run the Flash Policy Server as a daemon. Don't forget <b>sudo</b> in Linux
./inc | ./PassThruServer #Pipe <b>inc</b> out to the <b>PassThruServer</b></code></pre>
			</li>
			<li>Note that this will immediately start the PassThruServer receiving information from “inc”, so if you don’t get the client up in time, it may already be done counting and send you all the info at once (25 seconds).</li>
			<li>
				The PassThruServer will not end until one of the following conditions has been met:
				<ul>
					<li>The client has connected and the piped process is completed</li>
					<li>The client has connected and disconnected and the disconnect has been detected (when a packet send failed)</li>
					<li>It is manually killed through a signal</li>
				</ul>
			</li>
			<li>The Flash Policy Server daemon should probably just be left on indefinitely in the background (it only needs to be run once).</li>
		</ul>
	</li>
	<li>To run the client, open <a href="client.html">client.html</a> through a web server [i.e. <a href="http://httpd.apache.org/">Apache’s httpd</a>] in your web browser. <b>Don’t open the local file straight through your file system, it needs to be run through a web server for Flash to work correctly.</b></li>
	<li>Click “connect” (assuming you are running the PassThruServer already on localhost [the same computer]). You can click “connect” again every time a new PassThruServer is ran.</li>
</ol></div>