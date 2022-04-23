package longpoll

import (
	"context"
	"net/http"

	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/vkargs"
)

// Run runs longpoll listener and returns events channel.
// Channel will be closed when context is done.
func Run(ctx context.Context, vk *vkapi.Client, wait int) <-chan *Event {
	ch := make(chan *Event)
	go run(ctx, vk, wait, ch)
	return ch
}

type response struct {
	Updates []Event `json:"updates"`
	Ts      string  `json:"ts"`
	Failed  int     `json:"failed"`
}

func run(ctx context.Context, vk *vkapi.Client, wait int, ch chan<- *Event) {
	defer close(ch)

	args := vkapi.ArgsMap{
		"act":  "a_check",
		"wait": wait,
	}
	server := update(vk, args, true)

	for {
		var res response
		url := server + "?" + vkargs.Marshal(args).Encode()
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		status, canceled := vk.Do(req, &res)
		if canceled {
			return
		}

		if status != 200 {
			// TODO
		}

		switch res.Failed {
		default:
			args["ts"] = res.Ts
		case 2, 3:
			update(vk, args, res.Failed == 3)
		}

		for i := range res.Updates {
			select {
			case ch <- &res.Updates[i]:
			case <-ctx.Done():
				return
			}
		}
	}
}

func update(vk *vkapi.Client, args vkapi.ArgsMap, ts bool) string {
	s, err := vk.GetLongPollServer(vk.Self().ID)
	if err != nil {
		panic("longpoll update error\n" + err.String())
	}
	args["key"] = s.Key
	if ts {
		args["ts"] = s.Ts
	}
	return s.Server
}
