package longpoll

import (
	"context"
	"strconv"

	"github.com/Toffee-iZt/HwBot/common"
	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/vkhttp"
)

// New creates new longpoll instance.
func New(vk *vkapi.Client) *LongPoll {
	return &LongPoll{
		vk: vk,
	}
}

// LongPoll struct.
type LongPoll struct {
	vk   *vkapi.Client
	sync common.Sync
}

// Done returns a channel that closes when the longpoll is finished
// due to an error or context cancellation.
func (lp *LongPoll) Done() <-chan struct{} {
	return lp.sync.Done()
}

// Err returns longpoll error.
func (lp *LongPoll) Err() error {
	return lp.sync.Err()
}

// Run runs longpoll in new goroutine.
func (lp *LongPoll) Run(ctx context.Context, wait int) <-chan Event {
	if !lp.sync.Init() {
		return nil
	}

	ch := make(chan Event)
	go lp.run(ctx, ch, wait)

	return ch
}

func (lp *LongPoll) update() *vkapi.LongPollServer {
	s, err := lp.vk.GetLongPollServer(lp.vk.Self().ID)
	if err != nil {
		panic("longpoll: update error\n" + err.Error())
	}
	return s
}

func (lp *LongPoll) run(ctx context.Context, ch chan Event, wait int) {
	serv := lp.update()
	args := vkhttp.Args{
		"act":  "a_check",
		"wait": strconv.Itoa(wait),
		"key":  serv.Key,
		"ts":   serv.Ts,
	}

	for {
		var res struct {
			Ts      string  `json:"ts"`
			Updates []Event `json:"updates"`
			Failed  int     `json:"failed"`
		}
		ctxerr := lp.vk.Client.LongPoll(ctx, serv.Server, args, &res)
		if ctxerr != nil {
			lp.sync.ErrClose(ctxerr)
			return
		}

		switch res.Failed {
		default:
			args["ts"] = res.Ts
		case 2, 3:
			serv = lp.update()
			if res.Failed == 3 {
				args["ts"] = serv.Ts
			}
		}

		for _, u := range res.Updates {
			select {
			case ch <- u:
			case <-ctx.Done():
				lp.sync.ErrClose(ctx.Err())
				return
			}
		}
	}
}
