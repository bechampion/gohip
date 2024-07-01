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

          upgrade-nix-gohip-package = pkgs.writeShellScriptBin "upgrade-nix-gohip" ''
            echo "UPGRADING"
            OLD_HASH_FILE=`cat $WORKSPACE/gohip.vendor.hash.nix`
            echo "pkgs : pkgs.lib.fakeHash" > $WORKSPACE/gohip.vendor.hash.nix
            NEW_HASH=`nix build .#default 2>&1 | grep " got: " | awk '{print $2}'`
            echo "pkgs : \"$NEW_HASH\"" > $WORKSPACE/gohip.vendor.hash.nix
            NEW_HASH_FILE=`cat $WORKSPACE/gohip.vendor.hash.nix`
            if [[ "$OLD_HASH_FILE" != "$NEW_HASH_FILE" ]]; then
              git config --global user.name 'autobot'
              git config --global user.email 'gohip.nix@github.com'
              git commit -am "[skip ci] Automated nix hash update"
              git push
            fi
          '';
          custom-packages = {upgrade-nix-gohip = upgrade-nix-gohip-package;};

        in
        with pkgs;
        {
          devShells.default = (import ./shell.nix) pkgs custom-packages;
          apps = rec {
            gohip = flake-utils.lib.mkApp {
              drv = pkgs.writeShellScriptBin "gohip" ''
                cd `mktemp -d`
                "${gohip-package}"/bin/gohip
              '';
            };
            upgrade-nix-gohip = flake-utils.lib.mkApp {
              drv = upgrade-nix-gohip-package;
            };
            default = gohip;
          };
          packages = rec {
            gohip = gohip-package;
            upgrade-nix-gohip = upgrade-nix-gohip-package;
            default = gohip;
          };
        }
      );

}
