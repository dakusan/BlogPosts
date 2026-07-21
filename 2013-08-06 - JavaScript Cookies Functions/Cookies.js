function SetCookie(Name, Value, SecondsToExpire) { document.cookie=Name+"="+escape(Value)+"; expires="+new Date(new Date().getTime()+SecondsToExpire*1000).toUTCString(); }
function GetCookie(Name)
{
	var Match=document.cookie.match(new RegExp('(?:^|; ?)'+Name+'=(.*?)(?:;|$)'));
	return Match ? unescape(Match[1]) : null;
}