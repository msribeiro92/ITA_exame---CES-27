package general

import (
	"strings"
    "strconv"
	"math/rand"
	//"sync"
)

// RequestMessageArgs is the struct that hold data passed to
// RequestMessage RPC calls.
type RequestMessageArgs struct {
	Turn        int
	Me 			int
}

// RequestMessageReply is the struct that hold data returned by
// RequestMessage RPC calls.
type RequestMessageReply struct {
	Answers		map[string]int
}

// RequestMessage is called by other instances of General. It'll write the args received
// in the requestMessageChan.
func (rpc *RPC) RequestMessage(args *RequestMessageArgs, reply *RequestMessageReply) error {
	ans := make(map[string]int)
	for k, v := range rpc.General.Answers {
		if len(k) == args.Turn && !strings.Contains(k, strconv.Itoa(args.Me)){
			if rpc.General.Traitor {
				ans[k + strconv.Itoa(rpc.General.Me)] = rand.Intn(2)
			} else {
				ans[k + strconv.Itoa(rpc.General.Me)] = v
			}
		}
	}
	reply.Answers = ans
	return nil
}

// broadcastRequestMessage will send RequestMessage to all peers
func (general *General) broadcastRequestMessage() {
	//var mu = &sync.Mutex{}
	args := &RequestMessageArgs{
		Turn:		general.CurrentTurn,
		Me: 		general.Me,
	}

	for peerIndex := range general.peers {
		func(peer int) {
			//mu.Lock()
			reply := &RequestMessageReply{}
			ok := general.sendRequestMessage(peer, args, reply)
			if ok {
				general.Mutex.Lock()
				for k, v := range reply.Answers {
					general.Answers[k] = v
				}
				general.Mutex.Unlock()
			}
			//mu.Unlock()
		}(peerIndex)
	}
}

// sendRequestMessage will send RequestMessage to a peer
func (general *General) sendRequestMessage(peerIndex int, args *RequestMessageArgs, reply *RequestMessageReply) bool {
	err := general.CallHost(peerIndex, "RequestMessage", args, reply)
	if err != nil {
		return false
	}
	return true
}
