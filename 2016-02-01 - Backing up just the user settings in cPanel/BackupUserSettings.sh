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