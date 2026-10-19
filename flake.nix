{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "https://channels.nixos.org/nixpkgs-unstable/nixexprs.tar.zst";
  };

  outputs = inputs: {
    packages = builtins.mapAttrs (system: pkgs: {
      default =   inputs.self.packages.${system}.keyvalue;
      keyvalue =  import ./nix/packages/keyvalue.nix { inherit pkgs; };
      cli =       import ./nix/packages/cli.nix { inherit pkgs; };
    }) inputs.nixpkgs.legacyPackages;


    devShells = builtins.mapAttrs (system: pkgs: {
      default = import ./nix/shell.nix { inherit pkgs; };
    }) inputs.nixpkgs.legacyPackages;
  };
}
