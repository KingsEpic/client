package main

import (
	"encoding/json"
	"fmt"
	"github.com/KingsEpic/kinglib"
	"gopkg.in/v0/qml"
	"log"
	"os"
	"time"
)

var m Manager
var w *World
var win *qml.Window
var context *qml.Context
var quit chan bool = make(chan bool)
var requests chan func() = make(chan func(), 1000)
var packets chan *kinglib.Packet = make(chan *kinglib.Packet, 5000)

var deadcount int

var game *Game

var job_model = &JobModel{}
var inv_model = &InventoryModel{}
var craft_model = &CraftModel{}
var build_model = &BuildModel{}

// var messages chan string = make(chan string, 1000)

type Manager struct {
	State   int
	Address string
}

type Game struct {
	// qml.Object
	CSprite           *CSprite
	Map               qml.Object
	Delegated         bool
	SelectedArchetype string
	WindowZ           int // Stores the highest window z
}

func (g *Game) SetAddress(address string) {
	fmt.Printf("Changing address from %s to %s\n", m.Address, address)
	m.Address = address
}

func (g *Game) SetState(new_state int) {
	fmt.Printf("Attempting to change state from %d to %d\n", m.State, new_state)
	switch new_state {
	case 1:
		// Connect
		err := connect()
		if err != nil {
			log.Printf("Error connecting: %s\n", err)
		} else {
			m.State = new_state
			w = &World{}
			w.Init()

			connection.SendGob(&kinglib.Player{"player1"})

			// game.CSprite.CreateSprite(game.Map)
		}
	}
}

func (g *Game) SetMap(obj qml.Object) {
	g.Map = obj
}

func (g *Game) State() int {
	return m.State
}

func (g *Game) Address() string {
	return m.Address
}

func (g *Game) CreateSprite() qml.Object {
	cs := g.CSprite.CreateSprite(g.Map)

	return cs
}

func (g *Game) SetPlx(n int) {
}

func (g *Game) SetPly(n int) {
}

func (g *Game) NewWindowZ() int {
	g.WindowZ++
	return g.WindowZ
}

func (g *Game) RightClick(obj qml.Object) {
	// Create a move:
	am := &kinglib.ActionMove{obj.Int("tx"), obj.Int("ty")}
	connection.SendGob(am)
}

func (g *Game) CreateJob(name string, delegate bool, attributes *qml.Map) {
	var mp map[string]interface{}
	// fmt.Printf("Attributes: %+v\n", attributes)
	attributes.Convert(&mp)
	jcode, err := json.Marshal(mp)

	if err != nil {
		fmt.Printf("Could not issue command: %s", err)
	} else {
		// fmt.Printf("Creating job %s with attributes: %+v\n", name, jcode)
		j := &kinglib.Job{name, delegate, []byte(jcode)}
		connection.SendGob(j)
	}

}

func (g *Game) DestroyMe(obj qml.Object) {
	obj.Destroy()
}

func (g *Game) Plx() int {
	if WorldReady() && w.Player > 0 {
		e := w.Entities[w.Player]
		if e != nil {
			x := int(e.FloatX * 32.0)
			fmt.Printf("plx: %d\n", x)
			return x
		}
	}

	return 0
}

func (g *Game) Ply() int {
	if WorldReady() && w.Player > 0 {
		e := w.Entities[w.Player]
		if e != nil {
			y := int(e.FloatY * 32.0)
			fmt.Printf("ply: %d\n", y)
			return y
		}
	}
	return 0

}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

