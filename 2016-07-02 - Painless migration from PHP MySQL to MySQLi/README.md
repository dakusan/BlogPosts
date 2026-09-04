<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Painless_migration_from_PHP_MySQL_to_MySQLi">Painless migration from PHP MySQL to MySQLi</a><br>
Date: 2016-07-02<br>
Tags: Examples, Legacy Technology, Migrate, MySQL, MySQLi, PHP, Programming hacks, Quick Fixes, runkit, Servers<br><br>

<p>The PHP <a href="http://php.net/manual/en/book.mysql.php">MySQL extension</a> is being deprecated in favor of the <a href="http://php.net/manual/en/book.mysqli.php">MySQLi extension</a> in PHP 5.5, and removed as of PHP 7.0. MySQLi was first referenced in <a href="http://php.net/ChangeLog-5.php#5.0.0b4">PHP v5.0.0 beta 4</a> on 2004-02-12, with the first stable release in PHP 5.0.0 on 2004-07-13<sup><a href="https://www.quora.com/When-were-MySQLi-and-PDO-released">[1]</a></sup>. Before that, the PHP MySQL extension was by far the most popular way of interacting with MySQL on PHP, and still was for a very long time after. This website was opened only 2 years after the first stable release!</p><br>

<p>With the deprecation, problems from some websites I help host have popped up, many of these sites being very, very old. I needed a quick and dirty solution to monkey-patch these websites to use MySQLi without rewriting all their code. The obvious answer is to overwrite the functions with wrappers for MySQLi. The generally known way of doing this is with <a href="http://php.net/manual/en/book.apd.php">the Advanced PHP Debugger (APD)</a>. However, using this extension has a lot of requirements that are not appropriate for a production web server. Fortunately, another extension I recently learned of offers the renaming functionality; <a href="http://php.net/manual/en/book.runkit.php">runkit</a>. It was a super simple install for me.</p>
<ol>
	<li>From the command line, run <code>pecl install runkit</code></li>
	<li>Add <code>extension=runkit.so</code> and <code>runkit.internal_override=On</code> to the php.ini</li>
</ol><br>
<p>Besides the ability to override these functions with wrappers, I also needed a way to make sure this file was always loaded before all other PHP files. The simple solution for that is adding <code>auto_prepend_file=<kbd>/PATH/TO/FILE</kbd></code> to the “.user.ini” in the user’s root web directory.</p>
<p>The code for this script is as follows. It only contains a limited set of the MySQL functions, including some very esoteric ones that the web site used. This is not a foolproof script, but it gets the job done.</p><br>
</div>

