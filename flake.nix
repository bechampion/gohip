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
          gohip-vendor-hash = import ./gohip.vendor.hash.nix;
          gohip-package = (pkgs.buildGoModule {
            name = "gohip";
            pname = "gohip";

            src = ./.;
            vendorHash = gohip-vendor-hash pkgs;
            proxyVendor = true;
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
                "${gohip-package}"/bin/gohip
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
