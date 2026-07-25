<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Sending_URLs_as_a_file_in_an_HTML_form_using_AJAX">Sending URLs as a file in an HTML form using AJAX</a><br>
Date: 2015-10-03<br>
Tags: Data Passing, Examples, HTML, JavaScript, Standards, Web browsers, XSS<br><br>
</div>

It is common knowledge that you can use the FormData class to send a file via AJAX as follows:
```js
function NormalFileSend(PostVariableName, VariableData, PostFileVariableName, FileElement, PostFileName, YOUR_URL)
{
	var DataToSend=new FormData();
	DataToSend.append(PostVariableName, VariableData); //Send a normal variable
	DataToSend.append(PostFileVariableName, FileElement.files[0], PostFileName); //Send a file
	var xhr=new XMLHttpRequest();
	xhr.open("POST", YOUR_URL, true);
	xhr.send(DataToSend);
}
```

Something that is much less known, which doesn't have any really good full-process examples online (that I could find), is sending a URL's file as the posted file.<br>
This is doable by downloading the file as a Blob, and then directly passing that blob to the FormData. The 3rd parameter to the FormData.append should be the file name.<br><br>

The following code demonstrates downloading the file. I did not worry about adding error checking.
```js
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
```

And the following code demonstrates submitting that file
```js
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
```

<p>Unfortunately, due to cross site scripting (XSS) security settings, you can generally only use ajax to query URLs on the same domain. I use my <a href="https://www.github.com/dakusan/https://github.com/dakusan/HTTP-Forwarders">Cross site scripting solutions</a> and <a href="https://github.com/dakusan/HTTP-Forwarders">HTTP Forwarders</a> for this. Stackoverflow also has <a href="https://stackoverflow.com/questions/15005500/loading-cross-domain-html-page-with-jquery-ajax">a good thread</a> about it.</p></div>