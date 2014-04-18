package main

import (
	"fmt"
	"github.com/KingsEpic/kinglib"
	"gopkg.in/v0/qml"
	// "math"
)

type TilePosition struct {
	X, Y int
}

type ChunkIdentifier struct {
	ID    int
	Layer string
}

type World struct {
	Archetypes map[string]*Archetype
	Layers     map[string]*Layer
	Entities   map[int]*Entity
	Jobs       map[int]*kinglib.JobStatus
	Diaspora   map[int]*Entity
	Ready      bool
	Player     int // Entity ID of player
}

type Archetype struct {
	Name              string
	SimpleName        string
	Class             int
	Craftable         bool
	Buildable         bool
	CraftRequirements map[*Archetype]int
}

type Entity struct {
	EntityID       int
	FloatX         float32
	FloatY         float32
	Layer          string
	SimpleName     string
	CurrentChunk   *Chunk
	ContainerID    int
	ContainerIndex int
	Quantity       int
	Inventory      Inventory
	Visible        bool
	QObject        qml.Object
}

type Inventory struct {
	Slots []*Entity
}

type ArchetypeInstance struct {
	TP         TilePosition
	SimpleName string
	QObject    qml.Object
}

type Layer struct {
	Name   string
	Chunks map[int]*Chunk
}

func (a *Archetype) Init() {
	a.CraftRequirements = make(map[*Archetype]int)
}

func (l *Layer) Init() {
	l.Chunks = make(map[int]*Chunk)
}

type Chunk struct {
	ChunkID            int
	ArchetypeInstances map[TilePosition][]*ArchetypeInstance
	Entities           map[int]*Entity
}

func WorldReady() bool {
	if w != nil && w.Ready {
		return true
	}
	return false
}

func (a *ArchetypeInstance) SetQML(obj qml.Object) {
	a.QObject = obj
	obj.Set("simpleName", a.SimpleName)
	obj.Set("name", w.Archetypes[a.SimpleName].Name)
	qmlMoveUpdate(obj, float32(a.TP.X), float32(a.TP.Y), a.SimpleName, kinglib.MOVE_INSTANT, 0)

}

func (e *Entity) SetVisible(vis bool) {
	if vis != e.Visible {
		e.Visible = vis
		e.UpdateVisibility()
	}
}

func (e *Entity) UpdateVisibility() {
	if e.QObject != nil {
		e.QObject.Set("visible", e.Visible)
	}
}

func (e *Entity) SetQML(obj qml.Object) {
	e.QObject = obj
	obj.Set("simpleName", e.SimpleName)
	obj.Set("name", w.Archetypes[e.SimpleName].Name)
	obj.Set("entityID", e.EntityID)
	qmlMoveUpdate(obj, e.FloatX, e.FloatY, e.SimpleName, kinglib.MOVE_INSTANT, 0)
	e.UpdateVisibility()
}

func (e *Entity) SetContainer(container_id int, index int) {
	// Sometimes, an entity will be contained in another when the container is not yet known to
	// the client.  So we need to handle those cases.
	// fmt.Printf("Setting container_id %d index %d on %+v\n", container_id, index, e)

	// Delete from existing slot, if needed:
	if e.ContainerID > 0 {
		current_container := w.Entities[e.ContainerID]
		if current_container.Inventory.Slots[e.ContainerIndex] == e {
			current_container.Inventory.Slots[e.ContainerIndex] = nil
		}
	}

	if container_id == 0 {
		e.SetVisible(true)

		// Let's no longer wait to add this to container now that it has none:
		if w.Diaspora[e.EntityID] != nil {
			delete(w.Diaspora, e.EntityID)
		}

	} else {
		e.SetVisible(false)
		// Returns current container
		container := w.Entities[container_id]
		// fmt.Printf("* Container: %d, index: %d, entity: %+v\n", e, container_id, index)

		if container != nil {
			// fmt.Printf("  Adding to inventory since object found")
			container.Inventory.Slots[index] = e
		} else {
			// Not at all the best way to store I think, but I'm lazy at the moment
			// fmt.Printf("  Saving for later, since not found")
			w.Diaspora[e.EntityID] = e
		}

		// fmt.Printf("Container status: %+v\n", container)
	}

	e.ContainerID = container_id
	e.ContainerIndex = index

}

func qmlMoveUpdate(obj qml.Object, x, y float32, simple_name string, move_type int, move_time int) {
	// Not good for now:
	obj.Set("simpleName", simple_name)

	switch move_type {
	case kinglib.MOVE_WALK:
		// fmt.Printf("Setting move_walk\n")
		obj.Set("dsx", int(x*32.0))
		obj.Set("dsy", int(y*32.0))
		obj.Set("dtime", move_time)
		obj.Set("moving", false)
		obj.Set("moving", true)
	default:
		obj.Set("x", int(x*32.0))
		obj.Set("y", int(y*32.0))
	}

	z := GetArchetypeZ(simple_name, int(y), int(x))
	obj.Set("z", z)

	obj.Set("tx", x)
	obj.Set("ty", y)

}

