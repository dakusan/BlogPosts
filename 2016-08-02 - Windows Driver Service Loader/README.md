<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Windows_Driver_Service_Loader">Windows Driver Service Loader</a><br>
Date: 2016-08-02<br>
Tags: C, Command Line Utility, Driver, Kernel, Loader, Operating Systems, Service, Windows<br><br>

<p>Following is some C++ source code for a Windows kernel-driver service loader. It could be used to load other service types too by changing the <i>dwServiceType</i> flag on the <i>CreateService</i> call. I threw this together for another project I am currently working on. It is also used in the following post (posting soon).</p>
It works in the following way:<ul>
	<li>It is a command line utility which takes 3 arguments:<ol>
		<li>The service name. Hereby referred to as <kbd>SERVICE_NAME</kbd></li>
		<li>The service display name. Hereby referred to as <kbd>DISPLAY_NAME</kbd></li>
		<li>The driver path (to the .sys file). Hereby referred to as <kbd>DRIVER_PATH</kbd></li>
	</ol></li>
	<li>This program (most likely) requires administrative access. There are also some caveats regarding driver code signing requirements that are thoroughly explored elsewhere.</li>
	<li>It first checks to see if a service already exists with the given <kbd>SERVICE_NAME</kbd>. If it does:<ol>
		<li>If the <kbd>DISPLAY_NAME</kbd> matches, the service is kept as is.</li>
		<li>If the <kbd>DISPLAY_NAME</kbd> does not match, the user is prompted on if they want to delete the current service. If they do not, the program exits.</li>
	</ol></li>
	<li>If the service needs to be created (it did not already exist or was deleted), it creates the service with the given <kbd>SERVICE_NAME</kbd>, <kbd>DISPLAY_NAME</kbd>, and <kbd>DRIVER_PATH</kbd>. If the service is not created during this run, the <kbd>DRIVER_PATH</kbd> is ignored.<br>Note: The <kbd>DRIVER_PATH</kbd> must be to a direct local file system file. I have found that network links and symbolic links do not work.</li>
	<li>The service is started up:<ul>
		<li>If it is already running, the user is prompted on if they want to stop the currently running service. If they say no, the program exits.</li>
	</ul></li>
	<li>The program then waits for a final user input on if they want to close the service before exiting the program.</li>
	<li>If there was an error, the program reports the error, otherwise, it reports “Success”.</li>
	<li>The program pauses at the end until the user presses any key to exit.</li>
	<li>The program returns 0 on success, and 1 if an error occurred.</li>
</ul></div>

