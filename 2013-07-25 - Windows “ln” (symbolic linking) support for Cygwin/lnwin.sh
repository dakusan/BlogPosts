#!/bin/bash
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