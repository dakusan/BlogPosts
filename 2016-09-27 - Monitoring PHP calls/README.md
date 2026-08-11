<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Monitoring_PHP_calls">Monitoring PHP calls</a><br>
Date: 2016-09-27<br>
Tags: Diagnosing, Examples, PHP, Programming hacks, runkit<br><br>
</div>

<p>I recently had a Linux client that was, for whatever odd reason, making infinite recursive HTTP calls to a single script, which was making the server process count skyrocket. I decided to use the same module as I did in my <a href="../2016-07-02 - Painless migration from PHP MySQL to MySQLi">Painless migration from PHP MySQL to MySQLi</a> post, which is to say, overriding base functions for fun and profit using the <a href="http://php.net/manual/en/book.runkit.php">PHP runkit</a> extension. I did this so I could gather, for debugging, logs of when and where the calls that were causing this to occur.</p><br>
<p>The below code overrides all functions listed on the line that says “List of functions to intercept” <code>[Line 9]</code>. It works by first renaming these built in functions to “OVERRIDE_<b>$FuncName</b>” <code>[Line 12]</code>, and replacing them with a call to “<b>GlobalRunFunc</b>()” <code>[Line 13]</code>, which receives the original function name and argument list. The <b>GlobalRunFunc()</b>:
<ol>
	<li>Checks to see if it is interested in logging the call<ul>
		<li>In the case of this example, it will log the call if <code>[Line 20]</code>:<ul>
			<li>Line 21: <b>curl_setopt</b> is called with the CURLOPT_URL parameter (enum=10002)</li>
			<li>Line 22: <b>curl_init</b> is called with a first parameter, which would be a URL</li>
			<li>Line 23: <b>file_get_contents</b> or <b>fopen</b> is called and is not an absolute path<br>(Wordpress calls everything absolutely. Normally I would have only checked for http[s] calls).</li>
		</ul></li>
		<li>If it does want to log the call, it stores it in a global array (which holds all the calls we will want to log).<br>The logged data includes <code>[Line 25]</code>:<ul>
			<li>The function name</li>
			<li>The function parameters</li>
			<li>2 functions back of backtrace (which can often get quite large when stored in the log file)</li>
		</ul></li>
	</ul></li>
	<li>It then calls the original function, with parameters intact, and passes through the return <code>[Line 27]</code>.</li>
</ol><br>

<p>The “<b>GlobalShutdown</b>()” <code>[Line 30]</code> is then called when the script is closing <code>[Line 38]</code> and saves all the logs, if any exist, to “$<b>GlobalLogDir</b>/$<b>DATETIME</b>.srl”.</p>
<p>I have it using “serialize()” to encode the log data <code>[Line 25]</code>, as opposed to “json_encode()” or “print_r()” calls, as the latter were getting too large for the logs. You may want to have it use one of these other encoding functions for easier log perusal, if running out of space is not a concern.</p>

```php
<?
//The log data to save is stored here
global $GlobalLogArr, $GlobalLogDir;
$GlobalLogArr=Array();
$GlobalLogDir='./LOG_DIRECTORY_NAME';

//Override the functions here to instead have them call to GlobalRunFunc, which will in turn call the original functions
foreach(Array(
        'fopen', 'file_get_contents', 'curl_init', 'curl_setopt', //List of functions to intercept
) as $FuncName)
{
        runkit_function_rename($FuncName, "OVERRIDE_$FuncName");
        runkit_function_add($FuncName, '', "return GlobalRunFunc('$FuncName', func_get_args());");
}

//This optionally 
function GlobalRunFunc($FuncName, $Args)
{
        global $GlobalLogArr;
        if(
                ($FuncName=='curl_setopt' && $Args[1]==10002) || //CURLOPT enumeration can be found at https://curl.haxx.se/mail/archive-2004-07/0100.html
                ($FuncName=='curl_init' && isset($Args[0])) ||
                (($FuncName=='file_get_contents' || $FuncName=='fopen') && $Args[0][0]!='/')
        )
                $GlobalLogArr[]=serialize(Array('FuncName'=>$FuncName, 'Args'=>$Args, 'Trace'=>array_slice(debug_backtrace(), 1, 2)));

        return call_user_func_array("OVERRIDE_$FuncName", $Args);
}

function GlobalShutdown()
{
        global $GlobalLogArr, $GlobalLogDir;
        $Time=microtime(true);
        if(count($GlobalLogArr))
                file_put_contents($GlobalLogDir.date('Y-m-d_H:i:s.'.substr($Time-floor($Time), 2, 3), floor($Time)).'.srl', implode("\n", $GlobalLogArr));

}
register_shutdown_function('GlobalShutdown');
?>
```