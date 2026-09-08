/***Dakusan's Embedded Image Drawer - v1.0 http://www.castledragmire.com/Projects/Other_Web_Scripts#DrawImage ***/

var IsIE=(navigator.appName.indexOf("Microsoft")!=-1);

//---Render Object functions---
//ImageClicked(X, Y, Control_Value) is ran when the image is clicked. Control_Value is true if the right mouse button is used in IE, or if CTRL is held down in other browsers
//Create and return the object to render the image
function CreateRenderObject(CodeBase, Width, Height) //The 3 parameters are only used if an applet is needed when using IE. CodeBase is the directory the DrawImage.class is in.
{
	if(IsIE) //If IE, create the Java Applet
	{
		var RenderObject=document.createElement('APPLET');
		//Copy in the object attributes
		var ElementInfo={codebase:CodeBase, code:'DrawImage', width:Width, height:Height, VIEWASTEXT:true};
		for(var i in ElementInfo) if(ElementInfo.hasOwnProperty(i))
			RenderObject[i]=ElementInfo[i];
		return RenderObject;
	}
	
	//If not IE, create an image object
	var RenderObject=document.createElement('IMG');
	RenderObject.onclick=OnImageClickJS;
	return RenderObject;
}
//Functions for when a render object is clicked
function OnImageClickJS(event) { ImageClicked(event.offsetX?event.offsetX:event.pageX-event.target.offsetLeft, (event.offsetY?event.offsetY:event.pageY-event.target.offsetTop), event.ctrlKey); }
function OnImageClickApplet(X, Y, RightButtonClicked) { ImageClicked(parseInt(X), parseInt(Y), parseInt(RightButtonClicked)); }

//---Image output functionality---
function DrawImage(ImageObject, Width, Height, PixelData) //Outputs the image to ImageObject. PixelData is an array of pixel colors
{
	//IE ends the data to the Java applet
	if(IsIE)
		return ImageObject.DrawFrame(Width, Height, PixelData.join(','));

	//Other browsers render to an IMG object
	//Turn each 24bit pixel value into 3 seperate 8 bit color channel values
	var NewPixelData=[];
	for(var i=0;i<PixelData.length;i++)
	{
		NewPixelData.push(PixelData[i]&0xFF);
		NewPixelData.push((PixelData[i]&0xFF00)>>8);
		NewPixelData.push((PixelData[i]&0xFF0000)>>16);
	}

	OutputBitmapByteArrayToImage(CreateBitmapFile(Width, Height, 24, NewPixelData), ImageObject); //Render the bitmap file to an IMG
}

//Create a byte array of the bitmap file
function CreateBitmapFile(Width, Height, BitCount, Data) //BitCount must be divisible by 8, Data is a Numeric Byte Array
{
	var BytesPerPixel=(BitCount/8); //BitCount must be divisible by 8
	var BytesPerRow=Width*BytesPerPixel; //Current Bytes per Row, this is padded to a multiple of 4
	var ExtraBytes=BytesPerRow%4; //Extra bytes on the row that need to be padded
	if(ExtraBytes) //If row byte lengths are not multiples of 4, they need to be padded
	{
		Data=Data.slice(); //Copy the data array for manipulation

		//Determine the number of bytes to pad
		var PadBytes=[], NumPadBytes=4-ExtraBytes;
		for(var i=0;i<NumPadBytes;i++)
			PadBytes.push(0);

		//Pad each row with the end bytes
		BytesPerRow+=NumPadBytes;
		for(var i=0;i<Height;i++)
			Data.splice.apply(Data, [(i+1)*BytesPerRow-NumPadBytes, 0].concat(PadBytes));
	}

	return [].concat( //Concat all the byte elements into a single byte array and return, starting with the bitmap file header
		DecBytes(0x4D42, 2),		//bfType:	Specifies the file type. It must be BM.
		DecBytes(14+40+Data.length, 4),	//bfSize:	Specifies the size, in bytes, of the bitmap file. -- -- BITMAPFILEHEADER(14) + BITMAPINFOHEADER(40) + sizeof(DATA)
		DecBytes(0, 2),			//bfReserved1
		DecBytes(0, 2),			//bfReserved2
		DecBytes(14+40, 4),		//bfOffBits:	Specifies the offset, in bytes, from the BITMAPFILEHEADER structure to the bitmap bits. -- BITMAPFILEHEADER(14) + BITMAPINFOHEADER(40)
		DecBytes(40, 4),		//biSize:	Specifies the number of bytes required by the structure -- BITMAPINFOHEADER(40)
		DecBytes(Width, 4),		//biWidth:	Specifies the width of the bitmap, in pixels.
		DecBytes(0xFFFFFFFF-Height+1,4),//biHeight:	Specifies the height of the bitmap, in pixels. -- Height is made negative (using 2s compliment) for a top-down bitmap
		DecBytes(1, 2),			//biPlanes:	Specifies the number of planes for the target device. This value must be set to 1.
		DecBytes(BitCount, 2),		//biBitCount:	Specifies the number of bits per pixel. -- 8,16,24,32 (0,1,4 are not supported here)
		DecBytes(0, 4),			//biCompression:Specifies the type of compression for a compressed bottom-up bitmap -- Must be BI_RGB=0
		DecBytes(Data.length, 4),	//biSizeImage:	Specifies the size, in bytes, of the image. This may be set to zero for BI_RGB bitmaps.
		DecBytes(0, 4),			//biXPelsPerMeter
		DecBytes(0, 4),			//biYPelsPerMeter
		DecBytes(0, 4),			//biClrUsed:	Specifies the number of color indexes in the color table that are actually used by the bitmap. If this value is zero, the bitmap uses the maximum number of colors corresponding to the value of the biBitCount member for the compression mode specified by biCompression.
		DecBytes(0, 4),			//biClrImportant:Specifies the number of color indexes that are required for displaying the bitmap. If this value is zero, all colors are required.
		Data);				//The actual bitmap data
}

