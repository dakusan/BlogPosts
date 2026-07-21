<?
$Fields=Array('ClientNumber'=>'Client #', 'FirstName'=>'First Name', 'LastName'=>'Last Name', ...); //Field list
print '<form method=post action=index.php><table>'; //Form action needs to point to the current file
foreach($Fields as $Name => $Value) //Output search text boxes
	print "<tr><td>$Value</td><td><input name=\"$Name\" style='width:200px;' value=\"".
		(isset($_POST[$Name]) ? htmlentities($_POST[$Name], ENT_QUOTES) : '').'"></td></tr>';//Text boxes w/ POSTed values,if set
print '</table><input type=submit value=Search></form>';

if(!isset($_POST[key($Fields)])) //If search data has not been POSTed, stop here
	return;
	
$SearchArray=Array('1=1'); //Search parameters are stored here. 1=1 is passed in case no POSTed search parameter are ...
                           //... requested so there is at least 1 WHERE parameter, and is optimized out with the MySQL preprocessor anyways.
foreach($Fields as $Name => $Value) //Check each POSTed search parameter
	if(trim($_POST[$Name])!='') //If the POSTed search parameter is empty, do not use it as a search parameter
	{
		$V=mysql_escape_string($_POST[$Name]); //Prepare for SQL insertion
		$SearchArray[]=$Name.(strpos($Value, '#')===FALSE ? " LIKE '%$V%'" : "='$V'"); //Pound sign in the Viewable Name=exact ...
			//... value, otherwise, just a partial patch
	}
//Get data from MySQL
mysql_connect('SQL_HOST', 'SQL_USERNAME', 'SQL_PASSWORD');
mysql_select_db('SQL_DATABASE');
$q=mysql_query('SELECT * FROM TABLENAME WHERE '.implode(' AND ', $SearchArray));

//Output retrieved data
$i=0;
while($d=mysql_fetch_assoc($q)) //Iterate through found rows
{
	if(!($i++)) //If this is the first row found, output header
	{
		print '<table border=1 cellpadding=0 cellspacing=0><tr><td>Num</td>'; //Start table and output first column header (row #)
		foreach($Fields as $Name => $Value) //Output the rest of the column headers (Viewable Names)
			print "<td>$Value</td>";
		print '</tr>'; //Finish header row
	}
	print '<tr bgcolor='.($i&1 ? 'white' : 'gray')."><td>$i</td>"; //Start the data field's row. Row's colors are alternating white and gray.
	foreach($Fields as $Name => $Value) //Output row data
		print '<td>'.$d[$Name].'</td>';
	print '</tr>'; //End data row
}

print ($i==0 ? 'No records found.' : '</table>'); //If no records are found, output an error message, otherwise, end the data table
?>