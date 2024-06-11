{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          overlays = [];
          pkgs = import nixpkgs {
            inherit system overlays;
          };
          gohip-package = (pkgs.buildGoModule {
            name = "gohip";
            pname = "gohip";

            src = ./.;
            vendorHash = "sha256-kikKkz2XnAo6MVVbu4VOMyW9xUqlxjTrjhzAaGph1CY=";
            excludedPackages = ["osdata" "others" "systemd" "types"];
          });

        in
        with pkgs;
        {
          devShells.default = (import ./shell.nix) pkgs;
          apps = rec {
            gohip = flake-utils.lib.mkApp {
              drv = pkgs.writeShellScriptBin "gohip" ''
                cd `mktemp -d`
                "${gohip-package}"/bin/hip
              '';
            };
            default = gohip;
          };
          packages = rec {
            gohip = gohip-package;
            default = gohip;
          };
        }
      );

}
