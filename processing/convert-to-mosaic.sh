#!/bin/bash

RANDOM=$(cat /dev/urandom | LC_CTYPE=C tr -cd '0-9' | head -c 5)
FN="images/GROVER_MOSIAC_$RANDOM.png"
ffmpeg $1 -filter_complex hstack=inputs=$2 "$3/$FN"
echo $FN