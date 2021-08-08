package longpoll

import (
	"HwBot/common"
	"HwBot/shttp"
	"HwBot/vkapi"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	vk   *vkapi.Client
	serv *vkapi.LongPollServer
	ch   chan Event
	wait string
	sync common.Sync
}

func (lp *LongPoll) update(updateTS bool) error {
	s, err := lp.vk.Groups.GetLongPollServer()
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

	lp.ch = make(chan Event)
	go lp.run(ctx)

	return lp.ch
}

type response struct {
	Ts      string  `json:"ts"`
	Updates []Event `json:"updates"`
	Failed  int     `json:"failed"`
}

func (lp *LongPoll) run(ctx context.Context) {
	uri := shttp.NewURIBuilder(lp.serv.Server)

	args := new(shttp.Query)
	args.Set("act", "a_check")
	args.Set("wait", lp.wait)
	args.Set("key", lp.serv.Key)

	for {
		args.Set("ts", lp.serv.Ts)

		req, resp := shttp.New(shttp.GETStr, uri.Build(args))
		err := lp.vk.HTTP().DoContext(ctx, req, resp)
		body := resp.Body()
		shttp.Release(req, resp)
		if err != nil {
			if err != context.Canceled {
				err = fmt.Errorf("longpoll: %w", err)
			}
			lp.sync.ErrClose(err)
			return
		}

		var res response
		err = json.Unmarshal(body, &res)
		if err != nil {
			lp.sync.ErrClose(fmt.Errorf("longpoll json: %w\n%s", err, string(body)))
			return
		}

		switch res.Failed {
		case 0, 1:
			lp.serv.Ts = res.Ts
		case 2:
			err = lp.update(false)
			args.Set("key", lp.serv.Key)
		case 3:
			err = lp.update(true)
			args.Set("key", lp.serv.Key)
		default:
			panic("BUG: failed has unexpected value\n" + string(body))
		}
		if err != nil {
			lp.sync.ErrClose(fmt.Errorf("longpoll update: %w", err))
			return
		}

		for _, u := range res.Updates {
			select {
			case lp.ch <- u:
			case <-ctx.Done():
				lp.sync.ErrClose(ctx.Err())
				return
			}
		}
	}
}
