{ pkgs ? import <nixpkgs> { } }:
pkgs.stdenv.mkDerivation {
  name = "KVHelper";
  version = "1.0.0";
  src = ../../.;

  nativeBuildInputs = [ pkgs.go ];

  buildPhase = ''
    export GOPATH=$TMPDIR
    export GOCACHE=$TMPDIR/go-cache
    go build -o KVHelper ./cmd/cli
  '';

  installPhase = ''
    mkdir -p $out/bin
    cp KVHelper $out/bin/KVHelper
    chmod +x $out/bin/KVHelper
  '';
}