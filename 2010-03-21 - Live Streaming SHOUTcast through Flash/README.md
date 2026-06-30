<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Live_Streaming_SHOUTcast_through_Flash">Live Streaming SHOUTcast through Flash</a><br>
Date: 2010-03-21<br>
Tags: Data Passing, Examples, Flash, PHP, Quick Fixes, SHOUTCast<br><br>

<p>A client of mine wanted their website to have an applet that played streaming music from a <a href="http://www.shoutcast.com/">SHOUTcast</a> server. The easy solution would have been to just embed a Windows Media Player applet into the page, but that would only work for IE.</p>
<p>I thoroughly searched the web and was unable to find a Flash applet (or other solution) that already did this (and actually worked). Most of the information I was finding was people having problems getting this kind of thing working in Flash with no answer provided. After giving up on finding a resolution online, I decided to load up Flash and see what I could find from some tinkering.</p>
<p>Quite frankly, I’m shocked people were having so many problems with this. I started an ActionScript 2.0 project and put in the following code, and it worked right away in Flash CS3 (v9.0) with no problem:</p>
</div>

```as
var URL="http://example.shoutcast.castledragmire.com:1234/" //The URL to the SHOUTcast server
var MySound:Sound=new Sound(this);
MySound.loadSound(URL,true);
```

<div>
<br>

<p>Unfortunately, once I exported the Flash applet and loaded it up in my browsers, it was no longer working. After a few minutes of poking around, I had a hunch that the SHOUTcast host might be sending different data depending on the [Browser’s] User Agent. I changed Firefox’s User Agent to “Flash” through a Firefox add-on (<a href="http://chrispederick.com/work/user-agent-switcher/">User Agent Switcher</a>), and it worked 😁.</p>
<p>Once again, unfortunately, this was not a viable solution because I couldn’t have every user who visited the client’s web page change their browser User Agent string 🙂. The quickest solution at this point to the problem was to just create a passthrough script that grabbed the live stream on their server and passed it to the client. The following is the PHP script I used for this:
</div>

```php
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
```

<div>
<br>

The final two steps to get this working were:
<ol>
	<li>Setting the Flash Applet’s <b>URL</b> variable to the PHP file</li>
	<li>Turning off PHP output buffering for the file. This can only be done through Apache or the php.ini depending on the server setup. This is very important, as if it’s on, the data will never get sent to the user.</li>
</ol><br>

<p>The only problem with this method is that it taxes the server that is passing the data through, especially since it uses PHP... This kind of thing could very easily be done in C though (as a matter of fact, I will be writing a post on something very close to that very soon).</p></div>