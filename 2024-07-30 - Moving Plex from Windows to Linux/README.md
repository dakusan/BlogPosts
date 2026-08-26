<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Moving_Plex_from_Windows_to_Linux">Moving Plex from Windows to Linux</a><br>
Date: 2024-07-30<br>
Tags: Linux, Linux Migration, Plex, SQLite<br><br>

<p>Yesterday I <a href="https://support.plex.tv/articles/201370363-move-an-install-to-another-system/">moved my Plex installation</a> from a Windows machine to a Linux machine. The <a href="https://support.plex.tv/articles/202915258-where-is-the-plex-media-server-data-directory-located/">primary data folders</a> that needed to be copied over were <code>Media</code>, <code>Metadata</code>, <code>Plug-ins</code>, and <code>Plug-in Support</code>. It doesn't hurt to copy over some of the folders in <code>Cache</code> too. It's possible there may be some more data that needs to be moved over, but I don't have documented what.</p>
<p>After moving all the data, I updated the paths in the database for the new machine. Doing this allowed me to keep everything as it was and no new refreshes/scans needed to be done.</p>
<p>The <a href="https://support.plex.tv/articles/repair-a-corrupted-database/#toc-0">location of the Plex modified SQLite</a> for me was at <code>/usr/lib/plexmediaserver/Plex SQLite</code>. So the following is the bash commands to stop the plex server and open SQLite editor on Linux Mint 21.3.</p>

```bash
service plexmediaserver stop
/usr/lib/plexmediaserver/Plex\ SQLite '/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Plug-in Support/Databases/com.plexapp.plugins.library.db'
```

<p>And the following is the SQL I used to replace a <code>D:</code> drive with a path of <code>/home/plex/drives/d/</code>. You can replace these strings in the below code with your custom drives and paths.</p>

```sql
#Replace backslashes with forward slash in paths - This only needs to be run once
UPDATE media_parts       SET file     =REPLACE(file     , '\\', '/');
UPDATE section_locations SET root_path=REPLACE(root_path, '\\', '/');
UPDATE media_streams     SET url      =REPLACE(url      , '\\', '/') WHERE url LIKE 'file://%';

#Replace root paths - This will need to be run for every path replacement you want to make. This example remaps D:\ to /home/plex/drives/d/
UPDATE media_parts       SET file     =REPLACE(file     , 'D:/',         '/home/plex/drives/d/');
UPDATE section_locations SET root_path=REPLACE(root_path, 'D:/',         '/home/plex/drives/d/');
UPDATE media_streams     SET url      =REPLACE(url      , 'file://D:/',  'file:///home/plex/drives/d/') WHERE url LIKE 'file://%';
UPDATE media_streams     SET url      =REPLACE(url      , 'file:///D:/', 'file:///home/plex/drives/d/') WHERE url LIKE 'file://%';
UPDATE metadata_items    SET guid     =REPLACE(guid     , 'file:///D:/', 'file:///home/plex/drives/d/');
```