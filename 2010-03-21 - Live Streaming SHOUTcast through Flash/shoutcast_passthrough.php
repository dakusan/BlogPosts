<?php
$streamname='example.shoutcast.castledragmire.com';
$port      ='1234';
$path      ='/';

header('Content-type: audio/mpeg');
$sock=fsockopen($streamname,$port);
fputs($sock, "GET $path HTTP/1.0\r\n");
fputs($sock, "Host: $streamname\r\n");
fputs($sock, "User-Agent: WinampMPEG/2.8\r\n");
fputs($sock, "Accept: */*\r\n");
fputs($sock, "Connection: close\r\n\r\n");
fpassthru($sock);
fclose($sock);