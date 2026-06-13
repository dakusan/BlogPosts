<?php
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