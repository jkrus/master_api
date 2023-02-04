package ping_controller

import (
	"net/http"
)

func (p *PingController) Ping(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}
	msg := r.Form.Get("message")

	res, err := p.bl.Ping.Ping.Ping(r.Context(), msg)
	if err != nil {
		return nil, err
	}

	return res, nil
}
