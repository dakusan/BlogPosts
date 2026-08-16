<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Opening_IntelliJ_via_the_Symfony_ide_setting">Opening IntelliJ via the Symfony ide setting</a><br>
Date: 2018-08-30<br>
SubTitle: Nasty Escaping Problems<br>
Tags: Escaping, Examples, IntelliJ, PHP, PHPStorm, Quick Fixes, Symfony, URL Protocol Handler, Web, Web browsers, Windows, YAML<br><br>

<p>I wanted a simple setup in Symfony where the programmer could define their ide in the parameters file. Sounds simple, right? Just add something like <code>ide_url: 'phpstorm'</code> to <code>parameters.yml->parameters</code> and <code>ide: '%ide_url%'</code> to <code>config.yml->framework</code>. And it worked great, however, my problem was much more convoluted.</p>
<p>I am actually running the Symfony server on another machine and am accessing the files via NFS on Windows. So, it would try to open PHPStorm with the incorrect path. <a href="https://symfony.com/doc/3.4/reference/configuration/framework.html#ide">Symfony suggests the solution to this is writing your own custom URL handler</a> with <code>%f</code> and <code>%l</code> to fill in the filename and line, and use some weird formatting to do string replaces. So I wrote in <code>'idea://%%f:%%l&/<kbd>PROJECT_PATH_ON_SERVER</kbd>/><kbd>DRIVE_LETTER</kbd>:/<kbd>PATH_ON_WINDOWS</kbd>/'</code> (note the double parenthesis for escaping) directly in the <code>config.yml</code> and that worked, kind of. The URL was perfect, but IntelliJ does not seem to register the <code>idea://</code> protocol handler like PHPStorm theoretically does (according to some online threads) with <code>phpstorm://</code>. So I had to write my own solution.</p>
<p><a href="https://stackoverflow.com/questions/389204/how-do-i-create-my-own-url-protocol-e-g-so/7277284#7277284">This answer on stackoverflow has the answer on how to register a protocol handler in Windows</a>. But the problem now was that the first parameter passed to IntelliJ started with the <code>idea://</code> which broke the command-line file-open. So I ended up writing a script to fix this, which is at the bottom.</p>
<p>OK, so we’re almost there; I just had to paste the string I came up with back into the <code>parameters.yml</code>, right? I wish. While this was now working properly in a Symfony error page, a new problem arose. The Symfony <code>bin/console debug:config framework</code> command was failing with <code>You have requested a non-existent parameter "f:"</code>. The darn thing was reading the unescaped string as <code>'idea://%f:%l&...'</code> and it thought <code>%f:%</code> was supposed to be a variable. Sigh.</p>
<p>So the final part was to double escape the strings with 4 percent signs. <code>'idea://%%%%f:%%%%l&...'</code>. Except now the URL on the error pages gave me <code>idea://%<kbd>THE_PATH</kbd>:%<kbd>THE_LINE_NUMBER</kbd></code>. It was adding an extra parenthesis before both values. This was simple to resolve in the script I wrote, so I was finally able to open scripts directly from the error page. Yay.</p><br><br>

So here is the final set of data that has to be added to make this work:

<ins>Registry</ins>:<br>
	<code>HKCR/idea/(default)</code> = <code>URL:idea Protocol</code><br>
	<code>HKCR/idea/URL Protocol</code> = <code>""</code><br>
	<code>HKCR/idea/shell/open/command</code> = <code>"<kbd>PATH_TO_PHP</kbd>" -f "<kbd>PATH_TO_SCRIPT</kbd>" "%1" "%2" "%3" "%4" "%5" "%6" "%7" "%8" "%9"</code><br><br>

<pre><ins>parameters.yml</ins>:
	parameters:
		ide_url: 'idea://%%%%f:%%%%l&/<kbd>PROJECT_PATH_ON_SERVER</kbd>/><kbd>DRIVE_LETTER</kbd>:/<kbd>PATH_ON_WINDOWS</kbd>/'
</pre>

<pre><ins>config.yml:</ins>
	framework:
		ide: '%ide_url%'
</pre>

<ins>PHP_SCRIPT_FILE</ins>:

```php
<?php
function DoOutput($S)
{
	//You might want to do something like output the error to a file or do an alert here
	print $S;
}

if(!isset($argv[1]))
	return DoOutput('File not given');
if(!preg_match('~^idea://(?:%25|%)?([a-z]:[/\\\\][^:]+):%?(\d+)/?$~i', $argv[1], $MatchData))
	return DoOutput('Invalid format: '.$argv[1]);

$FilePath=$MatchData[1];
if(!file_exists($FilePath))
	return DoOutput('Cannot find file: '.$FilePath);

$String='"C:\Program Files\JetBrains\IntelliJ IDEA 2018.1.6\bin\idea64.exe" --line '.$MatchData[2].' '.escapeshellarg($FilePath);
DoOutput($String);
shell_exec($String);
?>
```