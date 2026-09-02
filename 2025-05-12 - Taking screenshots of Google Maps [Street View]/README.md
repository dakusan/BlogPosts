<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Taking_screenshots_of_Google_Maps_%5BStreet_View%5D">Taking screenshots of Google Maps [Street View]</a><br>
Date: 2025-05-12<br>
Tags: Google Maps, JavaScript, Quick Fixes, Screenshot, Web browsers<br><br>
</div>

The following are the steps to take a screenshot without overlays inside Google Maps Street View or Google Maps.
<ol>
	<li>Collapse the left side window if visible with the arrow</li>
	<li>Press F11 to go full screen</li>
	<li>Open up the dev console [F12]<br>
	Optional: You may need to open the dev console in another window (undocked). Leaving it docked with the window could cause problems.
	</li>
	<li>Run the following command in the dev console:<br>
</li></ol>

```js
document.querySelectorAll('*').forEach(el => el.style.display='none');
for(let c=document.getElementsByTagName('canvas')[0];c.parentNode;c=c.parentNode) c.style.display='';
```

<ol>
	This should be pretty foolproof, but you may need to change the <code>[0]</code> to another number if the primary canvas element isn't the first canvas.<br><br>
	If that doesn't work, you can also try the following, which is a bit more specialized, but may not catch all the elements on the screen (including the Google logo).
</ol>

```js
document.querySelectorAll('[jsaction^="mouseover:"],[jsaction^="mouseout:"],button,[role="navigation"],.scene-footer-container,#omnibox-container,[aria-label="Sign in"],[aria-label="Google apps"]')
	.forEach(el => el.style.display = 'none');
document.querySelectorAll('[aria-label="Interactive map"]').forEach(el => el.parentNode.style.display='none');
```

<ol start=5>
	<li>Close the dev console (or Alt+tab back to the chrome window if undocked)</li>
	<li>Move the mouse to the sky (or a portion that won't trigger possible popups)</li>
	<li>Wait for a few seconds for icons to disappear</li>
	<li>Take screenshot</li>
</ol></div>