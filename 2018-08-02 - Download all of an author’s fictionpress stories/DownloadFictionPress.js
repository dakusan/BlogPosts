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