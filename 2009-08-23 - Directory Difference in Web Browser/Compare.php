<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "https://www.w3.org/TR/html4/loose.dtd">
<html>
<head>
<title>Diff</title>
<style type="text/css">
.Dif { color: red; }			/* A file containing differences */
.DifD { color: red; }			/* A directory containing differences, this must be separated from files so when only showing "same" files, directories are not hidden */
.Equ { color: black; }			/* An item containing no differences */
.OnL { color: blue !important; }	/* An item only found on the left side */
.OnR { color: green !important; }	/* An item only found on the right side */
.Dir { font-weight: bold; }		/* A directory */
.File { font-weight: normal;}		/* A file */
.Info { color: purple; }		/* An item's extra info */

.Legend { display: block !important; }	/* A legend example */
.ToggleLink { text-decoration:none; }	/* A directory display toggle */
.HeadSection span { font-weight: bold; text-decoration:underline; } 
.HeadSection ul { margin-top:0; } 
</style>
</head>
<body>
<div class=HeadSection>
	<span>Legend:</span>
	<ul>
		<li>Differences:
			<ul>
				<li class="Legend Equ">No differences detected</li>
				<li class="Legend Dif">File is different, or directory contains differences</li>
				<li class="Legend OnL">File is only found on left side, when this happens to a directory, it is only counted as 1 difference</li>
				<li class="Legend OnR">File is only found on right side, when this happens to a directory, it is only counted as 1 difference</li>
			</ul>
		</li>
		<li>Type:
			<ul>
				<li class="Legend Dir">This is a directory with sub items. After the directory name, it lists as info the number of differences over the total number of sub items</li>
				<li class="Legend File">This is a file. If both sides contain the file but are different, it lists as info the different strings</li>
			</ul>
		</li>
		<li class="Legend Info">Info about the file/directory listed in between parenthesis.</li>
	</ul>
</div>

<div class=HeadSection>
	<span>Options:</span>
	<ul>
		<li><label><input type=radio name=DisplayType value=All checked onclick="ChangeDisplay(this);">Display all items</label></li>
		<li><label><input type=radio name=DisplayType value=Same onclick="ChangeDisplay(this);">Display only same files (Displays all directories)</label></li>
		<li>
			<label><input type=radio name=DisplayType value=Different onclick="ChangeDisplay(this);">Display different items (Does not display synced directories)</label>
			<ul>
				<li><label><input type=checkbox name=DoNotDisplayLeftOnly onclick="ChangeDisplay(this);">Do not display left only items</label></li>
				<li><label><input type=checkbox name=DoNotDisplayRightOnly onclick="ChangeDisplay(this);">Do not display right only items</label></li>
			</ul>
		</li>
		<li><input type=button value="Collapse all items" onclick="ToggleAll(RootList, false);"></li>
		<li><input type=button value="Expand all items" onclick="ToggleAll(RootList, true);"></li>
	</ul>
</div>

<?
//The 2 lists of files to compare
$File1="./Dir1.txt";
$File2="./Dir2.txt";

//Output starting info
print "Left Side: $File1<br>Right Side: $File2<br>";

//Other Options
$FileSeparator="\n"; //Should be "\r\n" for Windows text files
$DirectorySeparator="/"; //Should be "\\" for Windows
$StringSeparator="\t"; //The file's data string separator

//Extract the directory/file contents structures from both lists
$Contents=Array(file_get_contents($File1), file_get_contents($File2));
$Structures=Array();
for($i=0;$i<2;$i++)
{
	$MyStructure=Array(); //Create the contents structure root
	$MyFiles=explode($FileSeparator, $Contents[$i]); //Extract the list of files and their confirmation string
	foreach($MyFiles as $File) //Add each file to the contents structure
	{
		$FileInfo=explode("\t", $File); //Separate the path and the file confirmation string
		if(count($FileInfo)!=2) //Bypass invalid lines
			continue;

		$FilePath=explode($DirectorySeparator, $FileInfo[0]); //Extract the parent directories and the filename
		$NumDirectories=count($FilePath)-1; //The number of directories is 1 minus the number of entries because the last entry is the file name
		$CurrentDirectory=&$MyStructure; //The top level array is the top parent
		for($ParentNum=0;$ParentNum<$NumDirectories;$ParentNum++) //Add the parent directories to the contents structure
		{
			$CurrentDirectory=&$CurrentDirectory[$FilePath[$ParentNum]]; //Get the next parent directory in the path list
			if(!isset($CurrentDirectory)) //If next parent directory does not already exist, create it
				$CurrentDirectory=Array();
		}
		
		$CurrentDirectory[$FilePath[$NumDirectories]]=$FileInfo[1]; //Set the filename in the contents structure to its confirmation string
	}

	$Structures[$i]=$MyStructure; //Remember the structure
}

//Output directory differences
$Root=DirDiff($Structures[0], $Structures[1]);
print "Total Differences: $Root[NumDifferences]<br>Total Items: $Root[NumItems]<br><div id=Root>$Root[DiffString]";

