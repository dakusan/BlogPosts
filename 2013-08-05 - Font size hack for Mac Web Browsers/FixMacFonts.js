$(document).ready(function() {
	if(/Macintosh/.test(navigator.userAgent))
		$.each(document.styleSheets, function(Indx, SS) {
			var rules=SS.cssRules || SS.rules;
			for(var i=0;i<rules.length;i++)
				if(rules[i].style && rules[i].style.fontSize!='')
					rules[i].style.fontSize=(parseInt(rules[i].style.fontSize, 10)-2)+'px';
});