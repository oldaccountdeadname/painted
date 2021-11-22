{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { self, nixpkgs }: let pkgs = import nixpkgs { system = "x86_64-linux"; }; in {
    defaultPackage.x86_64-linux = pkgs.buildGoModule {
      name = "painted";
      version = "v0.1.0";

      src = builtins.filterSource
        (path: type: baseNameOf path != "contrib")
        ./.;

      vendorSha256 = "sha256-m6uYoW7a/PuqAjEyDvDrS4tDIZP9SsuhBGWWYznKWjc=";
    };

    devShell.x86_64-linux = pkgs.mkShell {
      buildInputs = [ pkgs.go pkgs.libnotify ];
      shellHook = ''
        ln -sf ../../.githooks/pre-commit .git/hooks/pre-commit
      '';
    };
  };
}
