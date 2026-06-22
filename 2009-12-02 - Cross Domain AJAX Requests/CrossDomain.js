function GetCookiesFromURL(Domains) //Return all the cookies for Domains specified in the Domains array
{
	var cookieManager = Components.classes["@mozilla.org/cookiemanager;1"].getService(Components.interfaces.nsICookieManager); //Requires privileges, which is granted in CrossDomainGetURL
	var iter=cookieManager.enumerator, CookieList=[], cookie; //The object used to find all cookies, the final list of cookies, and a temporary object
	while(iter.hasMoreElements()) //Loop through all cookies
		if(((cookie=iter.getNext()) instanceof Components.interfaces.nsICookie) && Domains.indexOf(cookie.host)!=-1) //If a cookie whose host matches one of our domains
			CookieList.push(cookie.name+'='+cookie.value); //Add it to our final list
	return CookieList.join("; "); //Return the cookie list for the specified domains
}

function CrossDomainGetURL(URL, Data, CookieDomains, ExtraOptions) //See AJAX Library GetURL function. CookieDomains is an array specifying what domains cookies are pulled from for the AJAX call. 
{
	//Access universal privileges in Firefox (Required to get cookies for other domains, and to use AJAX with other domains). This functionality is lost as soon as this function loses scope.
	try { netscape.security.PrivilegeManager.enablePrivilege("UniversalXPConnect"); }
	catch(e) { return alert('Cannot access browser privileges'); }

	if(CookieDomains instanceof Array) //If an array of domains is passed to get cookies from...
	{	
		ExtraOptions=((ExtraOptions instanceof Object) ? ExtraOptions : {}); //Make sure extra options is an object
		ExtraOptions.AdditionalHeaders=((ExtraOptions.AdditionalHeaders instanceof Object) ? ExtraOptions.AdditionalHeaders : {}); //Make sure extra options has an additional headers object
		ExtraOptions.AdditionalHeaders.Cookie=GetCookiesFromURL(CookieDomains); //Get cookies for the domains
	}
	
	return GetURL(URL, Data, ExtraOptions); //Do the AJAX Call
}