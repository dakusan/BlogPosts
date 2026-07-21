<div>
Original Post: <a href="https://www.castledragmire.com/Posts/JavaScript_problems_when_crossing_tabs_or_windows">JavaScript problems when crossing tabs or windows</a><br>
Date: 2010-01-18<br>
Tags: Data Passing, Examples, JavaScript, Web browsers<br><br>

<p>I was doing some research around April of 2009 on JavaScript interaction between web browser windows. I was doing this because web browsers are starting to split off each tab/window into separate processes/threads (Firefox is lagging in this), which can lead to some useful new implementations in the browser world, including multithreading. I wanted to explore the interaction between these windows to make sure there were no caveats that might creep up if I decided to take advantage of this.</p>
<p>The first one I found was that each browser window has its own instance of all of the base object classes, so prototypes do not carry over, and instanceof will not work as expected.</p>
<p>For example, if in <b>WindowOne</b>, you add a prototype to the Array class called <b>IsArray</b>, it is only accessible by arrays created in <b>WindowOne</b>. If you pass an array created in <b>WindowOne</b> into a second window, the prototype is still available on that one array (IIRC this was not true of some of the browsers at the time, but I tested again today, and it worked for IE8, Firefox3, and Google Chrome). Also, since the base object class in <b>Window1</b> and other windows are not the same, an object created in <b>Window1</b> and passed to another window will return false in a <b>instanceof Object</b> operation in that other window.</p><br>
Here is some example code to help show what I’m talking about.<br><br>

<b><a href="LocalWindow.html">LocalWindow.html</a> </b> [run this one]
</div>

```html
<html><body>
<input type=button onclick="RunTest();" value='Click me when the second window has opened to run the test'>
<script type="text/javascript">
Array.prototype.IsArray=true;
var NewWindow=window.open('RemoteWindow.html'); //Spawn the second window
function RunTest() { NewWindow.RunTest({}, [], new ExampleObject()); }; //Send the test data to remote window
function ExampleObject() { } //An example class
</script></body></html>
```

<div>
<br>

<b><a href="RemoteWindow.html">RemoteWindow.html</a></b> [do not run this one, it is opened as a popup from LocalWindow.html]
</div>

```html
<html><body><script type="text/javascript">
function RunTest(AnObject, AnArray, AnExampleObject)
{
	var MyTests=[
		'AnObject instanceof Object',
		'AnObject.IsArray',                               //Object.prototype does not have this (Array.prototype does)
		'AnArray instanceof Object',
		'AnArray instanceof Array',
		'AnArray.IsArray',                                //This was added to the Array.prototype in the parent window
		'AnArray instanceof opener.Array',                //This example does not work in IE7 because opener.x cannot be properly accessed
		'AnExampleObject instanceof opener.ExampleObject',//This example does not work in IE7 because opener.x cannot be properly accessed
		'AnExampleObject instanceof ExampleObject'        //This test should error because "ExampleObject" does not exist in this window
	];
	
	for(var i=0;i<MyTests.length;i++) //This runs each test like the following: alert("TEST: "+(TEST));
		try {
			eval('alert("'+MyTests[i]+': "+('+MyTests[i]+'));');
		} catch(e) {
			alert('Error on test "'+MyTests[i]+'": '+(e.hasOwnProperty('message') ? e.message : e.toString()));
		}
}
</script></body></html>
```

<a href="../2010-01-21 - JavaScript problems when crossing windows Part 2">Continued here</a>

<div>
</div>