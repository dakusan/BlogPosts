<div>
Original Post: <a href="https://www.castledragmire.com/Posts/FoxPro_Table_Memo_Corruption">FoxPro Table Memo Corruption</a><br>
Date: 2008-06-13<br>
Tags: Diagnosing, Reverse Engineering<br><br>

<p>My <s><a href="http://www.eyecareteam.com">father’s optometric practice</a></s> has been using an old DOS database called “Eyecare” since the (I believe) early 80’s.  For many years, he has been programming a new, very customized, database up from scratch in Microsoft Access which is backwards compatible with “Eyecare”, which uses a minor variant of FoxPro databases.  I’ve been helping him with minor things on it for a number of years, and more recently I’ve been giving a lot more help in getting it secured and migrated from Microsoft Access databases (.mdb) into MySQL.</p>
<p>A recent problem cropped up in that one of the primary tables started crashing Microsoft Access when it was opened (through a FoxPro ODBC driver). Through some tinkering, he discovered that the memo file (.fpt) for the table was corrupted, as trying to view any memo fields is what crashed Access.  He asked me to see if I could help in recovering the file, which fortunately I can do at my leisure, as he keeps paper backups of everything for just such circumstances. He keeps daily backups of everything too… but for some reason that’s not an option.</p><br>

<p>I went about trying to recover it through the easiest means first, namely, trying to open and export the database through FoxPro, which only recovered 187 of the ~9000 memo records.  Next, I tried finding a utility online that did the job, and the first one I found that I thought should work was called “FoxFix”, but it failed miserably.  There are a number of other Shareware utilities I could try, but I decided to just see how hard it would be to fix myself first.</p><br>

<p>I opened the memo file up in a HEX editor, and after some very quick perusing and calculations, it was quite easy to determine the format:</p>
<ul>
	<li>512 byte header</li>
	<li>Record Data (Padded to 32 bytes)<ul>
		<li>4 byte Record Type (Little endian, 1=Memo)</li>
		<li>4 byte Record Length (Little endian)</li>
		<li>Character data</li></ul></li>
</ul>
So I continued on the path of seeing what I could do to fix the file.
<ul>
<li>First, I had it jump to the header of each record and just get the record data length, and I very quickly found multiple invalid record lengths.</li>
<li>Next, I had it attempt to fix each of these by determining the real length of the memo by searching for the first null terminator (“\0”) character, but I quickly discovered an oddity.  There are weird sections in many of the memo fields in the format BYTE{0,0,0,1,0,0,0,1,x}, which is 2 little endian DWORDS which equal 1, and a final byte character (usually 0).</li>
<li>I added to the algorithm to include these as part of a memo record, and many more original memo lengths then agreed with my calculated memo lengths.</li>
<li>The final thing I did was determine how many invalid (non keyboard) characters there were in the memo data fields.  There were ~3500 0x8D characters, which were usually always followed by 0xA, so I assume these were supposed to be line breaks (Windows line breaks are denoted by [0xD<b>/</b>new line<b>/</b>\r],[0xA<b>/</b>carriage return<b>/</b>\n]). There were only 5 other invalid characters, so I just changed these to question marks ‘?’.</li>
</ul><br>

<p>Unfortunately, Microsoft Access still crashed when I tried to access the comments fields, so I will next try to just recover the data, tie it to its primary keys (which I will need to determine through the table file [.dbf]), and then rebuild the table. I should be making another post when I get around to doing this.</p><br>

The following code which “fixes” the table’s memo file took about 2 hours to code up.
</div>

