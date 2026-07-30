;Get the last values
IniPath=%A_ScriptDir%\AutoHotKey.ini
IniRead,PadText,%IniPath%,CenterString,PadText,-
IniRead,NewLength,%IniPath%,CenterString,NewLength,10
IniRead,TheString,%IniPath%,CenterString,TheString,The String

;Get the input
InputBox,PadText,Center String,Pad Character,,,,,,,,%PadText%
InputBox,NewLength,Center String,New Length,,,,,,,,%NewLength%
InputBox,TheString,Center String,String To Center,,,,,,,,%TheString%

;Cancel on blank pad or invalid number
if StrLen(PadText)==0
{
	MsgBox,Pad text cannot be blank
	return
}
if NewLength is not integer
{
	MsgBox,New length must be an integer
	return
}

;Save the last values
IniWrite,%PadText%,%IniPath%,CenterString,PadText
IniWrite,%NewLength%,%IniPath%,CenterString,NewLength
IniWrite,%TheString%,%IniPath%,CenterString,TheString

;Initial padding
PadStrLen:=StrLen(PadText)
PadLen:=NewLength-StrLen(TheString)
NewString:=""
Loop
{
	if StrLen(NewString)>=Ceil(PadLen/2)
		break
	NewString.=PadText
}

;Truncate initial padding to at least half
NewString:=Substr(NewString, 1, Ceil(PadLen/2))

;Add the string
NewString.=TheString

;Final padding
Loop
{
	if StrLen(NewString)>=NewLength
		break
	NewString.=PadText
}

;Truncate to proper length
NewString:=Substr(NewString, 1, NewLength)

;Output to console
Sleep,100
Send %NewString%
return