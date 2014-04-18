package main

import (
	"gopkg.in/v0/qml"
)

type CraftReq struct {
	Archetype *Archetype
	Quantity  int
}

type CraftElement struct {
	Archetype *Archetype
	Ready     bool // True if the player has everything in place to craft this right now
	reqs      []CraftReq
	Len       int
}

type CraftModel struct {
	Len        int
	Craftables []*CraftElement
}

// func (ce *CraftElement) Requirements() int {
// 	// TODO: This is less than ideal to update every time this is called rather than just when a requirement is added:
// 	ce.update_requirements()
// 	return len(ce.reqs)
// }

func (ce *CraftElement) Get(index int) CraftReq {
	return ce.reqs[index]
}

func (ce *CraftElement) UpdateRequirements() {
	ce.reqs = make([]CraftReq, len(ce.Archetype.CraftRequirements))
	count := 0
	for na, quantity := range ce.Archetype.CraftRequirements {
		ce.reqs[count] = CraftReq{na, quantity}
		count++
	}

	ce.Len = len(ce.reqs)
	qml.Changed(ce, &ce.Len)
}

func (cm *CraftModel) AddArchetype(a *Archetype) {
	qml.Lock()

	cm.Craftables = append(cm.Craftables, &CraftElement{Archetype: a, Ready: false})
	cm.Len = len(cm.Craftables)

	qml.Changed(cm, &cm.Len)
	qml.Unlock()
}

func (cm *CraftModel) UpdateArchetype(a *Archetype) {
	for _, ce := range cm.Craftables {
		if ce.Archetype == a {
			ce.UpdateRequirements()
		}
	}
}

func (cm *CraftModel) Get(index int) *CraftElement {
	return cm.Craftables[index]
}
