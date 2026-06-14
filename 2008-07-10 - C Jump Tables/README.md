<div>
Original Post: <a href="https://www.castledragmire.com/Posts/C_Jump_Tables">C Jump Tables</a><br>
Date: 2008-07-10<br>
SubTitle: The unfortunate reality of different feature sets in different language implementations<br>
Tags: Assembly, C, Examples, Jump Table, Programming hacks<br><br>

<p>I was thinking earlier today how it would be neat for C/C++ to be able to get the address of a jump-to label to be used in jump tables, specifically, for an <a href="https://www.castledragmire.com/Projects/Hynes">emulator</a>. A number of seconds after I did a Google query, I found out it is possible in gcc (the open source native Linux compiler) through the “label value operator” “&amp;&amp;”. I am crushed that MSVC doesn’t have native support for such a concept 🙁.</p>
<p>The reason it would be great for an emulator is for emulating the CPU, in which, usually, each first byte of a CPU instruction’s opcode [see ASM] gives what the instruction is supposed to do. An example to explain the usefulness of a jump table is as follows:</p>
</div>

```c
void DoOpcode(int OpcodeNumber, <b style="color:red">...</b>)
{
	void *Opcodes[]={&&ADD, &&SUB, &&JUMP, &&MUL}; //assuming ADD=opcode 0 and so forth
	goto *Opcodes[OpcodeNumber];
  	ADD:
		//<b style="color:red">...</b>
	SUB:
		//<b style="color:red">...</b>
	JUMP:
		//<b style="color:red">...</b>
	MUL:
		//<b style="color:red">...</b>
}
```

<div>
<br>
<p>Of course, this could still be done with virtual functions, function pointers, or a switch statement, but those are <i>theoretically</i> much slower. Having them in separate functions would also remove the possibility of local variables.</p>
<p>Although, again, theoretically, it wouldn’t be <i>too</i> bad to use, I believe, the _fastcall function calling convention with function pointers, and modern compilers SHOULD translate switches to jump tables in an instance like this, but modern compilers are so obfuscated you never know what they are really doing.</p>
<p>It would probably be best to try and code such an instance so that all 3 methods (function pointers, switch statement, jump table) could be utilized through compiler definitions, and then profile for whichever method is fastest and supported.</p>
</div>

```c
//Define the switch for which type of opcode picker we want
#define UseSwitchStatement
//#define UseJumpTable
//#define UseFunctionPointers

//Defines for how each opcode picker acts
#if defined(UseSwitchStatement)
	#define OPCODE(o) case OP_##o:
#elif defined(UseJumpTable)
	#define OPCODE(o) o:
	#define GET_OPCODE(o) &&o
#elif defined(UseFunctionPointers)
	#define OPCODE(o) void Opcode_##o()
	#define GET_OPCODE(o) (void*)&Opcode_##o
	//The above GET_OPCODE is actually a problem since the opcode functions aren't listed until after their ...
	//address is requested, but there are a couple of ways around that I'm not going to worry about going into here.
#endif

enum {OP_ADD=0, OP_SUB}; //assuming ADD=opcode 0 and so forth
void DoOpcode(int OpcodeNumber, <b style="color:red">...</b>)
{
	#ifndef UseSwitchStatement //If using JumpTable or FunctionPointers we need an array of the opcode jump locations
		void *Opcodes[]={GET_OPCODE(ADD), GET_OPCODE(SUB)}; //assuming ADD=opcode 0 and so forth
	#endif
	#if defined(UseSwitchStatement)
		switch(OpcodeNumber) { //Normal switch statement
	#elif defined(UseJumpTable)
		goto *Opcodes[OpcodeNumber]; //Jump to the proper label
	#elif defined(UseFunctionPointers)
		*(void(*)(void))Opcodes[OpcodeNumber]; //Jump to the proper function
		} //End the current function
	#endif

	//For testing under "UseFunctionPointers" (see GET_OPCODE comment under "defined(UseFunctionPointers)")
	//put the following OPCODE sections directly above this "DoOpcode" function
	OPCODE(ADD)
	{
		//<b style="color:red">...</b>
	}
	OPCODE(SUB)
	{
		//<b style="color:red">...</b>
	}

	#ifdef UseSwitchStatement //End the switch statement
	}
	#endif

#ifndef UseFunctionPointers //End the function
}
#endif
```

<div>
<br>

After some tinkering, I did discover through assembly insertion it was possible to retrieve the offset of a label in MSVC, so with some more tinkering, it could be utilized, though it might be a bit messy.
</div>

```c
void ExamplePointerRetreival()
{
	void *LabelPointer;
	TheLabel:
	_asm mov LabelPointer, offset TheLabel
}
```

<div>
</div>