//Output a bitmap byte array to an image object
function OutputBitmapByteArrayToImage(BitmapByteArray, ImageObject) { ImageObject.src='data:image/bmp;base64,'+Encode64(BitmapByteArray); }

//---String/math conversion functions---
//Base64 quick conversion array
var B64Convert=[]; //A-Z, a-z, 0-9, +, /
for(var c=0;c<26;c++) //A-Z
	B64Convert.push(String.fromCharCode(c+65)); //65='A'
for(var c=0;c<26;c++) //a-z
	B64Convert.push(String.fromCharCode(c+97)); //97='a'
for(var c=0;c<=9;c++) //0-9
	B64Convert.push(String.fromCharCode(c+48)); //48='0'
B64Convert.push('+');
B64Convert.push('/');

//Convert a binary (octet) byte array to a base 64 string
function Encode64(ByteArray)
{
	var Output=[], Len=ByteArray.length, LeftOver;

	for(var i=0; i<Len; i+=3) //3 octet characters get converted to 4 base64 characters
	{
		//First character
		var Octet=ByteArray[i];
		Output.push(B64Convert[(Octet&0xFC)>>2]);
		LeftOver=Octet&0x03;

		//Make sure string is not complete
		if(i+1>=Len)
			break;

		//Second character
		Octet=ByteArray[i+1];
		Output.push(B64Convert[(LeftOver<<4)|((Octet&0xF0)>>4)]);
		LeftOver=Octet&0x0F;

		//Make sure string is not complete
		if(i+2>=Len)
			break;

		//Third character
		Octet=ByteArray[i+2];
		Output.push(B64Convert[(LeftOver<<2)|((Octet&0xC0)>>6)]);
		Output.push(B64Convert[(Octet&0x3F)]); //LeftOver
	}

	if(ByteArray.length%3) //If there is LeftOver characters not in a tuplet
	{
		var NumLeft=ByteArray.length%3;
		Output.push(B64Convert[LeftOver<<(NumLeft==1 ? 4 : 2)]+(NumLeft==1 ? '==' : '='));
	}

	return Output.join(''); //Return the combined output of the base64 characters into a single string
}

//function HexToBinary(HexString) { return unescape(HexString.replace(/(..)/g, '%$1')) } //Helper function to convert a hex pair string to a binary string - Use this to convert the embedded hex output of an image file

//Convert a number to a little endian numeric byte array
function DecBytes(Num, NumBytes)
{
	Num=(typeof(Num)!='number' || Num<0 ? 0 : Num); //Force Num to be a number
	Num=Math.floor(Num) % Math.pow(2,8*NumBytes); //Clamp number to positive integer and of proper bit size to avoid numeric overflow problems
	var Output=[];
	for(var i=NumBytes-1;i>0;i--) //Go over bytes backwards for proper mod and divide action
	{
		var DivideBy=Math.pow(2,8*i); //Determine the value needed to extract the byte on a divide
		var CurChar=Math.floor(Num/DivideBy); //Extract the byte
		Output.unshift(CurChar); //Push bytes onto array in little endian order
		Num-=CurChar*DivideBy; //Remove the byte from the number
	}
	Output.unshift(Num); //Push smallest byte onto the array
	return Output;
}