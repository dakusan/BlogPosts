<div>
Original Post: <a href="https://www.castledragmire.com/Posts/SymLink_Fix_for_Combining_Android_Project_Versions">SymLink Fix for Combining Android Project Versions</a><br>
Date: 2011-01-03<br>
Tags: Examples, Symlinks<br><br>

<p>Since I found out that <a href="../2010-12-04 - Symlinks in a Windows programming environment">NTFS now has semi-working native symlinks</a>, I have updated the symlinking script used in the <a href="https://www.castledragmire.com/Posts/Combining_an_Android_Project's_Versions">Combining an Android Project's Versions</a> post. This script creates relative symlinks now through Perl instead of absolute hard links through Bash. It is as follows:</p>
</div>

```perl
#!/usr/bin/perl
#Run this file to install links to shared files into all branches
use warnings;
use strict;

#Configuration
my $SharedDirectoryName="Shared";
my $NonProjectDirectories="^\\.(|/\\.git|/$SharedDirectoryName)\$"; #Non Project directories (., .git, $SharedDirectoryName)
my $IsWindows=(index(lc(`uname`), 'cygwin')!=-1);

#Create a symlink
sub MakeLink
{
	my ($LinkTarget, $LinkName, $IsWindows, $IsDirectory)=@_;

	#Create the target directory if it does not exist
	my $LinkDirectory=$LinkName;
	$LinkDirectory =~ s/\/[^\/]+$//;
	if(!-e $LinkDirectory) {
		print "Creating directory: $LinkDirectory\n";
		`mkdir -p "$LinkDirectory"`;
	}
	
	#If the link already exists, issue a warning
	if(-l $LinkName) {
		print "Link already exists: $LinkName\n";
		return;
	}

	#Create the relative symlink
	my $RelativePathFromLinkToTarget=('../' x ($LinkName =~ tr/\///)).$LinkTarget; #Determine the relative path between the link and its target
	my $Command;
	if(!$IsWindows) { #Create the Linux command
		$Command="ln -s \"$RelativePathFromLinkToTarget\" \"$LinkName\"";
	}
	else #Create the Windows command
	{
		#Replace /s in path with \s
		$RelativePathFromLinkToTarget =~ s/\//\\/g;
		$LinkName =~ s/\//\\/g;
		
		$Command='cmd /c mklink'.($IsDirectory ? ' /d' : '')." \"$LinkName\" \"$RelativePathFromLinkToTarget\"";
	}

	print "$Command\n";
	`$Command`;
}

#Find required information from file searches
my @LocalBranches=grep(!/$NonProjectDirectories/, `find -maxdepth 1 -type d`); #Find version folders by ignoring Non Project directories
my @Files=split(/\n?^$SharedDirectoryName\//m, substr(`find $SharedDirectoryName -type f`, 0, -1)); shift @Files; #Find shared files

#Propagate shared files into different versions
foreach my $LocalBranch (@LocalBranches) {
	$LocalBranch=substr($LocalBranch, 2, -1); #Remove ./ and new line separator
	foreach my $File (@Files) {
		MakeLink("$SharedDirectoryName/$File", "$LocalBranch/$File", $IsWindows, 0);
	}
}
```