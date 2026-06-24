;See https://www.castledragmire.com/Posts/My_New_Boot_Loader for more information

ORG 0x7C00 ;Origin location
jmp main ;Bypass starting data secction

;Screen dimensions
ScreenWidth equ 80
ScreenHeight equ 25
ScreenCharacters equ ScreenWidth*ScreenHeight

;String to display in center of screen
DisplayString db 'RABiD BUNNY FEVER',0 ;The string
DisplayStringLen equ $-DisplayString-1 ;The string's length
DisplayStringCenter equ ((ScreenWidth-DisplayStringLen)/2+(ScreenHeight/2)*ScreenWidth)*2 ;The point in the buffer to start displaying the string at

main: ;Bypass starting data section

;Prepare the frame buffer
mov ax,1000h
mov es,ax
;Read local data by default
xor ax,ax
mov ds,ax

;Hide the cursor
mov ah,0x1; ;Hide the cursor
mov cx,0x2706 ;Scanline start > Scanline end, Bit 5 Set=Hide cursor
int 0x10 ;BIOS Video call

;Prepare the random number generator
mov si,22 ;Random number seed
mov di,0x5DE7 ;Random number generator

DrawFrame: ;Draw a frame
	mov ecx,ScreenCharacters ;Screen buffer size
	PickCharacter: ;Add a Character to the frame buffer
		;Choose the next character randomly
		mov ax,di
		mul si
		add ax,0xB
		ror ax,13
		mov si,ax

		;If character is a non displayable character, choose another it
		cmp al,0
		je PickCharacter
		cmp al,0x7
		je PickCharacter
		cmp al,0x8
		je PickCharacter
		cmp al,0xA
		je PickCharacter
		cmp al,0xD
		je PickCharacter
		cmp al,0xFF
		je PickCharacter

		mov [es:ecx*2-2],al ;Save character to the frame buffer

		PickColor: ;Choose the color for the character
			;Choose the next color randomly
			mov ax,di
			mul si
			add ax,0xB
			ror ax,13
			mov si,ax

			and al,0xF ;Only provide foreground color (Remove to also choose a background color)

			;Make sure the fore and background colors dont match (This could be changed to just make sure al!=0 since we are no longer changing background too)
			mov ah,al
			shr ah,4
			mov bh,al
			and bh,0xF
			cmp ah,bh
			je PickColor

			mov [es:ecx*2-1],al ;Save color to the frame buffer

		dec ecx ;Decrement character counter
		jnz PickCharacter; If more characters to choose, pick the next one

	;Add the display string to the center of the screen
	mov ecx,DisplayStringLen
	MoveDisplayCharacter:
		mov al,[DisplayString+ecx-1] ;Grab the current character
		cmp al,32 ;Test against space
		je BypassDisplayCharacter ;Do not display spaces
		mov [es:ecx*2-2+DisplayStringCenter],al ;Save the character to the frame buffer
;		mov byte [es:ecx*2-1+DisplayStringCenter],0xf ;Statically set the color (Decided to comment out)
		BypassDisplayCharacter:
		dec ecx ;Decrement character counter
		jnz MoveDisplayCharacter; If more characters, set the next one

	;Draw the frame buffer (minus one character so a new line doesnt occur at the end)
	mov ah,0x13 ;Draw the character to the screen
	mov al,3 ;bit 0=Update cursor after writing=true, bit 1=String contains alternating characters and attributes=true, bits 2-7=Reserved (0).
	mov bh,0 ;Page number
	mov bl,0xF; color
	mov cx,ScreenCharacters-1 ;Screen buffer size
	mov dx,0 ;row,column at which to start writing.
	int 0x10 ;BIOS Video call

	;Write the final character
	mov ah,0x9 ;Write character+color
	xor ecx,ecx ;ecx=0
	mov al,[es:ecx+ScreenCharacters*2-2] ;The character at the end of the buffer
	mov bl,[es:ecx+ScreenCharacters*2-1] ;The color at the end of the buffer
	mov cx,1 ;Write character once
	int 0x10 ;BIOS Video call

;	;Wait for some time to pass before starting next frame
;	call ReadTimer
;	mov bx,ax ;Remember the start timer
;	mov cx,500 ;Number of PTI ticks to wait
;	WaitForNextFrame: ;Wait a fixed amount of time until starting the next frame
;		hlt ;Do nothing for a bit
;		call ReadTimer ;Read the millisecond count
;		sub ax,bx ;eax=Ellapsed time
;		cmp ax,cx ;Test against required wait time
;		jl WaitForNextFrame ;If not enough time has ellapsed, keep waiting

;	;Wait for keyboard input to repeat
;	mov ah,0
;	int 0x16

	jmp DrawFrame ;Draw the next frame

;;Read the Programmable Interval Timer
;ReadTimer:
;	pushfd
;	cli
;	mov al, 00000000b    ; al = channel in bits 6 and 7, remaining bits clear
;	out 0x43, al         ; Send the latch command
;	in al, 0x40          ; al = low byte of count
;	mov ah, al           ; ah = low byte of count
;	in al, 0x40          ; al = high byte of count
;	rol ax, 8            ; al = low byte, ah = high byte (ax = current count)
;	popfd
;	ret

;Fill out the rest of the 512 byte sector with 0s followed by the boot loader signature
times 510-($-$$) db 0 ;Fill the rest with zeros
dw 0xAA55 ;Boot loader signature