#!/bin/bash

sqlite3 '/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Plug-in Support/Databases/com.plexapp.plugins.library.db' <<'SQL'
SELECT MDI.id, file, title, MDI."index", hints, GROUP_CONCAT(T.text) AS taggings
FROM media_parts AS MP
INNER JOIN media_items AS MI ON MI.id=MP.media_item_id
INNER JOIN metadata_items AS MDI ON MDI.id=MI.metadata_item_id
LEFT JOIN taggings as T ON T.metadata_item_id=MDI.id AND T.text IN ("credits", "intro")
GROUP BY MDI.id;
SQL