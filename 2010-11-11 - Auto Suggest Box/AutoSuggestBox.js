/***Dakusan's AutoSuggestBox - v1.0.1 http://www.castledragmire.com/Projects/Other_Web_Scripts#AutoSuggestBox ***/

var SuggestBoxes=[];
function AutoSuggestBox(InputBox, Items, MaxSuggestions, InheritWidthDiff)
{
	var Me=this;
	$CP(Me, {InputBox:InputBox, Items:Items, Visible:false, Num:SuggestBoxes.length, SelectedIndex:-1, CurrentSuggestions:[], LastDisplayedValue:'', MaxSuggestions:(MaxSuggestions===undefined ? 5 : MaxSuggestions)}); //Fill in properties
	SuggestBoxes.push(Me); //Remember all auto suggestion boxes for lookup

	//Need a parent for the suggest box for layout reasons
	Me.ParentBox=$CE('DIV', {
			'style.position':'relative',
			'style.display':'none',
			className:'SuggestionBoxParent'});
	InputBox.parentNode.insertBefore(Me.ParentBox, InputBox);

	//Create the suggest box
	Me.Box=$CE('DIV', {
		'style.position':'absolute',
		className:'SuggestionBox',
		'style.width':(InheritWidthDiff!==undefined ? InputBox.offsetWidth+$I(InheritWidthDiff)+'px' : '')
		}, Me.ParentBox);

	//Prepare the InputBox to use the AutoSuggestBox
	InputBox.onfocus=InputBox.onkeyup=function(event) { ASB_ShowSuggestions(Me.Num, event); }
	InputBox.onkeypress=function(event) { return ASB_CheckKey(Me.Num, event); }
	InputBox.onblur=function() { ASB_HideSuggestions(Me.Num); }
	InputBox.autocomplete='off';

	return Me;
}

function ASB_CheckKey(Num, event)
{
	var SB=SuggestBoxes[Num]; //Get the AutoSuggestBox
	event=$EV(event);
	if(event.keyCode!=13 || SB.SelectedIndex==-1) //If not enter or no selection, nothing to do
		return;

	//Fill in the selected suggestion
	ASB_DoSuggest(Num, SB.Items[SB.CurrentSuggestions[SB.SelectedIndex]]);

	//Cancel the event
	event.cancelBubble=true;
	return false;
}

function ASB_ShowSuggestions(Num, event)
{
	var SB=SuggestBoxes[Num]; //Get the AutoSuggestBox
	event=$EV(event);
	
	if(event.keyCode==38 || event.keyCode==40) //Up or down
		return ASB_SelectIndex(SB, SB.SelectedIndex+(event.keyCode==38 ? -1 : 1)); //Move the selection up or down
	if(SB.LastDisplayedValue==SB.InputBox.value || $EV(event).keyCode==13) //If value has not changed, nothing to do
		return;
	
	SB.ParentBox.style.display='none'; //Hide the suggestbox

	//Find the suggestions
	var SearchString=SB.InputBox.value.toLowerCase();
	var Suggestions=Array();
	if(SearchString.length>0)
		for(var i in SB.Items)
			if(SB.Items.hasOwnProperty(i) && SB.Items[i].toLowerCase().indexOf(SearchString)!=-1)
			{
				Suggestions.push(i);
				if(Suggestions.length>=SB.MaxSuggestions)
					break;
			}

	//Delete the old elements from the suggest box
	while(SB.Box.childNodes.length)
		SB.Box.removeChild(SB.Box.childNodes[0]);
	SB.SelectedIndex=-1;

	//Add the new suggestions to the box
	for(var i=0;i<Suggestions.length;i++)
	{
		//Create the item parent for the suggestion box
		var SuggestionItem=$CE('DIV', {}, SB.Box);
		var MySuggestion=SB.Items[Suggestions[i]];
		ASB_SetItemInteractivity(SuggestionItem, Num, MySuggestion);

		//Output the suggestion into the item, with the search strings bolded
		var SplitText=MySuggestion.toLowerCase().split(SearchString), StartOn=0;
		for(var n=0;n<SplitText.length-1;n++)
		{
			SuggestionItem.appendChild(document.createTextNode(MySuggestion.substr(StartOn, SplitText[n].length)));
			StartOn+=SplitText[n].length;
			$CE('span', {className:'SuggestionBoxHighlight'}, SuggestionItem).appendChild(document.createTextNode(MySuggestion.substr(StartOn, SearchString.length)));
			StartOn+=SearchString.length;
		}
		SuggestionItem.appendChild(document.createTextNode(MySuggestion.substr(StartOn, SplitText[n].length)));
	}
	SB.CurrentSuggestions=Suggestions;
	SB.LastDisplayedValue=SearchString;

	if(!Suggestions.length) //If no suggestions, exit prematurely
		return;

	SB.ParentBox.style.display=''; //Show the suggestbox
	ASB_SelectIndex(SB, 0); //Select the first item
}

function ASB_SetItemInteractivity(El, Num, Suggestion) { El.onmousedown=function() { ASB_DoSuggest(Num, Suggestion); } }

function ASB_DoSuggest(Num, Suggestion)
{
	SuggestBoxes[Num].InputBox.value=Suggestion;
	ASB_HideSuggestions(Num);
}

function ASB_HideSuggestions(Num)
{
	$CP(SuggestBoxes[Num], {SelectedIndex:-1, LastDisplayedValue:'', 'ParentBox.style.display':'none'});
}

function ASB_SelectIndex(SB, ItemNum)
{
	if(SB.SelectedIndex!=-1) //Unselect old
	{
		var OldItem=SB.Box.childNodes[SB.SelectedIndex];
		OldItem.className=OldItem.className.replace(/ SuggestBoxSelected\b/, ''); //Remove the "SuggestBoxSelected" class name
	}

	SB.SelectedIndex=ItemNum=(ItemNum<0 ? SB.CurrentSuggestions.length-1 : (ItemNum>=SB.CurrentSuggestions.length ? 0 : ItemNum)); //Wrap ItemNum around if it reaches its bounds
	SB.Box.childNodes[ItemNum].className+=" SuggestBoxSelected"; //Add the "SuggestBoxSelected" class name
}

//Common functions (See http://www.castledragmire.com/Projects/Other_Web_Scripts#CommonJavaScript)
var IsIE=(navigator.appName.indexOf("Microsoft")!=-1), IEVer=(IsIE?parseFloat(navigator.appVersion.replace(/^.*MSIE(.*?);.*/i,"$1")):0); //For browser specific functions
var $   =function(O)		{return document.getElementById(O);};
var $GE =function(O)		{return typeof(O)=="string"?$(O):O;};
var $CP =function(O,P)		{for(var i in P)if(P.hasOwnProperty(i))eval("O."+i+"=P[i];");return O;};
var $CE =function(T,P,A)	{T=document.createElement(IsIE&&typeof(P)=="object"&&'name' in P?'<'+T+' name="'+P.name+'">':T);if(P)$CP(T,P);if(A)$GE(A).appendChild(T);return T;};
var $EV =function(E)		{return IsIE?window.event:E;};
var $I  =function(N)		{return parseInt(N, 10);};