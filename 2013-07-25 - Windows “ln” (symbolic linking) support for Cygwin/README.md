<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Windows_%E2%80%9Cln%E2%80%9D_%28symbolic_linking%29_support_for_Cygwin">Windows “ln” (symbolic linking) support for Cygwin</a><br>
Date: 2013-07-25<br>
SubTitle: Cygwin requires lots of tweaks to really be usable<br>
Tags: Command Line Utility, Cygwin, Quick Fixes, Symlinks, Windows<br><br>

The “ln -s” command in Cygwin creates a fake symbolic link only supported in Cygwin. I whipped up the following script to create a true windows symbolic link that is supported in both Windows and Cygwin (it shows up as a symlink in Cygwin).
</div>

```bash
TARGET=`echo $1 | perl -pe 's/\\//\\\\/g'`; #Determine the target from the first parameter (where we are linking to). Change forward slashes to back slashes
LINK=`echo $2 | perl -pe 's/\\//\\\\/g'` #Determine the link name from the second parameter (where the symlink is made). Change forward slashes to back slashes
cmd /c mklink $3 $4 $5 $6 "$LINK" "$TARGET" #Perform the windows mklink command with optional extra parameters
```

<div>

Note that for the link name, you have to include the full filename, not just specify a directory.<br>
See <a href="http://technet.microsoft.com/en-us/library/cc753194%28v=ws.10%29.aspx">here for information on mklink and its switches.</a> Specifically:
<ul>
	<li>/d : directory symlink (more hard)</li>
	<li>/j : directory junction (more soft)</li>
</ul>
This handles the problem a little more directly than my <a href="../2010-12-04 - Symlinks in a Windows programming environment">other post on the topic ("Symlinks in a Windows programming environment")</a>.

<br><br><b><u>[Edit on 2016-01-12 @ 12:34am]</u></b><br>
And once again, I have a new version of the code. This version has the following advantages:
<ul>
	<li>No longer limited to just 4 extra parameters (dynamic instead of static)</li>
	<li>Can now just specify your directory as the link location, and the filename will be automatically filled in</li>
	<li>Handles spaces better</li>
</ul>
Do note, if you want to link directly to a hard drive letter, you must use "c:/" instead of "/cygdrive/c/"<br><br>

</div>

```bash
#Get the target and link
TARGET="$1"
shift
LINK="$1"
shift

#If the link is already a directory, append the filename to the end of it
if [ -d "$LINK" ]; then
    #Get the file/directory name without the rest of the path
    ELEMENT_NAME=`echo "$TARGET" | perl -pe 's/^.*?\/([^\/]+)\/?$/$1/'`

    #Append the file name to the target, making sure there is only 1 separating "/"
    LINK=`echo "$LINK" | perl -pe 's/^(.*?)\/?$/$1/'`
    LINK="$LINK"/"$ELEMENT_NAME"
fi

#Replace forward slashes with back slashes
TARGET=`echo $TARGET | perl -pe 's/\\//\\\\/g'`
LINK=`echo $LINK | perl -pe 's/\\//\\\\/g'`

#Perform the windows mklink command with optional extra parameters
cmd /c mklink "$@" "$LINK" "$TARGET"
```
