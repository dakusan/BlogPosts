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
