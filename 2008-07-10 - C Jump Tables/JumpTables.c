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