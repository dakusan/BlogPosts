//** Double stars denotes changes for custom circumstance [The About Window Mask]
#include <windows.h>
#include <stdio.h>
#include <conio.h>

//Our encoding functions
int ErrorOut(char* Error, FILE* HandleToClose); //If an error occurs, output
UINT Encode(UCHAR* Input, UCHAR* Output, UINT Width, UINT Height); //Encoding process
UCHAR NumBitsRequired(UINT Num); //Tests how many bits are required to represent a number
void WriteToBits(UCHAR* StartPointer, UINT BitStart, UINT Value); //Write Value to Bit# BitStart after StartPointer - Assumes more than 8 bits are never written

//Decode
void Decode(UCHAR* Input, UCHAR* Output); //Decoding process
UCHAR _fastcall ReadBits(UCHAR* StartPointer, UINT BitStart, UCHAR BitSize); //Read value from Bit# BitStart after StartPointer - Assumes more than 8 bits are never read

//Confirm Decoding
BOOL CheckDecode(UCHAR* Input1, UCHAR* Input2, UINT Width, UINT Height); //Confirm Decoding

//Optimized DeCode
void DeCode(UCHAR *Input, UCHAR *Output);

//Program constants
const UCHAR BytesPerPixel=3, TranspMask=255; //24 bits per pixel, and white = transparent background color

//Encoding file header
typedef struct
{
	USHORT DataSize; //Data size in bits - **Should be UINT
	UCHAR Width, Height; //**Should be USHORTs
	UCHAR TranspSize, TranslSize; //Largest number of bits required for a run length for Transp[arent] and Transl[ucent]
	UCHAR NumIndexes, Indexes[0]; //Number and list of indexes
} EncodedFileHeader;

int main()
{
	UCHAR *InputBuffer, *OutputBuffer; //Where we will hold our input and output data
	FILE *File; //Handle to current input or output file
	UINT FileSize; //Holds input and output file sizes

	//The bitmap headers tell us about its contents
	BITMAPFILEHEADER BitmapFileHead;
	BITMAPINFOHEADER BitmapHead;

	//Read in bitmap header and confirm file type
	File=fopen("AboutWindow-Mask.bmp", "rb"); //Normally you'd read in the filename from passed arguments (argv)
	if(!File) //Confirm file open
		return ErrorOut("Cannot open file for reading", NULL);
	fread(&BitmapFileHead, sizeof(BITMAPFILEHEADER), 1, File);
	if(BitmapFileHead.bfType!=*(WORD*)"BM" || BitmapFileHead.bfReserved1 || BitmapFileHead.bfReserved2) //Confirm we are opening a bitmap
		return ErrorOut("Not a bitmap", File);

	//Read in the rest of the data
	fread(&BitmapHead, sizeof(BITMAPINFOHEADER), 1, File);
	if(BitmapHead.biPlanes!=1 || BitmapHead.biBitCount!=24 || BitmapHead.biCompression!=BI_RGB) //Confirm bitmap type - this code would probably have been simpler if I did an 8 bit indexed file instead... oh well, NBD.  **It has also been programmed for easy transition to 8 bit indexed files via the "BytesPerPixel" constant.
		return ErrorOut("Bitmap must be in 24 bit RGB format", File);
	FileSize=BitmapFileHead.bfSize-sizeof(BITMAPINFOHEADER)-sizeof(BITMAPFILEHEADER); //Size of the data portion
	InputBuffer=malloc(FileSize);
	fread(InputBuffer, FileSize, 1, File);
	fclose(File);

	//Run Encode
	OutputBuffer=malloc(FileSize); //We should only ever need at most FileSize space for output (output should always be smaller)
	memset(OutputBuffer, 0, FileSize); //Needs to be zeroed out due to how writing of data file is non sequential
	FileSize=Encode(InputBuffer, OutputBuffer, BitmapHead.biWidth, BitmapHead.biHeight); //Encode the file and get the output size

	//Write encoded data out
	File=fopen("Output.msk", "wb");
	fwrite(OutputBuffer, FileSize, 1, File);
	fclose(File);
	printf("File %d written with %d bytes\n", 1, FileSize);

	{
		//Run Decode
		UCHAR* O2=(BYTE*)malloc(BitmapFileHead.bfSize);
		Decode(OutputBuffer, O2); //Change to DeCode for optimized

/*		//If writing back out to a 24 bit windows bitmap, this adds the row padding back in
		File=fopen("output.bmp", "wb");
		fwrite(&BitmapFileHead, sizeof(BITMAPFILEHEADER), 1, File);
		fwrite(&BitmapHead, sizeof(BITMAPINFOHEADER), 1, File);
		fwrite(O2, BitmapFileHead.bfSize-sizeof(BITMAPINFOHEADER)-sizeof(BITMAPFILEHEADER), 1, File);*/

		printf(CheckDecode(InputBuffer, O2, BitmapHead.biWidth, BitmapHead.biHeight) ? "good" : "bad");
	}

	//Free up memory and wait for user input
	free(InputBuffer);
	free(OutputBuffer);
	getch(); //Pause for user input
	return 0;
}

