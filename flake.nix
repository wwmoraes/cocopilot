{
  description = "cocopilot fetches API tokens to use GitHub Copilot with any tool";

  inputs = {
    flake-utils = {
      inputs.systems.follows = "systems";
      url = "github:numtide/flake-utils";
    };
    gomod2nix = {
      inputs.flake-utils.follows = "flake-utils";
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:tweag/gomod2nix";
    };
    nixpkgs.url = "github:NixOS/nixpkgs/25.05";
    nur = {
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-parts.follows = "flake-parts";
      url = "github:nix-community/NUR";
    };
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:numtide/treefmt-nix";
    };
    unstable.url = "github:NixOS/nixpkgs?rev=e38c80c027d6bbdfa2a305fc08e732b9fac4928a";
  };

  nixConfig = {
    substituters = [
      "https://wwmoraes.cachix.org/"
      "https://nix-community.cachix.org/"
      "https://cache.nixos.org/"
    ];
    trusted-public-keys = [
      "wwmoraes.cachix.org-1:N38Kgu19R66Jr62aX5rS466waVzT5p/Paq1g6uFFVyM="
      "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
      "cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="
    ];
  };

  outputs =
    inputs@{
      self,
      flake-parts,
      gomod2nix,
      nixpkgs,
      nur,
      # sops-nix,
      systems,
      treefmt-nix,
      unstable,
      ...
    }:
    (flake-parts.lib.mkFlake { inherit inputs; } {
      flake = {
        overlays = {
          default =
            final: prev:
            prev.lib.recursiveUpdate prev {
              lib.maintainers.wwmoraes = {
                email = "nixpkgs@artero.dev";
                github = "wwmoraes";
                githubId = 682095;
                keys = [ { fingerprint = "32B4 330B 1B66 828E 4A96  9EEB EED9 9464 5D7C 9BDE"; } ];
                matrix = "@wwmoraes:hachyderm.io";
                name = "William Artero";
              };
            }
            // {
              inherit (self.packages.${prev.system}) cocopilot;
            };
        };
      };

      perSystem =
        { pkgs, system, ... }:
        let
          treefmt = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
        in
        {
          _module.args.pkgs = import nixpkgs {
            inherit system;
            overlays = [
              gomod2nix.overlays.default
              nur.overlays.default
              self.overlays.default
              (final: prev: {
                goEnv = prev.mkGoEnv { pwd = ./.; };
                unstable = import unstable { inherit (prev) system; };
              })
            ];
            config = { };
          };

          checks = {
            formatting = treefmt.config.build.check self;
          };

          formatter = treefmt.config.build.wrapper;

          devShells = import ./shell.nix { inherit pkgs; } // {
            treefmt = treefmt.config.build.devShell;
          };

          packages = rec {
            default = import ./default.nix { inherit pkgs; };
            cocopilot = default;
          };
        };

      systems = import systems;
    });
}
