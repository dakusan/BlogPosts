#!/bin/bash
LIBRARY_ID=1 #Set this from the list of libraries

#Get list of libraries
'/usr/lib/plexmediaserver/Plex Media Scanner' --list

#Force library intro dection
'/usr/lib/plexmediaserver/Plex Media Scanner' --analyze --manual --server-action intros --section $LIBRARY_ID

#Force library credit dection
'/usr/lib/plexmediaserver/Plex Media Scanner' --analyze --manual --server-action credits --section $LIBRARY_ID