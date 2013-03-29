      ___
    C(o o)O     C H U N K Y   M O N K E Y
      (.)   \
     w====m==|       minecraft server
            /

Chunky Monkey is a Minecraft Beta multiplayer server.  It is licensed under
the MIT open source license, please see the LICENSE file for more information.

Website: http://github.com/b1naryth1ef/chunkymonkey

Status
------

This fork of Chunky Monkey is more or less the third revision. The last version was abandonded around Minecraft 1.1.
The goal of this fork is to update the source to work with Minecraft 1.5 and then start a period of consistent releases
and upgrades to maintain this source. The overal goal is a decent replacement of the vanilla minecraft server that implements
a smart and easy to use scripting system, benchmarks faster, and performs better.


Features include:

*   Compatibility with Minecraft Beta clients.
*   Blocks may be dug, items picked up, and blocks placed.
*   Crafting using the 2x2 or 3x3 crafting grids and the furnace.
*   Partial item physics.
*   World persistency.
*   Basic world generation.

Currently missing features include:

*   Block physics.
*   Complete item physics (there is a minimal implementation in place).
*   Mob behaviour.
*   Many block interactions.
*   Decent world generation.

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

Serve up a single player world:

    $ bin/chunkymonkey ~/.minecraft/saves/World1
    2010/10/03 16:32:13 Listening on  :25565

Record/replay
-------------

For debugging it is often useful to record a player's actions and replay them
one or more times later.  This makes it possible to simulate multiplayer games
without having real people logging in.

To record a session, run the intercept proxy:

    $ bin/intercept -record player.log localhost:25567 localhost:25565

Which will accept client connections on localhost port 25567 and relay the
connection to the server at localhost port 25565. This has the side effect of
display packets that pass through to stderr.

Connect your Minecraft client to localhost:25567, and a record of the clients
actions will be stored to player.log-1, player.log-2 etc. (one file per client
connection).

To replay a session:

    $ bin/replay localhost:25565 player.log-1


[1]: http://golang.org/doc/install.html          "Go toolchain installation"
[2]: http://code.google.com/p/godag/wiki/Install "Godag builder"
[3]: https://github.com/huin                     "Huin on Github"
[4]: http://code.google.com/p/gomock/            "GoMock mocking library"
