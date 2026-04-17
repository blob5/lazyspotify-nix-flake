final: prev: {
  lazyspotify-librespot = final.callPackage ./lazyspotify/librespot.nix { };
  lazyspotify = final.callPackage ./lazyspotify/default.nix { };
}
