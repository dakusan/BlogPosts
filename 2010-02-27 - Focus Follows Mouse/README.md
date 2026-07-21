<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Focus_Follows_Mouse">Focus Follows Mouse</a><br>
Date: 2010-02-27<br>
SubTitle: Active Window Tracking in Windows<br>
Tags: Focus, Mouse, Quick Fixes, Windows<br><br>

<p>A friend of mine who mainly works in Linux is always frustrated when he has to deal with working in Windows and doesn’t have FFM (Focus Follows Mouse). FFM means focusing a window when the mouse cursor moves over it (rollover), preferably without raising the window’s z-order position. I told him I’d throw together a program that did this, but my original approach was problematic.</p>
<p>My original approach was to use LowLevelMouseProc (WinAPI:SetWindowsHookEx), which is the same basic approach as my <a href="https://www.castledragmire.com/Projects/HalfKey">HalfKey project</a>. Whenever a non-focused window (WinAPI:GetActiveWindow) is moused over (WinAPI:WindowFromPoint), it would activate it (gain the focus) (WinAPI:SetWindowLong/WinAPI:SetForegroundWindow). Unfortunately, as I found out, Windows is specifically designed to make an activated window go to the top of the z-order, and there is no way around this. The only solution I could find was recording the original position of the activated window (WinAPI:GetParent) and restoring it after the operation. This was however less than optimal 🙁.</p>
<p>After some additional research on the topic, I found out the FFM functionality is actually built into Windows through a little known “System Parameter” (WinAPI:SystemParametersInfo). Microsoft calls it “Active Window Tracking”. Below is the code to the <a href="FocusFollowsMouse.exe">FocusFollowsMouse.exe (command line executable)</a> to activate this system feature (the executable can be run without the command line and it will use the default options).</p><br>

<a href="FocusFollowsMouse.exe">FocusFollowsMouse</a> program information:<br>
<b>Turns on active window tracking in Windows</b><br>
Parameters:
<ul>
	<li>-h The help screen (cancels further actions)</li>
	<li>-i installs new settings permanently</li>
	<li>-d disables active window tracking</li>
	<li>-r raise the window when focusing</li>
</ul><br>

</div>

```c
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
```