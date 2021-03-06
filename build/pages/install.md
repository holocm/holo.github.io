# How to install Holo

Whenever possible, install Holo as a package. Package sources are listed below by distribution. (If you are a packager,
please [contact us](https://github.com/holocm/holo/issues/new) if you're working on getting Holo packaged in your
distribution.)

## From source

Alternatively, install from source by cloning the corresponding Git repos and looking at the `README.md`:

- the core executable: [holo](https://github.com/holocm/holo) (including the default set of plugins)
- the package building tool: [holo-build](https://github.com/holocm/holo-build)

## Arch Linux packages

Packages are available from the [AUR](https://aur.archlinux.org/packages/?K=holo&SeB=n). If you prefer pre-compiled
packages, you can use our repository (only available for x86\_64 at the moment). Add the following to your
`/etc/pacman.conf`:

```
[holo]
Server = https://repo.holocm.org/archlinux/$arch
```

The packages are signed with GPG, so you need to import and trust the signing key. The repo contains a keyring package
that you can use for this, but make sure to validate the checksum before installing it.

```
$ wget https://repo.holocm.org/archlinux/x86_64/holo-keyring-20201009.1-1-any.pkg.tar.xz
...
$ sha256sum < holo-keyring-20201009.1-1-any.pkg.tar.xz
dec378054732fad0109eeff5da3933cefb70eaeb14217f20a33510f3772aea95  -
$ sudo pacman -U holo-keyring-20201009.1-1-any.pkg.tar.xz
...
```

To get a list of all packages in the repo, use:

```
$ sudo pacman -Syl holo
```

To give feedback about these packages, please [use GitHub Issues](https://github.com/majewsky/holo-pacman-repo/issues/new).

## Parabola packages

Parabola includes a Holo package in its \[PCR\] repository, which is enabled by default:

```
$ sudo pacman -S holo
```

To give feedback on this package, use [Parabola's bug tracker](https://labs.parabola.nu/projects/issue-tracker/issues).

However, Parabola does not include holo-build. Follow the Arch Linux section above to install it.
