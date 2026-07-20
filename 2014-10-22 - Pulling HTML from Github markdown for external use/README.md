<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Pulling_HTML_from_Github_markdown_for_external_use">Pulling HTML from Github markdown for external use</a><br>
Date: 2014-10-22<br>
SubTitle: Although, converting to markdown is a time consuming pain<br>
Tags: Convert, Examples, Formats, Github, HTML, Markdown, PHP<br><br>

<p>So I started getting on the Github bandwagon <i>FINALLY</i>. I figured that while I was going to the trouble of remaking readme files for the projects into github markdown files, I might as well duplicate the compiled HTML for my website.</p>
<p>The below code is a simple PHP script to pull in the converted HTML from Github’s API and then do some more modifications to facilitate directly inserting it into a website.</p><br>
Usage:
<ul>
	<li>The variables that can be updated are all at the top of the file.</li>
	<li>The script will always output the finished result to the user’s browser, but can also optionally save it to an external file by setting the <b>$SaveFileName</b> variable.</li>
	<li>Stylesheet:<ul>
		<li>The script automatically includes a specified stylesheet from the <b>$StylesheetLocation</b> variable.</li>
		<li>The stylesheet I used is from <a href="https://gist.github.com/somebox/1082608">https://gist.github.com/somebox/1082608</a>. I’m not too happy with its coloring scheme, but it’ll do for now.</li>
		<li>The required modifications that need to be made to the css are to change “body” to “.GHMarkdown”, and then add “.GHMarkdown” before all other rules.</li>
		<li><a href="https://www.castledragmire.com/layoutx/github-markdown.css">This</a> is the one I am currently using for my website, but it also has a few modifications made specifically for my layouts.</li>
	</ul></li>
	<li>Modifications<ul>
		<li>In my markdowns, I like to link to internal sections by first creating a bookmark as “&lt;div name="BOOKMARK_NAME"&gt;...&lt;/div&gt;” and then linking via “[LinkName](#BOOKMARK_NAME)”. While this works on github, the bookmark’s names are actually changed to something like “user-content-BOOKMARK-NAME”, which is not useable outside of github. The first <b>$RegexModifications</b> item therefore updates the bookmarks back to their original name, and turns them into &lt;span&gt;s (which github does not support).</li>
		<li>The second rule just removes the “aria-hidden” attributes, which my W3C checking scripts throw a warning on.</li>
	</ul></li>
	<li>Note that sometimes, the script may return an error of “transfer closed with XXX bytes remaining to read”. This means that github denied the request (probably due to too many requests in too short a timespan), but the input is too large so github prematurely terminated the connection. If this happens, try sending a tiny input and see if you get back a proper error.</li>
</ul><br>

<p><b>Edit on 2026-07-20</b>: This has basically been superseded by my <a href="https://github.com/dakusan/MarkDown-To-Nexus-BBCode-and-HTML">MarkDown To Nexus BBCode and HTML</a> project.</p>

</div>

```php
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
```