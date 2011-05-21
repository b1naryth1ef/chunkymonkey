package shardserver_external

import (
	"rand"

	"chunkymonkey/entity"
	. "chunkymonkey/types"
	"chunkymonkey/interfaces"
)

// ITransmitter is the interface by which shards communicate packets to
// players.
type ITransmitter interface {
	TransmitPacket(packet []byte)
}

// IShardConnection is the interface by which shards can be communicated to by
// player frontend code.
type IShardConnection interface {
	SubscribeChunk(chunkLoc ChunkXz)
	UnsubscribeChunk(chunkLoc ChunkXz)

	// TODO better method to send events to chunks from player frontend.
	Enqueue(fn func())

	// Removes connection to shard, and removes all subscriptions to chunks in
	// the shard. Note that this does *not* send packets to tell the client to
	// unload the subscribed chunks.
	Disconnect()
}

// IShardConnecter is used to look up shards and connect to them.
type IShardConnecter interface {
	// Must currently be called from with the owning IGame's Enqueue:
	ShardConnect(entityId EntityId, player ITransmitter, shardLoc ShardXz) IShardConnection

	// TODO Eventually remove these methods - everything should go through
	// IShardConnection.
	EnqueueAllChunks(fn func(chunk IChunk))
	EnqueueOnChunk(loc ChunkXz, fn func(chunk IChunk))
}

// TODO remove this interface when Enqueue* removed from IShardConnection
type IChunk interface {
	// Safe to call from outside of the shard's goroutine.:
	GetLoc() *ChunkXz // Do not modify return value

	// Everything below must be called from within the containing shard's
	// goroutine.

	// Called from game loop to run physics etc. within the chunk for a single
	// tick.
	Tick()

	// Intended for use by blocks/entities within the chunk.
	GetRand() *rand.Rand
	AddSpawner(spawner entity.ISpawn)
	// Tells the chunk to take posession of the item/mob.
	TransferSpawner(e entity.ISpawn)
	// Tells the chunk to take posession of the item/mob.
	GetBlock(subLoc *SubChunkXyz) (blockType BlockId, ok bool)
	PlayerBlockHit(player interfaces.IPlayer, subLoc *SubChunkXyz, digStatus DigStatus) (ok bool)
	PlayerBlockInteract(player interfaces.IPlayer, target *BlockXyz, againstFace Face)

	// Register players to receive information about the chunk. When added,
	// a player will immediately receive complete chunk information via
	// their TransmitPacket method, and changes thereafter via the same
	// mechanism.
	AddPlayer(entityId EntityId, player ITransmitter)
	// Removes a previously registered player to updates from the chunk. If
	// sendPacket is true, then an unload-chunk packet is sent.
	RemovePlayer(entityId EntityId, sendPacket bool)

	MulticastPlayers(exclude EntityId, packet []byte)

	// Tells the chunk about the position of a player in/near the chunk. pos =
	// nil indicates that the player is no longer nearby.
	SetPlayerPosition(player interfaces.IPlayer, pos *AbsXyz)

	// Get packet data for the chunk
	SendUpdate()
}
