package proto

import (
    "bytes"
    "reflect"
    "testing"

    . "chunkymonkey/types"
    "nbt"
    te "testencoding"
)

const (
    Float32One = "\x3f\x80\x00\x00"
    Float32Two = "\x40\x00\x00\x00"

    Float64One   = "\x3f\xf0\x00\x00\x00\x00\x00\x00"
    Float64Two   = "\x40\x00\x00\x00\x00\x00\x00\x00"
    Float64Three = "\x40\x08\x00\x00\x00\x00\x00\x00"
    Float64Four  = "\x40\x10\x00\x00\x00\x00\x00\x00"
)

func testPacketSerial(t *testing.T, fromClient bool, outputPkt IPacket, expectedSerialization te.IBytesMatcher) {
    ps := new(PacketSerializer)

    // Test reading.
    input := new(bytes.Buffer)
    expectedSerialization.Write(input)
    if inputPkt, err := ps.ReadPacket(input, fromClient); err != nil {
        t.Errorf("Unexpected error reading packet: %v", err)
    } else {
        if !reflect.DeepEqual(outputPkt, inputPkt) {
            t.Errorf("Packet did not read expected value:\n  expected: %#v\n    result: %#v", outputPkt, inputPkt)
        }
    }

    // Test writing.
    output := new(bytes.Buffer)
    if err := ps.WritePacket(output, outputPkt); err != nil {
        t.Errorf("Unexpected error writing packet: %v\n  %#v\v", err, outputPkt)
    } else {
        if err := te.Matches(expectedSerialization, output.Bytes()); err != nil {
            t.Errorf("Output of writing packet did not match: %v\n%#v", err, outputPkt)
        }
    }
}

func Test_PacketLogin(t *testing.T) {
    testPacketSerial(
        t,
        true,
        &PacketLogin{
            VersionOrEntityId: 5,
            Username:          "username",
            MapSeed:           123,
            LevelType:         "DEFAULT",
            GameMode:          1,
            Dimension:         DimensionNormal,
            Difficulty:        GameDifficultyNormal,
            WorldHeight:       128,
            MaxPlayers:        12,
        },
        te.LiteralString("\x01"+
            "\x00\x00\x00\x05"+ // Version/EntityID
            "\x00\x08\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e"+ // Username
            "\x00\x00\x00\x00\x00\x00\x00\x7b"+ // MapSeed
            "\x00\x07\x00D\x00E\x00F\x00A\x00U\x00L\x00T"+ // LevelType
            "\x00\x00\x00\x01"+ // GameMode
            "\x00"+ // Dimension
            "\x02"+ // Difficulty
            "\x80"+ // WorldHeight
            "\x0c"), // MaxPlayers
    )
}

func Test_PacketHandshake(t *testing.T) {
    // Test long username.
    testPacketSerial(
        t,
        true,
        &PacketHandshake{
            UsernameOrHash: "username1username2username3username4username5" +
                "username6username7username8username9",
        },
        te.LiteralString("\x02"+
            "\x00\x51\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e\x001"+
            "\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e\x002"+
            "\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e\x003"+
            "\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e\x004"+
            "\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e\x005"+
            "\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e\x006"+
            "\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e\x007"+
            "\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e\x008"+
            "\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e\x009"),
    )

    // Test non-ASCII
    testPacketSerial(
        t,
        true,
        &PacketHandshake{
            UsernameOrHash: "üßərnáme",
        },
        te.LiteralString("\x02"+
            "\x00\x08\x00\xfc\x00\xdf\x02\x59\x00r\x00n\x00\xe1\x00m\x00e"),
    )
}

func Test_PacketUseEntity(t *testing.T) {
    testPacketSerial(
        t,
        true,
        &PacketUseEntity{
            User:      2,
            Target:    5,
            LeftClick: true,
        },
        te.LiteralString("\x07"+
            "\x00\x00\x00\x02"+
            "\x00\x00\x00\x05"+
            "\x01"),
    )
}

func Test_PacketPlayerPosition(t *testing.T) {
    testPacketSerial(
        t,
        true,
        &PacketPlayerPosition{
            X:  1, Y: 2, Stance: 3, Z: 4,
            OnGround: true,
        },
        te.LiteralString("\x0b"+
            Float64One+
            Float64Two+
            Float64Three+
            Float64Four+
            "\x01"),
    )
}

func Test_PacketPlayerPositionLook(t *testing.T) {
    testPacketSerial(
        t,
        true,
        &PacketPlayerPositionLook{
            X:  1, Y1: 2, Y2: 3, Z: 4,
            Look:     LookDegrees{Yaw: 1, Pitch: 2},
            OnGround: true,
        },
        te.LiteralString("\x0d"+
            Float64One+
            Float64Two+
            Float64Three+
            Float64Four+
            Float32One+
            Float32Two+
            "\x01"),
    )
}

