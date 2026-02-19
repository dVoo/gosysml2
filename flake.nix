{
  description = "Cline environment";

  inputs = {
    # Use your preferred nixpkgs channel
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    # OpenCode as flake (tracks its latest main by default)
    opencode-flake.url = "github:AodhanHayter/opencode-flake"; # wraps SST OpenCode CLI [web:1]

    # Claude Code hourly-updated flake
    claude-code.url = "github:sadjow/claude-code-nix"; # self-updating Claude Code package [web:9][web:12]

    # numtide’s ai tools (includes codex-cli)
    nix-ai-tools.url = "github:numtide/nix-ai-tools"; # provides a codex CLI output [web:7][web:10]
  };

  outputs =
    {
      self,
      nixpkgs,
      opencode-flake,
      claude-code,
      nix-ai-tools,
      ...
    }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        packages = [
          # OpenCode CLI binary from the opencode flake
          opencode-flake.packages.${system}.default # usually exposes `opencode` [web:1]

          # Claude Code CLI (the flake bundles its own Node runtime) [web:9][web:12]
          claude-code.packages.${system}.default # usually exposes `claude`

          # Codex CLI from nix-ai-tools (often called `codex-cli` or `codex`) [web:7][web:10]
          nix-ai-tools.packages.${system}.codex

          pkgs.antlr
          pkgs.openjdk
        ];
      };
    };
}
