# Contributing

## Writing Patches
Painted is built and developed with Nix to make everything deterministic: what
works on your machine will work on mine, so long as you do your development in
a Nix-aware fashion. You'll usually want to your development inside a Nix
development shell, where you have access to everything on your host system, plus
libnotify of a known-good version, the Go toolchain, and a Git hook for
formatting your code.

To get access to a development shell (and thus, all necessary prerequisites for
building) run the below on a machine with a version of Nix supporting Flakes
installed:

```bash
nix develop
```

## Submitting Patches (Merge Requests)
Development is done through GitLab, *not* GitHub! Fork, push, and submit a merge
request there.
