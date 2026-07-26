#!/bin/bash
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