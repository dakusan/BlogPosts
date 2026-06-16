<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Data_Format_Conversion">Data Format Conversion</a><br>
Date: 2008-07-14<br>
SubTitle: Moving from Point A to Point B<br>
Tags: Data Passing, Examples, MySQL, PHP<br><br>

<p>I am often asked to transfer data sets into MySQL databases, or other formats. In this case, I’ll use a Microsoft Excel file without line breaks in the fields to MySQL as an example. While there are many programs out there to do this kind of thing, this method doesn’t take too long and is a good example use of <a href="https://www.castledragmire.com/Posts/Regular_Expression_Examples">regular expressions</a>.</p><br>

<p>First, select all the data in Excel (ctrl+a) and copy (ctrl+c) it to a text editor with regular expression support. I recommend <a href="http://www.editpadpro.com/">EditPad Pro</a> as a very versatile and powerful text editor.</p>
<p>Next, we need to turn each row into the format “('FIELD1','FIELD2','FIELD3',<b style="color:red">...</b>),”. Four regular expressions are needed to format the data:</p>
<link href="layoutx/RegularExpressions.css" rel=stylesheet type="text/css">
<table class=SpaceTable><tr><th>Search</th><th>Replace</th><th>Explanation</th></tr>
<tr><td>'</td><td><span class=REReplaceSp>\\</span>'</td><td>Escape single quotes</td></tr>
<tr><td><span class=REEscape>\t</span></td><td>','</td><td>Separate fields and quote as strings</td></tr>
<tr><td><span class=RESpecial>^</span></td><td>('</td><td>Start of row</td></tr>
<tr><td><span class=RESpecial>$</span></td><td>'),</td><td>End of row</td></tr>
</table>
From there, there are only 2 more steps to complete the query.
<ul>
	<li>Add the start of the query: “INSERT INTO <b>TABLENAME</b> VALUES”</li>
	<li>End the query by changing the last row's comma “,” at the very end of the line to a semi-colon “;”.</li>
</ul><br>

For example:
</div>

```
a	b	c
d	e	f
g	h	i
```

<div>

would be converted to
</div>

```sql
INSERT INTO MyTable VALUES
('a','b','c'),
('d','e','f'),
('h','h','i');
```

<div>
<br>

<p>Sometimes queries may get too long and you will need to separate them by performing the “2 more steps to complete the query” from above.</p><br>

<p>After doing one of these conversions recently, I was also asked to make the data searchable, so I made a very simple PHP script for this.</p>
<p>This script lets you search through all the fields and lists all matches. The fields are listed on the 2nd line in an array as "SQL_FieldName"=>"Viewable Name".  If the “Viewable Name” contains a pound sign “#” it is matched exactly, otherwise, only part of the search string needs to be found.</p>
</div>

```php
<?
$Fields=Array('ClientNumber'=>'Client #', 'FirstName'=>'First Name', 'LastName'=>'Last Name', <b style="color:red">...</b>); //Field list
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
mysql_connect('<b>SQL_HOST</b>', '<b>SQL_USERNAME</b>', '<b>SQL_PASSWORD</b>');
mysql_select_db('<b>SQL_DATABASE</b>');
$q=mysql_query('SELECT * FROM <b>TABLENAME</b> WHERE '.implode(' AND ', $SearchArray));

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
```

<div>
</div>