package interfaces

import (
    "io"
    "os"
    "rand"

    "chunkymonkey/entity"
    "chunkymonkey/physics"
    "chunkymonkey/slot"
    .   "chunkymonkey/types"
)

// Subset of player methods used by chunks.
type IChunkSubscriber interface {
    GetEntityId() EntityId
    TransmitPacket(packet []byte)
    // Offers an item to the subscriber. If the subscriber takes it into their
    // inventory, it returns true. This function is called from the item's
    // parent chunk's goroutine, so all methods are safely accessible.
    OfferItem(IItem) (taken bool)
}

type IPlayer interface {
    // Safe to call from outside of player's own goroutine.
    GetEntityId() EntityId
    GetEntity() *entity.Entity // Only the game mainloop may modify the return value
    GetName() string           // Do not modify return value
    LockedGetChunkPosition() *ChunkXz
    TransmitPacket(packet []byte)
    OfferItem(IItem) (taken bool)

    Enqueue(f func(IPlayer))

    // Everything below must be called from within Enqueue

    SendSpawn(writer io.Writer) (err os.Error)
    IsWithin(p1, p2 *ChunkXz) bool
}

type IItem interface {
    // Safe to call from outside of chunk's own goroutine
    GetEntity() *entity.Entity // Only the game mainloop may modify the return value
    GetSlot() *slot.Slot

    // Item methods must be called from the goroutine of their parent chunk.
    // Note that items move between chunks.
    GetPosition() *AbsXyz
    GetItemType() ItemId
    // Note that most code assumes that an item's count is 1.
    GetCount() ItemCount
    SetCount(count ItemCount)
    SendSpawn(writer io.Writer) (err os.Error)
    SendUpdate(writer io.Writer) (err os.Error)
    Tick(blockQuery physics.BlockQueryFn) (leftBlock bool)
}

type IBlockType interface {
    Destroy(chunk IChunk, blockLoc *BlockXyz) bool
    IsSolid() bool
    GetName() string
    GetTransparency() int8
}

type IChunk interface {
    // Safe to call from outside of Enqueue:
    GetLoc() *ChunkXz // Do not modify return value

    Enqueue(f func(IChunk))

    // Everything below must be called from within Enqueue

    // Called from game loop to run physics etc. within the chunk for a single
    // tick.
    Tick()

    // Intended for use by blocks/entities within the chunk.
    GetRand() *rand.Rand
    AddItem(item IItem)
    // Tells the chunk to take posession of the item.
    TransferItem(item IItem)
    GetBlock(subLoc *SubChunkXyz) (blockType BlockId, ok bool)
    DestroyBlock(subLoc *SubChunkXyz) (ok bool)

    // Register subscribers to receive information about the chunk. When added,
    // a subscriber will immediately receive complete chunk information via
    // their TransmitPacket method, and changes thereafter via the same
    // mechanism.
    AddSubscriber(subscriber IChunkSubscriber)
    // Removes a previously registered subscriber to updates from the chunk. If
    // sendPacket is true, then an unload-chunk packet is sent.
    RemoveSubscriber(subscriber IChunkSubscriber, sendPacket bool)

    // Tells the chunk about the position of a player in/near the chunk. pos =
    // nil indicates that the player is no longer nearby.
    SetSubscriberPosition(subscriber IChunkSubscriber, pos *AbsXyz)

    // Get packet data for the chunk
    SendUpdate()
}

type IChunkManager interface {
    // Must currently be called from with the owning IGame's Enqueue:
    Get(loc *ChunkXz) (chunk IChunk)
    ChunksInRadius(loc *ChunkXz) <-chan IChunk
    ChunksActive() <-chan IChunk
}

type IGame interface {
    // Safe to call from outside of Enqueue:
    GetStartPosition() *AbsXyz      // Do not modify return value
    GetChunkManager() IChunkManager // Respect calling methods on the return value within Enqueue
    GetBlockTypes() map[BlockId]IBlockType

    Enqueue(f func(IGame))

    // Everything below must be called from within Enqueue

    AddEntity(entity *entity.Entity)
    AddPlayer(player IPlayer)
    RemovePlayer(player IPlayer)
    MulticastPacket(packet []byte, except interface{})
    SendChatMessage(message string)
}
