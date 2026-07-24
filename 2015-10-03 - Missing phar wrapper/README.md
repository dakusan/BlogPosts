<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Missing_phar_wrapper">Missing phar wrapper</a><br>
Date: 2015-10-03<br>
Tags: Cygwin, Examples, phar, PHP, PHP wrapper, Quick Fixes<br><br>

<p><a href="https://secure.php.net/manual/en/intro.phar.php">Phar files</a> are PHP’s way of distributing an entire PHP solution in a single package file. I recently had a problem on my Cygwin PHP server that said “Unable to find the wrapper "phar" - did you forget to enable it when you configured PHP?”. I couldn’t find any solution for this online, so I played with it a bit.</p>
<p>The quick and dirty solution I came up with is to include the phar file like any normal PHP file, which sets your current working directory inside of the phar file. After that, you can include files inside the phar and then change your directory back to where you started. Here is the code I used:</p>

</div>

```php
if(preg_match('/^(?:win|cygwin)/i', PHP_OS))
{
	$CWD=getcwd();
	require_once('Scripts/PHPExcel.phar');
	require_once('PHPExcel/IOFactory.php');
	chdir($CWD);
}
else
	require_once('phar://Scripts/PHPExcel.phar/PHPExcel/IOFactory.php');
```