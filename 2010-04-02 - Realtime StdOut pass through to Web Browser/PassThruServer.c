#include <stdio.h>
#include <malloc.h>
#include <fcntl.h>
#include <sys/types.h> 
#include <sys/socket.h>
#include <netinet/in.h>
#include <signal.h>
#include <unistd.h> //_exit

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
