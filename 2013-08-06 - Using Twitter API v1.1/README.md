<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Using_Twitter_API_v1.1">Using Twitter API v1.1</a><br>
Date: 2013-08-06<br>
Tags: Examples, OAuth, PHP, Twitter<br><br>

<p>Twitter recently turned off their v1.0 API which broke a multitude of applications, many of which are still broken today. I had the need to immediately update at least one website to the new Twitter v1.1 API, but I could not find a simple bare bones example on the internet. Pretty much all the code/libraries out there I could find for the new API version were hundreds to thousands of lines long, or didn’t work 😣; . So anywho, here is the simple PHP code I put together to make API calls.</p>
<p>This code requires certain consumer/authentication keys and secrets. You can find how to generate them elsewhere online.</p>
</div>

```php
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
```

<br>If you don’t have an AccessToken and AccessSecret yet, you can get them through the following code:
```php
$OAuthResult=SendTwitterRequest('https://twitter.com/oauth/request_token');
parse_str($OAuthResult, $OauthRet);
if(!isset($OauthRet['oauth_token']))
	throw new Exception("OAuth error: $OAuthResult");
$AccessToken=$OauthRet['oauth_token'];
$AccessSecret=$OauthRet['oauth_token_secret'];
```

<br>Here is an example to pull the last 4 tweets from a user:
```php
$UserName='TheUserName'; //FILL ME IN
$Result=json_decode(SendTwitterRequest('https://api.twitter.com/1.1/statuses/user_timeline.json', Array('screen_name'=>$UserName, 'count'=>4)));
if(isset($Result->{'errors'}))
	throw new Exception($Result->{'errors'}[0]->{'message'});
$Tweets=Array();
foreach($Result as $Tweet)
	$Tweets[]=$Tweet->text;

print implode('<br>', $Tweets);
```