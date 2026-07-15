<div>
Original Post: <a href="https://www.castledragmire.com/Posts/JavaScript_Cookies_Functions">JavaScript Cookies Functions</a><br>
Date: 2013-08-06<br>
Tags: Cookies, Examples, JavaScript<br><br>

<p>Just throwing these up here for reference. Simple JavaScript scripts to get and set cookies.  Not particularly foolproof, robust, or fully featured.</p>
</div>

```js
function SetCookie(Name, Value, SecondsToExpire) { document.cookie=Name+"="+escape(Value)+"; expires="+new Date(new Date().getTime()+SecondsToExpire*1000).toUTCString(); }
function GetCookie(Name)
{
	var Match=document.cookie.match(new RegExp('(?:^|; ?)'+Name+'=(.*?)(?:;|$)'));
	return Match ? unescape(Match[1]) : null;
}
```

<div>
</div>