<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Auto_Locking_Windows_on_Login">Auto Locking Windows on Login</a><br>
Date: 2018-12-07<br>
Tags: Automated, Boot, Boot Loader, Login, Quick Fixes, Security, Windows<br><br>

<p>On my primary computer (whose harddrive is encrypted) I always have Windows auto logging in to help with the bootup time. However, my bootup time can be rather slow; and if I needed to have my computer booted but locked, I had to wait for the login to complete so I could lock the computer. This has been becoming a nuisance lately when I need to get out of my house quickly in the morning.</p>
<p>For the solution I created a windows boot entry that auto locks the computer after logging the user in. This also requires a batch file, to run for the user on startup, to detect when this boot entry was selected. Here are the steps to create this setup:</p><br>

<ol>
	<li>
		<b>Create the new boot entry</b>:
		In the windows command line, run: <code>bcdedit /copy {current} /d "Lock on Startup"</code><br>
		This creates a new boot option, duplicated from your currently selected boot option, in the boot menu labeled “<b>Lock on Startup</b>”.
	</li>
	<li>
		<b>(Optional) Change the bootup timeout</b>:
		In the windows command line, run: <code>bcdedit /timeout <kbd>5</kbd></code><br>
		Where 5 is a 5 second timeout.
	</li>
	<li>
		<b>Create a batch file to run on login</b>:
		In your start menu’s startup folder, add a batch file. You can name it anything as long as the extension is “.bat”.<br>
		Add the following to the file: <code>bcdedit /enum {current} | findstr /r /c:"description  *Lock on Startup" &amp;&amp; rundll32.exe user32.dll,LockWorkStation</code><br>
		Note that there are 2 spaces in the description search string to replicate the regular expression's 1 or more quantifier “+”, since <b>findstr</b> only supports the 0 or more quantifier “*”.
	</li>
</ol></div>