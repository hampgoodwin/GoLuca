{
  description = "Goluca App";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils?ref=v1.0.0"; # I'm not sure if this is proper tag ref
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-25.05";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.goluca = pkgs.mkShell {
          name = "goluca";
          buildInputs = [
            # languages+toolings
            ## go
            pkgs.go_1_24
            pkgs.gopls
            pkgs.delve
            pkgs.gofumpt
            pkgs.golangci-lint
            ## bash
            pkgs.shellcheck
            pkgs.bash-language-server
            ## nix
            pkgs.nix
            pkgs.nixfmt-rfc-style
            pkgs.nixd

            # tooling
            ## protoencoding
            pkgs.buf
            ## openapi
            pkgs.redocly

            # containerization
            pkgs.colima
            pkgs.docker_28 # used for docker tools, not the runtime/engine
          ];

          shellHook = ''
            echo "📚 goluca"
          '';
        };
      }
    );
}
