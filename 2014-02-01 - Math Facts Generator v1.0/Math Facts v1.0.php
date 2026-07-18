<? header('Content-Type: text/html; charset=utf-8'); ?>
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "https://www.w3.org/TR/html4/loose.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<style type="text/css">
	/* Main table styles */
	.DataTable { font-family:monospace; font-size:15px; float:left; border-collapse:collapse; }
	.DataTable tr>* { border:1px solid black; padding:4px; }
	.AnswerTable { margin-left:100px; }

	/*Column stles*/
	.ColHead { } /* The row number column */
	.Col0 { background-color:#FFFFFF; }
	.Col1 { background-color:#A8A8A8; }
	.Col2 { background-color:#828282; }
	.Col3 { background-color:#A8A8A8; }
	.Col4 { background-color:#FFFFFF; }

	/*Row Styles*/
	.Row0 { border-top:3px solid black !important; }
	.Row1 {} /*Further rows can be defined*/

	/* Other styles */
	@media print { form { display:none; } } /*Do not display the generation form on print*/
	.GenerateText { font-size:12px; font-weight:bold; color:red; }
	.Title { font-weight:bold; }

</style>
<?
//Get integer parameter's value. If it isn't set, use the default
function GetParam($ParamName, $Default) { return isset($_REQUEST[$ParamName]) && is_string($_REQUEST[$ParamName]) ? $_REQUEST[$ParamName] : $Default; }
foreach(Array('Range1Start'=>1, 'Range1End'=>10, 'Range2Start'=>1, 'Range2End'=>10, 'Iterations'=>1, 'RowClassRepeat'=>3, 'NumColumns'=>5) as $VarName => $Default)
	${$VarName}=intval(GetParam($VarName, $Default));
foreach(Array('NumColumns', 'RowClassRepeat') as $Name) //Remove possible division by 0s on integers
	if(${$Name}==0)
		${$Name}=1;

//Process other parameters
$IsGenerated=isset($_REQUEST['Generate']);
foreach(Array('MakeRandom', 'IncludeAnswers', 'IncludeTitle') as $VarName) //Booleans (Default true)
	${$VarName}=(!$IsGenerated || isset($_REQUEST[$VarName]));
$FullTypeNames=Array('Add'=>'Addition', 'Sub'=>'Subtraction', 'Mul'=>'Multiplication', 'Div'=>'Division');
$MathType=GetParam('MathType', 'Add');
$MathType=(!isset($FullTypeNames[$MathType]) ? 'Add' : $MathType); //If MathType is invalid, set it to Add
$Title=$FullTypeNames[$MathType]." (${Range1Start}-$Range1End)x(${Range2Start}-$Range2End)"; //Title

//Output the form to generate tables
function EHTML($S) { return htmlspecialchars($S, ENT_QUOTES, 'UTF-8'); } //Escape for HTML
function TextBox($Name) { global ${$Name}; print "<input type=text name=$Name value='".EHTML(${$Name})."'>"; } //Output a text input
?>
<title>Math Facts<?=($IsGenerated ? ": $Title" : '')?></title>
</head>
<body>
<form action="<?=$_SERVER['SCRIPT_NAME']?>" method="GET">
	<table>
		<tr><td>Range 1 Start</td><td><?=TextBox('Range1Start')?></td></tr>
		<tr><td>Range 1 End</td><td><?=TextBox('Range1End')?></td></tr>
		<tr><td>Range 2 Start</td><td><?=TextBox('Range2Start')?></td></tr>
		<tr><td>Range 2 End</td><td><?=TextBox('Range2End')?></td></tr>
		<tr><td>Iterations</td><td><?=TextBox('Iterations')?></td></tr>
		<tr><td>Math Type</td><td>
<? foreach($FullTypeNames as $Name => $Label)
	print "\t\t\t<input type=radio name=MathType value=$Name".($MathType==$Name ? ' checked' : '').">$Label<br>\n";
?>
		</td></tr>
		<tr><td>Make Random</td><td><input type=checkbox name=MakeRandom <?=($MakeRandom ? 'checked' : '')?>></td></tr>
		<tr><td>Include Answers</td><td><input type=checkbox name=IncludeAnswers <?=($IncludeAnswers ? 'checked' : '')?>></td></tr>
		<tr><td>Include Title</td><td><input type=checkbox name=IncludeTitle <?=($IncludeTitle ? 'checked' : '')?>></td></tr>
		<tr><td>Number of Columns</td><td><?=TextBox('NumColumns')?></td></tr>
		<tr><td>How often row class repeats</td><td><?=TextBox('RowClassRepeat')?></td></tr>
	</table>
	<input type=submit name=Generate value=Generate> <span class=GenerateText>Make sure to confirm your page setup (portrait/landscape, prints backgrounds [including colors], margins, etc). Confirm before first print with a print preview.</span>
</form>
<?
//If parameters are not yet given, do not generate the tables
if(!$IsGenerated)
	exit;

//Output the title
if($IncludeTitle)
	print "<div class=Title>$Title</div>";

//Functions to solve the math types
function Add($x, $y) { return $x+$y; }
function Sub($x, $y) { return $x-$y; }
function Mul($x, $y) { return $x*$y; }
function Div($x, $y) { return $x/$y; }

//Function to check if a math combination is valid
function CheckAdd($x, $y) { return true; }
function CheckSub($x, $y) { return $x>=$y; }
function CheckMul($x, $y) { return true; }
function CheckDiv($x, $y) { return $x%$y==0; }

//Static parameters
$MathSymbols=Array('Add'=>'+', 'Sub'=>'-', 'Mul'=>'x', 'Div'=>'÷');
$ShowMathSymbol=' '.$MathSymbols[$MathType].' ';
$CheckFunc="Check$MathType";

//Create the valid facts
$Nums=Array(); //[]=Array(DisplayFact, Result)
for($x=0;$x<$Iterations;$x++)
	for($i=$Range1Start;$i<=$Range1End;$i++)
		for($n=$Range2Start;$n<=$Range2End;$n++)
			if($CheckFunc($i, $n))
				$Nums[]=Array(FixNum($i).$ShowMathSymbol.FixNum($n), $MathType($i, $n));
function FixNum($S) { return strtr(str_pad($S, 2, ' ', STR_PAD_LEFT), Array(' '=>'&nbsp;')); } //Pad a number to 2 spaces

//Randomize the facts
$NewNums=Array();
if($MakeRandom)
	while(count($Nums))
		$NewNums=array_merge($NewNums, array_splice($Nums, rand(0, count($Nums)-1), 1));
else //If not randomizing
	$NewNums=$Nums;

//Output a table - No need to duplicate code since both the facts and answer table are generated the same
function OutputTable($TableNum) 
{
	global $NewNums, $NumColumns, $RowClassRepeat;
	$RowNum=-1;
	foreach($NewNums as $i => $N)
	{
		if($i%$NumColumns==0)
		{
			$RowNum++;
			print ($i!=0 ? '</tr><tr>' : '').'<td class="ColHead Row'.($RowNum%$RowClassRepeat).'">'.chr(ord('A')+($i/$NumColumns)).')</td>';
		}
		print '<td class="Col'.($i%$NumColumns).' Row'.($RowNum%$RowClassRepeat).'">'.$N[$TableNum].($TableNum==0 ? ' = ___' : '').'</td>';
	}
	$ColsToFillIn=$NumColumns-(count($NewNums)%$NumColumns);
	if($ColsToFillIn!=$NumColumns)
		for($i=0;$i<$ColsToFillIn;$i++)
			print '<td>&nbsp;</td>';
	print '</tr></table>';
}

//Output the facts
print '<table class=DataTable><tr>';
OutputTable(0);

//Output the answers
if($IncludeAnswers)
{
	print '<table class="DataTable AnswerTable"><tr>';
	OutputTable(1);
}
?>
<div style="clear:both"></div>
<script type="text/javascript">window.print()</script>
</body></html>