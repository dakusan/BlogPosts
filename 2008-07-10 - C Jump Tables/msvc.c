void ExamplePointerRetreival()
{
	void *LabelPointer;
	TheLabel:
	_asm mov LabelPointer, offset TheLabel
}
