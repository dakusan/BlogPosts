<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Visual_Studio_IDE_Tab_Order">Visual Studio IDE Tab Order</a><br>
Date: 2010-11-17<br>
SubTitle: Microsoft fails at usability<br>
Tags: Examples, Microsoft, Quick Fixes, Visual Studio<br><br>

<p>I’ve been really annoyed for a while by the unintuitive IDE tab ordering in Visual Studio 2005+. When you type [shift+]alt+tab, you don’t get the next/previous tab in the list as would be the <b>OBVIOUS</b> way to do it (which probably all other IDEs do this right). No, it switches between tabs in an arbitrary hidden order related to the last acces order of the tabs.</p>
<p><a href="http://stackoverflow.com/questions/21027/changing-ctrltab-behavior-for-moving-between-document-in-visual-studio">Searching the internet for a solution to this was pretty fruitless</a>, so I tried to tackle the problem myself. I dug through all the possible structures I could find in the Visual Studio IDE macro explorer, and was unfortunately unable to find where the tab order was kept in a window pane (if it is even accessible to the user). I thought I had the solution at one point, but realized it also just switches tabs in the order they were originally opened 🙁. This is the VB macro code I came up with to do at least that, which uses “DTE.ActiveWindow.Document.Collection” for the tab-open order.</p>
</div>

```vb
Public Sub TabDirection(ByVal Direction As Integer)
	'Find the index of the current tab
	Dim i As Integer
	Dim Index As Integer
	Dim Count As Integer
	Count = DTE.ActiveWindow.Document.Collection.Count
	For i = 1 To Count
		If DTE.ActiveWindow.Document.Collection.Item(i).ActiveWindow.Equals(DTE.ActiveWindow) Then Index = i
	Next i

	'Determine the new index
	Index = Index + Direction
	If Index > Count Then Index = 1
	If Index = 0 Then Index = Count

	DTE.ActiveWindow.Document.Collection.Item(Index).Activate() 'Activate the next tab in the proper direction
End Sub

Public Sub TabForward()
	TabDirection(1)
End Sub

Public Sub TabBackward()
	TabDirection(-1)
End Sub
```

<div>
</div>