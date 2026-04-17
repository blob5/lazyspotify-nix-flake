{
  description = "Nix packaging for lazyspotify and patched librespot";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    {
      self,
      nixpkgs,
      ...
    }:
    let
      systems = [ "x86_64-linux" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f system);
    in
    {
      overlays.default = import ./pkgs/overlay.nix;

      packages = forAllSystems (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [ self.overlays.default ];
          };
        in
        {
          lazyspotify-librespot = pkgs.lazyspotify-librespot;
          lazyspotify = pkgs.lazyspotify;
          default = pkgs.lazyspotify;
        }
      );
    };
}
