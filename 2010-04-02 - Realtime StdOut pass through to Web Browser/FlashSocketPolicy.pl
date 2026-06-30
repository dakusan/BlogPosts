#!/usr/bin/perl
use warnings;
use strict;

#Listen for kill signals
$SIG{'QUIT'}=$SIG{'INT'}=$SIG{__DIE__} = sub
{
	close Server;
	print "Socket Policy Server Ended: $_[0]\n";
	exit;
};

#Start the server:
use Socket;
use IO::Handle;
my $FlashPolicyPort=843;
socket(Server, PF_INET, SOCK_STREAM, getprotobyname('tcp')) or die "'socket' call: $!"; #Open the socket
setsockopt(Server, SOL_SOCKET, SO_REUSEADDR, 1) or die "'setsockopt' call: $!"; #Allow reusing of port/address if in TIME_WAIT state
bind(Server, sockaddr_in($FlashPolicyPort,INADDR_ANY)) or die "'bind' call: $!"; #Listen on port $FlashPolicyPort for connections from any INET adapter
listen(Server,SOMAXCONN) or die "'listen' call: $!"; #Start listening for connections
Server->autoflush(1); #Do not buffer output

#Infinite loop that accepts connections
$/ = "\0"; #Reset terminator from new line to null char
while(my $paddr=accept(Client,Server))
{
	Client->autoflush(1); #Do not buffer IO
	if(<Client> =~ /.*policy\-file.*/i) { #If client requests policy file...
		print Client '<cross-domain-policy><allow-access-from domain="*" to-ports="*" /></cross-domain-policy>'.$/; #Output policy info: Allow any flash applets from any domain to connect
	}
	close Client; #Close the client
}