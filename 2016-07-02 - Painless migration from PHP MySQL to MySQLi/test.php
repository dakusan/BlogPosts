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