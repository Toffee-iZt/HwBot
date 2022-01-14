package longpoll

import (
	"context"
	"encoding/json"
	"fmt"
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

func (lp *LongPoll) update() (*vkapi.LongPollServer, error) {
	return lp.vk.GetLongPollServer(lp.vk.Self().ID)
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

func (lp *LongPoll) run(ctx context.Context, ch chan Event, wait int) {
	serv, err := lp.update()
	if err != nil {
		lp.sync.ErrClose(fmt.Errorf("longpoll init server: %w", err))
	}

	builder := vkhttp.NewRequestsBuilder(serv.Server)
	args := vkhttp.Args{
		"act":  "a_check",
		"wait": strconv.Itoa(wait),
		"key":  serv.Key,
		"ts":   serv.Ts,
	}

	for {
		req := builder.Build(args)

		status, body, err := lp.vk.HTTP().DoContext(ctx, req)
		if err != nil || status != vkhttp.StatusOK {
			if err != context.Canceled {
				err = fmt.Errorf("longpoll: %w", err)
			}
			lp.sync.ErrClose(&vkapi.Error{
				Message:    err.Error(),
				HTTPStatus: status,
				Body:       body,
			})
			return
		}

		var res struct {
			Ts      string  `json:"ts"`
			Updates []Event `json:"updates"`
			Failed  int     `json:"failed"`
		}

		err = json.Unmarshal(body, &res)
		if err != nil {
			lp.sync.ErrClose(fmt.Errorf("longpoll json: %w\n%s", err, string(body)))
			return
		}

		switch res.Failed {
		default:
			args["ts"] = res.Ts
		case 2, 3:
			serv, err = lp.update()
			if err != nil {
				lp.sync.ErrClose(fmt.Errorf("longpoll update: %w", err))
				return
			}
			args["key"] = serv.Key
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
