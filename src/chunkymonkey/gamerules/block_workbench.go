package gamerules

import (
    . "chunkymonkey/types"
)

func makeWorkbenchAspect() (aspect IBlockAspect) {
    return &InventoryAspect{
        name:                 "Workbench",
        createBlockInventory: createWorkbenchInventory,
    }
}

// Creates a new tile entity for a chest. UnmarshalNbt and SetChunk must be
// called before any other methods.
func NewWorkbenchTileEntity() ITileEntity {
    return createWorkbenchInventory(nil)
}

func createWorkbenchInventory(instance *BlockInstance) *blockInventory {
    inv := NewWorkbenchInventory()
    return newBlockInventory(
        instance,
        inv,
        true,
        InvTypeIdWorkbench,
    )
}
