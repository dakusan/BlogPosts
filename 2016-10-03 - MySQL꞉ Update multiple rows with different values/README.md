<div>
Original Post: <a href="https://www.castledragmire.com/Posts/MySQL%3A_Update_multiple_rows_with_different_values">MySQL: Update multiple rows with different values</a><br>
Date: 2016-10-03<br>
Tags: Batch, Benchmarks, Examples, MySQL, Query, Update<br><br>

There are 3 different methods for updating multiple rows at once in MySQL with different values:
<ol>
	<li>
		<b>INSERT</b>: INSERT with ON DUPLICATE KEY UPDATE<br>
</li></ol>

```sql
INSERT INTO FooBar (ID, foo)
VALUES (1, 5), (2, 8), (3, 2)
ON DUPLICATE KEY UPDATE foo=VALUES(foo);
```

<ol start="2">
	<li>
		<b>TRANSACTION</b>: Where you do an update for each record within a transaction (InnoDB or other DBs with transactions)<br>
</li></ol>

```sql
START TRANSACTION;
UPDATE FooBar SET foo=5 WHERE ID=1;
UPDATE FooBar SET foo=8 WHERE ID=2;
UPDATE FooBar SET foo=2 WHERE ID=3;
COMMIT;
```

<ol start="3">
	<li>
		<b>CASE</b>: In which you a case/when for each different record within an UPDATE
</li></ol>

```sql
UPDATE FooBar SET foo=CASE ID
	WHEN 1 THEN 5
	WHEN 2 THEN 8
	WHEN 3 THEN 2
END
WHERE ID IN (1,2,3);
```

<p>I feel knowing the speeds of the 3 different methods is important.</p>
<p>All of the following numbers apply to InnoDB.</p><br>

<p>I just tested this, and the <b>INSERT</b> method was <kbd>6.7x</kbd> faster for me than the <b>TRANSACTION</b> method. I tried on a set of both 3,000 and 30,000 rows and got the same results.</p>
<p>The <b>TRANSACTION</b> method still has to run each individually query, which takes time, though it batches the results in memory, or something, while executing. The <b>TRANSACTION</b> method is also pretty expensive in both replication and query logs.</p>
<p>Even worse, the <b>CASE</b> method was <kbd>41.1x</kbd> slower than the <b>INSERT</b> method w/ 30,000 records (<kbd>6.1x</kbd> slower than <b>TRANSACTION</b>). And <kbd>75x</kbd> slower in MyISAM. <b>INSERT</b> and <b>CASE</b> methods broke even at ~1,000 records. Even at 100 records, the <b>CASE</b> method is BARELY faster.</p>
<p>So in general, I feel the <b>INSERT</b> method is both best and easiest to use. The queries are smaller and easier to read and only take up 1 query of action. This applies to both InnoDB and MyISAM.</p>

<ins>Bonus stuff:</ins>
<p>Using the <b>INSERT</b> method, there can be a problem in which NON-NULL fields with no default (in other words, required fields) are not being updated. You will get an error like <code>Field 'fieldname' doesn't have a default value</code>. The solution is to temporarily turn off STRICT_TRANS_TABLES and STRICT_ALL_TABLES in the SQL mode: <code>SET SESSION sql_mode=REPLACE(REPLACE(@@SESSION.sql_mode,"STRICT_TRANS_TABLES",""),"STRICT_ALL_TABLES","")</code>. Make sure to save the sql_mode first if you plan on reverting it.</p><br>

<p>As for other comments I’ve seen that say the auto_increment goes up using the <b>INSERT</b> method, I tested that too and it seems to not be the case.</p><br>

<ins>Code to run the tests is as follows:</ins> (It also outputs .SQL files to remove PHP interpreter overhead)<br>

```php
<?php
//Variables
$NumRows=30000;

//These 2 functions need to be filled in
function InitSQL()
{

}
function RunSQLQuery($Q)
{

}

//Run the 3 tests
InitSQL();
for($i=0;$i<3;$i++)
    RunTest($i, $NumRows);

function RunTest($TestNum, $NumRows)
{
    $TheQueries=Array();
    $DoQuery=function($Query) use (&$TheQueries)
    {
        RunSQLQuery($Query);
        $TheQueries[]=$Query;
    };

    $TableName='Test';
    $DoQuery('DROP TABLE IF EXISTS '.$TableName);
    $DoQuery('CREATE TABLE '.$TableName.' (i1 int NOT NULL AUTO_INCREMENT, i2 int NOT NULL, primary key (i1)) ENGINE=InnoDB');
    $DoQuery('INSERT INTO '.$TableName.' (i2) VALUES ('.implode('), (', range(2, $NumRows+1)).')');

    if($TestNum==0)
    {
        $TestName='Transaction';
        $Start=microtime(true);
        $DoQuery('START TRANSACTION');
        for($i=1;$i<=$NumRows;$i++)
            $DoQuery('UPDATE '.$TableName.' SET i2='.(($i+5)*1000).' WHERE i1='.$i);
        $DoQuery('COMMIT');
    }

    if($TestNum==1)
    {
        $TestName='Insert';
        $Query=Array();
        for($i=1;$i<=$NumRows;$i++)
            $Query[]=sprintf("(%d,%d)", $i, (($i+5)*1000));
        $Start=microtime(true);
        $DoQuery('INSERT INTO '.$TableName.' VALUES '.implode(', ', $Query).' ON DUPLICATE KEY UPDATE i2=VALUES(i2)');
    }

    if($TestNum==2)
    {
        $TestName='Case';
        $Query=Array();
        for($i=1;$i<=$NumRows;$i++)
            $Query[]=sprintf('WHEN %d THEN %d', $i, (($i+5)*1000));
        $Start=microtime(true);
        $DoQuery("UPDATE $TableName SET i2=CASE i1\n".implode("\n", $Query)."\nEND\nWHERE i1 IN (".implode(',', range(1, $NumRows)).')');
    }

    print "$TestName: ".(microtime(true)-$Start)."<br>\n";

    file_put_contents("./$TestName.sql", implode(";\n", $TheQueries).';');
}
```