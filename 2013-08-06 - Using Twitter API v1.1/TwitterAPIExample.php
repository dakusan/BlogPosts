<?php
//OAuth Request parameters - FILL THESE IN
$ConsumerKey='FILL ME IN';
$ConsumerSecret='FILL ME IN';
$AccessToken='FILL ME IN';
$AccessTokenSecret='FILL ME IN';

function EncodeParam($input) { return strtr(rawurlencode($input), Array('+'=>' ', '%7E'=>'~')); }
function SendTwitterRequest($RequestURL, $Params=Array())
{
	//Compile the OAuth parameters
	global $ConsumerKey, $ConsumerSecret, $AccessToken, $AccessTokenSecret;
	$Params=array_merge(
		$Params,
		Array('oauth_version'=>'1.0', 'oauth_nonce'=>mt_rand(), 'oauth_timestamp'=>time(), 'oauth_consumer_key'=>$ConsumerKey, 'oauth_signature_method'=>'HMAC-SHA1'),
		isset($AccessToken) ? Array('oauth_token'=>$AccessToken) : Array()
	);
	uksort($Params, 'strcmp'); //Must be sorted to determine the signature
	foreach($Params as $Key => &$Val) //Create the url encoded parameter list
		$Val=EncodeParam($Key).'='.EncodeParam($Val);
	$Params=implode('&', $Params); //Combine the parameter list
	$Params.='&oauth_signature='.EncodeParam(base64_encode(hash_hmac('sha1', 'GET&'.EncodeParam($RequestURL).'&'.EncodeParam($Params), EncodeParam($ConsumerSecret).'&'.EncodeParam($AccessTokenSecret), TRUE)));

	//Do the OAuth request
	$CurlObj=curl_init();
	foreach(Array(CURLOPT_URL=>"$RequestURL?$Params", CURLOPT_SSL_VERIFYHOST=>0, CURLOPT_SSL_VERIFYPEER=>0, CURLOPT_RETURNTRANSFER=>1) as $Key => $CurlVal)
		curl_setopt($CurlObj, $Key, $CurlVal);
	$Result=curl_exec($CurlObj);
	curl_close($CurlObj);
	return $Result;
}

//Here is an example to pull the last 4 tweets from a user:
$UserName='TheUserName'; //FILL ME IN
$Result=json_decode(SendTwitterRequest('https://api.twitter.com/1.1/statuses/user_timeline.json', Array('screen_name'=>$UserName, 'count'=>4)));
if(isset($Result->{'errors'}))
	throw new Exception($Result->{'errors'}[0]->{'message'});
$Tweets=Array();
foreach($Result as $Tweet)
	$Tweets[]=$Tweet->text;

print implode('<br>', $Tweets);