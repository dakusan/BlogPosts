<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Weird_filename_encoding_issues_on_windows">Weird filename encoding issues on windows</a><br>
Date: 2016-02-08<br>
Tags: Convert, Cygwin, Encoding, NTFS, Quick Fixes, UTF16, UTF8, Windows<br><br>
</div>

<p>So somehow all of the file names in my Rammstein music directory, and some in my Daft Punk, had characters with diacritics replaced with an invalid character. I pasted one of such filenames into a hex editor to evaluate what the problem was. First, I should note that Windows encodes its filenames (and pretty much everything) in UTF16. Everything else in the world (mostly) has settled on UTF8, which is a much better encoding for many reasons. So during some file copy/conversion at some point in the directories’ lifetime, the file names had done a freakish (utf16*)(utf16->utf8) rename, or something to that extent. I had noticed that all I needed to do was to replace the first 2 bytes of the diacritic character with a different byte. Namely “EF 8x” to “Cx”, and the rest of the bytes for the character were fine. So if anyone ever needs it, here is the bash script.</p>

```bash
LANG=;
IFS=$'\n'
for i in `find -type f | grep -P '\xEF[\x80-\x8F]'`; do
	FROM="$i";
	TO=$(echo "$i" | perl -pi -e 's/\xEF([\x80-\x8F])/pack("C", ord($1)+(0xC0-0x80))/e');
	echo Renaming "'$FROM'" to "'$TO'"
	mv "$FROM" "$TO"
done
```

<p>I may need to expand the range beyond the x80-x8F range, but am unsure at this point. I only confirmed the range x82-x83.</p>