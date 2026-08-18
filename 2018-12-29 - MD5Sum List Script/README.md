<div>
Original Post: <a href="https://www.castledragmire.com/Posts/MD5Sum_List_Script">MD5Sum List Script</a><br>
Date: 2018-12-29<br>
Tags: Bash, Command Line Utility, Examples, Hash, MD5<br><br>
</div>

```bash
#This script takes a newline delimited file list from STDIN for md5 hashing
#This script requires the `md5sum`, `pv`, `paste`, `bc`, and 'numfmt' commands

#The output of the md5s are stored in the file specified by the first parameter
#The format for each md5 hash to the output file is "$FileName\t$Hash\n"

#File sizes are always output in megabytes with 3 decimal places
#While calculating the hashes the script keeps the user informed of the progress of both the current file and all the files as follows:
#1) Before file starts: "Hashing: $FileName ($FileSize MiB)\n"
#2) During transfer: The progress of the hash of the current file ran through `pv`
#3) During transfer: The progress of the hashing of all the files, ran through `pv`
#4) After transfer: "Finished $TotalProgressPercent% ($ProcessedBytes/$TotalBytes MiB)\n\n"

#Get $Outfile from the first argument and the $FileList from STDIN (newline delimited)
OutFile="$1";
FileList=`cat /dev/stdin`

#Format a byte count in MegaBytes with comma grouping and 3 decimal places
MbFmtNoExt ()
{
	echo "scale=3; $1/1024/1024" | bc | echo -n `xargs numfmt --grouping`
}

#Add " MiB" to the end of MbFmtNoExt
MbFmt ()
{
	echo `MbFmtNoExt $1`" MiB"
}

#Calculate and output the total size of the file list
echo -n "Calculating total size: "
TotalSize=`echo "$FileList" | xargs -d"\n" stat --printf="%s\n" | paste -s -d+ | bc`
MbFmt $TotalSize
echo #Add an extra newline

#Create a fifo to keep track of the total complete
TotalDoneFifo=$(mktemp)
TotalDoneBG=0
rm "$TotalDoneFifo"
mkfifo "$TotalDoneFifo"
cat > "$TotalDoneFifo" & #Do not close read of fifo
Cleanup() {
	rm "$TotalDoneFifo"
	kill $TotalDoneBG
	exit 0
}
trap Cleanup SIGTERM SIGINT

#Start the TOTAL line
tail -f "$TotalDoneFifo" | pv -s $TotalSize -F  "%b %t %p %e" > /dev/null &
TotalDoneBG=$!

#Run over the list (newline delimited)
CalculatedBytes=0
IFS=$'\n'
for FileName in `echo "$FileList"`
do
	#Output the file size and name to STDOUT
	FileSize=`stat --printf="%s" "$FileName"`
	echo "Hashing: $FileName ("`MbFmt $FileSize`")"

	#Output the filename to $OutFile
	echo -n $FileName$'\t' >> $OutFile

	#Run the md5 calculation with `pv` progress
	#Output the hash to $OutFile after the FileName and a tab
	cat "$FileName" | pv -s $FileSize -c | tee -a "$TotalDoneFifo" | md5sum | awk '{print $1}' >> $OutFile

	#Output the current progress for the entire file list
	#Format: "Finished $TotalProgressPercent% ($ProcessedBytes/$TotalBytes MiB)\n\n"
	CalculatedBytes=$(($CalculatedBytes+$FileSize))
	echo -n "Finished "
	printf "%.3f" `echo "scale=4; $CalculatedBytes*100/$TotalSize" | bc`
	echo "% ("`MbFmtNoExt $CalculatedBytes`"/"`MbFmt $TotalSize`$')\n'
done

Cleanup
```