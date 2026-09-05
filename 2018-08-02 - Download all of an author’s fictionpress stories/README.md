<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Download_all_of_an_author%E2%80%99s_fictionpress_stories">Download all of an author’s fictionpress stories</a><br>
Date: 2018-08-02<br>
Tags: Examples, Fanfiction, Fictionpress, GrabAll, HTML, JavaScript, Quick Fixes, Web, Web browsers<br><br>

<p>I was surprised in my failure to find a script online to download all of an author’s stories from <a href="https://www.fictionpress.com">Fiction Press</a> or <a href="http://fanfiction.net">Fan Fiction.Net</a>, so I threw together the below.</p>
<p>If you go to an author’s page in a browser (only tested in Chrome) it should have all of their stories, and you can run the following script in the console (F12) to grab them all. Their save name format is <code><kbd>STORY_NAME_LINK_FORMAT</kbd> - <kbd>CHAPTER_NUMBER</kbd>.html</code>. It works as follows:</p>
<ol>
	<li>Gathers all of the names, chapter 1 links, and chapter counts for each story.</li>
	<li>Converts this information into a list of links it needs to download. The links are formed by using the chapter 1 link, and just replacing the chapter number.</li>
	<li>It then downloads all of the links to your current browser’s download folder.</li>
</ol><br>
<p>Do note that chrome should prompt you to answer “This site is attempting to download multiple files”. So of course, say yes. The script is also designed to detect problems, which would happen if fictionpress changes their html formatting.</p>

```js
//Gather the story information
const Stories=[];
$('.mystories .stitle').each((Index, El) =>
	Stories[Index]={Link:$(El).attr('href'), Name:$(El).text()}
);
$('.mystories .xgray').each((Index, El) =>
	Stories[Index].NumChapters=/ - Chapters: (\d+) - /.exec($(El).text())[1]
);

//Get links to all stories
const LinkStart=document.location.protocol+'//'+document.location.host;
const AllLinks=[];
$.each(Stories, (_, Story) => {
	if(typeof(Story.NumChapters)!=='string' || !/^\d+$/.test(Story.NumChapters))
		return console.log('Bad number of chapters for: '+Story.Name);
	const StoryParts=/^\/s\/(\d+)\/1\/(.*)$/.exec(Story.Link);
	if(!StoryParts)
		return console.log('Bad link format for stories: '+Story.Name);
	for(let i=1; i<=Story.NumChapters; i++)
		AllLinks.push([LinkStart+'/s/'+StoryParts[1]+'/'+i+'/'+StoryParts[2], StoryParts[2]+' - '+i+'.html']);
});

//Download all the links
$.each(AllLinks, (_, LinkInfo) =>
	$('a').attr('download', LinkInfo[1]).attr('href', LinkInfo[0])[0].click()
);

jQuery('.blurb.group .heading a[href^="/works"]').map((_, El) => jQuery(El).text()).toArray().join('\n');
```