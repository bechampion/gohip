pkgs :

let
  upgrade-nix-gohip = pkgs.writeShellScriptBin "upgrade-nix-gohip" ''
    echo "UPGRADING"
    echo "pkgs : pkgs.lib.fakeHash" > $WORKSPACE/gohip.vendor.hash.nix
    NEW_HASH=`nix build .#default 2>&1 | grep " got: " | awk '{print $2}'`
    echo "pkgs : \"$NEW_HASH\"" > $WORKSPACE/gohip.vendor.hash.nix

  '';
in
  pkgs.mkShell {
    buildInputs = [ pkgs.go upgrade-nix-gohip ];
  }
