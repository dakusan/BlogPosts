//User Variables
var DownloadURL="https://www.castledragmire.com/layout/PopupBG.png";
var PostURL="https://www.castledragmire.com/ProjectContent/WebScripts/Default_PHP_Variables.php";
var PostFileVariableName="MyFile";
var OutputFileName="Example.jpg";
//End of User Variables

DownloadFile(DownloadURL, function(DownloadedFileBlob) {
	//Get the data to send
	var Data=new FormData();
	Data.append(PostFileVariableName, DownloadedFileBlob, OutputFileName);

	//Function to run on completion
	var CompleteFunction=function(ReturnData) {
		//Add your code in this function to handle the ajax result
		var ReturnText=(ReturnData.responseText ? ReturnData : this).responseText;
		console.log(ReturnText);
	}

	//Normal AJAX example
	var Req=new XMLHttpRequest();
	Req.onload=CompleteFunction; //You can also use "onreadystatechange", which is required for some older browsers
	Req.open("POST", PostURL, true);
	Req.send(Data);

	//jQuery example
	$.ajax({type:'POST', url:PostURL, data:Data, contentType:false, processData:false, cache:false, complete:CompleteFunction});
});

function DownloadFile(
	FileURL,     //http://...
	Callback,    //The function to call back when the file download is complete. It receives the file Blob.
	ContentType) //The output Content-Type for the file. Example=image/jpeg
{
	var Req=new XMLHttpRequest();
	Req.responseType='arraybuffer';
	Req.onload=function() {
		Callback(new Blob([this.response], {type:ContentType}));
	};
	Req.open("GET", FileURL, true);
	Req.send();
}

function NormalFileSend(PostVariableName, VariableData, PostFileVariableName, FileElement, PostFileName, YOUR_URL)
{
	var DataToSend=new FormData();
	DataToSend.append(PostVariableName, VariableData); //Send a normal variable
	DataToSend.append(PostFileVariableName, FileElement.files[0], PostFileName); //Send a file
	var xhr=new XMLHttpRequest();
	xhr.open("POST", YOUR_URL, true);
	xhr.send(DataToSend);
}