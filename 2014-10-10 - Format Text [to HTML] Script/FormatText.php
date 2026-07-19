<? header('Content-Type: text/html; charset=utf-8'); ?>
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "https://www.w3.org/TR/html4/loose.dtd">
<html>
<head>
<title>Format Text</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<?
$AllowRenderText=true; //Set to true only if this is in a secure environment, as directly outputting a given value can lead to XSS
if(isset($_REQUEST['RenderText']))
	return print '</head><body>'.($AllowRenderText ? $_REQUEST['RenderText'] : 'Rendering of text not allowed').'</body></html>';
?>
<style type="text/css">
html, body { width:100%; height:100%; margin:0; padding:0; }
.HalfScreen { display:block; width:calc(100% - 2px); height:calc(50% - 2px - 30px/2); margin:0; border:1px solid black; }
#RenderForm { overflow:hidden; }
#RenderText { margin:0; border:0; width:100%; height:100%; }
#RenderHTML { overflow-x:hidden; overflow-y:scroll; }
.TopBar { height:30px; background-color:grey; }
.Hide { position:absolute; visibility:hidden; top:-10000px; }
</style>

<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
<script type="text/javascript">$(document).ready(function() {

//History for undoing
var UndoBuf=[], RedoBuf=[];
function Undo()
{
	if(!UndoBuf.length)
		return;
	RedoBuf.push(UndoBuf.pop());
	$('#RenderText').val(UndoBuf[UndoBuf.length-1]);
	$('#RenderHTML').html(UndoBuf[UndoBuf.length-1]);
}
function Redo()
{
	if(!RedoBuf.length)
		return;
	$('#RenderText').val(RedoBuf[RedoBuf.length-1]);
	$('#RenderHTML').html(RedoBuf[RedoBuf.length-1]);
	UndoBuf.push(RedoBuf.pop());
}
$('#Undo').click(function(e) { e.preventDefault(); Undo(); });
$('#Redo').click(function(e) { e.preventDefault(); Redo(); });

//Render HTML
function Render() {
	//Do the render
	var MyVal=$('#RenderText').val();
	$('#RenderHTML').html(MyVal);

	//Save current value to the history
	//*Better history functionality here would be real nice (using smart currentTarget.selectionStart/End calculations), along with an undo/redo button, but not within the scope of this project
	if(RedoBuf.length) //Empty redo buffer
		RedoBuf=[];
	UndoBuf.push(MyVal);
	if(UndoBuf.length>100) //Limit history buffer
		UndoBuf.shift();
}
$('#RenderText').on('keypress paste', function(e) { setTimeout(Render, 1); }); //Automatic update on paste requires a timeout

//Open in new page
$('#OpenInNewPage').click(function(e) {
	e.preventDefault();
	$('#RenderForm').submit();
});

//Escape HTML
$('#EscapeHTML').click(function(e) {
	e.preventDefault();
	$('#RenderText').val(function(index, value) {
		$.each({"&amp;":/&/g, "&lt;":/</g, "&gt;":/>/g, "&quot;":/"/g, "&#039;":/'/g}, function(HTMLStr, ReplStr) {
			value=value.replace(ReplStr, HTMLStr); });
		return value;
	});
	Render();
});

//Listize based on tabbing
//If a successive line is tabbed over beyond the current, it is made inside a new nested list.
//Tabbing over more than once on a successive line will create multiple nests
//Having @@@ at the beginning of a line will include it in the previous line item, no matter the tabbing
//Make sure to have @@@ blank lines tabbed over to the proper nested level
$('#Listize').click(function(e) {
	//Get the text to replace
	e.preventDefault();
	var T=$('#RenderText').val();

	//Go over each line and if the next line is tabbed beyond it, make it a new nested list. Blank
	var CurTabLevel=0, NewLines=[]; //NewLines is 2 items per line: the original string and the new html tags
	$.each(T.split(/\r?\n/), function(Index, Str) {
		//Check for a continued line item
		if(Str.substr(0, 3)=='@@@')
			return NewLines.push('<br>', Str.substr(3));

		//In/de-dent as needed
		var Tags='';
		var NewTabLevel=/^\t*/.exec(Str)[0].length, PreLevel=CurTabLevel; //Get the nested level
		for(;NewTabLevel>CurTabLevel;CurTabLevel++)
			Tags+='<ul><li>';
		for(;NewTabLevel<CurTabLevel;CurTabLevel--)
			Tags+='</li></ul>';

		//Fill out the rest of the line
		if(NewTabLevel==0) //Breaks between top level new lines
			Tags+=(Index && PreLevel==0 ? '<br>' : '');
		else if(PreLevel>=NewTabLevel) //If previous item needs to be ended (new level is not greater and not 0)
			Tags+='</li><li>';

		NewLines.push(Tags, Str);
	});

	//Finish de-dent as needed
	var Final=[NewLines.shift()];
	var EndLine='';
	while(CurTabLevel--)
		EndLine+='</li></ul>';
	NewLines.push(EndLine);

	//Combine each line with the tags
	for(var i=0;i<NewLines.length;i+=2)
		Final.push(NewLines[i+0]+NewLines[i+1]);



	//Update from the replaced text
	$('#RenderText').val(Final.join("\n"));
	Render();
});

});</script>

</head>
<body>
	<div class=TopBar>
		<input type=button id=EscapeHTML value="Escape HTML">
		<input type=button id=Listize value="Listize">
		<? if($AllowRenderText) { ?> <input type=button id=OpenInNewPage value="Open In New Page"> <? } ?>
		<input type=button id=Undo value="Undo">
		<input type=button id=Redo value="Redo">
	</div>
	<form action="FormatText.php" method=post id=RenderForm target="_blank" class=HalfScreen>
		<textarea id=RenderText name=RenderText></textarea>
		<input type=submit class=Hide>
	</form>
	<div id=RenderHTML class=HalfScreen></div>
</body>
</html>