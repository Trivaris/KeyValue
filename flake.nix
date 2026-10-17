{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "https://channels.nixos.org/nixpkgs-unstable/nixexprs.tar.zst";
  };

  outputs = inputs: {
    packages = builtins.mapAttrs (system: pkgs: {
      default = inputs.self.packages.${system}.keyvalue;
      keyvalue = pkgs.stdenv.mkDerivation {
        name = "KeyValue";
        version = "1.0.0";
        src = ./.;

        nativeBuildInputs = [ pkgs.go ];

        buildPhase = ''
          export GOPATH=$TMPDIR
          export GOCACHE=$TMPDIR/go-cache
          go build -o KeyValue ./cmd/server
        '';

        installPhase = ''
          mkdir -p $out/bin
          cp KeyValue $out/bin/KeyValue
          chmod +x $out/bin/KeyValue
        '';
      };
    }) inputs.nixpkgs.legacyPackages;


    devShells = builtins.mapAttrs (system: pkgs: {
      keyvalue = pkgs.mkShell {
        packages = [
          pkgs.go
          pkgs.gopls
          pkgs.jq
          pkgs.docker
        ];

      shellHook = ''
        mkdir -p .vscode
        [ ! -f .vscode/settings.json ] && echo "{}" > .vscode/settings.json
        jq '.["go.goroot"] = "${pkgs.go}/share/go"' .vscode/settings.json > .vscode/settings.json
      '';
      };

      default = inputs.self.devShells.${system}.keyvalue;
    }) inputs.nixpkgs.legacyPackages;
  };
}
