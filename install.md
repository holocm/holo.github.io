---
title: "Installation"
---

<header>
 <div id="header-buttons">
  [[]{.logo .logo-twitter}](https://twitter.com/holocm "Follow on Twitter")
  [[]{.logo .logo-github}](https://github.com/holocm "Fork on GitHub")
 </div>
 <div id="small-logo">
  [![](img/holo-logo.svg)](./index.html)
 </div>
</header>
<section>

How to install Holo
===================

Whenever possible, install Holo as a package. Package sources are listed
below by distribution. (If you are a packager, please [contact
us](https://github.com/holocm/holo/issues/new) if you\'re working on
getting Holo packaged in your distribution.)

From source
-----------

Alternatively, install from source by cloning the corresponding Git
repos and looking at the `README.md`:

-   the core execuable: [holo](https://github.com/holocm/holo)
-   optional plugins for Holo:
    [holo-users-groups](https://github.com/holocm/holo-users-groups),
    [holo-ssh-keys](https://github.com/holocm/holo-ssh-keys) and
    [holo-run-scripts](https://github.com/holocm/holo-run-scripts)
-   the package building tool:
    [holo-build](https://github.com/holocm/holo-build)

Arch Linux packages
-------------------

Packages are available from the
[AUR](https://aur.archlinux.org/packages/?K=holo&SeB=n). If you prefer
pre-compiled packages, you can use our repository (only available for
x86\_64 at the moment). Add the following to your `/etc/pacman.conf`:

    [holo]
    Server = https://repo.holocm.org/archlinux/$arch

The packages are signed with GPG, so you need to import and trust the
signing key:

    $ sudo pacman-key -r 0xF7A9C9DC4631BD1A
    $ sudo pacman-key -f 0xF7A9C9DC4631BD1A
    #     check the key fingerprint before continuing, it should be:
    #     2A53 49F6 B4D7 305A 85DE  D8D4 F7A9 C9DC 4631 BD1A
    $ sudo pacman-key --lsign-key 0xF7A9C9DC4631BD1A

To get a list of all packages in this repo, use:

    $ sudo pacman -Syl holo

To give feedback about these packages, please [use GitHub
Issues](https://github.com/majewsky/holo-pacman-repo/issues/new).

Debian/Ubuntu packages
----------------------

We offer packages that have been cross-compiled on Arch Linux. Add the
following to your `/etc/apt/sources.list`:

    deb https://repo.holocm.org/debian stable main

The repository is signed with GPG, so you need to import and trust the
signing key:

    $ sudo apt-key adv --keyserver pool.sks-keyservers.net --recv-keys 0xD6019A3E17CA2D96
    $ sudo apt-key adv --fingerprint 0xD6019A3E17CA2D96
    #     check the key fingerprint before continuing, it should be:
    #     BB41 E373 DB03 9091 6D23  2BFF D601 9A3E 17CA 2D96

To give feedback about these packages, please [use GitHub
Issues](https://github.com/majewsky/holo-foreign-repo/issues/new).

</section>