int ErrorOut(char* Error, FILE* HandleToClose) //If an error occurs, output
{
	if(HandleToClose)
		fclose(HandleToClose);
	printf("%s\n", Error);
	getch(); //Pause for user input
	return 1;
}

UINT Encode(UCHAR* Input, UCHAR* Output, UINT Width, UINT Height) //Encoding process
{
	UCHAR Indexes[256], NumIndexes, IndexSize, RowPad; //The index re-mappings, number of indexes, number of bits an index takes in output data, padding at input row ends for windows bitmaps
	USHORT TranspSize, TranslSize; //Largest number of bits required for a run length for Transp[arent] (zero) and Transl[ucent] (non zero) - should be UCHAR's, but these are used as explained under "CurTranspLen" below
	UINT BitSize, x, y, ByteOn, NumPixels; //Current output size in bits, x/y coordinate counters, current byte location in Input, number of pixels in mask

	//Calculate some stuff
	NumPixels=Width*Height; //Number of pixels in mask
	RowPad=4-(Width*BytesPerPixel%4); //Account for windows DWORD row padding - see declaration comment
	RowPad=(RowPad==4 ? 0 : RowPad);

	{ //Do a first pass to find number of different mask values, run lengths, and their encoded values
		const UCHAR UnusedIndex=255; //In our index list, unused indexes are marked with this constant
		USHORT CurTranspLen, CurTranslLen; //Keep track of the lengths of the current transparent & translucent runs - TranspSize and TranslSize are temporarily used to hold the maximum run lengths
		//Zero out all index references and counters
		memset(Indexes, UnusedIndex, 256);
		NumIndexes=0;
		TranspSize=TranslSize=CurTranspLen=CurTranslLen=0;
		//Start gathering data
		for(y=ByteOn=0;y<Height;y++) //Column
		{
			for(x=0;x<Width;x++,ByteOn+=BytesPerPixel) //Row
			{
				UCHAR CurMask=Input[ByteOn]; //Curent alpha mask
				if(CurMask!=TranspMask) //Translucent value?
				{
					//Determine if index has been used yet
					if(Indexes[CurMask]==UnusedIndex) //We only need to check 1 byte per pixel as they are all the same for gray-scale **This would need to change if using non 24-bit or non gray-scale
					{
						((EncodedFileHeader*)Output)->Indexes[NumIndexes]=CurMask; //Save mask number in the index header
						Indexes[CurMask]=NumIndexes++; //Save index number to the mask
					}

					//Length of current transparent run
					TranspSize=(CurTranspLen>TranspSize ? CurTranspLen : TranspSize); //Max(CurTranspLen, TranspSize)
					CurTranspLen=0;

					//Length of current translucent run
					CurTranslLen++;
				}
				else //Transparent value?
				{
					//Length of current translucent run
					TranslSize=(CurTranslLen>TranslSize ? CurTranslLen : TranslSize);  //Max(CurTranslLen, TranslSize)
					CurTranslLen=0;

					//Length of current transparent run
					CurTranspLen++;
				}
			}

			ByteOn+=RowPad; //Account for windows DWORD row padding
		}
		//Determine number of required bits per value
		printf("Number of Indexes: %d\nLongest Transparent Run: %d\nLongest Translucent Run: %d\n", NumIndexes,
			TranspSize=CurTranspLen>TranspSize ? CurTranspLen : TranspSize, //Max(CurTranspLen, TranspSize)
			TranslSize=CurTranslLen>TranslSize ? CurTranslLen : TranslSize  //Max(CurTranslLen, TranslSize)
			);
		IndexSize=NumBitsRequired(NumIndexes);
		TranspSize=NumBitsRequired(TranspSize); //**This is currently overwritten a few lines down
		TranslSize=NumBitsRequired(TranslSize); //**This is currently overwritten a few lines down
		printf("Bit Lengths of - Indexes, Trasparent Run Length, Translucent Run Length: %d, %d, %d\n", IndexSize, TranspSize, TranslSize);
	}

	//**Modify run sizes (custom) - this function could be run multiple times with different TranspSize and TranslSize until the best values are found - the best values would always be a weighted average
	TranspSize=6;
	TranslSize=3;

	//Start processing data
	BitSize=(sizeof(EncodedFileHeader)+NumIndexes)*8; //Skip the file+bitmap headers and measure in bits
	x=ByteOn=0;
	do
	{
		//Transparent run
		UINT CurRun=0;
		while(Input[ByteOn]==TranspMask && x<NumPixels && CurRun<(UINT)(1<<TranspSize)-1) //Final 2 checks are for EOF and capping run size to max bit length
		{
			x++;
			CurRun++;
			ByteOn+=BytesPerPixel;
			if(x%Width==0) //Account for windows DWORD row padding
				ByteOn+=RowPad;
		}
		WriteToBits(Output, BitSize, CurRun);
		BitSize+=TranspSize;

		//Translucent run
		CurRun=0;
		BitSize+=TranslSize; //Prepare to start writing masks first
		while(x<NumPixels && Input[ByteOn]!=TranspMask && CurRun<(UINT)(1<<TranslSize)-1) //Final 2 checks are for EOF and and capping run size to max bit length
		{
			WriteToBits(Output, BitSize+CurRun*IndexSize, Indexes[Input[ByteOn]]);
			x++;
			CurRun++;
			ByteOn+=BytesPerPixel;
			if(x%Width==0) //Account for windows DWORD row padding
				ByteOn+=RowPad;
		}
		WriteToBits(Output, BitSize-TranslSize, CurRun); //Write the mask before the indexes
		BitSize+=CurRun*IndexSize;
	} while(x<NumPixels);

	{ //Output header
		EncodedFileHeader *OutputHead;
		OutputHead=(EncodedFileHeader*)Output;
		OutputHead->DataSize=BitSize-(sizeof(EncodedFileHeader)+NumIndexes)*8; //Length of file in bits not including header
		OutputHead->Width=Width;
		OutputHead->Height=Height;
		OutputHead->TranspSize=(UCHAR)TranspSize;
		OutputHead->TranslSize=(UCHAR)TranslSize;
		OutputHead->NumIndexes=NumIndexes;
	}
	return BitSize/8+(BitSize%8 ? 1 : 0); //Return entire length of file in bytes
}

