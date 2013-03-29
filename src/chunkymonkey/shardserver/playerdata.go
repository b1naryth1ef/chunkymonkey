package shardserver

import (
    "chunkymonkey/gamerules"
    "chunkymonkey/proto"
    . "chunkymonkey/types"
)

const (
    // Assumed values for size of player axis-aligned bounding box (AAB).
    playerAabH = AbsCoord(0.75) // Each side of player.
    playerAabY = AbsCoord(2.00) // From player's feet position upwards.
)

// playerData represents a Chunk's knowledge about a player. Only one Chunk has
// this data at a time. This data is occasionally updated from the frontend
// server.
type playerData struct {
    entityId   EntityId
    name       string
    position   AbsXyz
    look       LookBytes
    heldItemId ItemTypeId
    // TODO Armor data.
}

func (player *playerData) SpawnPackets(pkts []proto.IPacket) []proto.IPacket {
    return append(pkts, &proto.PacketNamedEntitySpawn{
        EntityId:    player.entityId,
        Username:    player.name,
        Position:    *player.position.ToAbsIntXyz(),
        Rotation:    player.look,
        CurrentItem: player.heldItemId,
    })
    // TODO Armor packet(s).
}

func (player *playerData) UpdatePackets(pkts []proto.IPacket) []proto.IPacket {
    return append(pkts, &proto.PacketEntityTeleport{
        EntityId: player.entityId,
        Position: *player.position.ToAbsIntXyz(),
        Look:     player.look,
    })
}

func (player *playerData) OverlapsItem(item *gamerules.Item) bool {
    // TODO note that calling this function repeatedly is not as efficient as it
    // could be.

    minX := player.position.X - playerAabH
    maxX := player.position.X + playerAabH
    minZ := player.position.Z - playerAabH
    maxZ := player.position.Z + playerAabH
    minY := player.position.Y
    maxY := player.position.Y + playerAabY

    pos := item.Position()

    return pos.X >= minX && pos.X <= maxX && pos.Y >= minY && pos.Y <= maxY && pos.Z >= minZ && pos.Z <= maxZ
}
