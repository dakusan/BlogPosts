<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Symlinks_in_a_Windows_programming_environment">Symlinks in a Windows programming environment</a><br>
Date: 2010-12-04<br>
SubTitle: Windows will get it right one day<br>
Tags: Cygwin, Examples, GIT, NTFS, Symlinks, TortiseGIT, Windows<br><br>

<p>I have been having some problems regarding <a href="http://en.wikipedia.org/wiki/Symbolic_link">symlinks (symbolic links)</a> for <a href="https://www.castledragmire.com/Posts/A_weird_thing_happened_today">a project that I’ve been working on recently</a> which is requiring work in at least 5 very different operating systems (and about a dozen programming languages). Not many programs support symlinks properly that I have the need to because support for it wasn’t added for NTFS until Windows Vista, and it still has some problems.</p>
<p>It is really great that Windows Vista and Windows 7 now support native symlinks so they can be utilized by programs out of the box. For example, one such instance where I need this a lot is directory relinking in Apache. While Apache’s mod_alias can duplicate the functionality of symlinks for many needs, creating special cases for this one piece of software when distributing a code repository is just not practical, and having proper symlinks natively followed without the program knowing they aren’t the actual file/directory is really the best solution so everything works without special cases.</p>
<p>The way to create NTFS symlinks in Windows Vista+ is through the “mklink” command, which is unfortunately implemented directly in the Window’s command shell, and not a separate executable, so it is not accessible to Cygwin. Further, Cygwin has made a stance to only support reading NTFS symlinks, and not creating them, because they can only be created by administrators, and require specification as to whether the link’s target is a directory or file. Cygwin itself in Windows has had support for symlinks for a long time, but these are not compatible with any program run outside of the Cygwin environment.</p>
<p>Now, my real problem started occurring when trying to use these NTFS symlinks with <a href="http://en.wikipedia.org/wiki/Git_%28software%29">GIT</a>. While GIT natively supports symlinks, <a href="https://code.google.com/p/tortoisegit/">TortoiseGIT</a> doesn’t really support them at all, and throws errors when they are encountered. This is still a big problem that I am going to have to think about 😕. Fortunately, when working with GIT in Cygwin they still work, with caveats. As previously mentioned, only reading the NTFS symlinks in Cygwin work, so when you fetch/pull from a repository and it creates Cygwin style symlinks, Windows still does not read them properly. The following is a script I wrote to change the Cygwin style symlinks into NTFS style symlinks. It can be run from the root folder of the GIT project.</p>
</div>

```bash
#!/bin/bash
IFS=$'\n' #Spaces do not count as new delimiters

function makewinlink
{
	LINK=$1
	OPTIONS=$2
	TARGET=`find $LINK -maxdepth 0 -printf %l`
	LASTMODTIME=`find $LINK -maxdepth 0 -printf "%t"`
	LINKDIR=`find $LINK -maxdepth 0 -printf %h`
	TARGET=`echo $LINKDIR/$TARGET`
	rm -f $LINK
	cmd /c mklink $OPTIONS "$(cygpath -wa $LINK)" "$(cygpath -wa $TARGET)"
	touch -h -d $LASTMODTIME $LINK
}

#Relink all directories
FILES=`find -type l -print0 | xargs -0 -i find -L {} -type d -maxdepth 0`
for f in $FILES
do
	makewinlink $f /D
done

#Relink all files
FILES=`find -type l -print0 | xargs -0 -i find -L {} -type f -maxdepth 0`
for f in $FILES
do
	makewinlink $f
done
```

<div>
<br>
<p>Make sure when committing symlinks in a GIT repository in Windows to use Cygwin with Cygwin style symlinks instead of TortoiseGIT. Also, as previously mentioned, after running this script, TortoiseGIT will show these symlinks as modified 😕. If this is a problem, you can always reverse the process in Cygwin by changing the “cmd /c mklink $OPTIONS” line to a “ln -s” in the above script (note that “target” and “symlink’s name” need to be switched) along with a few other changes.</p><br>

<b>[EDIT ON 2011-01-03 @ 6:30am]</b> See <a href="../2011-01-03 - SymLink Fix for Combining Android Project Versions">here</a> for a better example of symlinking in Windows that uses relative paths. <b>[/EDIT]</b></div>