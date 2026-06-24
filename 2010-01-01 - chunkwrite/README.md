<div>
Original Post: <a href="https://www.castledragmire.com/Posts/chunkwrite">chunkwrite</a><br>
Date: 2010-01-01<br>
SubTitle: Another useful little Unix like utility for command line<br>
Tags: chunk, Command Line Utility, Data Passing<br><br>

<p>This is a modification of the utility I made in <a href="../2009-12-31 - chunk">yesterday’s post, <b>chunk</b></a>, and most of the information mentioned in that post applies to this one too. This utility instead writes bytes from STDIN to a certain byte offset of a preexisting file. Below is the <a href="chunkwrite.cpp">source code</a> for the result, which I call <b>chunkwrite</b> (<a href="chunkwrite.exe">Windows Executable</a>).</p>
<blockquote>
Chunkwrite writes bytes to a file at a given offset from STDIN. The parameters are:<br>
1) The file to write to<br>
2) The byte offset to write at (hex is supported like 0xA)
</blockquote><br>

The source is as follows:
</div>

```cpp
//See http://www.castledragmire.com/Posts/chunkwrite for more information

#define __LARGE64_FILES
#include <stdio.h>
#include <stdlib.h> //strtoull

#ifdef WIN32 //STDIN only needs to be set to binary mode in windows
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
	if(argc!=3)
		return fprintf(stderr, "Chunkwrite writes bytes to a file at a given offset from STDIN. The parameters are:\n1) The file to write to\n2) The byte offset to write at (hex is supported like 0xA)\n") & 0;

	//Open the file to output to
	FILE *TheFile=fopen64(argv[1], "r+b");
	if(TheFile==NULL)
		return fprintf(stderr, "File not found or cannot open file\n") & 0;

	//Determine the requested start offset
	UINT64 Offset=GetNumberFromString(argv[2]);

	//Write the data 10MB at a time from STDIN to the file
	char *Buffer=new char[MaxSizeToRead];
	fseeko64(TheFile, Offset, SEEK_SET); //Seek to the beginning write offset of our file
	#ifdef WIN32 //STDIN only needs to be set to binary mode in windows
		_setmode(_fileno(stdin), _O_BINARY);
	#endif
	do
	{
		size_t AmountRead=fread(Buffer, 1, MaxSizeToRead, stdin); //Read the data from STDIN
		fwrite(Buffer, AmountRead, 1, TheFile); //Write the data to the file
	}
	while(!feof(stdin)); //Keep reading and writing until STDIN is complete

	//Cleanup
	delete[] Buffer;
	fclose(TheFile);
	return 1;
}
```

<div>
</div>