#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>
#include <string.h>

FILE *SideLog; //This logs all events to the file "log"
FILE *RecordLeft; //This logs the exact packets that the left side sends to the file "left"
int IsChrome=1; //True for roboform chrome compatibility, false to just pass through data as is
int DisplayStatus=1; //Whether to show the "Status: x" updates
int AddRoboformAnswerHacks=0; //Whether to force send fake responses that the rf-chrome server is not responding to
const char* HomePath="/home/USERNAME/.wine/drive_c/Program Files (x86)/Siber Systems/AI RoboForm/9.6.1.1/"; //This must be set.
const char* ExePath="/home/USERNAME/.wine/drive_c/Program Files (x86)/Siber Systems/AI RoboForm/9.6.1.1/rf-chrome-nm-host.exe"; //This must be set. You can make sure rf-chrome-nm-host.exe is running with a 'ps'

//Send back a custom packet to the left. This will show in the log file BEFORE the packet it is responding to.
void SendLeft(char *str)
{
	char Buffer[255];
	int StrSize=strlen(str);
	*(int*)Buffer=StrSize;
	strcpy(Buffer+4, str);

	write(STDOUT_FILENO, Buffer, StrSize+4);
	fprintf(SideLog, "OVERRIDE REQUEST: %s\n\n", Buffer+4);
	fflush(SideLog);
}

//Forward data from left to right or vice versa. Also saves to logs and "left" file
void ForwardSide(int InFileHandle, int OutFileHandle, fd_set* ReadFDS, int IsLeft)
{
	//Exit here if no data
	if(!FD_ISSET(InFileHandle, ReadFDS))
		return;

	//Create a static 1MB+1K buffer - Max packet with chrome is 1MB
	const int BufferLen=1024*(1024+1);
	static char* Buffer=0;
	if(!Buffer)
		Buffer=malloc(BufferLen);

	//If not chrome, just pass the data as is
	char* Side=(IsLeft ? "left" : "right");
	if(!IsChrome) {
		int ReadSize=read(InFileHandle, Buffer, BufferLen);
		write(OutFileHandle, Buffer, ReadSize);
		if(IsLeft)
			fwrite(Buffer, ReadSize, 1, RecordLeft);
		Buffer[ReadSize]=0;
		fprintf(SideLog, "%s (%d): %s\n\n", Side, ReadSize, Buffer);
		fflush(SideLog);
		return;
	}

	//Read the 4 byte packet size and store it at the beginning of the buffer
	unsigned int PacketSize;
	read(InFileHandle, &PacketSize, 4);
	*(unsigned int*)Buffer=PacketSize;

	//Read in the packet and zero it out at the end for the string functions
	read(InFileHandle, Buffer+4, PacketSize);
	Buffer[PacketSize+4]=0;

	//Send fake product-info packet since rf-chrome-nm-host.exe was not responding to it
	if(AddRoboformAnswerHacks && IsLeft && strstr(Buffer+4, "\"name\":\"getProp\"") && strstr(Buffer+4, "\"args\":\"product-info\"")) {
		//The return packet with header at the front
		char VersionInfo[]="{\"callbackId\":\"2\",\"result\":\"{\\\"version\\\":\\\"9-6-1-1\\\",\\\"haveReportAnIssue\\\":true,\\\"haveBreachMon\\\":true,\\\"haveLoginIntoAccount\\\":true}\"}";

		//Simplistic version counter hack since chrome always sends 2 version info requests at the beginning
		static int VersionCount=2;
		VersionInfo[15]='0'+VersionCount;
		VersionCount++;

		SendLeft(VersionInfo);
	//Send fake initialization info packet since rf-chrome-nm-host.exe was not responding to it
	} else if(AddRoboformAnswerHacks && IsLeft && strstr(Buffer+4, "\"name\":\"Initialize2\"")) {
		SendLeft("{\"callbackId\":\"4\",\"result\":\"rf-api\"}");
	//Forward the packet to the other side and store in the "left" file if left side
	} else {
		write(OutFileHandle, Buffer, PacketSize+4);
		if(IsLeft)
			fwrite(Buffer, PacketSize+4, 1, RecordLeft);
	}

	//Output the packet to the log
	fprintf(SideLog, "%s (%d): %s\n\n", Side, PacketSize, Buffer+4);
	fflush(SideLog);
}

int main(void) {
	//Create pipes
	int pipe1[2]; //Parent writes to child
	int pipe2[2]; //Child writes to parent
	if(pipe(pipe1)==-1 || pipe(pipe2)==-1) {
		perror("pipe");
		exit(EXIT_FAILURE);
	}

	//Fork the current process
	pid_t pid = fork();
	if(pid==-1) {
		perror("fork");
		exit(EXIT_FAILURE);
	}

	//New (child) process
	if(pid == 0) {
		//Close unused ends of the pipes
		close(pipe1[1]); // Close write end of pipe1
		close(pipe2[0]); // Close read end of pipe2

		//Redirect stdin to the read end of pipe1
		dup2(pipe1[0], STDIN_FILENO);
		close(pipe1[0]);

		//Redirect stdout to the write end of pipe2
		dup2(pipe2[1], STDOUT_FILENO);
		close(pipe2[1]);

		//Move to the roboform home directory
		if(IsChrome) {
			if(chdir(HomePath) == -1) {
				perror("chdir");
				exit(EXIT_FAILURE);
			}
		}

		//Execute a command that reads from stdin and writes to stdout. The default is the chrome command. If not in chrome, you can fill in the exe and parameter you wish to use
		char *cmd[] = {"/usr/bin/wine", (char*)ExePath, "chrome-extension://pnlccmojcmeohlpggmfnbbiapkmbliob/", "--parent-window=0", NULL};
		if(!IsChrome) {
			cmd[0]="/usr/bin/php";
			cmd[1]="echo.php";
		}
		execvp(cmd[0], cmd);
		perror("execlp");
		exit(EXIT_FAILURE);
	}

	//---Parent process - forwards both sides---
	//Close unused ends of the pipes
	close(pipe1[0]); // Close read end of pipe1
	close(pipe2[1]); // Close write end of pipe2

	//Open the log files
	SideLog = fopen("./log", "w+");
	RecordLeft = fopen("./left", "w+");

	//Run the main loop
	int max_fd=pipe2[0]+1; //Other pipe is STDIN which is 0
	while(1) {
		//Create the structures needed for select
		fd_set read_fds;
		FD_ZERO(&read_fds);
		FD_SET(STDIN_FILENO, &read_fds);
		FD_SET(pipe2[0], &read_fds);

		struct timeval timeout;
		timeout.tv_sec = 10;
		timeout.tv_usec = 0;
	
		//Listen for an update
		int status = select(max_fd, &read_fds, NULL, NULL, &timeout);

		//Display "Status: x" if its setting is true
		if(DisplayStatus) {
			fprintf(SideLog, "Status: %d\n", status);
			if(status==0)
				fflush(SideLog);
		}

		//Exit on bad status
		if (status==-1) {
			perror("select");
			break;
		}

		//Check both sides to see if they need to forward a packet
		ForwardSide(STDIN_FILENO, pipe1[1], &read_fds, 1);
		ForwardSide(pipe2[0], STDOUT_FILENO, &read_fds, 0);
	}

	//Close pipes
	close(pipe1[1]);
	close(pipe2[0]);

	//Wait for the child process to finish
	wait(NULL);

	return EXIT_SUCCESS;
}