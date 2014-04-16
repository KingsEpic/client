/*
Build model keeps a list of all the things a player can build.  Ultimately, the list will
be derived from both inventory and stockpiles.
*/

package main

import (
	// "github.com/KingsEpic/kinglib"
	// "fmt"
	"gopkg.in/v0/qml"
)

// This is currently not optimised to be as fast as it could be.  When something changes, everything gets redrawn.  When an index is requested, the object is recreated

type BuildSlot struct {
	Index     int
	Archetype *Archetype
	Quantity  int
}

type BuildModel struct {
	Size   int
	Entity *Entity
	Slots  []*BuildSlot
}

func (bs *BuildSlot) ImageName() string {
	if bs.Archetype == nil {
		return "blank32"
	} else {
		return bs.Archetype.SimpleName
	}
}

func (bm *BuildModel) Update(entity_id int) {
	qml.Lock()
	bm.Entity = w.Entities[entity_id]

	record := make(map[*Archetype]int)

	if bm.Entity != nil {
		// These first two lines ensure that a complete redraw is performed:
		bm.Size -= 1
		qml.Changed(bm, &bm.Size)

		// Check how many inventory items are buildable:
		count := 0
		for _, e := range bm.Entity.Inventory.Slots {
			if e != nil {
				a := w.Archetypes[e.SimpleName]
				if a.Buildable {
					record[a] += e.Quantity
					count++
				}
			}
		}

		bm.Size = len(record)
		bm.Slots = make([]*BuildSlot, len(record))

		count = 0
		for a, q := range record {
			bm.Slots[count] = &BuildSlot{count, a, q}
			count++
		}

		// for _, e := range bm.Entity.Inventory.Slots {
		// 	if e != nil && w.Archetypes[e.SimpleName].Buildable {
		// 		bm.Slots[count] = &BuildSlot{count, w.Archetypes[e.SimpleName], e.Quantity}
		// 		count++
		// 	}
		// }
	} else {
		bm.Size = 0
		bm.Slots = make([]*BuildSlot, 0)
	}

	qml.Changed(bm, &bm.Size)
	qml.Unlock()
}

func (bm *BuildModel) ImageName(index int) string {
	if index < bm.Size {
		is := bm.Get(index)
		return is.ImageName()
	} else {
		return "blank32"
	}
}

func (bm *BuildModel) Get(index int) *BuildSlot {
	return bm.Slots[index]
}