func Test_PacketPlayerBlockInteract(t *testing.T) {
    testPacketSerial(
        t,
        true,
        &PacketPlayerBlockInteract{
            Block: BlockXyz{1, 2, 3},
            Face:  2,
            Tool: ItemSlot{
                ItemTypeId: 1,
                Count:      2,
                Data:       3,
            },
        },
        te.LiteralString("\x0f"+
            "\x00\x00\x00\x01"+
            "\x02"+
            "\x00\x00\x00\x03"+
            "\x02"+
            "\x00\x01"+
            "\x02"+
            "\x00\x03"),
    )

    // Test with last two fields missing (no tool used).
    testPacketSerial(
        t,
        true,
        &PacketPlayerBlockInteract{
            Block: BlockXyz{1, 2, 3},
            Face:  2,
            Tool: ItemSlot{
                ItemTypeId: -1,
            },
        },
        te.LiteralString("\x0f"+
            "\x00\x00\x00\x01"+
            "\x02"+
            "\x00\x00\x00\x03"+
            "\x02"+
            "\xff\xff"),
    )

    // Test with an item that can have NBT data, but which is empty.
    testPacketSerial(
        t,
        true,
        &PacketPlayerBlockInteract{
            Block: BlockXyz{1, 2, 3},
            Face:  4,
            Tool: ItemSlot{
                ItemTypeId: 0x111,
                Count:      1,
                Data:       0x30,
            },
        },
        te.LiteralString("\x0f"+
            "\x00\x00\x00\x01"+
            "\x02"+
            "\x00\x00\x00\x03"+
            "\x04"+
            "\x01\x11"+
            "\x01"+
            "\x00\x30"+
            "\xff\xff"),
    )

    // Test with an item that has NBT data.
    testPacketSerial(
        t,
        true,
        &PacketPlayerBlockInteract{
            Block: BlockXyz{1, 2, 3},
            Face:  4,
            Tool: ItemSlot{
                ItemTypeId: 0x111,
                Count:      1,
                Data:       0x30,
                Nbt: nbt.Compound{
                    "ench": &nbt.List{
                        TagType: nbt.TagCompound,
                        Value: []nbt.ITag{
                            nbt.Compound{
                                "id":  &nbt.Short{3},
                                "lvl": &nbt.Short{2},
                            },
                        },
                    },
                },
            },
        },
        te.LiteralString("\x0f"+
            "\x00\x00\x00\x01"+
            "\x02"+
            "\x00\x00\x00\x03"+
            "\x04"+
            "\x01\x11"+
            "\x01"+
            "\x00\x30"+
            "\x00\x37"+
            // TODO This really should use gzip library to read the output data.
            // Literal is somewhat fragile to underlying harmless changes.
            "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe2\x62\x60\xe0\x64\x60\x49\xcd"+
            "\x4b\xce\x00\xb2\x18\x18\x99\x18\x98\x73\xca\x72\x18\x98\x98\x18\x98\x32"+
            "\x53\x18\x98\x19\x18\x00\x01\x00\x00\xff\xff\x76\xca\x64\xf3\x20\x00\x00"+
            "\x00"),
    )
}

func Test_PacketEntityMetadata(t *testing.T) {
    testPacketSerial(
        t,
        false,
        &PacketEntityMetadata{
            EntityId: 5,
            Metadata: EntityMetadataTable{
                EntityMetadata{0, 0, byte(5)},
            },
        },
        te.LiteralString("\x28"+
            "\x00\x00\x00\x05"+
            "\x00\x05"+
            "\x7f"),
    )
}

func Test_PacketMapChunk(t *testing.T) {
    testPacketSerial(
        t,
        false,
        &PacketMapChunk{
            Corner: BlockXyz{16, 0, 32},
            Data: ChunkData{
                Size:       ChunkDataSize{0, 1, 2},
                Blocks:     []byte{1, 2, 3, 4, 5, 6},
                BlockData:  []byte{1, 2, 3},
                BlockLight: []byte{4, 5, 6},
                SkyLight:   []byte{7, 8, 9},
            },
        },
        te.InOrder(
            te.LiteralString("\x33"+
                "\x00\x00\x00\x10"+
                "\x00"+
                "\x00\x00\x00\x20"),
            // TODO This really should use zlib library to read the output data.
            // Literal is somewhat fragile to underlying harmless changes.
            te.LiteralString(""+
                "\x00\x01\x02"),
            te.LiteralString(""+
                "\x00\x00\x00\x17"+
                "\x78\x9c\x62\x64\x62\x66\x61\x65\x83\x90\xec\x1c"+
                "\x9c\x80\x00\x00\x00\xff\xff\x01\xa9\x00\x43"),
        ),
    )
}

func Test_PacketMultiBlockChange(t *testing.T) {
    testPacketSerial(
        t,
        false,
        &PacketMultiBlockChange{
            ChunkLoc: ChunkXz{1, 2},
            Changes: MultiBlockChanges{
                Coords:    []int16{5, 7, 9},
                TypeIds:   []byte{1, 2, 3},
                BlockData: []byte{4, 5, 6},
            },
        },
        te.LiteralString("\x34"+
            "\x00\x00\x00\x01\x00\x00\x00\x02"+
            "\x00\x03"+
            "\x00\x05\x00\x07\x00\x09"+
            "\x01\x02\x03"+
            "\x04\x05\x06"),
    )
}

