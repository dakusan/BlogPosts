<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Backing_up_just_the_user_settings_in_cPanel">Backing up just the user settings in cPanel</a><br>
Date: 2016-02-01<br>
Tags: Backup, cPanel, Examples, MySQL, r1soft, Servers<br><br>
</div>

<p>One of the companies I work for recently moved one of our cPanel servers to a new collocation, still running cPanel. We decided to use a new backup solution called <a href="http://www.r1soft.com/">r1soft</a>, which so far has been working spectacularly. I’d love to use it for my personal computers, except the licenses, which are geared towards enterprise business, are way too costly.</p>
<p>However, since r1soft only backs up files (on the incrementally block level, yay) you can’t use it to restore a cPanel account. It can only restore things like the user’s home directory and SQL databases. Because of this, when we had need to restore an entire account today, and found out there is no easy/quick way to do it, we were up a creek. The obvious future solution for this would be to use cPanel’s backup (or legacy backup) systems, but unfortunately, you can’t easily set them to not backup the user’s databases and home directory, which can be very large, and are already taken care of by r1soft. I ended up adding the following script, ran nightly via cron, to back up user account settings.</p>
<p>It saves all the user settings under the backup path in their own directory, uncompressed, and named cpmove-<kbd>USERNAME</kbd>. It is best to do it this way so r1soft’s incremental backups don’t have much extra work if anything changes. Make sure to change line 3 in the following script to the path where you want backups to occur.</p>

```bash
#!/bin/bash
#Create and move to backup directory
BACKUPDIR=/backup/userbackup
mkdir -p $BACKUPDIR #Make sure the directory exists
cd $BACKUPDIR

#Remove old backups
rm -rf cpmove-*

#Loop over accounts
for USER in `/usr/sbin/whmapi1 listaccts | grep -oP '(?<=user: )\w+$' | sort -u`; do
  #Backup the account
  /scripts/pkgacct --nocompress --skipbwdata --skiphomedir --skiplogs --skipmysql --skipmailman $USER ./

  #Extract from and remove the tar container file
  tar -xvf cpmove-$USER.tar
  rm -f cpmove-$USER.tar

  #Save MySQL user settings
  mysqldump --compact -fnt -w "User LIKE '$USER""_%'" mysql user db tables_priv columns_priv procs_priv proxies_priv \
  | perl -pe "s~('|NULL)\),\('~\1),\n('~ig" \
  > cpmove-$USER/mysql-users.sql
done;
```

<p>This script skips a few backup items that need to be noted. Mailman, logs, homedir, and bandwidth data should all be easy 1:1 copy over restores from r1soft. I excluded them because those can take up a lot of room, which we want r1soft to handle. The same goes for MySQL, except that your MySQL users are not backed up to your account, which is why I added the final section.</p>
<p>Do note, for the final section, the line starting with <code>| perl</code> is optional. It is there to separate the insert rows into their own lines. A very minor warning though; it would also pick up cases where the last field in MySQL’s user table ends in <code>NULL),&#8203;(</code>. This would only happen if someone is trying to be malicious and knew about this script, and even then, it couldn’t harm anything.</p>
<p>Bonus note: To restore a MySQL database which does not use a shared-file (like InnoDB does by default), you could actually stop the MySQL server, copy over the binary database files, and start the server back up.</p></div>