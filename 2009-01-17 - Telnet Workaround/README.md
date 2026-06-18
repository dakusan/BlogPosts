<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Telnet_Workaround">Telnet Workaround</a><br>
Date: 2009-01-17<br>
Tags: Examples, Quick Fixes, Servers, SSH<br><br>

<p>I recently had to do some work on a system where I was not allowed SSH/telnet access. Trying to do work strictly across ftp can take hours, especially when you have thousands of files to transfer, so I came up with a quick solution in PHP for simple command line access.</p>
</div>

```php
<form method=post action="exec.php">
	<table style="height:100%;width:100%">
		<tr><td height="10%">
			<textarea name=MyAction style="height:100%;width:100%"><?=(isset($_REQUEST['MyAction']) ? $_REQUEST['MyAction'] : '')?></textarea>
		</td></tr><tr><td height="90%">
			<textarea name=Output style="height:100%;width:100%"><?
if(isset($_REQUEST['MyAction']))
{
	$MyAction=preg_split("/\\r?\\n/", $_REQUEST['MyAction']);
	foreach($MyAction as $Action)
	{
		exec($Action, $MyOutput);
		print htmlentities(implode("\n", $MyOutput), ENT_QUOTES, 'ISO8859-1')."\n-----------------------\n";
	}
}
			?></textarea>
		</td></tr><tr><td height=1>
			<input type=submit>
		</td></tr>
	</table>
</form>
```

<div>

<p>This code allows you to enter commands on separate lines in the top box, and after the form is submitted, the output of each command is entered into the bottom box separated by dashed lines.</p>
<p>Note that  between each command the environment is reset, so commands like "cd" which change the current directory are not useable 🙁. You must also change the line 'action="exec.php"' to reflect the name you give the file.</p><br>

<p>A more suitable solution would be possible through AJAX and a program that redirected console output from a persistent session, but this was just meant as quick fix 🙂.</p></div>