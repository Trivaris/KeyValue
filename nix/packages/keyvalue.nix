{ pkgs ? import <nixpkgs> { } }:
pkgs.stdenv.mkDerivation {
  name = "KeyValue";
  version = "1.0.0";
  src = ../../.;

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
}