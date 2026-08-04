#!/bin/bash
LANG=;
IFS=$'\n'
for i in `find -type f | grep -P '\xEF[\x80-\x8F]'`; do
	FROM="$i";
	TO=$(echo "$i" | perl -pi -e 's/\xEF([\x80-\x8F])/pack("C", ord($1)+(0xC0-0x80))/e');
	echo Renaming "'$FROM'" to "'$TO'"
	mv "$FROM" "$TO"
done