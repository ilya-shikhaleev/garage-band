package game

import (
	"log"
)

type Room struct {
	name string

	// Registered connections.
	connectedPlayers map[*ConnectedPlayer]bool

	// New action added.
	AddAction chan *ConnectedPlayer

	// Register requests from the connections.
	Join chan *ConnectedPlayer

	// Unregister requests from connections.
	Leave chan *ConnectedPlayer
}

func (r *Room) Name() string {
	return r.name
}

// Run the room in goroutine
func (r *Room) run() {
	for {
		log.Println("run room select")
		select {
		case cp := <-r.Leave:
			log.Println("r.leave catch", cp)
			delete(r.connectedPlayers, cp)
			r.OnPlayerLeave(cp)
		case cp := <-r.Join:
			log.Println("r.join catch")
			r.connectedPlayers[cp] = true
			r.OnPlayerJoined(cp)
		case cp := <-r.AddAction:
			log.Println("r.addAction catch", cp)
			r.OnActionAdded(cp)
		}
	}

}

func (r *Room) OnPlayerJoined(newPlayer *ConnectedPlayer) {
	for cp := range r.connectedPlayers {
		cp.SendMessage("{\"event\":\"join\",\"player\":\"" + newPlayer.name + "\",\"instrument\":\"" + newPlayer.instrument.Name() + "\"}")
	}
}

func (r *Room) OnPlayerLeave(newPlayer *ConnectedPlayer) {
	for cp := range r.connectedPlayers {
		cp.SendMessage("{\"event\":\"leave\",\"player\":\"" + newPlayer.name + "\"}")
	}
}

func (r *Room) OnActionAdded(p *ConnectedPlayer) {
	name := p.Name()
	audio, err := p.PlayedAudio()
	if err != nil {
		log.Println("Bad command from ", name, "command:", p.Action())
		return
	}

	for cp := range r.connectedPlayers {
		cp.SendMessage("{\"event\":\"play\",\"player\":\"" + name + "\",\"audio\":\"" + audio + "\"}")
	}
}

func (r *Room) GetFreeInstruments() map[InstrumentType]string {
	instruments := AvailableInstruments()
	for cp := range r.connectedPlayers {
		delete(instruments, cp.instrument.Type())
	}
	return instruments
}

func NewRoom(name string) *Room {
	r := &Room{
		name:             name,
		connectedPlayers: make(map[*ConnectedPlayer]bool),
		AddAction:        make(chan *ConnectedPlayer),
		Join:             make(chan *ConnectedPlayer),
		Leave:            make(chan *ConnectedPlayer),
	}

	go r.run()
	return r
}
