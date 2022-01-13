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

## Submitting Patches
Painted attempts to have a linear history with few massive merge commits. Make
sure to commit atomic changes that always build and work properly with *good
commit messages*. Good commit messages look like this:

```
use the imperative to describe your change

The first line is in the imperative with no capital or punctuation.
Following paragraphs details why you made the change and what it does.
These should be mostly hard-wrapped at 72 characters.

More paragraphs may provide additional detail. If you include any code,
indent it with one tab, and precede it with a newline:

	fmt.Println("\thello world!");

If text follows a code block, precede it with a blank line.
```

## Updating Dependencies
Painted has two build systems, Go and Nix. Nix dependencies are managed through
Flakes, and this repo has a workflow to update them weekly. Assuming no breaking
compiler changes, the PRs created by that workflow can be merged with minimal
checks. Go dependencies are a bit more complicated, despite the relatively small
number of them. Dependabot updates them when needed, but the CI pipelines will
always fail. Why? The pipelines don't build through Go, they build through Nix.
Nix is deterministic, and, if the vendorSha256 isn't what Nix expects (as
happens when dependencies change) the build will fail.

The procedure to fix, assuming dependency updates don't break anything extra, is
as follows:
+ Update the venderSha256 (set it to pkgs.lib.fakeSha256 and copy the output of
  the failing build)
+ Make sure `nix build` works.
+ Ammend dependabot's commit to contain it.
+ Mark yourself as the author (you're the one who updated the hash).
+ Replace the first line of the commit message with `bump $dep: $oldVer ->
  $newVer`.
+ Note that you updated the sha.
+ Replace the last line's `Signed-off-by` with `co-authored-by` to indicate that
  this is dependabot's work, too.
+ Push.
