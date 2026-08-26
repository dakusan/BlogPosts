#!/bin/bash
service plexmediaserver stop
/usr/lib/plexmediaserver/Plex\ SQLite '/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Plug-in Support/Databases/com.plexapp.plugins.library.db'
