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