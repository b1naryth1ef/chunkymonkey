// The window package handles windows for inventories.
package window

import (
    "chunkymonkey/gamerules"
    "chunkymonkey/proto"
    . "chunkymonkey/types"
)

// IInventory is the interface that windows require of inventories.
type IInventory interface {
    NumSlots() SlotId
    Click(click *gamerules.Click) (txState TxState)
    SetSubscriber(subscriber gamerules.IInventorySubscriber)
    GetProtoSlots(slots proto.ItemSlotSlice)
}

// IWindow is the interface on to types that represent a view on to multiple
// inventories.
type IWindow interface {
    WindowId() WindowId
    Click(click *gamerules.Click) (txState TxState)
    PacketWindowOpen() *proto.PacketWindowOpen
    PacketWindowItems() *proto.PacketWindowItems
    Finalize(sendClosePacket bool)
}

// IWindowViewer is the required interface of types that wish to receive packet
// updates from changes to inventories viewed inside a window. Typically
// *player.Player implements this.
type IWindowViewer interface {
    SendPacket(packet proto.IPacket)
}

// inventoryView provides a single mapping between a window view onto an
// inventory at a particular slot range inside the window.
type inventoryView struct {
    window    *Window
    inventory IInventory
    startSlot SlotId
    endSlot   SlotId
}

func (iv *inventoryView) Init(window *Window, inventory IInventory, startSlot SlotId, endSlot SlotId) {
    iv.window = window
    iv.inventory = inventory
    iv.startSlot = startSlot
    iv.endSlot = endSlot
    iv.inventory.SetSubscriber(iv)
}

func (iv *inventoryView) Resubscribe() {
    iv.inventory.SetSubscriber(iv)
}

func (iv *inventoryView) Finalize() {
    iv.inventory.SetSubscriber(nil)
}

// Implementing IInventorySubscriber - relays inventory changes to the viewer
// of the window.
func (iv *inventoryView) SlotUpdate(slot *gamerules.Slot, slotIndex SlotId) {
    iv.window.viewer.SendPacket(slot.UpdatePacket(iv.window.windowId, iv.startSlot+slotIndex))
}

func (iv *inventoryView) ProgressUpdate(prgBarId PrgBarId, value PrgBarValue) {
    iv.window.viewer.SendPacket(&proto.PacketWindowProgressBar{
        WindowId: iv.window.windowId,
        PrgBarId: prgBarId,
        Value:    value,
    })
}

// Window represents the common base behaviour of an inventory window. It acts
// as a view onto multiple Inventories.
type Window struct {
    windowId  WindowId
    invTypeId InvTypeId
    viewer    IWindowViewer
    views     []inventoryView
    title     string
    numSlots  SlotId
}

// NewWindow creates a Window as a view onto the given inventories.
func NewWindow(windowId WindowId, invTypeId InvTypeId, viewer IWindowViewer, title string, inventories ...IInventory) (w *Window) {
    w = &Window{}
    w.Init(windowId, invTypeId, viewer, title, inventories...)
    return
}

// Init is the same as NewWindow, but allows for direct embedding of the Window
// type.
func (w *Window) Init(windowId WindowId, invTypeId InvTypeId, viewer IWindowViewer, title string, inventories ...IInventory) {
    w.windowId = windowId
    w.invTypeId = invTypeId
    w.viewer = viewer
    w.title = title

    w.views = make([]inventoryView, len(inventories))
    startSlot := SlotId(0)
    for index, inv := range inventories {
        endSlot := startSlot + inv.NumSlots()
        w.views[index].Init(w, inv, startSlot, endSlot)
        startSlot = endSlot
    }
    w.numSlots = startSlot

    return
}

func (w *Window) WindowId() WindowId {
    return w.windowId
}

// Finalize cleans up, subscriber information so that the window can be
// properly garbage collected. This should be called when the window is thrown
// away.
func (w *Window) Finalize(sendClosePacket bool) {
    for index := range w.views {
        w.views[index].Finalize()
    }
    if sendClosePacket {
        w.viewer.SendPacket(&proto.PacketWindowClose{w.windowId})
    }
}

// PacketWindowOpen creates a packet describing the window to the writer.
func (w *Window) PacketWindowOpen() *proto.PacketWindowOpen {
    return &proto.PacketWindowOpen{
        WindowId:  w.windowId,
        Inventory: w.invTypeId,
        Title:     w.title,
        // Note that the window size is the number of slots in the first inventory,
        // not including the player inventories.
        NumSlots: byte(w.views[0].inventory.NumSlots()),
    }
}

// WriteWindowItems writes a packet describing the window contents to the
// writer. It assumes that any required locks on the inventories are held.
func (w *Window) PacketWindowItems() *proto.PacketWindowItems {
    items := make(proto.ItemSlotSlice, w.numSlots)

    for i := range w.views {
        view := &w.views[i]
        view.inventory.GetProtoSlots(items[view.startSlot:view.endSlot])
    }

    return &proto.PacketWindowItems{
        WindowId: w.windowId,
        Slots:    items,
    }
}

func (w *Window) Click(click *gamerules.Click) TxState {
    if click.SlotId >= 0 {
        for _, inventoryView := range w.views {

            if click.SlotId >= inventoryView.startSlot && click.SlotId < inventoryView.endSlot {
                invClick := *click
                invClick.SlotId = click.SlotId - inventoryView.startSlot

                result := inventoryView.inventory.Click(&invClick)

                click.Cursor = invClick.Cursor

                return result
            }
        }
    }

    return TxStateRejected
}
