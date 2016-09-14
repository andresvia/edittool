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
    # EDITTOL GENERATED DO NOT EDIT
    1.1.1.1 my.special.server.example.com
    # EDITTOL GENERATED DO NOT EDIT

Run the same command as many times as you want there should be no side effects.

Return values
-------------

 - 0 - nothing happened
 - 2 - changes were made
 - anything else - error

Potential uses
--------------

Edit files before real provisioning kicks-in, like:

 - `/etc/fstab`
 - `/etc/hosts`
 - `/etc/resolv.conf`
 - `$HOME/.ssh/authorized_keys`
 - etc...
