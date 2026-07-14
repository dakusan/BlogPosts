<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Converting_Videos_to_HTML5_Formats">Converting Videos to HTML5 Formats</a><br>
Date: 2013-08-05<br>
Tags: Convert, Examples, HTML5<br><br>

<p>The following is a Linux script I threw together to convert an mpeg4 (.mp4, .mpv, .mpeg4, .mpg) file into ogg vorbis (.ogg, .ogv) and flash video (.flv) maintaining the same bit rate. It also doesn’t hurt to have a .webm format for some older devices. This assumes you already have the appropriate packages installed, including ffmpeg. The script also wasn’t made to be foolproof, and might not be able to read the bitrate on some videos. The script takes a list of mpeg4 files to convert.</p>

</div>

```bash
for filename in "$@"
do
  extension="${filename##*.}"
  filename="${filename%.*}"

  export BITRATE=`ffmpeg -i $filename.$extension 2>&1 | grep -oP 'Video.*\d+ kb/s' | grep -oP '\d+ kb/s' | grep -oP '\d+'`
  ffmpeg2theora -V $BITRATE -o $filename.ogv $filename.$extension
  ffmpeg -i $filename.$extension -ar 44100 -b ${BITRATE}k -f flv $filename.flv #Add an optional "-threads #" to make this faster
done
```

<div>
</div>