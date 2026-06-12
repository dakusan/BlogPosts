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