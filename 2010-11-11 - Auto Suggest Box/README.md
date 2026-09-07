<div>
Original Web Script: <a href="https://www.castledragmire.com/Projects/Other_Web_Scripts#AutoSuggestBox">AutoSuggestBox v1.0.1</a><br>
Date: 2010-11-11<br>
</div><br>

A dropdown box for HTML text boxes that gives case insensitive suggestions from a predefined list for what the user is typing. The search phrase is highlighted in the search results. Search results can be selected through keyboard or mouse.

See <a href="https://www.castledragmire.com/Projects/Other_Web_Scripts#AutoSuggestBox">project page</a> for live example.

## Information
<ul>
	<li>Global variables:
	<ul>
		<li><a id=ASB_SuggestBoxes></a><kbd>SuggestBoxes</kbd>: This holds a list of all created auto suggest boxes.</li>
	</ul></li>
	<li>Functions:
	<ul>
		<li><b>AutoSuggestBox</b>
		<ul>
			<li>Description: Creates an “AutoSuggestBox” object. Make sure to call with “new”.</li>
			<li>Returns: The newly created AutoSuggestBox object.</li>
			<li><a id=ASB_PassedParameters></a>Parameters:
			<ul>
				<li><a id=ASB_InputBox></a>1) <kbd>InputBox</kbd>: The input text box</li>
				<li><a id=ASB_Items></a>2) <kbd>Items</kbd>: The array of auto suggest item stings</li>
				<li><a id=ASB_MaxSuggestions></a>[3]) <kbd>MaxSuggestions</kbd>: The maximum number of suggestions. Default is 5.</li>
				<li><a id=ASB_InheritWidthDiff></a>[4]) <kbd>InheritWidthDiff</kbd>: If set, the width of the AutoSuggestBox is the width of the <a href="#ASB_InputBox">InputBox</a> plus the given value.</li>
			</ul></li>
			<li>Object Members:
			<ul>
				<li><kbd>InputBox</kbd>, <kbd>Items</kbd>, <kbd>MaxSuggestions</kbd>: <a href="#ASB_PassedParameters">The passed parameters to the function</a></li>
				<li><a id=ASB_Visible></a><kbd>Visible</kbd>: Whether the auto suggest box is currently visible</li>
				<li><a id=ASB_Num></a><kbd>Num</kbd>: The index of the object in the <a href="#ASB_SuggestBoxes">SuggestBoxes</a> global variable</li>
				<li><a id=ASB_ParentBox></a><kbd>ParentBox</kbd>: A parent container that holds the auto suggest box. It is inserted directly before the <a href="#ASB_InputBox">InputBox</a> in the DOM. It has the class of “<a href="#ASB_SuggestionBoxParent">SuggestionBoxParent</a>”.</li>
				<li><a id=ASB_Box></a><kbd>Box</kbd>: The auto suggest box element held in <a href="#ASB_ParentBox">ParentBox</a>. It has the class of “<a href="#ASB_SuggestionBox">SuggestionBox</a>”. Suggestions are held in sub-DIVs, with selected text in SPANs with the class name “<a href="#ASB_SuggestBoxSelected">SuggestBoxSelected</a>”.</li>
				<li><a id=ASB_SelectedIndex></a><kbd>SelectedIndex</kbd>: The index of the currently selected item. -1 if no selection.</li>
				<li><a id=ASB_CurrentSuggestions></a><kbd>CurrentSuggestions</kbd>: The currently showing suggestions.</li>
			</ul></li>
		</ul></li>
	</ul></li>
	<li>CSS Classes: (Default styles)
	<ul>
		<li><a id=ASB_SuggestionBoxParent></a><kbd>.SuggestionBoxParent</kbd> :: See AutoSuggestBox.<a href="#ASB_ParentBox">ParentBox</a>
		<br><code>position:relative; display:none; /*required styles set in JavaScript*/</code></li>
		<li><a id=ASB_SuggestionBox></a><kbd>.SuggestionBox</kbd> :: See AutoSuggestBox.<a href="#ASB_Box">Box</a>
		<pre><code>position:absolute; /*required style set in JavaScript*/
top:20px; z-index:100; background-color:lightgrey; width:200px;</code></pre>
		</li>
		<li><a id=ASB_SuggestionBoxItem></a><kbd>.SuggestionBox div</kbd> :: An individual suggestion in the auto suggest box
		<br><code>color:black; overflow-x:hidden; white-space:nowrap;</code></li>
		<li><a id=ASB_SuggestionBoxItemHover></a><kbd>.SuggestionBox div:hover</kbd> :: When a suggestion is moused over in the auto suggest box. This may not work in Internet Explorer without a compatible doctype.
		<br><code>background-color:blue;</code></li>
		<li><a id=ASB_SuggestionBoxHighlight></a><kbd>.SuggestionBoxHighlight</kbd> :: Matched search text in a suggestion
		<br><code>font-weight:bold;</code></li>
		<li><a id=ASB_SuggestBoxSelected></a><kbd>.SuggestBoxSelected</kbd> :: The currently selected item in a suggest box
		<br><code>background-color:darkgrey;</code></li>
	</ul></li>
	<li>Other Information:
	<ul>
		<li>Only some browsers support turning off real autocomplete on a per-textbox basis. You may need to add “autocomplete=off” to your form to turn off browser autocomplete.</li>
	</ul></li>
</ul>
<br><br>[Whitespace left at bottom so inner document links place the correct line at the top]<br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br>