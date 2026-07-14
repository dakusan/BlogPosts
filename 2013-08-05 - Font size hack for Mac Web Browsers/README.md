<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Font_size_hack_for_Mac_Web_Browsers">Font size hack for Mac Web Browsers</a><br>
Date: 2013-08-05<br>
Tags: Apple, Examples, Fonts, JavaScript, Quick Fixes, Web browsers<br><br>

<p>Mac renders fonts differently than windows, which is a problem when trying to make cross-system compatible web pages. I've read in many places that using "em" instead of pixels fixes this problem, but sometimes using pixel based measurements is required for a project.</p>
<p>I found when testing a recent project's web page out on Apple's OSX that no matter what browser I used, the fonts were always 2px bigger than on windows, which threw off my layouts. So I used the simple JavaScript solution below (jQuery required). Note that this assumes all font sizes are in pixels. It might not hurt to add a check for that if you mix font size types.</p>
</div>

```js
$(document).ready(function() {
	if(/Macintosh/.test(navigator.userAgent))
		$.each(document.styleSheets, function(Indx, SS) {
			var rules=SS.cssRules || SS.rules;
			for(var i=0;i<rules.length;i++)
				if(rules[i].style && rules[i].style.fontSize!='')
					rules[i].style.fontSize=(parseInt(rules[i].style.fontSize, 10)-2)+'px';
});
```

<div>
</div>