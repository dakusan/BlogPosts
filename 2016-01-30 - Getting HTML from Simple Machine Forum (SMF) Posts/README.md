<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Getting_HTML_from_Simple_Machine_Forum_%28SMF%29_Posts">Getting HTML from Simple Machine Forum (SMF) Posts</a><br>
Date: 2016-01-30<br>
Tags: Convert, HTML, PHP, Quick Fixes, Simple Machines Forum, SMF, WebSite<br><br>
</div>

<p>When I first created my website 10 years ago, from scratch, I did not want to deal with writing a comment system with HTML markups. And in those days, there weren’t plugins for <i>everything</i> like there is today. My solution was setting up a forum which would contain a topic for every Project, Update, and Post, and have my pages mirror the linked topic’s posts.</p>
<p>I had just put in a quick hack at the time in which the pulled <a href="http://www.simplemachines.org/">SMF</a> message’s body had links converted from bbcode (there might have been 1 other bbcode I also hooked). I had done this with regular expressions, which was a nasty hack.</p>
<p>So anywho, I finally got around to writing a script that converts SMF messages’ bbcode to HTML and caches it. You can download it <a href="GetForumTopicPosts.php">here</a>, or see the code below. The script is optimized so that it only ever needs to load SMF code when a post has not yet been cached. Caching happens during the initial loading of an SMF post within the script’s main function, and is discarded if the post is changed.</p>
<p>The script requires that you run the query on line #3 of itself in your SMF database. Directly after that are 3 variables you need to set. The script assumes you are already logged in to the appropriate user. To use it, call <code>FTP\GetForumTopicPosts($ForumTopicID)</code>. I have the functions split up so you can do individual posts too if needed (requires a little extra code).</p><br>

```php
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
```