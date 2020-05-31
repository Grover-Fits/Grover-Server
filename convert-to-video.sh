#!/bin/bash

RANDOM=$(cat /dev/urandom | LC_CTYPE=C tr -cd '0-9' | head -c 5)
FN="videos/fitsMOV_$RANDOM.mkv"
cat $1 | ffmpeg -framerate 1 -probesize 42M -f image2pipe -i - "$2/$FN"
printf "$FN"