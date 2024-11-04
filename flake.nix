{
  description = "Nixity";

  inputs = {
    devenv.url = "github:cachix/devenv";
    nixities.url = "github:ereslibre/nixities";
    systems.url = "github:nix-systems/default";
  };

  outputs = {
    self,
    devenv,
    nixities,
    systems,
    ...
  } @ inputs: let
    eachSystem = nixities.nixpkgs.lib.genAttrs (import systems);
  in {
    # Fix issue with devenv-up missing with flakes: https://github.com/cachix/devenv/issues/756
    packages = eachSystem (system: {
      devenv-up = self.devShells.${system}.default.config.procfileScript;
    });
    devShells = eachSystem (system: let
      pkgs = import nixities.nixpkgs {inherit system;};
    in {
      # The default devShell
      default = devenv.lib.mkShell {
        inherit pkgs;
        inputs.nixpkgs = nixities.nixpkgs;
        modules = [
          ({pkgs, ...}: {
            env = {
              DOCKER_HOST = "unix:///var/run/docker.sock";
              WASMTIME_NEW_CLI = "1";
              TERM = "xterm";
            };
            languages = {
              c.enable = true;
              go.enable = true;
            };
            packages =
              (with pkgs; [
                alejandra
                bat
                fermyon-spin
                just
                wasm-tools
                wasmtime
              ])
              ++ pkgs.lib.lists.optionals pkgs.stdenv.isx86_64 (with nixities.packages.${system}; [
                wasi-sdk-20
              ]);
          })
        ];
      };
    });
  };
}
