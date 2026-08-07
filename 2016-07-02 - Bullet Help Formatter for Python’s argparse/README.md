<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Bullet_Help_Formatter_for_Python%E2%80%99s_argparse">Bullet Help Formatter for Python’s argparse</a><br>
Date: 2016-07-02<br>
Tags: Command Line Utility, Formatting, Python<br><br>
</div>

This is a <a href="https://docs.python.org/3/library/argparse.html#formatter-class">HelpFormatter</a> for Python’s <a href="https://docs.python.org/3/library/argparse.html">argparse</a> class which:
<ul>
	<li>Takes raw input and wraps long lines to indent against the current line start.</li>
	<li>When an indented/list line is encountered, which starts with spaces followed by a star "*", wrapped line’s indents will start 2 spaces after the star.</li>
	<li>Lines attempt to split at words of 10 characters or less (see .MinCharsInSplitWord).</li>
	<li>If a line needs to split along a word longer than this, a hyphen is inserted at the end of the line.</li>
</ul><br>

```python
import argparse
import re
class BulletHelpFormatter(argparse.HelpFormatter):
    def __init__(self, *args, **kwargs):
        super(BulletHelpFormatter, self).__init__(*args, **kwargs)
        self.MinCharsInSplitWord=10

    def _split_lines(self, text, width):
        #Split lines around line breaks and then modify each line
        Lines=[]
        for Line in text.splitlines():
            #Get the number of spaces to put at subsequent lines
            #0 if not a list item, oherwise, 2+list item start
            ListEl=re.match(r'^ *\*', Line)
            NumBeginningSpace=(0 if ListEl==None else ListEl.end()+1)

            #Add extra spaces at the beginning of each line to match the start of the current line, and go to a maxium of $width
            IsFirstPass=True
            SpacesToAdd=''
            NumSpacesToAdd=0
            while(True):
                #Get the word break points before and after where the line would end
                MaxLineLen=max(min(width-NumSpacesToAdd, len(Line)), 1)
                PrevWordBreak=CurWordBreak=0
                for WordBreak in re.finditer(r'(?<=\W).|\W|$', Line):
                    PrevWordBreak=CurWordBreak
                    CurWordBreak=WordBreak.start()
                    if CurWordBreak>=MaxLineLen:
                        if CurWordBreak==MaxLineLen:
                            PrevWordBreak=CurWordBreak
                        break

                #If previous wordbreak is more than MinCharsInSplitWord away from MaxLineLen, then split at the end of the line
                IsSplit=(PrevWordBreak<1 or CurWordBreak-PrevWordBreak>self.MinCharsInSplitWord)
                SplitPos=(MaxLineLen if IsSplit else PrevWordBreak)

                #Append the new line to the list of lines
                Lines.append(SpacesToAdd+Line[0:SplitPos]+('-' if IsSplit else ''))
                Line=Line[SplitPos:]

                #If this is the end, nothing left to do
                if len(Line)==0:
                    break

                #If this is the first pass, update line creation variables
                if IsFirstPass:
                    IsFirstPass=False
                    NumSpacesToAdd=NumBeginningSpace
                    SpacesToAdd=(' ' * NumSpacesToAdd)

        return Lines
```