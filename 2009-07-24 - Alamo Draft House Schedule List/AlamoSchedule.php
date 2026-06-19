<?php
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
?>