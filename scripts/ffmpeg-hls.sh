#!/bin/bash

INPUT_FILE="../sea_clip.mp4"
OUTPUT_DIR="../hls"
MASTER_PLAYLIST_NAME="master.m3u8"
CHILD_PLAYLIST_NAME="playlist_1080p.m3u8"

# 1. Verifica se o arquivo existe
if [ ! -f "$INPUT_FILE" ]; then
    echo "ERRO: O arquivo de entrada não foi encontrado: ${INPUT_FILE}"
    exit 1
fi

# 2. Cria pasta de saída
mkdir -p "$OUTPUT_DIR"

echo -e "\n\n-------------------------------------------------------------------"
echo "Streaming (Live Loop) de ${INPUT_FILE} para HLS (1 resolução)."
echo "Master Playlist: ${OUTPUT_DIR}/${MASTER_PLAYLIST_NAME}"
echo "Child Playlist:  ${OUTPUT_DIR}/${CHILD_PLAYLIST_NAME}"
echo "Usando loop infinito e segmentação de 2 segundos."
echo "-------------------------------------------------------------------\n\n"

# 3. Executa o FFmpeg
TZ='UTC' ffmpeg \
    -stream_loop -1 -re \
    -i "$INPUT_FILE" \
    -fflags +nobuffer \
    -c:v libx264 -preset slow -b:v 4M -maxrate 4M -bufsize 8M -g 48 -sc_threshold 0 -keyint_min 48 \
    -c:a aac -b:a 128k -ac 2 \
    -f hls \
    -hls_time 2 \
    -hls_list_size 6 \
    -hls_flags delete_segments+append_list+omit_endlist+program_date_time+independent_segments \
    -hls_segment_type mpegts \
    -master_pl_name "${MASTER_PLAYLIST_NAME}" \
    -hls_segment_filename "${OUTPUT_DIR}/segment_%03d.ts" \
    "${OUTPUT_DIR}/${CHILD_PLAYLIST_NAME}"

# 4. Gera master playlist apontando para a child
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=4500000,RESOLUTION=1920x1080,NAME="1080p"
${CHILD_PLAYLIST_NAME}
EOF

echo -e "\n\n-------------------------------------------------------------------"
echo "Streaming (Live Loop) encerrado."
echo "Arquivos gerados em: ${OUTPUT_DIR}"
echo "-------------------------------------------------------------------\n\n"
