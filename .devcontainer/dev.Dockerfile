FROM mcr.microsoft.com/devcontainers/miniconda:dev-3

# nix installation
RUN curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install linux \
  --init none \
  --no-confirm
ENV PATH="${PATH}:/nix/var/nix/profiles/default/bin"

# go and pnpm installation
RUN nix profile install "nixpkgs#go"
RUN nix profile install "nixpkgs#pnpm_9"
