<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Directory_Difference_in_Web_Browser">Directory Difference in Web Browser</a><br>
Date: 2009-08-23<br>
SubTitle: Another quick and dirty solution<br>
Tags: Diff, Examples, FileSync<br><br>

<p>Since my <a href="https://www.castledragmire.com/Projects/FileSync">FileSync Project</a> is still a long way to being where I want it to be and is in a state that makes it annoying to use, I decided to throw together a script that essentially emulates the primary functions of it for what I need. I find it a bit annoying that I’ve never been able to find another good project that does exactly what I want for quick syncing of files over networks 😕. I used to think rsync would be a good solution for it but it’s very... quirky and unstable in certain ways. Not as flexible as I would like. Alas.</p>
<p>Anywho, <a href="Compare.php">this PHP script</a> takes 2 file lists and gives you back their differences in a directory tree view. Each file has a data string after it used to tell if the files are different. The data string can be anything you want, but will usually probably be a timestamp or data checksum. The tree view lets you hide files/directories depending on whether each item is only on one side, different, or the same. An example is as follows (<a href="Dir1.txt">List1</a>, <a href="Dir2.txt">List2</a>):</p>
<table align=center><tr><td>
	<table width=300>
		<tr style="border-bottom:1px solid black;"><th align=left>Dir1.txt files</th><th>Data String (Timestamps)</th></tr>
		<tr><td>Same.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>Diff.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>LeftOnly.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>PartialDir/Same.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>PartialDir/Diff.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>PartialDir/LeftOnly.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>SameDir/Same.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>LeftOnlyDir/Left.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>LeftOnlyDir/Left2.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>DiffDir/Diff.txt</td><td>9999-99-99 99:99:99</td></tr>
	</table>
</td><td valign=top>
	<table width=300>
		<tr style="border-bottom:1px solid black;"><th align=left>Dir2.txt files</th><th>Data String (Timestamps)</th></tr>
		<tr><td>Same.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>Diff.txt</td><td>1111-11-11 11:11:11</td></tr>
		<tr><td>RightOnly.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>PartialDir/Same.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>PartialDir/Diff.txt</td><td>1111-11-11 11:11:11</td></tr>
		<tr><td>PartialDir/RightOnly.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>SameDir/Same.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>RightOnlyDir/Right.txt</td><td>9999-99-99 99:99:99</td></tr>
		<tr><td>&nbsp;</td><td>&nbsp;</td></tr>
		<tr><td>DiffDir/Diff.txt</td><td>1111-11-11 11:11:11</td></tr>
	</table>
</td></tr></table><br><br>

<b><a href="ExampleOutput.html">Example Output</a>, <a href="https://www.castledragmire.com/OtherContent/Compare/ExampleOutput.html">Example Output Live</a></b><br>
<img src=./Diff.png>

<br><b>Some example <a href="https://www.castledragmire.com/Posts/Useful_Bash_commands_and_scripts">bash commands</a> used to create file lists:</b><br>
```bash
#Output all files and their timestamps to “Dir1.txt”
find -type f -printf '%P\t%T+\n' > Dir1.txt

#Output all files and their md5sums to “Dir1.txt”:
find -type f -print0 | xargs -0 md5sum | perl -pe 's/^(.*?) (.*)$/$2\t$1/g' > Dir1.txt</li>
```