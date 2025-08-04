{
  description = "Goluca App";

  inputs = {
    inputs.flake-utils.url = "github:numtide/flake-utils?ref=1.0.0"; # I'm not sure if this is proper tag ref
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-25.05";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem(system:
    let pkgs = nixpkgs.legacyPackages.${system};
    in {
        devShells.goluca = pkgs.mkShell { buildInput  = [pkgs.go_1_24 pkgs.gopls pkgs.delve pkgs.gofumpt pkgs.golangci-lint ];};
      }
  )
}
