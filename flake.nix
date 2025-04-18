{
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    panfactum.url = "github:panfactum/stack/edge.25-04-03";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
  };

  outputs = { nixpkgs, panfactum, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem
    (system:
     let
        pkgs = import nixpkgs {
         inherit system;
         config = { allowUnfree = true; };
       };
     in
      {
        devShell = panfactum.lib.${system}.mkDevShell {
         packages = [
           pkgs.nodejs_22
           pkgs.go
         ];
        };
      }
    );
}