UCHAR NumBitsRequired(UINT Num) //Tests how many bits are required to represent a number
{
	UCHAR RetNum;
	_asm //Find the most significant bit
	{
		xor eax, eax //eax=0
		bsr eax, Num //Find most significant bit in eax
		mov RetNum, al
	}
	return RetNum+((UCHAR)(1<<RetNum)==Num ? 0 : 1); //Test if the most significant bit is the only one set, if not, at least 1 more bit is required
}

void WriteToBits(UCHAR* StartPointer, UINT BitStart, UINT Value) //Write Value to Bit# BitStart after StartPointer - Assumes more than 8 bits are never written
{
	*(WORD*)(&StartPointer[BitStart/8])|=Value<<(BitStart%8);
}

//Decoding
void Decode(UCHAR* Input, UCHAR* Output) //Decoding process
{
	EncodedFileHeader H=*(EncodedFileHeader*)Input; //Save header locally so we have quick memory lookups
	UCHAR Indexes[256], IndexSize=NumBitsRequired(H.NumIndexes); //Save indexes locally so we have quick lookups, use 256 index array so we don't have to allocate memory
	UINT BitOn=0; //Bit we are currently on in reading
	memcpy(Indexes, ((EncodedFileHeader*)Input)->Indexes, 256); //Save the indexes
	Input+=(sizeof(EncodedFileHeader)+H.NumIndexes); //Start reading input past the header

	//Unroll/unencode all the pixels
	do
	{
		UINT i, l; //index counter, length (transparent and then index)
		//Transparent pixels
		memset(Output, TranspMask, l=ReadBits(Input, BitOn, H.TranspSize)*BytesPerPixel);
		Output+=l;

		//Translucent pixels
		l=ReadBits(Input, BitOn+=H.TranspSize, H.TranslSize);
		BitOn+=H.TranslSize;
		for(i=0;i<l;i++) //Write the gray scale out to the 3 pixels, this should technically be done in a for loop, which would unroll itself anyways, but this way ReadBits+index lookup is only done once - ** Would need to be in a for loop if not using gray-scale or 24 bit output
			Output[i*BytesPerPixel]=Output[i*BytesPerPixel+1]=Output[i*BytesPerPixel+2]=Indexes[ReadBits(Input, BitOn+i*IndexSize, IndexSize)];
		Output+=l*BytesPerPixel;
		BitOn+=l*IndexSize;
	} while(BitOn<H.DataSize);

/*	{ //If writing back out to a 24 bit windows bitmap, this adds the row padding back in
		UINT i;
		UCHAR RowPad=4-(H.Width*BytesPerPixel%4); //Account for windows DWORD row padding
		RowPad=(RowPad==4 ? 0 : RowPad);
		Output-=H.Width*H.Height*BytesPerPixel; //Restore original output pointer
		for(i=H.Height;i>0;i--) //Go backwards so data doesn't overwrite itself
			memcpy(Output+(H.Width*BytesPerPixel+RowPad)*i, Output+(H.Width*BytesPerPixel)*i, H.Width*BytesPerPixel);
	}*/
}

