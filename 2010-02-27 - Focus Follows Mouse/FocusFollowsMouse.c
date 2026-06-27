//Toggles active window tracking, with options
#include <windows.h>
#include <stdio.h>

int main(int argc, char *argv[])
{
	BOOL SendUpdate=SPIF_SENDCHANGE; //How to send the update to the system. Default is to only send to the system, but not make the change permanent
	BOOL TurnOn=TRUE; //By default we are turning on the option
	BOOL RaiseWindow=FALSE; //Whether to raise the window that is gaining the focus
	int i;

	//Read in options
	for(i=1;i<argc;i++)
	{
		if(argv[i][0]!='-') //Parameters must start with a -
			continue;
		switch(argv[i][1]) //Execute a valid action on the first character. Characters after it are ignored
		{
			case 'h': //Help  screen
				printf("%s", "Turns on active window tracking in Windows\n"
					"-h this help screen (cancels further actions)\n"
					"-i installs new settings permanently\n"
					"-d disables active window tracking\n"
					"-r raise the window when focusing\n");
					return 0; //Cancel further actions
			case 'i': //Install permanently
				SendUpdate|=SPIF_UPDATEINIFILE; //Writes the new system-wide parameter setting to the user profile.
				break;
			case 'd': //Turn off
				TurnOn=FALSE;
				break;
			case 'r': //Raise the window
				RaiseWindow=TRUE;
				break;
		}
	}

	//Execute the system parameters change
	SystemParametersInfo(SPI_SETACTIVEWINDOWTRACKING, 0, (PVOID)TurnOn, SendUpdate);
	SystemParametersInfo(SPI_SETACTIVEWNDTRKZORDER, 0, (PVOID)RaiseWindow, SendUpdate);

	return 0;
}