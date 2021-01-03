package chromemobile

import (
	"github.com/b97tsk/chrome/service"
	"github.com/b97tsk/chrome/service/http"
	"github.com/b97tsk/chrome/service/http/goagent"
	"github.com/b97tsk/chrome/service/http/httpfs"
	"github.com/b97tsk/chrome/service/socks"
	"github.com/b97tsk/chrome/service/socks/shadowsocks"
	"github.com/b97tsk/chrome/service/socks/v2ray"
	"github.com/b97tsk/chrome/service/tcptun"
)

func newManager() *service.Manager {
	man := service.NewManager()
	man.Add(goagent.Service{})
	man.Add(http.Service{})
	man.Add(httpfs.Service{})
	man.Add(shadowsocks.Service{})
	man.Add(socks.Service{})
	man.Add(tcptun.Service{})
	man.Add(v2ray.Service{})

	return man
}
