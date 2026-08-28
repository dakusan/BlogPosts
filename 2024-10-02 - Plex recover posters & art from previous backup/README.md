<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Plex_recover_posters_%26_art_from_previous_backup">Plex recover posters &amp; art from previous backup</a><br>
Date: 2024-10-02<br>
Tags: Examples, Plex, Quick Fixes, Recover, SQLite<br><br>
</div>

<p>My Plex server decided to go through my database and screw with all the posters and art, so I made the below script which can be used to recover them. You should only use this if doing a full revert to backup is not a viable option. Do also note that if you have refreshed the metadata on any tv series since the backup, it could cause issues with the episode posters.</p><br>

```bash
#Run the following commands as root
sudo su

#User variables
SERVICE_NAME='plexmediaserver'
BACKUP_LOCATION='/ROOT_BACKUP_LOCATION/'
PLEX_USER='plex'

#Derived variables (Slashes at end of directory paths are required)
SQLITE3_PATH="/usr/lib/$SERVICE_NAME/Plex SQLite"
DB_PATH="/var/lib/$SERVICE_NAME/Library/Application Support/Plex Media Server/Plug-in Support/Databases/com.plexapp.plugins.library.db"
IMG_PATH="/var/lib/$SERVICE_NAME/Library/Application Support/Plex Media Server/Metadata/"
CACHE_PATH="/var/lib/$SERVICE_NAME/Library/Application Support/Plex Media Server/Cache/PhotoTranscoder/"

#Stop the plex service until we have updated everything
echo "Stopping Plex"
systemctl stop "$SERVICE_NAME"

#Set the (user_thumb_url, user_art_url, and updated_at) from the old database in the new database.
#The updated_at is required to match the cached images
echo "Updating database"
sudo -u "$PLEX_USER" "$SQLITE3_PATH" "$DB_PATH" -cmd "ATTACH DATABASE '$BACKUP_LOCATION$DB_PATH' as OLDDB" "
UPDATE
  metadata_items AS NMDI
SET
  user_thumb_url=OMDI.user_thumb_url,
  user_art_url=OMDI.user_art_url,
  updated_at=OMDI.updated_at
FROM
	OLDDB.metadata_items AS OMDI
WHERE
	NMDI.id=OMDI.id
"

#Copy the posters and art from the old directory to the new directory
echo "Copying posters"
cd "$BACKUP_LOCATION$IMG_PATH"
find -type d \( -name "art" -o -name "posters" \) | xargs -n100 bash -c 'rsync -a --relative "$@" "'"$IMG_PATH"'"' _

#Copy the cached images over
echo "Copying cache"
rsync -a "$BACKUP_LOCATION$CACHE_PATH" "$CACHE_PATH"

#Restart plex
echo "Starting Plex"
systemctl start "$SERVICE_NAME"
```

<b><ins>[Edit on 2024-10-17]</ins></b><br>
Adding a solution to my problem where movie posters kept disappearing in between plex server reboots.<br>
For some reason it was storing the movie poster metadata in `Metadata/*/*/*/Upload` and it expected them in `Metadata/*/*/*/Contents/_combined/`. I ran the following command in `/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Metadata` to fix the missing posters:
```bash
IFS=$'\n'; for i in `ls */*/* -d`; do cd "$i"; if [ -d "Uploads" ]; then echo "$i"; rsync -a Uploads/ Contents/_combined/; fi; cd -; done
```