function DirDiff($Dir1, $Dir2) //Returns Array(DiffString, NumItems, NumDifferences)
{
	//Determine all the file and subdirectory names in both directories
	$DirectoryNames=Array();
	$FileNames=Array();
	foreach(Array($Dir1, $Dir2) as $CurDir)
		foreach($CurDir as $ItemName => $Item)
			if(is_array($Item))
				$DirectoryNames[$ItemName]=0;
			else
				$FileNames[$ItemName]=0;
	
	//Sort directory and file lists into a natural case order
	$DirectoryNames=array_keys($DirectoryNames);
	$FileNames=array_keys($FileNames);
	natcasesort($DirectoryNames);
	natcasesort($FileNames);

	//Output list of items (directories then files) within this directory
	$NumDifferences=0;
	$NumItems=0;
	$RetVal='<ul>'; //Start a directory listing
	foreach($DirectoryNames as $Directory)
	{
		$OneSideOnly=(!isset($Dir1[$Directory]) ? ' OnR' : (!isset($Dir2[$Directory]) ? ' OnL' : '')); //Determine if only one of the directories has this item
		$SubInfo=DirDiff($OneSideOnly==' OnR' ? Array() : $Dir1[$Directory], $OneSideOnly==' OnL' ? Array() : $Dir2[$Directory]); //Recurse through subdirectories (use empty array if one of the sides does not contain this entry)
		$NumItems+=$SubInfo['NumItems']; //Add number of items in directory to total number of items
		$NumDifferences+=$MyNumDifferences=($OneSideOnly=='' ? $SubInfo['NumDifferences'] : 1); //Add number of differences in directory to total number of differences. If this item is only on one side, consider it only 1 difference.

		$RetVal.=
			'<li class="Dir '. //Start the directory item, and set it as a "Directory"
			($SubInfo['NumDifferences'] ? 'DifD' : 'Equ'). //Is the directory item "Equal" or "Different"
			$OneSideOnly. //Set if directory is a left or right side only
			'">'.htmlentities($Directory)." <span class=Info>($MyNumDifferences/$SubInfo[NumItems])</span> $SubInfo[DiffString]</li>"; //The directories: Name, Num items, num differences, sub content; and end of the directory item
	}
	foreach($FileNames as $File)
	{
		$OneSideOnly=(!isset($Dir1[$File]) ? ' OnR' : (!isset($Dir2[$File]) ? ' OnL' : '')); //Determine if only one of the directories has this item
		$IsDiff=($OneSideOnly || $Dir1[$File]!=$Dir2[$File]);
		$NumDifferences+=($IsDiff ? 1 : 0); //Add 1 to number of differences if files do not match
		
		$RetVal.=
			'<li class="File '. //Start the directory item, and set it as a "File"
			($IsDiff ? 'Dif' : 'Equ'). //Is the directory item "Equal" or "Different"
			$OneSideOnly. //Set if directory is a left or right side only
			'">'.htmlentities($File).(!$OneSideOnly && $IsDiff ? " <span class=Info>($Dir1[$File] :: $Dir2[$File])</span>" : '').'</li>'; //The files name, differences, and end of the directory item
	}
	$NumItems+=count($FileNames); //Add number of files to total number of items
	
	$RetVal.='</ul>'; //End a directory listing

	return Array('DiffString'=>$RetVal, 'NumItems'=>$NumItems, 'NumDifferences'=>$NumDifferences);
}
?>
</div>
<script type="text/javascript">
function ChangeDisplay(SwitchObject) //Hides/Shows items depending on their "differences" state
{
	if(SwitchObject.name=='DisplayType')
	{
		ModifyCSS('Dif', 'display', SwitchObject.value!='Same' ? '' : 'none');
		ModifyCSS('Equ', 'display', SwitchObject.value!='Different' ? '' : 'none');
		return;
	}
	
	ModifyCSS(SwitchObject.name=='DoNotDisplayLeftOnly' ? 'OnL' : 'OnR', 'display', SwitchObject.checked ? 'none' : '');
}
function ModifyCSS(Name, Rule, Value) //Modify a global CSS rule (1 stylesheet assumed)
{
	var MyStyleSheet=GetStyleSheet('Dif');
	for(var i=0;i< MyStyleSheet.length;i++)
		if(MyStyleSheet[i].selectorText=='.'+Name)
		{
			MyStyleSheet[i].style[Rule]=Value;
			break;
		}
}
function GetStyleSheet(FirstRule)
{
	var RulesName=(document.styleSheets[0].cssRules ? 'cssRules' : 'rules');
	for(var i=0;i< document.styleSheets.length;i++)
		if(document.styleSheets[i][RulesName][0].selectorText=="."+FirstRule)
			return document.styleSheets[i][RulesName];
}

function ToggleDirectory(TheDirectory, DoExpand) //Expand or contract a single directory
{
	DoExpand=(DoExpand===undefined ? TheDirectory.firstChild.innerHTML=='+' : DoExpand); //If DoExpand is not defined, swap the current direction
	TheDirectory.firstChild.innerHTML=(DoExpand ? '-' : '+');
	TheDirectory.lastChild.style.display=(DoExpand ? '' : 'none');
	return false; //Returns false so the anchor link is not followed
}

function ToggleAll(ListObject, DoExpand) //Expand or contract all directories (recursive)
{
	for(var i=0, ChildNode;i< ListObject.childNodes.length;i++)
		if((ChildNode=ListObject.childNodes[i]).className.indexOf('Dir')!='-1')
		{
			ToggleAll(ChildNode.lastChild, DoExpand);
			ToggleDirectory(ChildNode, DoExpand);
		}
}

function AddToggles(ListObject) //Add directory expanding/collapsing toggle link (recursive). This is done in the JavaScript phase to make the returned HTML document smaller.
{
	for(var i=0, ChildNode;i< ListObject.childNodes.length;i++)
		if((ChildNode=ListObject.childNodes[i]).className.indexOf('Dir')!='-1')
		{
			AddToggles(ChildNode.lastChild);
			var NewLink=document.createElement('span');
			NewLink.innerHTML='<a href="#" onclick="return ToggleDirectory(this.parentNode);" class="ToggleLink">-<'+'/a>';
			ChildNode.insertBefore(NewLink.firstChild, ChildNode.firstChild);
		}
}

//Prepare the document by adding toggle links and contracting all directories
var RootList=document.getElementById('Root').firstChild;
AddToggles(RootList);
ToggleAll(RootList, false);
</script>
</body>
</html>