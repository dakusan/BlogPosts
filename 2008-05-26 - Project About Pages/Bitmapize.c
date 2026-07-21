//Prepare to decode signature
//const UCHAR BytesPerPixel=4, TranspMask=255; //32 bits per pixel (for quicker copies and such - variable not used due to changing BYTE*s to DWORD*s), and white=transparent background color - also not used anymore since we directly write in the background color
//Load data from executable
HGLOBAL GetData=LoadResource(NULL, FindResource(NULL, "DakSig", "Sig")); //Load the resource from the executable
BYTE *Input=(BYTE*)LockResource(GetData); //Get the resource data

//Prepare header and decoding data
UINT BitOn=0; //Bit we are currently on in reading
EncodedFileHeader H=*(EncodedFileHeader*)Input; //Save header locally so we have quick memory lookups
DWORD *Output=Signature=new DWORD[H.Width*H.Height]; //Allocate signature memory

//Prepare the index colors
DWORD Indexes[17], IndexSize=NumBitsRequired(H.NumIndexes); //Save full color indexes locally so we have quick lookups, use 17 index array so we don't have to allocate memory (since we already know how many there will be), #16=transparent color
DWORD BackgroundColor=GetSysColor(COLOR_BTNFACE), FontColor=0x067906;
BYTE *BGC=(BYTE*)&BackgroundColor, *FC=(BYTE*)&FontColor;
for(UINT i=0;i<16;i++) //Alpha blend the indexes
{
	float Alpha=((EncodedFileHeader*)Input)->Indexes[i] / 255.0f;
	BYTE IndexColor[4];
	for(int n=0;n<3;n++)
		IndexColor[n]=(BYTE)(BGC[n]*Alpha + FC[n]*(1-Alpha));
	//IndexColor[3]=0; //Don't really need to worry about the last byte as it is unused
	Indexes[i]=*(DWORD*)IndexColor;
}
Indexes[16]=BackgroundColor; //Translucent background = window background color

//Unroll/unencode all the pixels
Input+=(sizeof(EncodedFileHeader)+H.NumIndexes); //Start reading input past the header
do
{
	UINT l; //Length (transparent and then index)
	//Transparent pixels
	memsetd(Output, Indexes[16], l=ReadBits(Input, BitOn, H.TranspSize));
	Output+=l;

	//Translucent pixels
	l=ReadBits(Input, BitOn+=H.TranspSize, H.TranslSize);
	BitOn+=H.TranslSize;
	for(i=0;i<l;i++) //Write the gray scale out to the 3 pixels, this should technically be done in a for loop, which would unroll itself anyways, but this way ReadBits+index lookup is only done once - ** Would need to be in a for loop if not using gray-scale or 24 bit output
		Output[i]=Indexes[ReadBits(Input, BitOn+i*IndexSize, IndexSize)];
	Output+=l;
	BitOn+=l*IndexSize;
} while(BitOn<H.DataSize);

//Output the signature
const BITMAPINFOHEADER MyBitmapInfo={sizeof(BITMAPINFOHEADER), 207, 42, 1, 32, BI_RGB, 0, 0, 0, 0, 0};
SetDIBitsToDevice(MyDC, x, y, MyBitmapInfo.biWidth, MyBitmapInfo.biHeight, 0, 0, 0, MyBitmapInfo.biHeight, Signature, (BITMAPINFO*)&MyBitmapInfo, DIB_RGB_COLORS);