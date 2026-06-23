<div>
Original Post: <a href="https://www.castledragmire.com/Posts/chunk">chunk</a><br>
Date: 2009-12-31<br>
SubTitle: Useful little Unix like utility for command line<br>
Tags: chunk, Command Line Utility, Data Passing<br><br>

<p>I needed a command line utility for Bash (for both Windows and Linux) that only outputs bytes between 2 points in a file to STDOUT. <i>head</i> and <i>tail</i> weren’t really cutting it so I figured I’d throw something together. Below is the <a href="chunk.cpp">source code</a> for the result, which I call <b>chunk</b> (<a href="chunk.exe">Windows Executable</a>).</p>
<p>I compiled the file as c++, but it should be c99 compatible. The file has been tested as compilable for: GCC4 on Slackware, GCC3 on Red Hat, and GCC3 on MingW (Windows [<i>WIN32</i> should be defined by the compiler]).</p><br>
<blockquote>
Chunk outputs bytes between 2 points in a file to STDOUT. The parameters are:<br>
1) The file<br>
2) The byte offset to start at (hex is supported like 0xA)<br>
3) The number of bytes to output. If not given, the end of the file is assumed.
</blockquote><br>

The source is as follows:
</div>

```cpp
//See http://www.castledragmire.com/Posts/chunk for more information

#define __LARGE64_FILES
#include <stdio.h>
#include <stdlib.h> //strtoull

#ifdef WIN32 //STDOUT only needs to be set to binary mode in windows
	#include <io.h> //_setmode
	#include <fcntl.h> //_O_BINARY
#endif

typedef unsigned long long UINT64;
const UINT64 MaxSizeToRead=1024*1024*10; //The maximum number of bytes to read at a time to our buffer (Must be < 2^31)

UINT64 GetNumberFromString(const char* S) //Extract both hexidecimal and decimal numbers from a string
{
	bool IsHex=S[0]=='0' && (S[1]|32=='x'); //If string starts as 0x, then is a hex number
	return strtoull(S+(IsHex ? 2 : 0), NULL, IsHex ? 16 : 10); //Hex number starts after 2 characters and uses base 16
}

int main(int argc, char *argv[], char *envp[])
{
	//Determine if proper number of parameters are passed, and if not, output help info
	if(argc!=3 && argc!=4)
		return fprintf(stderr, "Chunk outputs bytes between 2 points in a file to STDOUT. The parameters are:\n1) The file\n2) The byte offset to start at (hex is supported like 0xA)\n3) The number of bytes to output. If not given, the end of the file is assumed.\n") & 0;

	//Open the file and get its length
	FILE *TheFile=fopen64(argv[1], "rb");
	if(TheFile==NULL)
		return fprintf(stderr, "File not found or cannot open file\n") & 0;
	fseeko64(TheFile, 0, SEEK_END); //Get the length by seeking to the end
	UINT64 FileSize=ftello64(TheFile);

	//Determine the requested start offset
	UINT64 Offset=GetNumberFromString(argv[2]), SizeToOutput;
	if(Offset>=FileSize)
	{
		fprintf(stderr, "Offset is larger than file's size\n");
		fclose(TheFile);
		return 0;
	}

	//Determine the size to read
	if(argc==3) //If no final parameter, read to the end of the file
		SizeToOutput=FileSize-Offset;
	else //Determine from the 3rd parameter
	{
		SizeToOutput=GetNumberFromString(argv[3]);
		if(Offset+SizeToOutput>FileSize)
		{
			fprintf(stderr, "Requested size is larger than the file, truncating to end of file\n");
			SizeToOutput=FileSize-Offset;
		}
		else if(!SizeToOutput) //If nothing to output, exit prematurely
		{
			fclose(TheFile);
			return 1;
		}
	}

	//Output requested data 10MB at a time from the file to STDOUT
	char *Buffer=new char[SizeToOutput>MaxSizeToRead ? MaxSizeToRead : SizeToOutput]; //Only allocate as many bytes to our read buffer as is necessary
	fseeko64(TheFile, Offset, SEEK_SET); //Seek to the beginning read offset of our file
	#ifdef WIN32 //STDOUT only needs to be set to binary mode in windows
		_setmode(_fileno(stdout), _O_BINARY);
	#endif
	while(SizeToOutput) //Keep reading and outputting until requested data is complete
	{
		UINT64 SizeToRead=SizeToOutput>MaxSizeToRead ? MaxSizeToRead : SizeToOutput; //Number of bytes to read and write
		fread(Buffer, SizeToRead, 1, TheFile); //Read the data
		fwrite(Buffer, SizeToRead, 1, stdout); //Write the data to STDOUT
		SizeToOutput-=SizeToRead; //Decrease number of bytes we still need to read
	}

	//Cleanup
	delete[] Buffer;
	fclose(TheFile);
	return 1;
}
```

<div>
</div>