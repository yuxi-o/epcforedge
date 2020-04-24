// CertusNet Copyright  
package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func getAllQoSSubscriptions(cliCtx context.Context, afCtx *Context) (
	[]AsSessionWithOoSSub, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	tSubs, resp, err := cli.QoSSubGetAllAPI.QoSSubscriptionsGetAll(
		cliCtx, afCtx.cfg.AfID)

	if err != nil {
		return nil, resp, err
	}
	return tSubs, resp, nil

}

//GetAllQoSSubscriptions function
func GetAllQoSSubscriptions(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		tsResp     []AsSessionWithQoSSub
		resp       *http.Response
		tsRespJSON []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("AsSessionWithQoS Subscription get all: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	tsResp, resp, err = getAllQoSSubscriptions(cliCtx, afCtx)
	if err != nil {
		log.Errf("AsSessionWithQoS Subscriptions get all : %s", err.Error())
		if resp != nil {
			w.WriteHeader(resp.StatusCode)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Updating the Self Link in AF
	for key, v := range tsResp{
		var self string
		self, err = updateQoSSelfLink(afCtx.cfg, r, v)
		if err != nil {
			errRspHeader(&w, "GET ALL", err.Error(), http.StatusInternalServerError)
			return
		}
		v.Self = Link(self)
		tsResp[key] = v
	}

	tsRespJSON, err = json.Marshal(tsResp)
	if err != nil {
		log.Errf("AsSessionWithQoS Subscriptions get all: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(tsRespJSON); err != nil {
		log.Errf("AsSessionWithQoS Subscription get all: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