```c
//Usually included in windows.h
typedef unsigned long DWORD;
typedef unsigned char BYTE;

//Includes
#include <iostream.h> //cout
#include <stdio.h> //file io
#include <conio.h> //getch
#include <ctype.h> //isprint

//Memo file structure
#pragma warning(disable: 4200) //Remove zero-sized array warning
const MemoFileHeadLength=512;
const RecordBlockLength=32; //This is actually found in the header at (WORD*)(Start+6)
struct MemoRecord //Full structure must be padded at end with \0 to RecordBlockLength
{
	DWORD Type; //Type in little endian, 1=Memo
	DWORD Length; //Length in little endian
	BYTE Data[0];
};
#pragma warning(default: 4200)

//Input and output files
const char *InFile="EXAM.Fpt.old", *OutFile="EXAM.Fpt";

//Assembly functions
__forceinline DWORD BSWAP(DWORD n) //Swaps endianness
{
	_asm mov eax,n
	_asm bswap eax
	_asm mov n, eax
	return n;
}

//Main function
void main()
{
	//Read in file
	const FileSize=6966592; //This should actually be found when the file is opened...
	FILE* MyFile=fopen(InFile, "rb");
	BYTE *MyData=new BYTE[FileSize];
	fread(MyData, FileSize, 1, MyFile);
	fclose(MyFile);

	//Start checking file integrity
	DWORD FilePosition=MemoFileHeadLength; //Where we currently are in the file
	DWORD RecordNum=0, BadRecords=0, BadBreaks=0, BadChars=0; //Data Counters
	const DWORD OneInLE=0x01000000; //One in little endian
	while(FilePosition<FileSize) //Loop until EOF
	{
		FilePosition+=sizeof(((MemoRecord*)NULL)->Type); //Advanced passed record type (1=memo)
		DWORD CurRecordLength=BSWAP(*(DWORD*)(MyData+FilePosition)); //Pull in little endian record size
		cout << "Record #" << RecordNum++ << " reports " << CurRecordLength << " characters long. (Starts at offset " << FilePosition << ")" << endl; //Output record information

		//Determine actual record length
		FilePosition+=sizeof(((MemoRecord*)NULL)->Length); //Advanced passed record length
		DWORD RealRecordLength=0; //Actual record length
		while(true)
		{
			for(;MyData[FilePosition+RealRecordLength]!=0 && FilePosition+RealRecordLength<FileSize;RealRecordLength++) //Loop until \0 is encountered
			{
#if 1 //**Check for valid characters might not be needed
				if(!isprint(MyData[FilePosition+RealRecordLength])) //Makes sure all characters are valid
					if(MyData[FilePosition+RealRecordLength]==0x8D) //**0x8D maybe should be in ValidCharacters string? - If 0x8D is encountered, replace with 0xD
					{
						MyData[FilePosition+RealRecordLength]=0x0D;
						BadBreaks++;
					}
					else //Otherwise, replace with a "?"
					{
						MyData[FilePosition+RealRecordLength]='?';
						BadChars++;
					}
#endif
			}

			//Check for inner record memo - I'm not really sure why these are here as they don't really fit into the proper memo record format.... Format is DWORD(1), DWORD(1), BYTE(0)
			if(((MemoRecord*)(MyData+FilePosition+RealRecordLength))->Type==OneInLE && ((MemoRecord*)(MyData+FilePosition+RealRecordLength))->Length==OneInLE /*&& ((MemoRecord*)(MyData+FilePosition+RealRecordLength))->Data[0]==0*/) //**The last byte seems to be able to be anything, so I removed its check
			{ //If inner record memo, current memo must continue
				((MemoRecord*)(MyData+FilePosition+RealRecordLength))->Data[0]=0; //**This might need to be taken out - Force last byte back to 0
				RealRecordLength+=sizeof(MemoRecord)+1;
			}
			else //Otherwise, current memo is finished
				break;
		}
		if(RealRecordLength!=CurRecordLength) //If given length != found length
		{
			//Tell the user a bad record was found
			cout << "   Real Length=" << RealRecordLength << endl;
			CurRecordLength=RealRecordLength;
			BadRecords++;
			//getch();

			//Update little endian bad record length
			((MemoRecord*)(MyData+FilePosition-sizeof(MemoRecord)))->Length=BSWAP(RealRecordLength);
		}

		//Move to next record - Each record, including RecordLength is padded to RecordBlockLength
		DWORD RealRecordSize=sizeof(MemoRecord)+CurRecordLength;
		FilePosition+=CurRecordLength+(RealRecordSize%RecordBlockLength==0 ? 0 : RecordBlockLength-RealRecordSize%RecordBlockLength);
	}

	//Tell the user file statistics
	cout << "Total bad records=" << BadRecords << endl << "Total bad breaks=" << BadBreaks << endl << "Total bad chars=" << BadChars << endl;

	//Output fixed data to new file
	MyFile=fopen(OutFile, "wb");
	fwrite(MyData, FileSize, 1, MyFile);
	fclose(MyFile);

	//Cleanup and wait for user keystroke to end
	delete[] MyData;
	getch();
}
```

<div>
</div>