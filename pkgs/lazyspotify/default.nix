{
  buildGoModule,
  lib,
  fetchFromGitHub,
  pkg-config,
  alsa-lib,
  flac,
  libogg,
  libvorbis,
  lazyspotify-librespot,
}:

buildGoModule (finalAttrs: {
  pname = "lazyspotify";
  version = "0.2.3";

  src = fetchFromGitHub {
    owner = "dubeyKartikay";
    repo = "lazyspotify";
    rev = "v${finalAttrs.version}";
    hash = "sha256-JCXKSJmJ3NTwezKDEkb7fdcgbN0Bf1FKRpj7KBQzjMA=";
  };

  buildInputs = [
    alsa-lib
    flac
    libogg
    libvorbis
  ];

  nativeBuildInputs = [ pkg-config ];

  env.CGO_ENABLED = 1;

  subPackages = [ "cmd/lazyspotify" ];

  ldflags = [
    "-s"
    "-w"
    "-X github.com/dubeyKartikay/lazyspotify/buildinfo.Version=${finalAttrs.version}"
    "-X github.com/dubeyKartikay/lazyspotify/buildinfo.Commit=${finalAttrs.src.rev}"
    "-X github.com/dubeyKartikay/lazyspotify/buildinfo.BuildDate=unknown"
    "-X github.com/dubeyKartikay/lazyspotify/buildinfo.PackagedDaemonPath=${lib.getExe lazyspotify-librespot}"
  ];

  checkFlags = [ "./..." ];

  vendorHash = "sha256-Axdt3/3ZOZY9Z5VUI6Wh77oIREOO26ODMyEgtscTmn8=";

  meta = {
    description = "Terminal Spotify client";
    homepage = "https://github.com/dubeyKartikay/lazyspotify";
    license = with lib.licenses; [ mit ];
    mainProgram = "lazyspotify";
    platforms = lib.platforms.linux;
    maintainers = [ ];
  };
})
