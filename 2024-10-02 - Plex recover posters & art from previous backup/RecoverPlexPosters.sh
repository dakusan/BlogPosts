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