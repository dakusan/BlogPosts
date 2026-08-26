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