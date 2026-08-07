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