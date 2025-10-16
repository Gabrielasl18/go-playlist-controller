#!/bin/bash

# 1. Define o TZ para UTC e executa o FFmpeg
# O comando agora gera: #EXT-X-PROGRAM-DATE-TIME:YYYY-MM-DDTHH:MM:SS.sss+0000
TZ='UTC' ffmpeg -i sea_clip.mp4 \
    -filter_complex \
    "[0:v]split=3[v1][v2][v3]; [v1]scale=-2:1080[v1out]; [v2]scale=-2:720[v2out]; [v3]scale=-2:360[v3out]" \
    -map "[v1out]" -c:v:0 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:0 5M -maxrate:v:0 5M -minrate:v:0 5M -bufsize:v:0 10M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
    -map "[v2out]" -c:v:1 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:1 3M -maxrate:v:1 3M -minrate:v:1 3M -bufsize:v:1 6M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
    -map "[v3out]" -c:v:2 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:2 1M -maxrate:v:2 1M -minrate:v:2 1M -bufsize:v:2 2M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
    -map a:0 -c:a:0 aac -b:a:0 96k -ac 2 \
    -map a:0 -c:a:1 aac -b:a:1 96k -ac 2 \
    -map a:0 -c:a:2 aac -b:a:2 48k -ac 2 \
    -f hls \
    -hls_time 2 \
    -hls_flags independent_segments+program_date_time \
    -hls_segment_type mpegts \
    -hls_list_size 0 \
    -hls_segment_filename stream_%v/segment%02d.ts \
    -master_pl_name master.m3u8 \
    -var_stream_map "v:0,a:0,name:1080p v:1,a:1,name:720p v:2,a:2,name:360p" \
    playlist_%v.m3u8

