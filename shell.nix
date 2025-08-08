{
  pkgs,
  ...
}:
rec {
  default = pkgs.mkShell {
    nativeBuildInputs = [
      # keep-sorted start
      pkgs.cocogitto
      pkgs.editorconfig-checker
      pkgs.git
      pkgs.goEnv
      pkgs.goreleaser
      pkgs.gotestdox
      pkgs.grype
      pkgs.jq
      pkgs.moreutils
      pkgs.omnix
      pkgs.ripgrep
      pkgs.semgrep
      pkgs.unstable.go
      pkgs.unstable.golangci-lint
      # keep-sorted end
    ];
  };

  ci = default.overrideAttrs (
    final: prev: {
      nativeBuildInputs = [
        # keep-sorted start
        pkgs.go-junit-report
        pkgs.nur.repos.wwmoraes.codecov-cli-bin
        # keep-sorted end
      ] ++ prev.nativeBuildInputs;

      shellHook = ''
        export GOCACHE=$(go env GOCACHE)
        export GOMODCACHE=$(go env GOMODCACHE)
      '';
    }
  );

  terminal = default.overrideAttrs (
    final: prev: {
      nativeBuildInputs = [
        # keep-sorted start
        pkgs.curl
        pkgs.eclint
        pkgs.gomod2nix
        pkgs.nix-update
        pkgs.unstable.cocogitto
        pkgs.unstable.gotests
        pkgs.unstable.gotools
        pkgs.unstable.sarif-fmt
        # keep-sorted end
      ] ++ prev.nativeBuildInputs;

      shellHook = ''
        cog install-hook --all --overwrite
      '';
    }
  );
}
