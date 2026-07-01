<div>
Original Post: <a href="https://www.castledragmire.com/Posts/UTF8_BOM">UTF8 BOM</a><br>
Date: 2010-11-02<br>
Tags: BOM, Examples, Linux, PHP, Standards, Unicode, UTF8, Web browsers, Windows<br><br>

<p>While <a href="http://en.wikipedia.org/wiki/UTF-8">UTF-8</a> has almost universally been accepted as the de-facto standard for <a href="http://en.wikipedia.org/wiki/Unicode">Unicode</a> character encoding in most non-Windows systems (mmmmmm <a href="http://en.wikipedia.org/wiki/Plan9">Plan 9</a> ^_^), the <a href="http://en.wikipedia.org/wiki/Byte-order_mark">BOM (Byte Order Marker)</a> still has large adoption problems. While I have been allowing my text editors to add the UTF8 BOM to the beginning of all my text files for years, I have finally decided to rescind this practice for compatibility reasons.</p>
<p>While the UTF8 BOM is useful so that editors know for sure what the character encoding of a file is, and don’t have to guess, they are not really supported, for their reasons, in Unixland. Having to code solutions around this was becoming cumbersome. Programs like vi and pico/nano seem to ignore a file’s character encoding anyways and adopt the character encoding of the current terminal session.</p>
<p>The main culprit in which I was running into this problem a lot with is PHP. The funny thing about it too was that I had a solution for it working properly in Linux, but not Windows 🙂.</p>
<p>Web browsers do not expect to receive the BOM marker at the beginning of files, and if they encounter it, may have serious problems. For example, in a certain browser (*cough*IE*cough*) having a BOM on a file will cause the browser to not properly read the DOCTYPE, which can cause all sorts of nasty compatibility issues.</p>
<p>Something in my <a href="http://en.wikipedia.org/wiki/LAMP_%28software_bundle%29">LAMP</a> setup on my cPanel systems was removing the initial BOM at the beginning of outputted PHP contents, but through some preliminary research I could not find out why this was not occurring in Windows. However, both systems were receiving multiple BOMs at the beginning of the output due to PHP’s include/require functions not stripping the BOM from those included files. My solution to this was a simple overload of these include functions as follows (only required when called from any directly opened [non-included] PHP file):</p>
</div>

```php
<?
/*Safe include/require functions that make sure UTF8 BOM is not output
Use like: eval(safe_INCLUDETYPE($INCLUDE_FILE_NAME));
where INCLUDETYPE is one of the following: include, require, include_once, require_once
An eval statement is used to maintain current scope
*/

//The different include type functions
function safe_include($FileName)	{ return real_safe_include($FileName, 'include'); }
function safe_require($FileName)	{ return real_safe_include($FileName, 'require'); }
function safe_include_once($FileName)	{ return real_safe_include($FileName, 'include_once'); }
function safe_require_once($FileName)	{ return real_safe_include($FileName, 'require_once'); }

//Start the processing and return the eval statement
function real_safe_include($FileName, $IncludeType)
{
	ob_start();
	return "$IncludeType('".strtr($FileName, Array("\\"=>"\\\\", "'", "\\'"))."'); safe_output_handler();";
}

//Do the actual processing and return the include data
function safe_output_handler()
{
	$Output=ob_get_clean();
	while(substr($Output, 0, 3)=='?') //Remove all instances of UTF8 BOM at the beginning of the output
		$Output=substr($Output, 3);
	print $Output;
}
?>
```

<div>

<p>I would have like to have used PHP’s <a href="http://se.php.net/manual/en/outcontrol.configuration.php#ini.output-handler">output_handler</a> ini setting to catch even the root file’s BOM and not require include function overloads, but, as php.net puts it “Only built-in functions can be used with this directive. For user defined functions, use ob_start().”.</p>
<p>As a bonus, the following bash command can be used to find all PHP files in the current directory tree with a UTF8 BOM:</p>
</div>

```bash
find . -name '*.php' -print0 | xargs -0 -n1000 grep -l $'^\xef\xbb\xbf'
```

And to remove the BOMs (Cygwin compatible):
```bash
find . -name '*.php' -print0 | xargs -0 -n1000 grep -l $'^\xef\xbb\xbf' | xargs -i perl -i.bak -pe 'BEGIN{ @d=@ARGV } s/^\xef\xbb\xbf//; END{ unlink map "$_$^I", @d }' "{}"
```

Simpler remove BOMs (not Cygwin/Windows compatible):
```bash
find . -name '*.php' -print0 | xargs -0 -n1000 grep -l $'^\xef\xbb\xbf' | xargs -i perl -i -pe 's/^\xef\xbb\xbf//' "{}"
```