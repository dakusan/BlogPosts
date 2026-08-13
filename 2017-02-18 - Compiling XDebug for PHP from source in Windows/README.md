<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Compiling_XDebug_for_PHP_from_source_in_Windows">Compiling XDebug for PHP from source in Windows</a><br>
Date: 2017-02-18<br>
Tags: Compile, PHP, Tutorial, Windows, XDebug<br><br>

<p>I decided to try out the PHP IDE <a href="https://www.jetbrains.com/phpstorm/specials/phpstorm/phpstorm.html">PHPStorm</a> by <a href="https://www.jetbrains.com">jet brains</a> (you know, the people that made android studio) again recently but ran into a few problems. I tried it a while back and I abandoned it because I was using the browser connect method instead of xdebug and it was a major PITA. This time around, I tried xdebug and everything works much smoother and easier. However, there is a <a href="https://bugs.xdebug.org/view.php?id=842">major</a> <a href="https://bugs.xdebug.org/view.php?id=895">bug</a> with it in which single line statements inside control statements could not be stopped on in the debugger (breakpoint or stepping). This has been a problem since 2012 and was marked until very recently as unfixable. I am honestly shocked this is not a thing more cited. Do people just not use that form of formatting anymore? I find it much cleaner. Example:</p>

```php
if(true)
  $a=1; //This line could never be reached
```

So all code like this had to be changed to the following for debugging.
```php
if(true) {
  $a=1;
}
```

<p>However, in a <a href="https://bugs.xdebug.org/view.php?id=842#c3992">comment</a> by one of the developers on 2016-12-11 (just 2 months ago) he said “<span class="code short">This turned out to be a duplicate of 1165, for which I now have cooked up a fix, which will make it into 2.5.1.</span>”. Unfortunately, there has not yet been a release of this version, so I was forced to a compile a copy myself on Windows. This turned out to also be a major PITA. I have documented the process here. <a href="https://www.castledragmire.com/OtherContent/php_xdebug.dll">Here is the version I compiled (7.1.2 threadsafe x86) </a></p>
<ul>
	<li>To compile a php extension, you also have to compile php itself. So first <a href="http://windows.php.net/download">download</a> both the source and zip files for your chosen version. I used the latest version, <a href="http://windows.php.net/downloads/releases/php-7.1.2-Win32-VC14-x86.zip">7.1.2 VC14 x86 Thread Safe (2017-Feb-14 23:28:41)</a> [<a href="http://windows.php.net/downloads/releases/php-7.1.2-src.zip">7.1.2 source</a>]. The source is needed for obvious reasons. The binary is needed so we can extract some dlls not included with the source. Extract the php source root directory files to “<span class="code short">c:\php\src.</span>”</li>
	<li>You also need to grab the <a href="https://github.com/xdebug/xdebug">xdebug source [github]</a> via “<span class="code short"><span class="code short">git clone git://github.com/xdebug/xdebug.git</span></span>”. I used commit <a href="https://github.com/xdebug/xdebug/commit/62ac284bf36f7521e78d2b21492ce2899f7cc4ff">#62ac284bf36f7521e78d2b21492ce2899f7cc4ff</a> <strike>#<a href="https://github.com/xdebug/xdebug/commit/379ab7ab68d28d0af8f637aa28c7e9fab921f27a">379ab7ab68d28d0af8f637aa28c7e9fab921f27a</a>, but it actually had a bug in it which I fixed in <a href="https://github.com/dakusan/xdebug/commit/9fa7250f6f315d6e4558445fc4cd223266836081">my fork</a>. I sent a pull request so it should hopefully be integrated soon.</strike> Clone the git repo in “<span class="code short">c:\php\src-extensions</span>” (it will create the “<span class="code short">xdebug</span>” directory inside of it)</li>
	<li>You will also need <a href="https://www.visualstudio.com/downloads/">visual studio</a> <a href="https://www.visualstudio.com/post-download-vs/?sku=community&clcid=0x409&telem=ga#">14.0</a> (VS 2015) [<a href="http://download.microsoft.com/download/D/2/3/D23F4D0F-BA2D-4600-8725-6CCECEA05196/vs_community_ENU.exe">direct download</a>] installed.</li>
	<li>Once visual studio is installed, you need to open “<span class="code short">Developer Command Prompt for VS2015</span>”. It can be found in your start menu under “<span class="code short">Programs\Visual Studio 2015\Visual Studio Tools</span>”.</li>
	<li>Run the following commands to get php compiled:
<pre><code>cd\php\src
buildconf --add-modules-dir=..\src-extensions
:: You may want some different flags for the below command
configure "--enable-snapshot-build" "--enable-debug-pack" "--without-pdo-mssql" "--enable-com-dotnet=shared" "--with-mcrypt=static" "--without-analyzer"
:: At this point we are going to actually build php to get the phpize.bat
:: You may receive some warnings at the end that say something to the extent of
:: “'-9' is not recognized as an internal or external command”. You can ignore these
nmake snap
:: This next command will “install” php to “c:\php”.
nmake install</code></pre>
	</li>
	<li>Next we get to compile xdebug. But it requirements <a href="http://gnuwin32.sourceforge.net/packages/bison.htm">bison, which can be downloaded here</a> [<a href="http://downloads.sourceforge.net/gnuwin32/bison-2.4.1-bin.zip">Direct download</a>]. Extract bin/bison.exe and bin/m4.exe to the VC/bin directory, most likely at “<span class="code short">C:\Program Files (x86)\Microsoft Visual Studio 14.0\VC\bin</span>”.<br><br>
	For some reason, the first time I did this process I needed the following files also in VC/bin directory, all found at <a href="http://gnuwin32.sourceforge.net">http://gnuwin32.sourceforge.net</a>. However, during my second run to create this tutorial, I did not need them. iconv.exe, libcharset1.dll, libiconv2.dll, libintl3.dll, mspdb110.dll, regex2.dll. And the “<span class="code short">share</span>” folder from the bison zip file needed to be extracted to “<span class="code short">C:\Program Files (x86)</span>”.
	</li>
	<li>Some file edits in the xdebug directory:<ul>
		<li>Add the following line to the top of configure.js: “<span class="code short">PHP_SECURITY_FLAGS=true;</span>”</li>
		<li>If you are using the commit I mentioned, Add the following to line 471 of xdebug_code_coverage.c “<span class="code short">zend_op *base_address = &(opa-&gt;opcodes[0]);</span>”</li>
		<li>Also edit line 9 of template.rc by changing it from “<span class="code short">#include "main/php_version.h"</span>” to “<span class="code short">#include "../../src/main/php_version.h"</span>”</li>
	</ul></li>
	<li>And now to actually compile xdebug
	<pre class="code lang-bash">
cd ..\src-extensions\xdebug
..\..\SDK\phpize
configure --with-xdebug
nmake php_xdebug.dll
copy Release_TS\php_xdebug.dll ..\..\ext</pre>
	</li>
	<li>Your complete php install is now located in “<span class="code short">c:\php</span>”. But exclude the “<span class="code short">src</span>” and “<span class="code short">src-extensions</span>” folder. It still needs a few dlls that are found in the php zip file you downloaded earlier. Extract all the dlls from that root of that file to your root php folder, but do not overwrite anything that is already there. If you wanted to compile your own version of the apache dll, you can add one of “<span class="code short">--enable-apache2handler, --enable-apache2-2handler, --enable-apache2-4handler</span>” to the src/configure command.</li>
	<li>For some really odd reason, my version of php would always give me the following error “<span class="code short">Packets out of order. Expected 3 received 2. Packet size=22</span>” when trying to connect to mysql (both pdo and mysqli go through mysqlnd). So I just ended up using the php files provided through the zip file, with my newly compiled xdebug dll.</li>
	<li>Definitely have to mention that the following blog post was a big help to me in figuring this all out: <a href="http://blog.benoitblanchon.fr/build-php-extension-on-windows/">http://blog.benoitblanchon.fr/build-php-extension-on-windows/</a></li>
</ul></div>