// URLs dos vídeos
const hlsVideo = "/hls/playlist_1080p.m3u8";  // HLS
const dashVideo = "/dash/playlist_1080p.mpd";  // DASH

// Inicializa o player
const player = videojs("my-video");

// Função para carregar qualquer fonte
function loadVideo(src) {
  const type = src.endsWith(".m3u8")
    ? "application/x-mpegURL"
    : src.endsWith(".mpd")
    ? "application/dash+xml"
    : "video/mp4";

  player.src({ src, type });
  console.log(`Carregando vídeo: ${src} | type: ${type}`);
}

// Eventos para logar cada segmento (HLS/DASH)
player.on("loadedmetadata", () => {
  console.log("Metadados carregados");
});
player.on("loadeddata", () => {
  console.log("Primeiro frame carregado");
});
player.on("progress", () => {
  console.log("Progresso de download/segmentos");
});

// Botões para alternar vídeos
document.getElementById("play-hls").addEventListener("click", () => {
  loadVideo(hlsVideo);
});

document.getElementById("play-dash").addEventListener("click", () => {
  loadVideo(dashVideo);
});

// Carrega HLS por padrão
loadVideo(hlsVideo);