func Test_PacketExplosion(t *testing.T) {
    testPacketSerial(
        t,
        false,
        &PacketExplosion{
            Center: AbsXyz{1, 2, 3},
            Radius: 2,
            Blocks: BlocksDxyz{1, 2, 3, 4, 5, 6},
        },
        te.LiteralString("\x3c"+
            Float64One+Float64Two+Float64Three+
            Float32Two+
            "\x00\x00\x00\x02"+
            "\x01\x02\x03\x04\x05\x06"),
    )
}

func Test_PacketWindowItems(t *testing.T) {
    testPacketSerial(
        t,
        false,
        &PacketWindowItems{
            WindowId: 5,
            Slots: ItemSlotSlice{
                ItemSlot{ItemTypeId: -1},
                ItemSlot{ItemTypeId: 3, Count: 7, Data: 1},
            },
        },
        te.LiteralString("\x68"+
            "\x05"+
            "\x00\x02"+
            "\xff\xff"+
            "\x00\x03\x07\x00\x01"),
    )
}

func Test_PacketItemData(t *testing.T) {
    testPacketSerial(
        t,
        false,
        &PacketItemData{
            ItemTypeId: 10,
            MapId:      3,
            MapData: MapData{
                1, 10,
            },
        },
        te.LiteralString("\x83"+
            "\x00\x0a"+
            "\x00\x03"+
            "\x02"+
            "\x01\x0a"),
    )
}

func Test_PacketObjectSpawn(t *testing.T) {
    // Test for non-fireball data case.
    testPacketSerial(
        t,
        false,
        &PacketObjectSpawn{
            EntityId: 4,
            ObjType:  3,
            Position: AbsIntXyz{1, 2, 3},
        },
        te.LiteralString("\x17"+
            "\x00\x00\x00\x04"+
            "\x03"+
            "\x00\x00\x00\x01"+
            "\x00\x00\x00\x02"+
            "\x00\x00\x00\x03"+
            "\x00\x00\x00\x00"),
    )

    // Test for fireball data case.
    testPacketSerial(
        t,
        false,
        &PacketObjectSpawn{
            EntityId: 4,
            ObjType:  3,
            Position: AbsIntXyz{1, 2, 3},
            Fireball: FireballData{
                ThrowerId: 20,
                X:         4, Y: 5, Z: 6,
            },
        },
        te.LiteralString("\x17"+
            "\x00\x00\x00\x04"+
            "\x03"+
            "\x00\x00\x00\x01"+
            "\x00\x00\x00\x02"+
            "\x00\x00\x00\x03"+
            "\x00\x00\x00\x14"+
            "\x00\x04\x00\x05\x00\x06"),
    )
}

func Test_PacketPluginMessage(t *testing.T) {
    testPacketSerial(
        t,
        false,
        &PacketPluginMessage{
            Channel: "chan",
            Data:    []byte{1, 2, 3, 4, 5},
        },
        te.LiteralString("\xfa"+
            "\x00\x04\x00c\x00h\x00a\x00n"+
            "\x00\x05"+
            "\x01\x02\x03\x04\x05"),
    )
}

func Benchmark_Packet_ReadString16(b *testing.B) {
    var ps PacketSerializer
    data := []byte("\x08\x00u\x00s\x00e\x00r\x00n\x00a\x00m\x00e")

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        input := bytes.NewBuffer(data)
        _, _ = ps.readString16(input)
        input.Reset()
    }
}

func Benchmark_Packet_WriteString16(b *testing.B) {
    output := bytes.NewBuffer(make([]byte, 0, 1024))
    var ps PacketSerializer

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _ = ps.writeString16(output, "username")
        output.Reset()
    }
}

func benchmarkPacket(b *testing.B, pkt IPacket) {
    output := bytes.NewBuffer(make([]byte, 0, 1024))
    var ps PacketSerializer
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = ps.WritePacket(output, pkt)
        output.Reset()
    }
}

func Benchmark_WritePacketLogin(b *testing.B) {
    benchmarkPacket(b, &PacketLogin{
        VersionOrEntityId: 5,
        Username:          "username",
        MapSeed:           123,
        GameMode:          1,
        Dimension:         DimensionNormal,
        Difficulty:        GameDifficultyNormal,
        WorldHeight:       128,
        MaxPlayers:        12,
    })
}

func Benchmark_WritePacketKeepAlive(b *testing.B) {
    benchmarkPacket(b, &PacketKeepAlive{
        Id: 10,
    })
}

func Benchmark_WritePacketEntityMetadata(b *testing.B) {
    benchmarkPacket(b, &PacketEntityMetadata{
        EntityId: 5,
        Metadata: EntityMetadataTable{
            EntityMetadata{0, 0, byte(5)},
        },
    })
}
