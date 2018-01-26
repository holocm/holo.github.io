---
title: "Minimalistic Configuration Management"
---

<header>
 <div id="header-buttons">
  [[]{.logo .logo-twitter}](https://twitter.com/holocm "Follow on Twitter")
  [[]{.logo .logo-github}](https://github.com/holocm "Fork on GitHub")
 </div>

![Holo](./img/holo-logo.svg) Minimalistic configuration management
==================================================================
</header>

<section class="features flex">

<div>
### ![](img/feature-declarative.svg) Declarative

Document system configuration with a friendly DSL.
</div>

<div>
### ![](img/feature-integrated.svg) Integrated

Holo uses the system package manager to maintain the system state.
</div>

<div>
### ![](img/feature-serverless.svg) Serverless

No dedicated infrastructure, network access or permanently running agent
required.
</div>

</section>

<section class="audience flex">

<div>
This might be for you if\...
----------------------------

-   \...you are looking for a simple way to formalize the configuration
    of a small number of systems.
-   \...you want to get started with the basic ideas of configuration
    management before moving to an enterprise solution.
-   \...you ship configuration to your clients, and they are already
    doing configuration management. They can install your Holo-based
    configuration as system packages using their own configuration
    management tool.

<p class="buttons">
 [See example](./example.html){.button}
 [Installation](./install.html){.button}
</p>
</div>

<div>
This is not for you if\...
--------------------------

-   \...your OS does not have a system package manager (such as dpkg or
    RPM) that exerts authority over the system partition.
-   \...your applications are running in containers. You\'re better off
    baking your configuration into the container image from the start.
-   \...you need one of the [features that are not yet
    implemented.](https://github.com/holocm/holo/issues?q=is:issue+is:open+label:%22type:+feature%22)
    The biggest omissions right now are variable-based templating and
    verification of the full system state. (If in doubt, feel free to
    [ask a question](https://github.com/holocm/holo/issues/new) about
    your required features.)

<p class="buttons">
 [See alternatives](https://en.wikipedia.org/wiki/Comparison_of_open-source_configuration_management_software){.button}
</p>
</div>

</section>
