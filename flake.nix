{
  description = "AI work environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    opencode-flake.url = "github:anomalyco/opencode";  # Pin to matching version/commit from error
    claude-code.url = "github:sadjow/claude-code-nix";
    codex.url = "github:sadjow/codex-nix";
  };

  outputs =
    {
      self,
      nixpkgs,
      opencode-flake,
      claude-code,
      codex,
      ...
    }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        packages = [
          opencode-flake.packages.${system}.default
          claude-code.packages.${system}.default
          codex.packages.${system}.default

          pkgs.antlr
          pkgs.openjdk
        ];
      };
    };
}

