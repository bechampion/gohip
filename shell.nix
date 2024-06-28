pkgs :

let
  upgrade-nix-gohip = pkgs.writeShellScriptBin "upgrade-nix-gohip" ''
    echo "UPGRADING"
  '';
in
  pkgs.mkShell {
    buildInputs = [ pkgs.go upgrade-nix-gohip ];
  }
