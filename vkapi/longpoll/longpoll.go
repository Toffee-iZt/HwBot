package longpoll

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Toffee-iZt/HwBot/common"
	"github.com/Toffee-iZt/HwBot/shttp"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

// New creates new longpoll instance.
func New(vk *vkapi.Client, wait int) *LongPoll {
	return &LongPoll{
		vk:   vk,
		wait: strconv.Itoa(wait),
	}
}

// LongPoll struct.
type LongPoll struct {
	http shttp.Client
	vk   *vkapi.Client
	serv *vkapi.LongPollServer
	wait string
	sync common.Sync
}

func (lp *LongPoll) update(updateTS bool) error {
	s, err := lp.vk.GetLongPollServer(lp.vk.Self().ID)
	if err != nil {
		return err
	}
	if lp.serv == nil {
		lp.serv = s
		return nil
	}
	lp.serv.Key = s.Key
	if updateTS {
		lp.serv.Ts = s.Ts
	}
	return nil
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
func (lp *LongPoll) Run(ctx context.Context) <-chan Event {
	if !lp.sync.Init() {
		return nil
	}

	err := lp.update(true)
	if err != nil {
		lp.sync.ErrClose(fmt.Errorf("longpoll init server: %w", err))
		return nil
	}

	ch := make(chan Event)
	go lp.run(ctx, ch)

	return ch
}

func (lp *LongPoll) run(ctx context.Context, ch chan Event) {
	builder := shttp.NewRequestsBuilder(lp.serv.Server)

	args := new(shttp.Query)
	args.Set("act", "a_check")
	args.Set("wait", lp.wait)
	args.Set("key", lp.serv.Key)

	for {
		args.Set("ts", lp.serv.Ts)

		req := builder.Build(shttp.GETStr, args)
		status, body, err := lp.vk.HTTP().DoContext(ctx, req)
		if err != nil || status != shttp.StatusOK {
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
			lp.serv.Ts = res.Ts
		case 2, 3:
			err = lp.update(res.Failed == 3)
			if err != nil {
				lp.sync.ErrClose(fmt.Errorf("longpoll update: %w", err))
				return
			}
			args.Set("key", lp.serv.Key)
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
