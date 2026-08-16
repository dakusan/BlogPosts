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