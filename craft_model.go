package main

import (
	"gopkg.in/v0/qml"
)

type CraftElement struct {
	Archetype *Archetype
	Ready     bool // True if the player has everything in place to craft this right now
}

type CraftModel struct {
	Len        int
	Craftables []*CraftElement
}

// func (ce *CraftElement) Requirements() []*Archetype {
// 	reqs := make([]*Archetype, len(ce.Archetype.CraftRequirements))
// 	count := 0
// 	for a, _ := range ce.Archetype.CraftRequirements {
// 		reqs[count] = a
// 		count++
// 	}
// 	return reqs
// }

func (cm *CraftModel) AddArchetype(a *Archetype) {
	qml.Lock()

	cm.Craftables = append(cm.Craftables, &CraftElement{a, false})
	cm.Len = len(cm.Craftables)

	qml.Changed(cm, &cm.Len)
	qml.Unlock()
}

func (cm *CraftModel) Get(index int) *CraftElement {
	return cm.Craftables[index]
}