```php
<?php
//Override the MySQL functions
foreach(Array(
    'connect', 'error', 'fetch_array', 'fetch_row', 'insert_id', 'num_fields', 'num_rows',
    'query', 'select_db', 'field_len', 'field_name', 'field_type', 'list_dbs', 'list_fields',
    'list_tables', 'tablename'
) as $FuncName)
    runkit_function_redefine("mysql_$FuncName", '',
        'return call_user_func_array("mysql_'.$FuncName.'_OVERRIDE", func_get_args());');

//If a connection is not explicitely passed to a mysql_ function, use the last created connection
global $SQLLink; //The remembered SQL Link
function GetConn($PassedConn)
{
    if(isset($PassedConn))
        return $PassedConn;
    global $SQLLink;
    return $SQLLink;
}

//Override functions
function mysql_connect_OVERRIDE($Host, $Username, $Password) {
    global $SQLLink;
    return $SQLLink=mysqli_connect($Host, $Username, $Password);
}
function mysql_error_OVERRIDE($SQLConn=NULL) {
    return mysqli_error(GetConn($SQLConn));
}
function mysql_fetch_array_OVERRIDE($Result, $ResultType=MYSQL_BOTH) {
    return mysqli_fetch_array($Result, $ResultType);
}
function mysql_fetch_row_OVERRIDE($Result) {
    return mysqli_fetch_row($Result);
}
function mysql_insert_id_OVERRIDE($SQLConn=NULL) {
    return mysqli_insert_id(GetConn($SQLConn));
}
function mysql_num_fields_OVERRIDE($Result) {
    return mysqli_num_fields($Result);
}
function mysql_num_rows_OVERRIDE($Result) {
    return mysqli_num_rows($Result);
}
function mysql_query_OVERRIDE($Query, $SQLConn=NULL) {
    return mysqli_query(GetConn($SQLConn), $Query);
}
function mysql_select_db_OVERRIDE($DBName, $SQLConn=NULL) {
    return mysqli_select_db(GetConn($SQLConn), $DBName);
}
function mysql_field_len_OVERRIDE($Result, $Offset) {
    $Fields=$Result->fetch_fields();
    return $Fields[$Offset]->length;
}
function mysql_field_name_OVERRIDE($Result, $Offset) {
    $Fields=$Result->fetch_fields();
    return $Fields[$Offset]->name;
}
function mysql_field_type_OVERRIDE($Result, $Offset) {
    $Fields=$Result->fetch_fields();
    return $Fields[$Offset]->type;
}
function mysql_list_dbs_OVERRIDE($SQLConn=NULL) {
    $Result=mysql_query('SHOW DATABASES', GetConn($SQLConn));
    $Tables=Array();
    while($Row=mysqli_fetch_assoc($Result))
        $Tables[]=$Row['Database'];
    return $Tables;
}
function mysql_list_fields_OVERRIDE($DBName, $TableName, $SQLConn=NULL) {
    $SQLConn=GetConn($SQLConn);
    $CurDB=mysql_fetch_array(mysql_query('SELECT Database()', $SQLConn));
    $CurDB=$CurDB[0];
    mysql_select_db($DBName, $SQLConn);
    $Result=mysql_query("SHOW COLUMNS FROM $TableName", $SQLConn);
    mysql_select_db($CurDB, $SQLConn);
    if(!$Result) {
        print 'Could not run query: '.mysql_error($SQLConn);
        return Array();
    }
    $Fields=Array();
    while($Row=mysqli_fetch_assoc($Result))
        $Fields[]=$Row['Field'];
    return $Fields;
}
function mysql_list_tables_OVERRIDE($DBName, $SQLConn=NULL) {
    $SQLConn=GetConn($SQLConn);
    $CurDB=mysql_fetch_array(mysql_query('SELECT Database()', $SQLConn));
    $CurDB=$CurDB[0];
    mysql_select_db($DBName, $SQLConn);
    $Result=mysql_query("SHOW TABLES", $SQLConn);
    mysql_select_db($CurDB, $SQLConn);
    if(!$Result) {
        print 'Could not run query: '.mysql_error($SQLConn);
        return Array();
    }
    $Tables=Array();
    while($Row=mysql_fetch_row($Result))
        $Tables[]=$Row[0];
    return $Tables;
}
function mysql_tablename_OVERRIDE($Result) {
    $Fields=$Result->fetch_fields();
    return $Fields[0]->table;
}
```

And here is some test code to confirm functionality:
```php
<?php
global $MyConn, $TEST_Table;
$TEST_Server='localhost';
$TEST_UserName='...';
$TEST_Password='...';
$TEST_DB='...';
$TEST_Table='...';
function GetResult() {
    global $MyConn, $TEST_Table;
    return mysql_query('SELECT * FROM '.$TEST_Table.' LIMIT 1', $MyConn);
}
var_dump($MyConn=mysql_connect($TEST_Server, $TEST_UserName, $TEST_Password));
//Set $MyConn to NULL here if you want to test global $SQLLink functionality
var_dump(mysql_select_db($TEST_DB, $MyConn));
var_dump(mysql_query('SELECT * FROM INVALIDTABLE LIMIT 1', $MyConn));
var_dump(mysql_error($MyConn));
var_dump($Result=GetResult());
var_dump(mysql_fetch_array($Result));
$Result=GetResult(); var_dump(mysql_fetch_row($Result));
$Result=GetResult(); var_dump(mysql_num_fields($Result));
var_dump(mysql_num_rows($Result));
var_dump(mysql_field_len($Result, 0));
var_dump(mysql_field_name($Result, 0));
var_dump(mysql_field_type($Result, 0));
var_dump(mysql_tablename($Result));
var_dump(mysql_list_dbs($MyConn));
var_dump(mysql_list_fields($TEST_DB, $TEST_Table, $MyConn));
var_dump(mysql_list_tables($TEST_DB, $MyConn));
mysql_query('CREATE TEMPORARY TABLE mysqltest (i int auto_increment, primary key (i))', $MyConn);
mysql_query('INSERT INTO mysqltest VALUES ()', $MyConn);
mysql_query('INSERT INTO mysqltest VALUES ()', $MyConn);
var_dump(mysql_insert_id($MyConn));
mysql_query('DROP TEMPORARY TABLE mysqltest', $MyConn);
```