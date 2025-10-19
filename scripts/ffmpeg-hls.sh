#!/bin/bash

INPUT_FILE="../sea_clip.mp4"
OUTPUT_DIR="../hls"
MASTER_PLAYLIST_NAME="master.m3u8"

if [ ! -f "$INPUT_FILE" ]; then
    echo "ERRO: O arquivo de entrada não foi encontrado: ${INPUT_FILE}"
    exit 1
fi

mkdir -p "$OUTPUT_DIR"

echo -e "\n\n-------------------------------------------------------------------"
echo "Streaming (Live Loop) de ${INPUT_FILE} para HLS MBR (MPEG-TS)."
echo "Master Playlist: ${OUTPUT_DIR}/${MASTER_PLAYLIST_NAME}"
echo "Usando loop infinito e Segmentação por Time (2s)."
echo "Arquivos de Saída: ${OUTPUT_DIR}"
echo -e "-------------------------------------------------------------------\n\n"

# 3. Comando FFmpeg formatado para melhor leitura.
# Cada linha termina com '\' e SEM ESPAÇOS OU COMENTÁRIOS depois.

TZ='UTC' ffmpeg \
    -stream_loop -1 -re \
    -i "$INPUT_FILE" \
    -fflags +nobuffer \
    -filter_complex "[0:v]split=3[v1][v2][v3]; [v1]scale=-2:1080[v1out]; [v2]scale=-2:720[v2out]; [v3]scale=-2:360[v3out]" \
    -map "[v1out]" -c:v:0 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:0 5M -maxrate:v:0 5M -minrate:v:0 5M -bufsize:v:0 10M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
    -map a:0 -c:a:0 aac -b:a:0 96k -ac 2 \
    -map "[v2out]" -c:v:1 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:1 3M -maxrate:v:1 3M -minrate:v:1 3M -bufsize:v:1 6M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
    -map a:0 -c:a:1 aac -b:a:1 96k -ac 2 \
    -map "[v3out]" -c:v:2 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:2 1M -maxrate:v:2 1M -minrate:v:2 1M -bufsize:v:2 2M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
    -map a:0 -c:a:2 aac -b:a:2 48k -ac 2 \
    -f hls \
    -hls_time 2 \
    -hls_flags independent_segments+program_date_time \
    -hls_segment_type mpegts \
    -hls_list_size 3 \
    -hls_segment_filename "${OUTPUT_DIR}/stream_%v/segment%02d.ts" \
    -master_pl_name "${MASTER_PLAYLIST_NAME}" \
    -var_stream_map "v:0,a:0,name:1080p v:1,a:1,name:720p v:2,a:2,name:360p" \
    "${OUTPUT_DIR}/playlist_%v.m3u8"

echo -e "\n\n-------------------------------------------------------------------"
echo "Streaming (Live Loop) encerrado."
echo "Os arquivos estão localizados em: ${OUTPUT_DIR}"
echo -e "-------------------------------------------------------------------\n\n"