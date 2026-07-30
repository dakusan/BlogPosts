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