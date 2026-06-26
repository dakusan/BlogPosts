<div>
Original Post: <a href="https://www.castledragmire.com/Posts/JavaScript_problems_when_crossing_windows_Part_2">JavaScript problems when crossing windows Part 2</a><br>
Date: 2010-01-21<br>
SubTitle: IE being a pain in the butt like usual<br>
Tags: Data Passing, Examples, JavaScript, Standards, Web browsers<br><br>

<p>To continue the subject in my <a href="../2010-01-18 - JavaScript problems when crossing tabs or windows">last post</a>, these next cross-window bugs also derive from objects not being recognized properly when being passed between windows in JavaScript.</p>
<p>I needed the ability to dynamically run functions in the secondary window form the primary window where the parameters are taken from an array. Since a “function” from a secondary window is not seen as a function object from the primary window in IE, the <a href="http://www.webreference.com/js/column26/apply.html">apply member</a> was not working.</p>
<p>I have included a fix for this below in the “RunFunctionInRemoteWindow” function, which is just a wrapper function in the second window that calls the apply function. This function manually copies the array through a for loop, instead of using slice, because in IE7 (but not IE8), the passed arrays were not seen as valid JSObjects, so the <a href="http://www.w3schools.com/jsref/jsref_slice_array.asp">slice method</a> (which is a standard method used for copying arrays by value) was not working.</p><br>

<b><a href="LocalWindow.html">LocalWindow.html</a> </b> [run this one]
</div>

```html
<html><body>
<input type=button onclick="RunTest();" value='Click me when the second window has opened to run the test'>
<script type="text/javascript">

//Spawn the second window
var NewWindow=window.open('RemoteWindow.html');

//Run the test
function RunTest()
{
	LocalAlertString('This is an alert generated from the local window');
	NewWindow.RemoteAlertString('This is an alert generated from the remote window');
	alert('The local window alert function is of type function: '+(LocalAlertString instanceof Function));
	alert('The remote window alert function is of type function: '+(NewWindow.RemoteAlertString instanceof Function));
	LocalAlertString.apply(window, ['This is an alert generated from the local window through the APPLY member']);

	try {
	NewWindow.RemoteAlertString.apply(NewWindow.window, ['This is an alert generated from the remote window through the APPLY member. This will not work in IE because the remote window\'s function is not actually a function.']);
	} catch(e) { alert('The REMOTE APPLY failed: '+e.message); }

	NewWindow.RunFunctionInRemoteWindow('RemoteAlertString', ['This is an alert generated from the remote window through the FIXED APPLY function.']);
}

//Generate an alert in the local window
function LocalAlertString(TheString)
{
	alert('Local String: '+TheString);
}

</script></body></html>
```

<div>
<br>

<b><a href="RemoteWindow.html">RemoteWindow.html</a></b> [do not run this one, it is opened as a popup from LocalWindow.html]
</div>

```html
<html><body><script type="text/javascript">
//Generate an alert in the remote window
function RemoteAlertString(TheString)
{
	alert('Remote String: '+TheString);
}

//Call functions in this window remotely through the "apply" member
function RunFunctionInRemoteWindow(FunctionName, Parameters)
{
	//Manually copy the passed Parameters since "Parameters" may not be a valid JSObject anymore (this could be detected and array.slice used if it is still valid)
	var ParametersCopy=[];
	for(var i=0;i<Parameters.length;i++)
		ParametersCopy[i]=Parameters[i];
	
	window[FunctionName].apply(window, ParametersCopy);
}
</script></body></html>
```

<div>
</div>