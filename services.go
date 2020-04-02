package chromemobile

import (
	"github.com/b97tsk/chrome/service"
	"github.com/b97tsk/chrome/service/goagent"
	"github.com/b97tsk/chrome/service/http"
	"github.com/b97tsk/chrome/service/httpfs"
	"github.com/b97tsk/chrome/service/logging"
	"github.com/b97tsk/chrome/service/shadowsocks"
	"github.com/b97tsk/chrome/service/socks"
	"github.com/b97tsk/chrome/service/tcptun"
	"github.com/b97tsk/chrome/service/vmess"
)

func addServices(services *service.Manager) {
	services.Add(goagent.Service{})
	services.Add(http.Service{})
	services.Add(httpfs.Service{})
	services.Add(logging.Service{})
	services.Add(shadowsocks.Service{})
	services.Add(socks.Service{})
	services.Add(tcptun.Service{})
	services.Add(vmess.Service{})
}
