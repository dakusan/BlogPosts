<div>
Original Post: <a href="https://www.castledragmire.com/Posts/AutoHotKey_Scripts">AutoHotKey Scripts</a><br>
Date: 2015-12-01<br>
Tags: AutoHotKey, Examples<br><br>
</div><br>

In lieu of using <a href="http://www.castledragmire.com/Projects/HalfKey">my own custom C++ background services</a> to take care of hot key tasks in Windows, I started using <a href="https://www.autohotkey.com">AutoHotKey</a> a while back. While it’s not perfect, and it is missing a lot of Win32 API functionality, I am still able to mostly accomplish what I want in it. I was thinking I should add some of the simple scripts I use here.

### Center a string within padding characters and output as key-strokes
Example:
<ul>
	<li>PadText = ~*</li>
	<li>Length = 43</li>
	<li>Text = Example Text</li>
	<li>Result = ~*~*~*~*~*~*~*~*Example Text~*~*~*~*~*~*~*~</li>
</ul>
</div>

```ahk
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
```

<hr>

### Format rich clipboard text to plain text
```ahk
clipboard = %clipboard%
return
```

<hr>

### Force window to borderless full screen
Description: This takes the active window, removes all window dressing (titlebar, borders, etc), sets its resolution as 1920x1080, and positions the window at 0x0. In other words, this makes your current window take up the entirety of your primary monitor (assuming it has a resolution of 1920x1080).

```ahk
WinGetActiveTitle, WinTitle
WinSet, Style, -0xC40000, %WinTitle%
WinMove, %WinTitle%, , 0, 0, 1920, 1080
return
```

<hr>

### Continually press key on current window
Description: Saves the currently active window (by its title) and focused control object within the window; asks the user for a keypress interval and the key to press; starts to continually press the requested key at the requested interval in the original control (or top level window if an active control is not found); stops via the F11 key.<br>
Note: I had created this to help me get through the <a href="https://www.castledragmire.com/Posts/LISA_game_difficulty_level_save_hack">LISA</a> intro multiple times.

```ahk
;Get the current window and control
WinGetActiveTitle, TheTitle
ControlGetFocus FocusedControl, %TheTitle%
if(ErrorLevel)
	FocusedControl=ahk_parent

;Get the pause interval
InputBox,IntervalTime,Starting script with window '%TheTitle%',Enter pause interval in milliseconds. After submitted`, hold down the key to repeat,,,,,,,,200
if(ErrorLevel || IntervalTime=="") ;Cancel action if blank or cancelled
	return
IntervalTime := IntervalTime+0

;Get the key to keep pressing - Unfortunately, there is no other way I can find to get the currently pressed keycode besides polling all 255 of them
Sleep 500 ;Barrier to make sure one of the initialization keys is not grabbed
Loop {
	TestKey := 0
	Loop {
		SetFormat, INTEGER, H
		HexTextKey := TestKey
		SetFormat, INTEGER, D
		VirtKey = % "vk" . SubStr(HexTextKey, 3)
		if(GetKeyState(VirtKey)=1 || TestKey>255)
			break
		TestKey:=TestKey+1
	}
	if(TestKey<=255)
		break
	Sleep 500
}
VirtKey := GetKeyName(VirtKey)

;If a direction key, remap to the actual key
if(TestKey>=0x25 && TestKey<=0x28)
	VirtKey := SubStr(VirtKey, 7)

;Let the user know their key
MsgBox Received key: '%VirtKey%'. You may now let go of the key. Hold F11 to stop the script.

;Continually send the key at the requested interval
KeyDelay:=10
SetKeyDelay %KeyDelay% #Interval between up/down keys
IntervalTime-=%KeyDelay%
Loop {
	;Press the key
	ControlSend, %FocusedControl%, {%VirtKey% Up}{%VirtKey% Down}, %TheTitle%

	;Check for the cancel key
	if(GetKeyState("F11"))
		break

	;Wait the requested interval to press the key again
	Sleep, %IntervalTime%
}

;Let the user know the script has ended
MsgBox Ending script with window '%TheTitle%'
return
```