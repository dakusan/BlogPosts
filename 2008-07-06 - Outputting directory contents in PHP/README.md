<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Outputting_directory_contents_in_PHP">Outputting directory contents in PHP</a><br>
Date: 2008-07-06<br>
Tags: Examples, PHP<br><br>

A friend just asked me to write a PHP function to list all the contents of a directory and its sub-directories.<br>
Nothing special here... just a simple example piece of code and boredom...
</div>

```php
function ListContents($DirName)
{
	print '<ul>';
	$dir=opendir($DirName);
	while($file=readdir($dir))
		if($file!='.' && $file!='..')
		{
			$FilePath="$DirName/$file";
			$IsDir=is_dir($FilePath);
			print "<li>$file [".($IsDir ? 'D' : number_format(filesize($FilePath), 0, '.', ',')).']';
			if($IsDir)
				ListContents($FilePath);
			print '</li>';
		}
	closedir($dir);
	print '</ul>';
}
```

<div>

It wouldn’t be a bad idea to turn off PHP’s “output buffering” and on “implicit flush” when running something like this for larger directories.<br>
<span class=ToggleText>Example output for “ListContents('c:\\temp');”:</span><div class=HideText>
<ul><li>A.BMP [230]</li><li>Dir1 [D]<ul><li>codeblocks-1.0rc2_mingw.exe [13,597,181]</li><li>Dir1a [D]<ul><li>DEBUGUI.C [25,546]</li></ul></li></ul></li><li>Dir2 [D]<ul></ul></li><li>Dir3 [D]<ul><li>HW.C [12,009]</li><li>INIFILE.C [9,436]</li></ul></li><li>NTDETECT.COM [47,564]</li></ul>
</div><br><br>

I decided to make it a little nicer afterwards by <b>bolding</b> the directories, adding their total size, and changing sizes to a human readable format. This function is a lot more memory intensive because it holds data in strings instead of immediately outputting.<br>
</div>

```php
function HumanReadableSize($Size)
{
	$MetricSizes=Array('Bytes', 'KB', 'MB', 'GB', 'TB');
	for($SizeOn=0;$Size>=1024 && $SizeOn<count($MetricSizes)-1;$SizeOn++) //Loops until Size is < a binary thousand (1,024) or we have run out of listed Metric Sizes
		$Size/=1024;
	return preg_replace('/\\.?0+$/', '', number_format($Size, 2, '.', ',')).' '.$MetricSizes[$SizeOn]; //Forces to a maximum of 2 decimal places, adds comma at thousands place, appends metric size
}

function ListContents2($DirName, &$RetSize)
{
	$Output='<ul>';
	$dir=opendir($DirName);
	$TotalSize=0;
	while($file=readdir($dir))
		if($file!='.' && $file!='..')
		{
			$FilePath="$DirName/$file";
			if(is_dir($FilePath)) //Is directory
			{
				$DirContents=ListContents2($FilePath, $DirSize);
				$Output.="<li><b>$file</b> [".HumanReadableSize($DirSize)."]$DirContents</li>";
				$TotalSize+=$DirSize;
			}
			else //Is file
			{
				$FileSize=filesize($FilePath);
				$Output.="<li>$file [".HumanReadableSize($FileSize).']</li>';
				$TotalSize+=$FileSize;
			}
		}
	closedir($dir);
	$RetSize=$TotalSize;
	$Output.='</ul>';
	return $Output;
}
```

<div>

<span class=ToggleText>Example output for “print ListContents2('c:\\temp', $Dummy);”:</span><div class=HideText>
<ul><li>A.BMP [230 Bytes]</li><li><b>Dir1</b> [12.99 MB]<ul><li>codeblocks-1.0rc2_mingw.exe [12.97 MB]</li><li><b>Dir1a</b> [24.95 KB]<ul><li>DEBUGUI.C [24.95 KB]</li></ul></li></ul></li><li><b>Dir2</b> [0 Bytes]<ul></ul></li><li><b>Dir3</b> [20.94 KB]<ul><li>HW.C [11.73 KB]</li><li>INIFILE.C [9.21 KB]</li></ul></li><li>NTDETECT.COM [46.45 KB]</li></ul>
</div><br><br>

The memory problem can be rectified through a little extra IO by calculating the size of a directory before its contents is listed, thereby not needing to keep everything in a string.<br>
</div>

```php
function CalcDirSize($DirName)
{
	$dir=opendir($DirName);
	$TotalSize=0;
	while($file=readdir($dir))
		if($file!='.' && $file!='..')
			$TotalSize+=(is_dir($FilePath="$DirName/$file") ? CalcDirSize($FilePath) :  filesize($FilePath));
	closedir($dir);
	return $TotalSize;
}

function ListContents3($DirName)
{
	print '<ul>';
	$dir=opendir($DirName);
	$TotalSize=0;
	while($file=readdir($dir))
		if($file!='.' && $file!='..')
		{
			$FilePath="$DirName/$file";
			$IsDir=is_dir($FilePath);
			$FileSize=($IsDir ? CalcDirSize($FilePath) : filesize($FilePath));
			$TotalSize+=$FileSize;
			print '<li>'.($IsDir ? '<b>' : '').$file.($IsDir ? '</b>' : '').' ['.HumanReadableSize($FileSize).']';
			if($IsDir) //Is directory
				$TotalSize+=ListContents3($FilePath);
			print '</li>';
		}
	closedir($dir);
	print '</ul>';
}
```

<div>

<span class=ToggleText>Example output: for “ListContents3('c:\\temp');”:</span><div class=HideText>
<ul><li>A.BMP [230 Bytes]</li><li><b>Dir1</b> [12.99 MB]<ul><li>codeblocks-1.0rc2_mingw.exe [12.97 MB]</li><li><b>Dir1a</b> [24.95 KB]<ul><li>DEBUGUI.C [24.95 KB]</li></ul></li></ul></li><li><b>Dir2</b> [0 Bytes]<ul></ul></li><li><b>Dir3</b> [20.94 KB]<ul><li>HW.C [11.73 KB]</li><li>INIFILE.C [9.21 KB]</li></ul></li><li>NTDETECT.COM [46.45 KB]</li></ul>
</div><br><br>

Of course, after all this, my friend took the original advice I gave him before writing any of this code, which was that using <a href="https://www.castledragmire.com/Posts/Useful_Bash_commands_and_scripts">bash commands</a> might get him to his original goal much easier.</div>