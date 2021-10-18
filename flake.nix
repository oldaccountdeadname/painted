{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { self, nixpkgs }: {
    defaultPackage.x86_64-linux = let pkgs = import nixpkgs { system = "x86_64-linux"; }; in
      pkgs.buildGoModule {
        name = "painted";
        version = "v0.1.0";

        src = ./.;
        vendorSha256 = "sha256-Nsnw5er32WosaHUIqc13qwh+vnQ01LZ9wIXECIu2VXk=";
      };
  };
}
