<?php
//This SQL command must be ran before using the script
//ALTER TABLE smf_messages ADD body_html text, ADD body_md5 char(32) DEFAULT NULL;

namespace GFTP;

//Forum database variables
global $ForumInfo;
$ForumInfo=Array(
	'DBName'=>'YourDatabase_smf',
	'Location'=>'/home/YourUser/www',
	'MessageTableName'=>'smf2_messages',
);

function GetForumTopicPosts($ForumTopicID)
{
	//Change to the forum database
	global $ForumInfo;
	$CurDB=mysql_fetch_row(mysql_query('SELECT database()'))[0];
	if($CurDB!=$ForumInfo['DBName'])
		mysql_select_db($ForumInfo['DBName']);
	$OldEncoding=SetEncoding(true);

	//Get the posts
	$PostsInfos=Array();
	$PostsQuery=mysql_query('SELECT '.implode(', ', PostFields())." FROM $ForumInfo[MessageTableName] WHERE id_topic='".intval($ForumTopicID)."' AND approved=1 ORDER BY id_msg ASC LIMIT 1, 9999999");
	if($PostsQuery) //If query failed, do not process
		while(($PostInfo=mysql_fetch_assoc($PostsQuery)) && ($PostsInfos[]=$PostInfo))
			if(md5($PostInfo['body'])!=$PostInfo['body_md5']) //If the body md5s do not match, get new value, otherwise, use cached value
				ProcessPost($PostsInfos[count($PostsInfos)-1]); //Process the lastest post as a reference

	//Restore from the forum database
	if($CurDB!=$ForumInfo['DBName'])
		mysql_select_db($CurDB);
	SetEncoding(false, $OldEncoding);

	//Return the posts
	return $PostsInfos;
}

function ProcessPost(&$PostInfo) //PostInfo must have fields id_msg, body, body_md5, and body_html
{
	//Load SMF
	global $ForumInfo;
	if(!defined('SMF'))
	{
		global $context;
		require_once(rtrim($ForumInfo['Location'], DIRECTORY_SEPARATOR).DIRECTORY_SEPARATOR.'SSI.php');
		mysql_select_db($ForumInfo['DBName']);
		SetEncoding();
	}

	//Update the cached body_html field
	$ParsedCode=$PostInfo['body_html']=parse_bbc($PostInfo['body']);
	$EscapedHTMLBody=mysql_escape_string($ParsedCode);
	$BodyMD5=md5($PostInfo['body']);
	mysql_query("UPDATE $ForumInfo[MessageTableName] SET body_html='$EscapedHTMLBody', body_md5='$BodyMD5' WHERE id_msg=$PostInfo[id_msg]");
}

//The fields to select in the Post query
function PostFields() { return Array('id_msg', 'poster_time', 'id_member', 'subject', 'poster_name', 'body', 'body_md5', 'body_html'); }

//Swap character encodings. Needs to be set to utf8
function SetEncoding($GetOld=false, $NewSet=Array('utf8', 'utf8', 'utf8'))
{
	//Get the old charset if required
	$CharacterVariables=Array('character_set_client', 'character_set_results', 'character_set_connection');
	$OldSet=Array();
	if($GetOld)
	{
		//Fill in variables with default in case they are not found
		foreach($CharacterVariables as $Index => $Variable)
			$OldSet[$Variable]='utf8';

		//Query for the character sets and update the OldSet array
		$Query=mysql_query('SHOW VARIABLES LIKE "character_%"');
		while($VariableInfo=mysql_fetch_assoc($Query))
			if(isset($OldSet[$VariableInfo['Variable_name']]))
				$OldSet[$VariableInfo['Variable_name']]=$VariableInfo['Value'];

		$OldSet=array_values($OldSet); //Turn back into numerical array
	}

	//Change to the new database encoding
	$CompiledSets=Array();
	foreach($CharacterVariables as $Index => $Variable)
		$CompiledSets[$Index]=$CharacterVariables[$Index].'="'.mysql_escape_string($NewSet[$Index]).'"';
	mysql_query('SET '.implode(', ', $CompiledSets));

	//If requested, return the previous values
	return $OldSet;
}
?>