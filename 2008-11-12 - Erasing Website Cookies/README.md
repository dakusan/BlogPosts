<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Erasing_Website_Cookies">Erasing Website Cookies</a><br>
Date: 2008-11-12<br>
SubTitle: A quick useful code snippet because it takes way too long to do this through normal browser means<br>
Tags: Cookies, Examples, JavaScript, Web browsers<br><br>

This erases all cookies on the current domain (in the “ / ” path)<br><br>
JavaScript:
</div>

```js
function ClearCookies() //Clear all the cookies on the current website
{
	var MyCookies=document.cookie; //Remember the original cookie string since it will be changing soon
	var StartAt=0; //The current string pointer in MyCookies
	do //Loop through all cookies
	{
		var CookieName=MyCookies.substring(StartAt, MyCookies.indexOf('=', StartAt)).replace(/^ /,''); //Get the next cookie name in the list, and strip off leading white space
		document.cookie=CookieName+"=;expires=Thu, 01-Jan-1970 00:00:01 GMT; path=/"; //Erase the cookie
		StartAt=MyCookies.indexOf(';', StartAt)+1; //Move the string pointer to the end of the current cookie
	} while(StartAt!=0)
}
```

<div>
<br>

I went a little further with the script after finishing this to add a bit of a visual aspect.<br>
The following adds a textarea box which displays the current cookies for the site, and also displays the cookie names when they are erased.
</div>

```html
<input type=button value="Clear Cookies" onclick="ClearCookies()">
<input type=button value="View Cookies" onclick="ViewCookies()">
<textarea id=CookieBox style="width:100%;height:100%"></textarea>
<script type="text/javascript">
function ViewCookies() //Output the current cookies in the textbox
{
	document.getElementById('CookieBox').value=document.cookie.replace(/;/g,';\n\n');
}

function ClearCookies() //Clear all the cookies on the current website
{
	var CookieNames=[]; //Remember the cookie names as we erase them for later output
	var MyCookies=document.cookie; //Remember the original cookie string since it will be changing soon
	var StartAt=0; //The current string pointer in MyCookies
	do //Loop through all cookies
	{
		var CookieName=MyCookies.substring(StartAt, MyCookies.indexOf('=', StartAt)).replace(/^ /,''); //Get the next cookie name in the list, and strip off leading white space
		CookieNames.push(CookieName); //Remember the cookie name
		document.cookie=CookieName+"=;expires=Thu, 01-Jan-1970 00:00:01 GMT; path=/"; //Erase the cookie
		StartAt=MyCookies.indexOf(';', StartAt)+1; //Move the string pointer to the end of the current cookie
	} while(StartAt!=0)
	document.getElementById('CookieBox').value='Clearing: '+CookieNames.join("\nClearing: "); //Output the erased cookie names
}
</script>
```

<div>
<br>

<a href="LiveExample.html">Live Example HTML</a><br>
<a href="https://www.castledragmire.com/Posts/Erasing_Website_Cookies">Live Example</a><br>
</div>