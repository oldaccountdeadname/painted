{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { self, nixpkgs }:
    let pkgs = import nixpkgs { system = "x86_64-linux"; };
    in {
      defaultPackage.x86_64-linux = pkgs.buildGoModule {
        name = "painted";
        version = "v0.1.2";

        src = builtins.filterSource
          (path: type: baseNameOf path != "contrib")
          ./.;

        vendorSha256 = "sha256-8pkdIMgXduxrhB4LdEOyX81PGA6ya5b1VCJStptqmd0=";
      };

      devShell.x86_64-linux = pkgs.mkShell {
        buildInputs = [ pkgs.go pkgs.libnotify ];
        shellHook = ''
          ln -sf ../../.githooks/pre-commit .git/hooks/pre-commit
        '';
      };
    };
}