```cpp
//Compiler flags
#define WIN32_LEAN_AND_MEAN  //Include minimum amount of windows stuff
#ifndef _UNICODE //Everything in this script is unicode
	#define _UNICODE
#endif

#include <windows.h>
#include <stdio.h>
#include <conio.h>
#include <memory>

//Smart pointers
typedef std::unique_ptr<WCHAR, void(*)(WCHAR*)> SmartWinAlloc;
typedef std::unique_ptr<SC_HANDLE, void(*)(SC_HANDLE*)> SmartCloseService;
void Delete_SmartWinAlloc(WCHAR *p) { if(p) LocalFree(p); }
void Delete_SmartCloseService(SC_HANDLE *h) { if(h && *h) CloseServiceHandle(*h); }

//Function declarations
WCHAR* InitDriver(int argc, WCHAR *argv[]);
WCHAR* FormatError(WCHAR* Format, ...);
SmartWinAlloc GetLastErrorStr();
BOOLEAN AskQuestion(WCHAR* Question); //Returns if user answered yes

int wmain(int argc, WCHAR *argv[])
{
	//Run the init routine
	WCHAR* Ret=InitDriver(argc, argv);

	//If there is an error, report it, or otherwise, report success
	wprintf(L"%s\n", Ret ? Ret : L"Success");
	wprintf(L"%s\n", L"Press any key to exit");
	_getch();

	//Return if successful
	return (Ret ? 1 : 0);
}

WCHAR* InitDriver(int argc, WCHAR *argv[])
{
	//Confirm arguments
	if(argc<4)
		return FormatError(L"%s", L"3 arguments are required: Service Name, Display Name, Driver Path");
	const WCHAR* Param_ServiceName=argv[1];
	const WCHAR* Param_DisplayName=argv[2];
	const WCHAR* Param_DriverPath =argv[3];

	//Open the service manager
	wprintf(L"%s\n", L"Opening the service manager");
	SC_HANDLE HSCManager=OpenSCManager(nullptr, nullptr, SC_MANAGER_CREATE_SERVICE);
	if(!HSCManager)
		return FormatError(L"%s: %s", L"Error opening service manager", GetLastErrorStr());
	SmartCloseService FreeHSCManager(&HSCManager, Delete_SmartCloseService);

	//Check if the service already exists
	wprintf(L"%s\n", L"Checking previously existing service state");
	BOOL ServiceExists=false;
	{
		//Get the service name
		const DWORD NameBufferSize=255;
		WCHAR NameBuffer[NameBufferSize];
		WCHAR *NamePointer=NameBuffer;
		DWORD NamePointerSize=NameBufferSize;
		std::unique_ptr<WCHAR> Buf(nullptr); //May be swapped with a real pointer later
		for(INT_PTR i=0;i<2;i++)
		{
			//If we found the service, exit the lookup here
			if(GetServiceDisplayName(HSCManager, Param_ServiceName, NamePointer, &NamePointerSize))
			{
				ServiceExists=true;
				break;
			}

			//If the service does not exist, we can exit the lookup here
			if(GetLastError()==ERROR_SERVICE_DOES_NOT_EXIST)
				break;

			//If error is not insufficient buffer size, return the error
			if(GetLastError()!=ERROR_INSUFFICIENT_BUFFER)
				return FormatError(L"%s: %s", L"Could not query service information", GetLastErrorStr());

			//If second pass, error out
			if(i==1)
				return FormatError(L"%s: %s", L"Could not query service information", L"Second buffer pass failed");

			//Create a buffer of appropriate size (and make sure it will later be released)
			NamePointer=new WCHAR[++NamePointerSize];
			std::unique_ptr<WCHAR> Buf2(NamePointer);
			Buf.swap(Buf2);
		}

		//If the service already exists, confirm the service name matches, and if not, ask if user wants to delete the current service
		if(ServiceExists)
		{
			wprintf(L"%s\n", L"The service already exists");
			if(wcsncmp(NamePointer, Param_DisplayName, NamePointerSize+1))
			{
				//If the server names do not match, ask the user what to do
				wprintf(L"%s:\nCurrent: %s\nRequested: %s\n", L"The service names do not match", NamePointer, Param_DisplayName);

				//Make the request
				if(!AskQuestion(L"Would you like to replace the service? (y/n)")) //If user does not wish to replace the service
					return FormatError(L"%s", L"Cannot continue if service names do not match");

				//Delete the service
				wprintf(L"%s\n", L"Deleting the old service");
				ServiceExists=false;
				SC_HANDLE TheService=OpenService(HSCManager, Param_ServiceName, DELETE);
				if(!TheService)
					return FormatError(L"%s: %s", L"Could not open the service to delete it", GetLastErrorStr());
				SmartCloseService CloseTheService(&TheService, Delete_SmartCloseService); //Close the service handle
				if(!DeleteService(TheService))
					return FormatError(L"%s: %s", L"Could not delete the service", GetLastErrorStr());
				wprintf(L"%s\n", L"The service has been deleted");
			}
		}
	}

	//Create the service
	SC_HANDLE TheService;
	if(!ServiceExists)
	{
		//Confirm the driver path exists
		wprintf(L"%s\n", L"Checking the driver file");
		DWORD FileAttrs=GetFileAttributes(Param_DriverPath);
		if(FileAttrs==INVALID_FILE_ATTRIBUTES)
			return FormatError(L"%s: %s", L"Given path is invalid", GetLastErrorStr());
		if(FileAttrs&FILE_ATTRIBUTE_DIRECTORY)
			return FormatError(L"%s: %s", L"Given path is invalid", L"Path is a folder");

		//Create the service
		wprintf(L"%s\n", L"Creating the service");
		TheService=CreateService(
			HSCManager, Param_ServiceName, Param_DisplayName, 
			SERVICE_START|SERVICE_STOP, 
			SERVICE_KERNEL_DRIVER, SERVICE_DEMAND_START, SERVICE_ERROR_IGNORE,
			Param_DriverPath, nullptr, nullptr, nullptr, nullptr, nullptr);
		if(!TheService)
			return FormatError(L"%s: %s", L"Could not create the service", GetLastErrorStr());

	//Open the service if not creating
	} else {
		TheService=OpenService(HSCManager, Param_ServiceName, SERVICE_START|SERVICE_STOP);
		if(!TheService)
			return FormatError(L"%s: %s", L"Could not open the service", GetLastErrorStr());
	}
	SmartCloseService CloseTheService(&TheService, Delete_SmartCloseService); //Close the service on exit

	//Start the service
	wprintf(L"%s\n", L"Starting the service");
	for(INT_PTR i=0;i<2;i++)
	{
		if(StartService(TheService, 0, nullptr))
			break;

		//If not "service already running" error, or user does not want to stop the current service
		if(i==1 || GetLastError()!=ERROR_SERVICE_ALREADY_RUNNING || !AskQuestion(L"The service is already running. Would you like to stop it? (y/n)"))
			return FormatError(L"%s: %s", L"Could not start the service", GetLastErrorStr());

		//Stop the service
		SERVICE_STATUS ss;
		wprintf(L"%s\n", L"Stopping the current service");
		if(!ControlService(TheService, SERVICE_CONTROL_STOP, &ss))
			return FormatError(L"%s: %s", L"Could not stop the current service", GetLastErrorStr());
	}
	wprintf(L"%s\n", L"Started the service");

	//Ask if the user wants to close the service
	if(!AskQuestion(L"Would you like to stop the service before exit? (y/n)"))
		return nullptr;

	//Stop the service
	SERVICE_STATUS ss;
	if(!ControlService(TheService, SERVICE_CONTROL_STOP, &ss))
		return FormatError(L"%s: %s", L"Could not stop the service", GetLastErrorStr());
	if(ss.dwCurrentState!=SERVICE_STOP_PENDING && ss.dwCurrentState!=SERVICE_STOPPED)
		return FormatError(L"%s", L"The service does not appear to be closing");
	wprintf(L"%s\n", L"The service has been stopped");

	//Return success
	return nullptr;
}

WCHAR* FormatError(WCHAR* Format, ...)
{
	static WCHAR Err[255];
	va_list VAList;
	va_start(VAList, Format);
	vswprintf(Err, sizeof(Err)/sizeof(Err[0]), Format, VAList);
	return Err;
}

SmartWinAlloc GetLastErrorStr()
{
	LPWSTR MessageBuffer=nullptr;
	FormatMessage(
		FORMAT_MESSAGE_ALLOCATE_BUFFER|FORMAT_MESSAGE_FROM_SYSTEM|FORMAT_MESSAGE_IGNORE_INSERTS|FORMAT_MESSAGE_MAX_WIDTH_MASK,
		nullptr, GetLastError(), MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT), (LPWSTR)&MessageBuffer, 0, nullptr);
	return SmartWinAlloc(MessageBuffer, Delete_SmartWinAlloc);
}

BOOLEAN AskQuestion(WCHAR* Question)
{
	//Make the request and wait for an input character
	while(1)
	{
		//Ask the question and get the answer
		wprintf(L"%s:", Question);
		fflush(stdout);
		char InputChar=_getch();
		printf("\n");

		//Check for a valid answer
		if(InputChar=='n' || InputChar=='N')
			return FALSE;
		if(InputChar=='y' || InputChar=='Y')
			return TRUE;
	}
}
```