func (c *Chunk) Init() {
	c.ArchetypeInstances = make(map[TilePosition][]*ArchetypeInstance)
	c.Entities = make(map[int]*Entity)
}

func (c *Chunk) AddEntity(e *Entity) {
	c.Entities[e.EntityID] = e
	e.CurrentChunk = c
}

func (c *Chunk) RemoveEntity(e *Entity) {
	delete(c.Entities, e.EntityID)
	e.CurrentChunk = nil
}

func (w *World) Init() {
	w.Archetypes = make(map[string]*Archetype)
	w.Layers = make(map[string]*Layer)
	w.Entities = make(map[int]*Entity)
	w.Jobs = make(map[int]*kinglib.JobStatus)
	w.Diaspora = make(map[int]*Entity)

	w.Ready = true
}

func (w *World) UpdateJob(js kinglib.JobStatus) {

	fmt.Printf("Updating job\n")

	w.Jobs[js.ID] = &js

	if js.Finished {
		// fmt.Printf("Job finished.  Deleting j.ID: %d\n", js.ID)

		delete(w.Jobs, js.ID)
	}

	job_model.Update()

}

func (w *World) AddCraftableRequirement(cr *kinglib.CraftableRequirement) {
	a := w.Archetypes[cr.ArchetypeName]
	if a != nil {
		ra := w.Archetypes[cr.RequirementName]
		if ra == nil {
			fmt.Printf("WARNING!  Could not find craft requrement: %+v\n", cr)
		} else {
			a.CraftRequirements[ra] = cr.Quantity
			craft_model.UpdateArchetype(a)
		}
	} else {
		fmt.Printf("WARNING!  Could not find archetype to add requirement: %+v\n", cr)
	}
}

func (w *World) UpdateEntity(eu *kinglib.EntityUpdate) *Entity {
	var e *Entity
	update_inv_model := false

	if w.Entities[eu.EntityID] == nil {
		e = &Entity{
			EntityID:   eu.EntityID,
			FloatX:     eu.FloatX,
			FloatY:     eu.FloatY,
			SimpleName: eu.SimpleName,
			// ContainerID:    eu.ContainerID,
			// ContainerIndex: eu.ContainerIndex,
			Quantity:  eu.Quantity,
			Visible:   true,
			Inventory: Inventory{Slots: make([]*Entity, eu.InventorySize)},
		}
		w.Entities[eu.EntityID] = e
		if e.EntityID == w.Player {
			w.SetupPlayer()
		}
	} else {
		e = w.Entities[eu.EntityID]

		// Check before we change container ID:
		if e.ContainerID == w.Player {
			update_inv_model = true
		}
	}

	// fmt.Printf("* Entity: %+v, update: %+v\n", e, eu)

	// <-- Check and if needed update which chunk entity is in -->
	w.UpdateEntityChunk(e, ChunkIdentifier{eu.ChunkID, eu.Layer})

	e.FloatX = eu.FloatX
	e.FloatY = eu.FloatY
	if e.ContainerID != eu.ContainerID || e.ContainerIndex != eu.ContainerIndex {
		e.SetContainer(eu.ContainerID, eu.ContainerIndex)
	}

	// Terrible way to search.  I should have a map from container to containees, not other way!
	for _, ent := range w.Diaspora {
		if ent.ContainerID == e.EntityID {
			delete(w.Diaspora, ent.EntityID)
			ent.SetContainer(ent.ContainerID, ent.ContainerIndex)
		}
	}
	// e.ContainerIndex = eu.ContainerIndex
	e.Quantity = eu.Quantity

	// And check again after container ID changed:
	if e.ContainerID == w.Player {
		update_inv_model = true
	}

	if update_inv_model {
		// fmt.Printf("Updating inv_model\n")
		inv_model.SetContainer(w.Player)
	}

	return e

}

func (w *World) UpdateEntityChunk(e *Entity, cid ChunkIdentifier) {
	if e.CurrentChunk == nil {
		c := w.GetChunk(cid)
		c.AddEntity(e)
	}

	if e.CurrentChunk.ChunkID != cid.ID {
		e.CurrentChunk.RemoveEntity(e)
		c := w.GetChunk(cid)
		c.AddEntity(e)
	}
}

func (w *World) EntityDelete(ed *kinglib.EntityDelete) {
	e := w.Entities[ed.EntityID]

	w.entity_delete(e)
}

func (w *World) entity_delete(e *Entity) {
	if e != nil {
		if e.CurrentChunk != nil {
			e.CurrentChunk.RemoveEntity(e)
		}

		delete(w.Entities, e.EntityID)
	}

	// fmt.Printf("Destroying QML for entity %d\n", e.EntityID)
	e.QObject.Destroy()
}

