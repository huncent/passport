package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/controllers"
)

type Passport struct {
	ServDomain string
}

func (p *Passport) UserAdd(cellphone, email, pwd string) ([]byte, error) {
	scellphone, semail := strings.TrimSpace(cellphone), strings.TrimSpace(email)
	if scellphone == "" && semail == "" {
		return nil, fmt.Errorf("用户手机号码和邮箱地址同时为空.")
	}

	body, _ := json.Marshal(controllers.UserAdd{semail, scellphone, ""})
	return common.PostRequest(p.ServDomain+"/user/add", body, nil)
}