GLOOP:
	for {
		// Main game loop
		select {
		case _ = <-quit:
			break GLOOP
		case r := <-requests:
			r()
		case p := <-packets:
			switch p.SubType {
			case 1:
				ai := &kinglib.ArchetypeInstance{}
				connection.DecodeData(p.Data, ai)
				a := w.AddArchetypeInstance(ai)
				go func() {
					// By running this separately, significantly reduces pausing when new chunks load
					q := game.CreateSprite()
					// Need to synchronise access:
					requests <- func() { a.SetQML(q) }
				}()
				// fmt.Printf("AI: %+v\n", ai)
			case 2:
				cu := &kinglib.UnloadChunk{}
				connection.DecodeData(p.Data, cu)
				w.UnloadChunk(cu)
				// fmt.Printf("Chunk unload: %+v\n", cu)
			case 3:
				vs := &kinglib.Vessel{}
				connection.DecodeData(p.Data, vs)
				w.SetPlayer(vs.EntityID)
				new_x := int((win.Int("width") / 2) - 16 - game.Plx())
				new_y := int((win.Int("height") / 2) - 16 - game.Ply())
				fmt.Printf("Map new x, y: %d, %d\n", new_x, new_y)
				game.Map.Set("x", new_x)
				game.Map.Set("y", new_y)
				fmt.Printf("Vessel: %+v\n", vs)
			case 4:
				eu := &kinglib.EntityUpdate{}
				connection.DecodeData(p.Data, eu)
				e := w.UpdateEntity(eu)
				if e.QObject == nil {
					q := game.CreateSprite()
					e.SetQML(q)
				}
				// fmt.Printf("(P_ID: %d) Entity Update: %+v\n", w.Player, eu)
			case 5:
				m := &kinglib.Move{}
				connection.DecodeData(p.Data, m)
				w.MoveEntity(m)
				// fmt.Printf("Move: %+v\n", m)
			case 6:
				ti := &kinglib.TileInstantiation{}
				connection.DecodeData(p.Data, ti)
				w.InstantiateTile(ti)
				fmt.Printf("Tile instantiation: %+v\n", ti)
			case 7:
				ed := &kinglib.EntityDelete{}
				connection.DecodeData(p.Data, ed)
				w.EntityDelete(ed)
				fmt.Printf("Delete entity: %+v\n", ed)
			case 8:
				cr := &kinglib.CraftableRequirement{}
				connection.DecodeData(p.Data, cr)
				w.AddCraftableRequirement(cr)
				// fmt.Printf("Craftable requirement: %+v\n", cr)
			case 9:
				aa := &kinglib.AnimateAttack{}
				connection.DecodeData(p.Data, aa)
				fmt.Printf("Animate Attack: %+v\n", aa)
			case 10:
				at := &kinglib.Archetype{}
				connection.DecodeData(p.Data, at)
				a := &Archetype{Name: at.Name, SimpleName: at.SimpleName, Class: at.Class, Craftable: at.Craftable, Buildable: at.Buildable}
				a.Init()
				w.AddArchetype(a)
				// fmt.Printf("Archetype: %+v\n", at)
			case 14:
				js := kinglib.JobStatus{}
				connection.DecodeData(p.Data, &js)
				w.UpdateJob(js)
				// fmt.Printf("Job Status received: %+v\n", js)
			default:
				deadcount++
				log.Printf("Dead count: %d.  Unknown packet subtype: %s\n", deadcount, p.SubType)
			}
		default:
			time.Sleep(2 * time.Millisecond)
		}

	}
}

func run() error {
	qml.Init(nil)
	engine := qml.NewEngine()
	context = engine.Context()

	// qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{
	// 	Init: func(s *CSprite, obj qml.Object) {
	// 		s.Object = obj
	// 	},
	// }})

	game = &Game{}
	context.SetVar("game", game)

	context.SetVar("jobModel", job_model)
	context.SetVar("invModel", inv_model)
	context.SetVar("craftModel", craft_model)
	context.SetVar("buildModel", build_model)

	component, err := engine.LoadFile("main.qml")
	if err != nil {
		return err
	}

	cspriteComponent, err := engine.LoadFile("CSprite.qml")
	if err != nil {
		return err
	}

	csprite := &CSprite{Component: cspriteComponent}
	game.CSprite = csprite

	win = component.CreateWindow(nil)
	// win.On("widthChanged", func(width int) {
	// 	fmt.Printf("New width: %d\n", width)
	// })
	win.Show()

	// win.Wait()
	go func() {
		win.Wait()
		quit <- true
	}()

	return nil
}

type CSprite struct {
	Component qml.Object
	qml.Object
}

func (s *CSprite) CreateSprite(parent qml.Object) qml.Object {
	ns := s.Component.Create(nil)
	ns.Set("parent", parent)
	ns.Set("x", -2000)
	ns.Set("y", -2000)
	// ns.Set("simple_name", "goblin")

	return ns
}
