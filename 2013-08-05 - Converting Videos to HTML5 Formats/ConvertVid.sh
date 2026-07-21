#!/bin/bash
for filename in "$@"
do
  extension="${filename##*.}"
  filename="${filename%.*}"

  export BITRATE=`ffmpeg -i $filename.$extension 2>&1 | grep -oP 'Video.*\d+ kb/s' | grep -oP '\d+ kb/s' | grep -oP '\d+'`
  ffmpeg2theora -V $BITRATE -o $filename.ogv $filename.$extension
  ffmpeg -i $filename.$extension -ar 44100 -b ${BITRATE}k -f flv $filename.flv #Add an optional "-threads #" to make this faster
done