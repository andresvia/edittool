edittool
========

Dogfood implementation of [editlib](https://github.com/andresvia/editlib).

But also useful for bootstrap provisioning in a idempotent way.

Example use
-----------

    cat /tmp/extra_host_config
    1.1.1.1 my.special.server.example.com
    sudo edittool -edit /etc/hosts -ensure /tmp/extra_host_config
    cat /etc/hosts
    ##
    # Host Database
    #
    # localhost is used to configure the loopback interface
    # when the system is booting.  Do not change this entry.
    ##
    127.0.0.1 localhost
    255.255.255.255 broadcasthost
    ::1             localhost
    # EDITTOOL GENERATED DO NOT EDIT
    1.1.1.1 my.special.server.example.com
    # EDITTOOL GENERATED DO NOT EDIT

Run the same command as many times as you want there should be no side effects.

Use in scripting
----------------

Example 1.

    #!/bin/bash
    set -e -u
    # ...
    # set up /tmp/extra_host_config
    # ...
    sudo edittool -edit /etc/hosts -ensure /tmp/extra_host_config || [ $? -eq 2 ]

Example 2.

    #!/bin/bash
    set -e -u
    # ...
    # set up /tmp/my_service_extra_config
    # ...
    sudo edittool -edit /etc/my_service.conf -ensure /tmp/my_service_extra_config -reload 'sudo systemctl restart my_service'

Return values
-------------

Without the `-reload` flag.

 - 0 - nothing happened
 - 2 - changes were made
 - anything else - error

With the `-reload` flag.

 - 0 - nothing happened / reload succeed
 - anything else - error / reload return code

Potential uses
--------------

Edit files before [real](https://github.com/provisioningsucks/tools) provisioning kicks-in (especially services or configuration files without `conf.d` style configuration), like:

 - `/etc/fstab`
 - `/etc/hosts`
 - `/etc/resolv.conf`
 - `$HOME/.ssh/authorized_keys`
 - etc...
