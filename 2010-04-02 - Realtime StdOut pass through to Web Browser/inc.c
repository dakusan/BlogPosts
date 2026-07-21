#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[])
{
	int NumLoops=(argc>1 ? atoi(argv[1]) : 50); //Number of loops to run from passed argument 1. Default is 50 if not specified.
	int LoopWait=(argc>2 ? atoi(argv[2]) : 500); //Number of milliseconds to wait in between each loop from passed argument 2. Default is 500ms if not specified.
	LoopWait*=1000; //Convert to microseconds for usleep

	//Output an incremented number every half a second
	int i=0;
	while(i<NumLoops)
	{
		printf("%u\n", i++);
		fflush(stdout); //Force stdout flush
		usleep(LoopWait); //Wait for half a second
	};
	
	return 0;
}