UCHAR _fastcall ReadBits(UCHAR* StartPointer, UINT BitStart, UCHAR BitSize) //Read value from Bit# BitStart after StartPointer - Assumes more than 8 bits are never read
{
	return (*(WORD*)&StartPointer[BitStart/8]>>BitStart%8)&((1<<BitSize)-1);
}

BOOL CheckDecode(UCHAR* Input1, UCHAR* Input2, UINT Width, UINT Height) //Confirm Decoding
{
	UINT x,y,i;
	UCHAR RowPad=4-(Width*BytesPerPixel%4); //Account for windows DWORD row padding
	RowPad=(RowPad==4 ? 0 : RowPad);

	for(y=0;y<Height;y++)
		for(x=0;x<Width;x++)
			for(i=0;i<BytesPerPixel;i++)
				if(Input1[y*(Width*BytesPerPixel+RowPad)+x*BytesPerPixel+i]!=Input2[y*(Width*BytesPerPixel)+x*BytesPerPixel+i])
					return FALSE;
	return TRUE;
}

//Optimized DeCode
void DeCode(UCHAR *Input, UCHAR *Output)
{
	EncodedFileHeader *H=(EncodedFileHeader*)Input; //Save header locally so we have quick memory lookups
	UINT BitOn=0; //Bit we are currently on in reading
	UCHAR IndexSize;
	Input+=(sizeof(EncodedFileHeader)+H->NumIndexes); //Start reading input past the header 

	__asm
	{
		//Determine IndexSize
		xor ecx,ecx
		mov eax, H
		mov cl, [eax+6] //NumIndexes
		xor al, al
		xor bl, bl
		bsr eax, ecx
		bsf ebx, ecx
		cmp al,bl
		setne cl
		add al, cl
		mov IndexSize, al
	}

	//Unroll/unencode all the pixels
	do
	{
		//Transparent pixels
		UINT i, l=ReadBits(Input, BitOn, H->TranspSize)*BytesPerPixel; //index counter, length (transparent and then index)
		_asm
		{
			//memset(Output, TranspMask, l=*BytesPerPixel);
			mov edi, Output
			mov ecx, l
			mov al, TranspMask
			rep stosb
		}
		Output+=l;

		//Translucent pixels
		l=ReadBits(Input, BitOn+=H->TranspSize, H->TranslSize);
		BitOn+=H->TranslSize;
		for(i=0;i<l;i++) //Write the gray scale out to the 3 pixels
			Output[i*BytesPerPixel]=Output[i*BytesPerPixel+1]=Output[i*BytesPerPixel+2]=H->Indexes[ReadBits(Input, BitOn+i*IndexSize, IndexSize)];
		Output+=l*BytesPerPixel;
		BitOn+=l*IndexSize;
	} while(BitOn<H->DataSize);
}