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
