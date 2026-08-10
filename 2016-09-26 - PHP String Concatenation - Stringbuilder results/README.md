<div>
Original Post: <a href="https://www.castledragmire.com/Posts/PHP_String_Concatenation_-_Stringbuilder_results">PHP String Concatenation - Stringbuilder results</a><br>
Date: 2016-09-26<br>
Tags: Benchmarks, Examples, Memory Footprint, Optimization, PHP, StringBuilder, Strings<br><br>

<p>I wrote the code at the end of this post to test the different forms of string concatenation and they really are all almost exactly equal in both memory and time footprints.</p><br>

<p>The two primary methods I used are concatenating strings onto each other, and filling an array with strings and then imploding them. I did 500 string additions with a 1MB string in PHP 5.6 (so the result is a 500MB string). At every iteration of the test, all memory and time footprints were very very close (at ~$IterationNumber*1MB). The runtime of both tests was 50.398 seconds and 50.843 seconds consecutively which are most likely within acceptable margins of error.</p>

<p>Garbage collection of strings that are no longer referenced seems to be pretty immediate, even without ever leaving the scope. Since the strings are mutable, no extra memory is really required after the fact.</p><br>

<p><b>HOWEVER</b>, The following tests showed that there is a different in peak memory usage <i>WHILE</i> the strings are being concatenated.</p>

```php
$OneMB=str_repeat('x', 1024*1024);
$Final=$OneMB.$OneMB.$OneMB.$OneMB.$OneMB;
print memory_get_peak_usage();
```
Result=10,806,800 bytes (~10MB w/o the initial PHP memory footprint)<br><br>

```php
$OneMB=str_repeat('x', 1024*1024);
$Final=implode('', Array($OneMB, $OneMB, $OneMB, $OneMB, $OneMB));
print memory_get_peak_usage();
```
Result=6,613,320 bytes (~6MB w/o the initial PHP memory footprint)<br><br>

<p>So there is in fact a difference that could be significant in very very large string concatenations memory-wise (I have run into such examples when creating very large data sets or SQL queries).</p>

<p>But even this fact is disputable depending upon the data. For example, concatenating 1 character onto a string to get 50 million bytes (so 50 million iterations) took a maximum amount of 50,322,512 bytes (~48MB) in 5.97 seconds. While doing the array method ended up using 7,337,107,176 bytes (~6.8GB) to create the array in 12.1 seconds, and then took an extra 4.32 seconds to combine the strings from the array.</p><br>

<p>Anywho... the below is the benchmark code I mentioned at the beginning which shows the methods are pretty much equal. It outputs a pretty HTML table.</p>

```php
<?php
//Please note, for the recursion test to go beyond 256, xdebug.max_nesting_level needs to be raised.
//You also may need to update your memory_limit depending on the number of iterations

//Output the start memory
print 'Start: '.memory_get_usage()."B<br><br>Below test results are in MB<br>";

//Our 1MB string
global $OneMB, $NumIterations;
$OneMB=str_repeat('x', 1024*1024);
$NumIterations=500;

//Run the tests
$ConcatTest=RunTest('ConcatTest');
$ImplodeTest=RunTest('ImplodeTest');
$RecurseTest=RunTest('RecurseTest');

//Output the results in a table
OutputResults(
  Array('ConcatTest', 'ImplodeTest', 'RecurseTest'),
  Array($ConcatTest, $ImplodeTest, $RecurseTest)
);

//Start a test run by initializing the array that will hold the results and manipulating those results after the test is complete
function RunTest($TestName)
{
  $CurrentTestNums=Array();
  $TestStartMem=memory_get_usage();
  $StartTime=microtime(true);
  RunTestReal($TestName, $CurrentTestNums, $StrLen);
  $CurrentTestNums[]=memory_get_usage();

  //Subtract $TestStartMem from all other numbers
  foreach($CurrentTestNums as &$Num)
    $Num-=$TestStartMem;
  unset($Num);

  $CurrentTestNums[]=$StrLen;
  $CurrentTestNums[]=microtime(true)-$StartTime;

  return $CurrentTestNums;
}

//Initialize the test and store the memory allocated at the end of the test, with the result
function RunTestReal($TestName, &$CurrentTestNums, &$StrLen)
{
  $R=$TestName($CurrentTestNums);
  $CurrentTestNums[]=memory_get_usage();
  $StrLen=strlen($R);
}

//Concatenate 1MB string over and over onto a single string
function ConcatTest(&$CurrentTestNums)
{
  global $OneMB, $NumIterations;
  $Result='';
  for($i=0;$i<$NumIterations;$i++)
  {
    $Result.=$OneMB;
    $CurrentTestNums[]=memory_get_usage();
  }
  return $Result;
}

//Create an array of 1MB strings and then join w/ an implode
function ImplodeTest(&$CurrentTestNums)
{
  global $OneMB, $NumIterations;
  $Result=Array();
  for($i=0;$i<$NumIterations;$i++)
  {
    $Result[]=$OneMB;
    $CurrentTestNums[]=memory_get_usage();
  }
  return implode('', $Result);
}

//Recursively add strings onto each other
function RecurseTest(&$CurrentTestNums, $TestNum=0)
{
  Global $OneMB, $NumIterations;
  if($TestNum==$NumIterations)
    return '';

  $NewStr=RecurseTest($CurrentTestNums, $TestNum+1).$OneMB;
  $CurrentTestNums[]=memory_get_usage();
  return $NewStr;
}

//Output the results in a table
function OutputResults($TestNames, $TestResults)
{
  global $NumIterations;
  print '<table border=1 cellspacing=0 cellpadding=2><tr><th>Test Name</th><th>'.implode('</th><th>', $TestNames).'</th></tr>';
  $FinalNames=Array('Final Result', 'Clean');
  for($i=0;$i<$NumIterations+2;$i++)
  {
    $TestName=($i<$NumIterations ? $i : $FinalNames[$i-$NumIterations]);
    print "<tr><th>$TestName</th>";
    foreach($TestResults as $TR)
      printf('<td>%07.4f</td>', $TR[$i]/1024/1024);
    print '</tr>';
  }

  //Other result numbers
  print '<tr><th>Final String Size</th>';
  foreach($TestResults as $TR)
    printf('<td>%d</td>', $TR[$NumIterations+2]);
  print '</tr><tr><th>Runtime</th>';
    foreach($TestResults as $TR)
      printf('<td>%s</td>', $TR[$NumIterations+3]);
  print '</tr></table>';
}
?>
```