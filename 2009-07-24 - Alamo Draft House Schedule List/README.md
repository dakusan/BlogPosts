<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Alamo_Draft_House_Schedule_List">Alamo Draft House Schedule List</a><br>
Date: 2009-07-24<br>
SubTitle: Simple information reorganization example<br>
Tags: Alamo Draft House, Entertainment, Examples, PHP<br><br>

<p>After discovering the <a href="http://www.originalalamo.com/">Alamo Draft House’s</a> <a href="https://www.castledragmire.com/Posts/Star_Trek_at_The_Alamo">coolness</a> a few months ago, I’ve been trying to watch what they’re playing to make sure I catch anything I might want to see on the big screen. Unfortunately, it is not easy to get a good quick idea of all the movies playing from their calendar because it shows movies per day with showtimes, making the list repetitive and crowded with extra information.</p>
<p>I decided to throw together a real quick PHP script that would parse their data so I could organize it however I wanted. The final result can be viewed <a href="https://www.castledragmire.com/Misc/Alamo_Draft_House_Listings.php">here</a>. The code is as follows:</p>

</div>

```php
//The list of calendar pages in format TheaterName=>URL
$PagesToGrab=Array(
	'Ritz'=>'http://www.originalalamo.com/Calendar.aspx?l=2',
	'Village'=>'http://www.originalalamo.com/Calendar.aspx?l=3',
	'South Lamar'=>'http://www.originalalamo.com/Calendar.aspx?l=4'
);

foreach($PagesToGrab as $Name => $URL) //Grab the movies for each theater
{
	print "<b>$Name</b><br>"; //Output the theater name
	$TheHTML=file_get_contents($URL); //Grab the HTML
	$ShowList=Array(); //This will contain the final list of shows and what days they are on
	
	preg_match_all('/<td class="day">.*?<\/td>/', $TheHTML, $DayMatches); //Extract all table cells containing a day's info
	foreach($DayMatches[0] as $DayInfo) //Loop over each day's info
	{
		//Determine the day of month
		preg_match('/<a class=\"daynumber\" title=".*?, (.*?),/', $DayInfo, $DayOfMonth);
		$DayOfMonth=$DayOfMonth[1];
		
		//Determine all the shows for the day
		preg_match_all('/<span class="show"><a href=".*?">(.*?)<\/a>/', $DayInfo, $AllShows);
		foreach($AllShows[1] as $Show)
		{
			$Show=preg_replace('/^\s+|\s+$/', '', $Show); //Remove start and end of line whitespace
			if(!isset($ShowList[$Show])) //If show has not yet been added to overall list, add it
				$ShowList[$Show]=Array();
			$ShowList[$Show][]=$DayOfMonth; //Add this day as a time for the show
		}
	}
	
	//Output the shows and their days
	print '<table>';
	foreach($ShowList as $ShowName => $Days)
		print "<tr><td>$ShowName</td><td>".implode(', ', $Days).'</td></tr>';
	print '</table><br><br>';
}
```

<div>
</div>