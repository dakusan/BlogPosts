<div>
Original Post: <a href="https://www.castledragmire.com/Posts/MySQL_replication_ring_status_reporting_script">MySQL replication ring status reporting script</a><br>
Date: 2015-04-17<br>
Tags: Master, MySQL, Replication, Reporting, Ring, Servers, Status<br><br>

<p>I just threw together a quick script to report status on a MySQL replication ring. While replication rings have been the only real multi-master MySQL solution for replication (with the ability for nodes to go down without majorly breaking things) until recently, I have read that MariaDB (still not MySQL) now allows a slave to have multiple masters, meaning many replication topologies are now possible (star, mesh, etc). This script could easily be adapted for those circumstances too.</p>
<p>This script will report all the variables from “SHOW MASTER STATUS” and “SHOW SLAVE STATUS” from all servers in your replication ring, in a unified table. It also includes a “Pretty Status” row that lets you quickly see how things look. The possibilities for this row are:</p>
<ul>
	<li>
		<b>Bad state: ...</b><br>
		This shows if the Slave_IO_State is not “Waiting for master to send event”
	</li>
	<li>
		<b>Cannot determine master’s real position</b><br>
		This shows if the Position variable on the master could not be read
	</li>
	<li>
		<b>On old master file</b><br>
		This shows if the slave’s “Master_Log_File” variable does not equal the master’s “File” variable
	</li>
	<li>
		<b>Bytes behind: xxx</b><br>
		This shows if none of the above errors occurred. It subtracts the master’s “Position” from the slave’s “Read_Master_Log_Pos”. This should generally be at or around 0. A negative value essentially means 0 (this should only happen between the last and first server).
	</li>
</ul>
<p>The “Seconds_Behind_Master” variable can also be useful for determining the replication ring’s current replication status.</p>
<p>The code is below the example. The entire source file can also be found <a href="ReplicationStatus.php">here</a>. The 3 variables that need to be configured are at the top of the file. It assumes that all servers are accessible via the single given username and password.</p><br>
<div><b><u>Example:</u></b></div>

<table>
	<tr><td colspan=3 align=center><b>Master</b></td></tr>
	<tr><th align=left>Server Name</th><td><b>EXAMPLE1.MYDOMAIN.COM</b></td><td><b>EXAMPLE2</b></td></tr>
	<tr><th align=left>File</th><td>mysql-bin.000003</td><td>mysql-bin.000011</td></tr>
	<tr><th align=left>Position</th><td>25249746</td><td>3215834</td></tr>
	<tr><th align=left>Binlog_Do_DB</th><td>example_data,devexample_data</td><td>example_data,devexample_data</td></tr>
	<tr><th align=left>Binlog_Ignore_DB</th><td></td><td></td></tr>
	<tr><td colspan=3 align=center><b>Slave</b></td></tr>
	<tr><th align=left>Pretty Status</th><td><b>Bytes behind: 0</b></td><td><b>Bytes behind: 0</b></td></tr>
	<tr><th align=left>Slave_IO_State</th><td>Waiting for master to send event</td><td>Waiting for master to send event</td></tr>
	<tr><th align=left>Master_Host</th><td>EXAMPLE2</td><td>EXAMPLE1.MYDOMAIN.COM</td></tr>
	<tr><th align=left>Master_User</th><td>example_slave</td><td>example_slave</td></tr>
	<tr><th align=left>Master_Port</th><td>3306</td><td>3306</td></tr>
	<tr><th align=left>Connect_Retry</th><td>60</td><td>60</td></tr>
	<tr><th align=left>Master_Log_File</th><td>mysql-bin.000011</td><td>mysql-bin.000003</td></tr>
	<tr><th align=left>Read_Master_Log_Pos</th><td>3215834</td><td>25249746</td></tr>
	<tr><th align=left>Relay_Log_File</th><td>www-relay-bin.070901</td><td>www-relay-bin.071683</td></tr>
	<tr><th align=left>Relay_Log_Pos</th><td>252</td><td>252</td></tr>
	<tr><th align=left>Relay_Master_Log_File</th><td>mysql-bin.000011</td><td>mysql-bin.000003</td></tr>
	<tr><th align=left>Slave_IO_Running</th><td>Yes</td><td>Yes</td></tr>
	<tr><th align=left>Slave_SQL_Running</th><td>Yes</td><td>Yes</td></tr>
	<tr><th align=left>Replicate_Do_DB</th><td>example_data,devexample_data</td><td>example_data,devexample_data</td></tr>
	<tr><th align=left>Replicate_Ignore_DB</th><td></td><td></td></tr>
	<tr><th align=left>Replicate_Do_Table</th><td></td><td></td></tr>
	<tr><th align=left>Replicate_Ignore_Table</th><td></td><td></td></tr>
	<tr><th align=left>Replicate_Wild_Do_Table</th><td></td><td></td></tr>
	<tr><th align=left>Replicate_Wild_Ignore_Table</th><td></td><td></td></tr>
	<tr><th align=left>Last_Errno</th><td>0</td><td>0</td></tr>
	<tr><th align=left>Last_Error</th><td></td><td></td></tr>
	<tr><th align=left>Skip_Counter</th><td>0</td><td>0</td></tr>
	<tr><th align=left>Exec_Master_Log_Pos</th><td>3215834</td><td>25249746</td></tr>
	<tr><th align=left>Relay_Log_Space</th><td>552</td><td>552</td></tr>
	<tr><th align=left>Until_Condition</th><td>None</td><td>None</td></tr>
	<tr><th align=left>Until_Log_File</th><td></td><td></td></tr>
	<tr><th align=left>Until_Log_Pos</th><td>0</td><td>0</td></tr>
	<tr><th align=left>Master_SSL_Allowed</th><td>No</td><td>No</td></tr>
	<tr><th align=left>Master_SSL_CA_File</th><td></td><td></td></tr>
	<tr><th align=left>Master_SSL_CA_Path</th><td></td><td></td></tr>
	<tr><th align=left>Master_SSL_Cert</th><td></td><td></td></tr>
	<tr><th align=left>Master_SSL_Cipher</th><td></td><td></td></tr>
	<tr><th align=left>Master_SSL_Key</th><td></td><td></td></tr>
	<tr><th align=left>Seconds_Behind_Master</th><td>0</td><td>0</td></tr>
	<tr><th align=left>Master_SSL_Verify_Server_Cert</th><td>No</td><td>No</td></tr>
	<tr><th align=left>Last_IO_Errno</th><td>0</td><td>0</td></tr>
	<tr><th align=left>Last_IO_Error</th><td></td><td></td></tr>
	<tr><th align=left>Last_SQL_Errno</th><td>0</td><td>0</td></tr>
	<tr><th align=left>Last_SQL_Error</th><td></td><td></td></tr>
	<tr><th align=left>Replicate_Ignore_Server_Ids</th><td></td><td><div class=NotGiven>Not given</div></td></tr>
	<tr><th align=left>Master_Server_Id</th><td>2</td><td><div class=NotGiven>Not given</div></td></tr>
</table><br><br>

</div>

<b><u>Code:</u></b><br>
```php
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
```

<p>One final note. When having this script run, you might need to make sure none of the listed server IPs evaluates to localhost (127.x.x.x), as MySQL may instead then use the local socket pipe, which may not work with users who only have <a href="https://dev.mysql.com/doc/refman/5.0/en/privileges-provided.html#priv_replication-slave">REPLICATION permissions</a> and a wildcard host.</p>