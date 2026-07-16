<?php
$OAuthResult=SendTwitterRequest('https://twitter.com/oauth/request_token');
parse_str($OAuthResult, $OauthRet);
if(!isset($OauthRet['oauth_token']))
	throw new Exception("OAuth error: $OAuthResult");
$AccessToken=$OauthRet['oauth_token'];
$AccessSecret=$OauthRet['oauth_token_secret'];