func (w *World) MoveEntity(m *kinglib.Move) {
	e := w.Entities[m.EntityID]
	if e != nil && e.QObject != nil {
		// old_x := float32(e.QObject.Float64("x"))
		// old_y := float32(e.QObject.Float64("y"))
		// new_x := m.FloatX * 32.0
		// new_y := m.FloatY * 32.0
		// if (float32(math.Abs(float64(old_x-new_x))) > 16.0) || (float32(math.Abs(float64(old_y-new_y))) > 16.0) {
		// 	fmt.Printf("Large move.  old x,y: %0.2f, %0.2f, new x,y: %0.2f, %0.2f.  Move obj: %+v\n", old_x, old_y, new_x, new_y, m)
		// }
		// e.QObject.Set("x", new_x)
		// e.QObject.Set("y", new_y)
		qmlMoveUpdate(e.QObject, m.FloatX, m.FloatY, e.SimpleName, m.MoveType, int(m.Duration))
	} else {
		if e == nil {
			fmt.Printf("ERROR: Could not move entity because entity not found.  %+v\n", m)
		} else {
			fmt.Printf("ERROR: Could not find qml object for entity.  %+v\n", m)
		}
	}
}

func (w *World) InstantiateTile(it *kinglib.TileInstantiation) {
	// Delete all archetype instances for this tile.  That's it.  Maybe slow process for now:
	c := w.GetChunk(ChunkIdentifier{it.ChunkID, it.Layer})
	p := TilePosition{it.X, it.Y}
	// First count size of array:
	for _, ai := range c.ArchetypeInstances[p] {
		// Do what's needed to clean up these instances
		ai.QObject.Destroy()
	}

	// for _, l := range w.Layers {
	// 	for _, c := range l.Chunks {
	// 		for _, i := range c.ArchetypeInstances {
	// 			for _, a := range i {
	// 				if a.QObject != nil {
	// 					fmt.Printf("Destroying qobject\n")
	// 					a.QObject.Destroy()
	// 				} else {
	// 					fmt.Printf("QOobject is nil!\n")
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	c.ArchetypeInstances[p] = make([]*ArchetypeInstance, 0)
}

func (w *World) AddArchetype(a *Archetype) {
	w.Archetypes[a.SimpleName] = a
	if a.Craftable {
		craft_model.AddArchetype(a)
	}
}

func (w *World) AddArchetypeInstance(ai *kinglib.ArchetypeInstance) *ArchetypeInstance {
	c := w.GetChunk(ChunkIdentifier{ai.ChunkID, ai.Layer})
	p := TilePosition{ai.X, ai.Y}

	a := &ArchetypeInstance{TP: p, SimpleName: ai.SimpleName}

	c.ArchetypeInstances[p] = append(c.ArchetypeInstances[p], a)
	return a
}

func (w *World) UnloadChunk(ai *kinglib.UnloadChunk) {
	c := w.GetChunk(ChunkIdentifier{ai.ChunkID, ai.Layer})

	for _, i := range c.ArchetypeInstances {
		for _, a := range i {
			a.QObject.Destroy()
		}
	}

	for _, e := range c.Entities {
		w.entity_delete(e)
	}

	// Now separate this chunk from its layer so that it can be cleaned as part of GC...hopefully
	l := w.Layers[ai.Layer]
	if l != nil {
		delete(l.Chunks, ai.ChunkID)
	}
}

func (w *World) SetPlayer(id int) {
	w.Player = id
	w.SetupPlayer()
}

func (w *World) SetupPlayer() {
	// Called when we have the entity for the player, or when that changes
	fmt.Printf("Calling setup player\n")
	e := w.Entities[w.Player]
	if e != nil {
		inv_model.SetContainer(e.EntityID)
	}
}

func (w *World) GetChunk(cid ChunkIdentifier) (c *Chunk) {
	if w.Layers[cid.Layer] == nil {
		w.Layers[cid.Layer] = &Layer{Name: cid.Layer, Chunks: make(map[int]*Chunk)}
	}

	if w.Layers[cid.Layer].Chunks[cid.ID] == nil {
		c = &Chunk{ChunkID: cid.ID}
		c.Init()
		w.Layers[cid.Layer].Chunks[cid.ID] = c
	} else {
		c = w.Layers[cid.Layer].Chunks[cid.ID]
	}
	return
}

func GetArchetypeZ(simple_name string, y int, x int) int {
	y += 1 // Don't want to multiply by 0 for y = 0

	switch w.Archetypes[simple_name].Class {
	case 0:
		// foundation
		return 0
	case 1:
		// creature
		return 10 + (y * 100) + (x * 100)
	case 2:
		// block
		return 15 + (y * 100) + (x * 100)
	case 3:
		// ground covering
		return 1
	case 4:
		// Trees
		return 11 + (y * 100) + (x * 100)
	case 5:
		// Lower lying ground objects -- campfire, beds, tables, chairs, etc:
		return 8 + (y * 100) + (x * 100)
	default:
		return 10 + (y * 100) + (x * 100)
	}
}
