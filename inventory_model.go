package main

import (
	"github.com/KingsEpic/kinglib"
	// "fmt"
	"gopkg.in/v0/qml"
)

// This is currently not implemented well.  When a slot changes significantly, I cause the whole inventory to redraw instead of marking the specific item to redraw

type InventorySlot struct {
	Index  int
	ID     int
	Entity *Entity
}

type InventoryModel struct {
	Size   int
	Entity *Entity
	Slots  []*InventorySlot
}

func (is *InventorySlot) ImageName() string {
	if is.Entity == nil {
		return "blank32"
	} else {
		return is.Entity.SimpleName
	}
}

func (is *InventorySlot) Quantity() int {
	if is.Entity == nil {
		return 0
	} else {
		return is.Entity.Quantity
	}
}

func (im *InventoryModel) SetContainer(entity_id int) {
	// This function should only be called once whenever the vessel is set.  Vessel is the entity which represents the player
	qml.Lock()
	im.Entity = w.Entities[entity_id]
	if im.Entity != nil {
		// These first two lines ensure that a complete redraw is performed:
		im.Size -= 1
		qml.Changed(im, &im.Size)

		im.Size = len(im.Entity.Inventory.Slots)
		im.Slots = make([]*InventorySlot, im.Size)

		for i, _ := range im.Slots {
			im.Slots[i] = &InventorySlot{}
		}
	} else {
		im.Size = 0
		im.Slots = make([]*InventorySlot, 0)
	}
	qml.Changed(im, &im.Size)
	qml.Unlock()

	build_model.Update(entity_id)
}

func (im *InventoryModel) SlotUpdated(index int) {
	// For now we ignore that a slot is updated and just reset the whole thing.  Lazy, and slower, I know!
	qml.Lock()
	im.Size -= 1
	qml.Changed(im, &im.Size)
	im.Size += 1
	qml.Changed(im, &im.Size)
	qml.Unlock()

	build_model.Update(im.Entity.EntityID)
}

func (im *InventoryModel) ImageName(index int) string {
	if index < im.Size {
		is := im.Get(index)
		return is.ImageName()
	} else {
		return "blank32"
	}
}

func (im *InventoryModel) Get(index int) *InventorySlot {
	slot := im.Slots[index]

	e := im.Entity.Inventory.Slots[index]

	var the_id int
	if e != nil {
		the_id = e.EntityID
	}
	slot.Index = index
	slot.ID = the_id
	slot.Entity = e

	return slot
}

func (im *InventoryModel) Swap(index int, destination int) {
	connection.SendGob(kinglib.InventoryMove{index, destination})
}
