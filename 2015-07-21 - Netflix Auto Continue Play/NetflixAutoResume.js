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