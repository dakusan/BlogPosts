<?php
//Variables
$SaveFileName='Output.html'; //Optionally save output to a file. Comment out to not save
$InputFile='Input.md';
$StylesheetLocation='github-markdown.css';
$RegexModifications=Array(
	'/<div name="user-content-(.*?)"(.*?)<\/div>/s'=>'<span id="$1"$2</span>', //Change <div name="user-contentXXX ---TO--- <span name="XXX
	'/ ?aria-hidden="true"/'=>'' //Remove aria-hidden attribute
);

//Set the curl options
$CurlHandle=curl_init(); //Init curl
curl_setopt_array($CurlHandle, Array(
	CURLOPT_URL=>           'https://api.github.com/markdown/raw', //Markdown/raw takes and returns plain text input and output
	CURLOPT_FAILONERROR=>   false,
	CURLOPT_FOLLOWLOCATION=>1,
	CURLOPT_RETURNTRANSFER=>1, //Return result as a string
	CURLOPT_TIMEOUT=>       300,
	CURLOPT_POST=>          1,
	CURLOPT_POSTFIELDS=>    file_get_contents($InputFile), //Pull in the requested file
	CURLOPT_HTTPHEADER=>    Array('Content-type: text/plain'), //Github expects the given data to be plaintext
	CURLOPT_SSL_VERIFYPEER=>0, //In case there are problems with the PHP ssl chain (often the case in Windows), ignore the error
	CURLOPT_USERAGENT=>     'Curl/PHP' //Github now requires a useragent to process the request
));

//Pull in the html converted markdown from Github
$Return=curl_exec($CurlHandle);
if(curl_errno($CurlHandle)) //Check for error
	$Return=curl_error($CurlHandle);
curl_close($CurlHandle);

//Make regex modifications
$Return=preg_replace(array_keys($RegexModifications), array_values($RegexModifications), $Return);

//Generate the final HTML. It will also be output here if not saving to a file
header('Content-Type: text/html; charset=utf-8');
if(isset($SaveFileName)) //If saving to a file, buffer output
	ob_start();
?>
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>
<head>
<title>Markdown pull</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<link href="<?=$StylesheetLocation?>" rel=stylesheet type="text/css">
</head><body><div class=GHMarkdown>
<?=$Return?>
</div></body></html>
<?php
//Save to a file if requested
if(isset($SaveFileName))
	file_put_contents($SaveFileName, ob_get_flush()); //Actual output happens here too when saving to a file
?>