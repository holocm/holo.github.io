---
title: "Example"
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

How does Holo configure systems?
================================

Holo users ship configuration in packages, usually called \"holograms\".
These can be built with the regular package building tools (debuild,
rpmbuild, makepkg, etc.) or with Holo\'s own
[holo-build](https://github.com/holocm/holo-build) tool that offers a
much more pleasant syntax and process. Let\'s go through an example
hologram that installs and starts an OpenSSH server and tweaks some of
its configuration.

The package declaration
-----------------------

Package declarations for holo-build are
[TOML](https://github.com/toml-lang/toml) files.

    [package]
    name        = "hologram-openssh"
    version     = "1.0.0"
    description = "Start and configure OpenSSH"
    requires    = ["openssh"]

Packages can install files, directories and symlinks. For example, we
may want to start SSH only after the firewall is set up, so we create a
configuration file for systemd.

    [[file]]
    path    = "/etc/systemd/system/sshd.service.d/hardened.conf"
    content = """
      [Unit]
      After=firewall.service
    """

We also want to disable password authentication. This one is a bit more
tricky: We want to modify the configuration installed by the OpenSSH
package, but the configuration is also a package, so it may not install
the same file path. Instead, we install a script that Holo will later
find and execute to update the default configuration.

    [[file]]
    path    = "/usr/share/holo/files/10-openssh/etc/ssh/sshd_config.holoscript"
    mode    = "0755"
    content = """
      #!/bin/sh
      # stdin has the default config and stdout wants the updated config;
      # we just add a line at the bottom
      cat
      echo "PasswordAuthentication no"
    """

Any file below **/usr/share/holo** will imply a dependency on the
**holo** package.

When everything is set up, we start the daemon:

    [[action]]
    on     = "setup"
    script = """
      systemctl daemon-reload
      systemctl enable sshd
      systemctl restart sshd
    """

Rolling it out
--------------

Once the package declaration is complete, a system package (.deb, .rpm,
etc.) can be produced by
[holo-build](https://github.com/holocm/holo-build). No extra tools
needed.

    $ holo-build --debian hologram-openssh.pkg.toml

Since we had files below **/usr/share/holo**, Holo will be installed and
**holo apply** will be executed during installation:

    # dpkg -i hologram-openssh_1.0.0-1_any.deb
    ...
    Working on file:/etc/ssh/sshd_config
      store at /var/lib/holo/files/base/etc/ssh/sshd_config
      passthru /usr/share/holo/files/10-openssh/etc/ssh/sshd_config.holoscript
    ...

This tells us that the default configuration has been modified as
described by our holoscript. And indeed:

    $ tail -n1 /etc/ssh/sshd_config
    PasswordAuthentication no

Monitoring for changes
----------------------

When Holo provisions an entity (such as this config file), it will
always store a **base image** describing the original state of the
entity. If the entity is changed afterwards, Holo will be able to detect
this change:

    # sed -i '/PasswordAuthentication/ s/no/yes/' /etc/ssh/sshd_config
    # holo apply

    Working on file:/etc/ssh/sshd_config
      store at /var/lib/holo/files/base/etc/ssh/sshd_config
      passthru /usr/share/holo/files/10-openssh/etc/ssh/sshd_config.holoscript

    !! Entity has been modified by user (use --force to overwrite)

        diff --holo /var/lib/holo/files/provisioned/etc/ssh/sshd_config /etc/ssh/sshd_config
        --- /etc/ssh/sshd_config
        +++ /etc/ssh/sshd_config
        @@ -131,3 +131,3 @@
         #       ForceCommand cvs server

        -PasswordAuthentication no
        +PasswordAuthentication yes

But wait, there\'s more!
------------------------

With **plugins**, Holo can be taught to provision other things than
files. For example, there are plugins for [user accounts,
groups](https://github.com/holocm/holo-users-groups) or [SSH public
keys](https://github.com/holocm/holo-ssh-keys). You can easily write
your own plugins; they can be as small as [one shell
script](https://github.com/holocm/holo-run-scripts/tree/master/src/holo-run-scripts).

This example has demonstrated the holo-files plugin that ships with the
**holo** tool itself, but it can only scratch the surface. Check out the
[man pages](https://github.com/holocm/holo/tree/master/doc/) for the
full documentation. And don\'t forget to [install Holo](./install.html)
on your system, too.

</section>
