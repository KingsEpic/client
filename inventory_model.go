package main

import (
	"github.com/KingsEpic/kinglib"
	// "fmt"
	"gopkg.in/v0/qml"
)

// This is currently not optimised to be as fast as it could be.  When something changes, everything gets redrawn.  When an index is requested, the object is recreated

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
	qml.Lock()
	im.Entity = w.Entities[entity_id]
	if im.Entity != nil {
		// These first two lines ensure that a complete redraw is performed:
		im.Size -= 1
		qml.Changed(im, &im.Size)

		im.Size = len(im.Entity.Inventory.Slots)
		im.Slots = make([]*InventorySlot, im.Size)
	} else {
		im.Size = 0
		im.Slots = make([]*InventorySlot, 0)
	}
	qml.Changed(im, &im.Size)
	qml.Unlock()

	build_model.Update(entity_id)
}

// func (im *InventoryModel) SlotUpdated(index int) {
// 	qml.Changed(im, im[index])
// }

func (im *InventoryModel) ImageName(index int) string {
	if index < im.Size {
		is := im.Get(index)
		return is.ImageName()
	} else {
		return "blank32"
	}
}

func (im *InventoryModel) Get(index int) *InventorySlot {
	e := im.Entity.Inventory.Slots[index]
	var the_id int
	if e != nil {
		the_id = e.EntityID
	}

	is := &InventorySlot{index, the_id, e}
	// fmt.Printf("Index %d requested: %+v\n", index, is)

	return is
}

func (im *InventoryModel) Swap(index int, destination int) {
	connection.SendGob(kinglib.InventoryMove{index, destination})
}
