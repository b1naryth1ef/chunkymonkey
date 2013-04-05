      ___
    C(o o)O     C H U N K Y   M O N K E Y
      (.)   \
     w====m==|       minecraft server
            /

Chunky Monkey is a Minecraft multiplayer server.  It is licensed under
the MIT open source license, please see the LICENSE file for more information.

Website: http://github.com/b1naryth1ef/chunkymonkey

Status
------

This fork of Chunky Monkey is more or less the third revision. The last version was abandonded around Minecraft 1.1.
The goal of this fork is to update the source to work with Minecraft 1.5 and then start a period of consistent releases
and upgrades to maintain this source. The overal goal is a decent replacement of the vanilla minecraft server that implements
a smart and easy to use scripting system, benchmarks faster, and performs better.

Goals
-----

* 100% compatibility with Minecraft spec
* Optimized, fast, smart, and scriptable server
* Support for go-get

Requirements
------------

The [Go toolchain][1] must be installed. Note that chunkymonkey is developed
against the current stable release of the Go toolchain, so might not compile
against the weekly releases (gofix might be able to fix such cases).


Building & Testing
------------------

[Godag][2] is used to build chunkymonkey. Install it, and run:

    $ make

If you are developing, you are encouraged to run the unit tests with:

    $ make test

The unit tests require [GoMock][4] to be installed.


Running
-------

Serve up a single player world (Warning: dont do this on worlds you like,
we might destroy them by accident D:):

    $ bin/chunkymonkey ~/.minecraft/saves/World1
    2010/10/03 16:32:13 Listening on  :25565

[1]: http://golang.org/doc/install.html          "Go toolchain installation"
[2]: http://code.google.com/p/godag/wiki/Install "Godag builder"
[3]: https://github.com/huin                     "Huin on Github"
[4]: http://code.google.com/p/gomock/            "GoMock mocking library"
