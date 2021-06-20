package chromemobile

import (
	"github.com/b97tsk/chrome"
	"github.com/b97tsk/chrome/service/http"
	"github.com/b97tsk/chrome/service/http/goagent"
	"github.com/b97tsk/chrome/service/http/httpfs"
	"github.com/b97tsk/chrome/service/socks"
	"github.com/b97tsk/chrome/service/socks/shadowsocks"
	"github.com/b97tsk/chrome/service/socks/v2socks"
	"github.com/b97tsk/chrome/service/tcptun"
	"github.com/b97tsk/chrome/service/tcptun/dnstun"
	"github.com/b97tsk/chrome/service/v2ray"
)

func newManager() *chrome.Manager {
	var m chrome.Manager

	m.AddService(dnstun.Service{})
	m.AddService(goagent.Service{})
	m.AddService(http.Service{})
	m.AddService(httpfs.Service{})
	m.AddService(shadowsocks.Service{})
	m.AddService(socks.Service{})
	m.AddService(tcptun.Service{})
	m.AddService(v2ray.Service{})
	m.AddService(v2socks.Service{})

	return &m
}
