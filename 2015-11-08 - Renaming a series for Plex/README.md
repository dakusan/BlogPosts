<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Renaming_a_series_for_Plex">Renaming a series for Plex</a><br>
Date: 2015-11-08<br>
Tags: Examples, Plex, Quick Fixes, Symlinks<br><br>

<p>I was recently trying to upload a TV series into Plex and was having a bit of a problem with the file naming. While I will leave the show nameless, let’s just say it has a magic dog.</p>
<p>Each of the files (generally) contained 2 episodes and were named S##-E##-E## (Season #, First Episode #, Second Episode #). Plex really didn’t like this, as for multi-episode files, it only supports the naming convention of first episode number THROUGH a second episode number. As an example S02-E05-E09 is considered episodes 5 through 9 of season 2. So I wrote a quick script to fix up the names of the files to consider each file only 1 episode (the first one), and then create a second symlinked file, pointing to the first episode, but named for the second episode.</p>
<p>So, for the above example, we would get 2 files with the exact same original filenames, except with the primary file having “S02E05,E09” in place of the episode number information, and the linked file having “S02E09-Link” in its place.</p><br>
The following is the bash code for renaming/fixing a single episode file. It needs to be saved into a script file. This requires perl for regular expression renaming.
</div>

```bash
#Get the file path info and updated file name
FilePath=`echo "$1" | perl -pe 's/\/[^\/]*$//g'`
FileName=`echo "$1" | perl -pe 's/^.*\///g'`
UpdatedFileName=`echo "$FileName" | perl -pe 's/\b(S\d\d)-(E\d\d)-(E\d\d)\b/$1$2,$3/g'`

#If the file is not in the proper format, exit prematurely
if [ "$UpdatedFileName" == "$FileName" ]; then
    echo "Proper format not found: $FilePath/$FileName"
    exit 1
fi

#Rename the file
cd "$FilePath"
mv "$FileName" "$UpdatedFileName"

#Create a link to the file with the second episode name
NewLinkName=`echo "$FileName" | perl -pe 's/\b(S\d\d)-(E\d\d)-(E\d\d)\b/$1$3-Link/g'`
ln -s "$UpdatedFileName" "$NewLinkName"
```

If you save that to a file named “<i>RenameShow.sh</i>”, you would use this like `find /PATH/ -type f -print0 | xargs -0n 1 ./RenameShowl.sh`. For windows, make sure you use <a href="../2010-12-04 - Symlinks in a Windows programming environment">windows symlinks</a> with /H (to make them hard file links, as soft/symbolic link files really just don’t work in Windows).</p></div>