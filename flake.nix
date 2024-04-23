{
  description = "A flake for developing brain.test";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    nixpkgs-pw-v1_40.url =
      "github:nixos/nixpkgs/a3ed7406349a9335cb4c2a71369b697cecd9d351";
  };

  outputs = { self, nixpkgs, flake-utils, nixpkgs-pw-v1_40 }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        pkgs-pw-v1_40 = import nixpkgs-pw-v1_40 { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            #mdbook
            mdbook
            mdbook-mermaid
            # nodejs
            nodejs_21
            pkgs-pw-v1_40.playwright
            # golang
            go
            go-task
            air
            govulncheck
            gotestsum
            go-swag
            golangci-lint
            #
          ];
          PLAYWRIGHT_NODEJS_PATH = "${pkgs.nodejs_21}/bin/node";
          PLAYWRIGHT_BROWSERS_PATH =
            "${pkgs-pw-v1_40.playwright-driver.browsers}";
          PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD = 1;
          PLAYWRIGHT_LAUNCH_OPTIONS_EXECUTABLE_PATH =
            "${pkgs.playwright-driver.browsers}/chromium-1091/chrome-linux/chrome";
        };
      });
}
