<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Some_useful_Plex_commands">Some useful Plex commands</a><br>
Date: 2024-09-16<br>
Tags: Plex<br><br>
</div>

<p>The below commands are tailored for Linux, as file paths are different for Windows. They need to be ran with the <code>sudo -u plex</code> prefix for proper user access permissions.</p>
<ins>Get list of libraries</ins>

```bash
'/usr/lib/plexmediaserver/Plex Media Scanner' --list
```

<ins>Force library intro dection</ins>
<pre><code>'/usr/lib/plexmediaserver/Plex Media Scanner' --analyze --manual --server-action intros --section <kbd>LIBRARY_ID</kbd></code></pre>

<ins>Force library credit dection</ins>
<pre><code>'/usr/lib/plexmediaserver/Plex Media Scanner' --analyze --manual --server-action credits --section <kbd>LIBRARY_ID</kbd></code></pre><br>

<ins>Export media list</ins> (this is a bash command)

```sql
sqlite3 '/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Plug-in Support/Databases/com.plexapp.plugins.library.db' <<'SQL'
SELECT MDI.id, file, title, MDI."index", hints, GROUP_CONCAT(T.text) AS taggings
FROM media_parts AS MP
INNER JOIN media_items AS MI ON MI.id=MP.media_item_id
INNER JOIN metadata_items AS MDI ON MDI.id=MI.metadata_item_id
LEFT JOIN taggings as T ON T.metadata_item_id=MDI.id AND T.text IN ("credits", "intro")
GROUP BY MDI.id;
SQL
```

<div>This outputs the following (column separated by pipe “|”) for each individual media item:
<ol>
	<li>Metadata item ID</li>
	<li>File path</li>
	<li>Title</li>
	<li>Episode Number</li>
	<li>Hints (e.x. <code>=&episode=1&episodic=1&season=1&show=Ninja%20Turtles&show_name=Ninja%20Turtles&year=2002</code>)</li>
	<li>“credits” and “intros” tags (if the credits and intros timestamps have been calculated)</li>
</ol></div></div>