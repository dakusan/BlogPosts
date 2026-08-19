;Function to create lparam/wparam for SendMessage
CalculatePARAM(w1, w2)
{
	IfLess, w1, 0
		w1 := 65535 + w1 + 1
	IfLess, w2, 0
		w2 := 65535 + w2 + 1

	return (w2<<16 | w1)
}

;Send a mouse wheel action to a window
SendMouseWheel(WindowHWND, Steps, XPos, YPos)
{
	;Constants
	WM_MOUSEWHEEL := 0x20A
	WheelStepAmount := 120

	;Calculate and execute the message
	WinGetPos, ScreenX, ScreenY,,, ahk_id %WindowHWND%
	wparam := CalculatePARAM(0, Steps*WheelStepAmount)
	lparam := CalculatePARAM(XPos+ScreenX, YPos+ScreenY)
	SendMessage, %WM_MOUSEWHEEL%, %wparam%, %lparam%,, ahk_id %WindowHWND%
}

^!WheelUp::
ControlGet, ControlHWND, Hwnd,,AfxFrameOrView100s1,RME TotalMix
if ControlHWND
	SendMouseWheel(ControlHWND, 1, 36, 428)
return

^!WheelDown::
ControlGet, ControlHWND, Hwnd,,AfxFrameOrView100s1,RME TotalMix
if ControlHWND
	SendMouseWheel(ControlHWND, -1, 36, 428)
return