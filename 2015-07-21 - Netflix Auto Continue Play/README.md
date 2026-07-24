<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Netflix_Auto_Continue_Play">Netflix Auto Continue Play</a><br>
Date: 2015-07-21<br>
Tags: JavaScript, Netflix, Quick Fixes, Tampermonkey<br><br>

<p>Here is a little <a href="https://chrome.google.com/webstore/detail/tampermonkey/dhdgffkkebhmkfjojejmpbldmpobfkfo?hl=en-US">Tampermonkey</a> script for Chrome that automatically clicks the “Continue playing” button when it pops up on Netflix, pausing the current stream.</p><br>
</div>

```js
// ==UserScript==
// @name         Netflix auto continue play
// @namespace    https://www.castledragmire.com/Posts/Netflix_Auto_Continue_Play
// @version      1.0
// @description  When netflix pops up the "Continue play" button, this script auto-selects "Continue" within 1 second
// @author       Dakusan
// @match        http://www.netflix.com/
// @grant        none
// ==/UserScript==

setInterval(function() {
	var TheElements=document.getElementsByClassName('continue-playing');
	for(var i=0;i<TheElements.length;i++)
		if(/\bbutton\b/.test(TheElements[i].className))
		{
			console.log('"Continue Playing" Clicked');
			TheElements[i].click();
		}
}, 1000);
```

<p>Make sure to set the “User matches” in the settings page to include both “<b>http://www.netflix.com/*</b>” and “<b>https://www.netflix.com/*</b>”.</p></div>