{
  buildGoModule,
  lib,
  fetchFromGitHub,
  pkg-config,
  alsa-lib,
  flac,
  libogg,
  libvorbis,
}:

buildGoModule (finalAttrs: {
  pname = "lazyspotify-librespot";
  version = "0.7.1.1";

  src = fetchFromGitHub {
    owner = "dubeyKartikay";
    repo = "go-librespot";
    rev = "v${finalAttrs.version}";
    hash = "sha256-Hq9Qk8f8oKzpBwsbLNAvPO7qam3bh4L4RPUQC67/NZY=";
  };

  buildInputs = [
    alsa-lib
    flac
    libogg
    libvorbis
  ];

  nativeBuildInputs = [ pkg-config ];

  env.CGO_ENABLED = 1;

  subPackages = [ "cmd/daemon" ];

  ldflags = [
    "-s"
    "-w"
    "-X github.com/devgianlu/go-librespot.version=v${finalAttrs.version}"
  ];

  postInstall = ''
    mv "$out/bin/daemon" "$out/bin/lazyspotify-librespot"
  '';

  vendorHash = "sha256-5J5i2Wc0zHCdvJ3aUkftXeMKS5X8jWimup0Ir4HLuS8=";

  doCheck = false;

  meta = {
    description = "Patched go-librespot daemon for lazyspotify";
    homepage = "https://github.com/dubeyKartikay/go-librespot";
    license = with lib.licenses; [ gpl3Only ];
    mainProgram = "lazyspotify-librespot";
    platforms = lib.platforms.linux;
    maintainers = [ ];
  };
})
