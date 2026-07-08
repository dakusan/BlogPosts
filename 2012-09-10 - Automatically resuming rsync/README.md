<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Automatically_resuming_rsync">Automatically resuming rsync</a><br>
Date: 2012-09-10<br>
SubTitle: The old network file copy problem<br>
Tags: Command Line Utility, Examples, FileSync, SSH<br><br>

<p>Rsync is a spectacular bash utility for doing file copying and/or syncing operations. It has a multitude of switches to help optimize and handle any requirements for file copy operation over a local computer or network. However, sometimes networks are less than stable and stalls can happen during an rsync (or scp). This is quite the nuisance when doing very large (i.e. GB) transfers. To solve this, the following script can be used to auto resume a stalled rsync.</p><br>

</div>

```bash
export Result=1;
while [ $Result -ne 0 ]; do
  echo "STARTING ($Result) @" `date`;
  rsync -Pza --timeout=10 "$@";
  Result=$?;
  sleep 1;
done
```

<div>

<br><ul>
	<li>The -P switch is highly suggested as it activates:
	<ul>
		<li>--partial: This keeps a file even if it doesn’t finish transferring so it can be resumed when rsync is restarted. This is especially important if you have very large files.</li>
		<li>--progress: This shows you the progress of the current file copy operation.</li>
		</ul></li>
	<li>The -z switch turns on gzip compression during the file transfer, so may only be useful depending on the circumstances.</li>
	<li>The -a switch stands for “archive” and is generally a good idea to use. It includes the switches:
	<ul>
		<li>-r: Recurse into folders</li>
		<li>-t: Preserves file modification time stamp. This is highly recommended for this, and incremental backups, as rsync, by default, skips files whose file sizes and modification times match.</li>
		<li>-l and -D: Preserve file type (i.e. symlinks)</li>
		<li>-p: Preserve file (chmod) permissions</li>
		<li>-g and -o: Preserve file owners.</li>
	</ul></li>
	<li>The --timeout is the crux of the script, in that if an I/O timeout of 10 seconds occurs, the rsync exits prematurely so it can be restarted.</li>
</ul><br>

<p>For more useful switches and information, see the <a href="http://linux.die.net/man/1/rsync">rsync man page</a>.</p><br><br>

</div>

Script with comments:
```bash
export Result=1; #This will hold the result of the rsync. Set to 1 so the first loop check will fail.
while [ $Result -ne 0 ]; do #Loop until rsync result is successful
  echo "STARTING ($Result) @" `date`; #Inform the user of the time an rsync is starting and the last rsync failure code
  rsync -Pza --timeout=10 "$@"; #See rest of post for switch information. Pass script parameters through to rsync
  Result=$?; #Store the result of the rsync
  sleep 1; #This is an optional 1 second timeout between attempts
done
```