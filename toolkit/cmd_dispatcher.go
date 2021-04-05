package toolkit

import (
	"fmt"
)

type CmdHandler func([]string)
type HandlerInfo struct {
	Handler CmdHandler
	Info    string
	Min     int
	Max     int
}

type CmdDispatcher struct {
	mux map[string]*HandlerInfo
}

func NewCmdDispatcher() *CmdDispatcher {
	return &CmdDispatcher{
		mux: make(map[string]*HandlerInfo),
	}
}

func (t *CmdDispatcher) Register(key string, handler CmdHandler, min, max int, info string) {
	t.mux[key] = &HandlerInfo{
		Handler: handler,
		Info:    info,
		Min:     min,
		Max:     max}
}

func (t *CmdDispatcher) Process(args []string) {
	if len(args) > 0 {
		if h, pres := t.mux[args[0]]; pres && h != nil {
			an := len(args) - 1
			if an >= h.Min && an <= h.Max {
				h.Handler(args)
				return
			} else {
				fmt.Println("invalid number of args")

			}
		} else {
			fmt.Println("invalid command")

		}
	} else {
		fmt.Println("no args is found")
	}
	t.PrintHelp()
}

func (t *CmdDispatcher) PrintHelp() {

	for k, v := range t.mux {
		fmt.Println(k, "\t", v.Info)
	}
}
