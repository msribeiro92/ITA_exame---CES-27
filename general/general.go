package general

import (
	"fmt"
	"sync"
    "time"
    "strconv"
    "math/rand"
)

// eneral is the struct that hold all information that is used by this instance
// of general.
type General struct {
	sync.Mutex

	serv *server
	done chan struct{}

	peers map[int]string
	Me        int
    Turns     int
	Processes int
    Traitor   bool

    // Persistent state on all servers:
    CurrentTurn int
    Answers map[string]int
}

// NewGeneral create a new general object and return a pointer to it.
func NewGeneral(peers map[int]string, me int, turns int, processes int, traitor bool, cmdValue int, cmdTraitor bool) *General {
	var err error

	general := &General{
		done: make(chan struct{}),

		peers:      peers,
		Me:         me,
        Turns:      turns,
		Processes: processes,
        Traitor:    traitor,

        CurrentTurn: 0,
		Answers: make(map[string]int),
	}


    if cmdTraitor {
        general.Answers[strconv.Itoa(0)] = rand.Intn(2)
    } else {
        general.Answers[strconv.Itoa(0)] = cmdValue
    }

	general.serv, err = newServer(general, peers[me])
	if err != nil {
		panic(err)
	}

	time.Sleep(30 * time.Second)
    general.gossip()

	for ; general.CurrentTurn > 0; general.CurrentTurn -- {
		mapCount  := make(map[string]float64)
		for k, v := range general.Answers {
			if len(k) == general.CurrentTurn {
				mapCount[k[0: len(k)-1]] = mapCount[k[0: len(k)-1]] + float64(v);
			}
		}
		for k, v := range mapCount {
			mapCount[k] = v / float64(general.Processes)
		}
		for k, v := range mapCount {
			if v > float64(general.Processes) / 2.0 {
				general.Answers[k] = 1
			} else {
				general.Answers[k] = 0
			}
		}
	}
	if general.Answers[strconv.Itoa(0)] == 0 {
		fmt.Println("decisao: Recue!")
	} else {
		fmt.Println("decisao: Ataque!")
	}


	return general
}

// Done returns a channel that will be used when the instance is done.
func (general *General) Done() <-chan struct{} {
	return general.done
}

// All changes to General structure should occur in the context of this routine.
// This way it's not necessary to use synchronizers to protect shared data.
// To send data to each of the states, use the channels provided.
func (general *General) gossip() {

	err := general.serv.startListening()
	if err != nil {
		panic(err)
	}

    for general.CurrentTurn = 1; general.CurrentTurn < general.Turns; general.CurrentTurn++ {
        general.broadcastRequestMessage()
        time.Sleep(60 * time.Second)
    }

	fmt.Println("Mensagnes recebidas por ", general.Me)
	for k, v := range general.Answers {
		fmt.Println("caminho:", k, ", valor:", v)
	}

}
