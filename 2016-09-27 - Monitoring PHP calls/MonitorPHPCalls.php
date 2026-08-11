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