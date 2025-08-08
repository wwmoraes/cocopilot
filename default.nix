{
  pkgs,
  ...
}:

pkgs.buildGoApplication {
  pname = "cocopilot";
  version = "0.1.0";
  src =
    with pkgs.lib.fileset;
    toSource {
      root = ./.;
      fileset = unions [
        (fileFilter (file: file.hasExt "go") ./.)
        ./cmd/cocopilot
        ./go.mod
        ./go.sum
      ];
    };
  modules = ./gomod2nix.toml;
  subPackages = [ "cmd/cocopilot" ];
  meta = {
    description = "fetches API tokens to use GitHub Copilot with any tool";
    homepage = "https://github.com/wwmoraes/cocopilot";
    license = pkgs.lib.licenses.mit;
    maintainers = [ pkgs.lib.maintainers.wwmoraes ];
    mainProgram = "cocopilot";
  };
}
