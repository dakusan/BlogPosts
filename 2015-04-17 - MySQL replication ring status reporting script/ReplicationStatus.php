<?
//Configurations
$Servers=Array('SERVER1.YOURDOMAIN.COM', 'SERVER2.YOURDOMAIN.COM'); //List of host names to access mysql servers on. This must be in the order of the replication ring.
$SlaveUserName='SLAVE_RING_USERNAME'; //This assumes all servers are accessible via this username with the same password
$SlavePassword='SLAVE_RING_PASSWORD';

//Get the info for each server
$ServersInfo=Array(); //SERVER_NAME=>Array('Master'=>Array(Col1=>Val1, ...), 'Slave'=>Array(Col1=>Val1, ...)
$ColsNames=Array('Master'=>Array('Server Name'=>0), 'Slave'=>Array('Pretty Status'=>0)); //The column names for the 2 (master and slave) queries. Custom column names are also added here
$CustomFieldNames=array_merge($ColsNames['Master'], $ColsNames['Slave']); //Store the custom column names so they are not HTML escaped later
foreach($Servers as $ServerName)
{
	//Connect to the server
	$Link=@new mysqli($ServerName, $SlaveUserName, $SlavePassword);
	if($Link->connect_error)
		die(EHTML("Connection error to $ServerName server: $Link->connect_error"));

	//Get the replication status info from the server
	$MyServerInfo=$ServersInfo[$ServerName]=Array(
		'Master'=>$Link->Query('SHOW MASTER STATUS')->fetch_array(MYSQLI_ASSOC),
		'Slave'=>$Link->Query('SHOW SLAVE STATUS')->fetch_array(MYSQLI_ASSOC)
	);
	mysqli_close($Link); //Close the connection

	//Gather the column names
	foreach($ColsNames as $ColType => &$ColNames)
		foreach($MyServerInfo[$ColType] as $ColName => $Dummy)
			$ColNames[$ColName]=0;
}
unset($ColNames);

//Gather the pretty statuses
foreach($Servers as $Index => $ServerName)
{
	//Determine the pretty status
	$SlaveInfo=$ServersInfo[$ServerName]['Slave'];
	$MasterInfo=$ServersInfo[$Servers[($Index+1)%count($Servers)]]['Master'];
	if($SlaveInfo['Slave_IO_State']!='Waiting for master to send event')
		$PrettyStatus='Bad state: '.EHTML($SlaveInfo['Slave_IO_State']);
	else if(!isset($MasterInfo['Position']))
		$PrettyStatus='Cannot determine master’s real position';
	else if($SlaveInfo['Master_Log_File']!=$MasterInfo['File'])
		$PrettyStatus='On old master file';
	else
		$PrettyStatus='Bytes behind: '.($MasterInfo['Position']-$SlaveInfo['Read_Master_Log_Pos']);

	//Add the server name and pretty status to the output columns
	$ServersInfo[$ServerName]['Master']['Server Name']='<div class=ServerName>'.EHTML($ServerName).'</div>';
	$ServersInfo[$ServerName]['Slave']['Pretty Status']='<div class=PrettyStatus>'.EHTML($PrettyStatus).'</div>';
}

//Output the document
function EHTML($S) { return htmlspecialchars($S, ENT_QUOTES, 'UTF-8'); } //Escape HTML
?>
<!DOCTYPE html>
<html>
<head>
	<title>Replication Status</title>
	<meta charset="UTF-8">
	<style>
		table { border-collapse:collapse; }
		table tr>* { border:1px solid black; padding:3px; }
		th { text-align:left; font-weight:bold; }
		.ReplicationDirectionType { font-weight:bold; text-align:center; color:blue; }
		.ServerName { font-weight:bold; text-align:center; color:red; }
		.PrettyStatus { font-weight:bold; color:red; }
		.NotGiven { font-weight:bold; }
	</style>
</head>
<body><table>
<?
//Output the final table
foreach($ColsNames as $Type => $ColNames) //Process by direction type (Master/Slave) then columns
{
	print '<tr><td colspan='.(count($Servers)+1).' class=ReplicationDirectionType>'.$Type.'</td></tr>'; //Replication direction (Master/Server) type title column
	foreach($ColNames as $ColName => $Dummy) //Process each column name individually
	{
		print '<tr><th>'.EHTML($ColName).'</th>'; //Column name
		$IsHTMLColumn=isset($CustomFieldNames[$ColName]); //Do not escape HTML on custom fields
		foreach($ServersInfo as $ServerInfo) //Output the column for each server
			if($IsHTMLColumn) //Do not escape HTML on custom fields
				print '<td>'.$ServerInfo[$Type][$ColName].'</td>';
			else //If not a custom field, output the escaped HTML of the value. If the column does not exist for this server (different mysql versions), output "Not given"
				print '<td>'.(isset($ServerInfo[$Type][$ColName]) ? EHTML($ServerInfo[$Type][$ColName]) : '<div class=NotGiven>Not given</div>').'</td>';
		print '</tr>';
	}
}
?>
</table></body>
</html>