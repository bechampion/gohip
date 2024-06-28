pkgs : custom-packages :

let
in
  pkgs.mkShell {
    buildInputs = [ pkgs.go custom-packages.upgrade-nix-gohip ];
  }
