// CertusNet Copyright  

package af

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

func createQoSSubscription(cliCtx context.Context, qs AsSessionWithQoSSub,
	afCtx *Context) (AsSessionWithQoSSub, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	qsResp, resp, err := cli.QoSSubPostAPI.QoSSubscriptionPost(cliCtx,
		afCtx.cfg.AfID, qs)

	if err != nil {
		return AsSessionWithQoSSub{}, resp, err
	}
	return qsResp, resp, nil
}

// CreateSubscription function
func CreateQoSSubscription(w http.ResponseWriter, r *http.Request) {

	var (
		err            error
		qs				AsSessionWithQoSSub
		qsResp			AsSessionWithQoSSub
		resp           *http.Response
		url            *url.URL
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("AsSessionWithQoS Subscription create: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&qs); err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusBadRequest)
		return
	}

	qsResp, resp, err = createQoSSubscription(cliCtx, qs, afCtx)
	if err != nil {
		log.Errf("AsSessionWithQos Subscription create : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	if url, err = resp.Location(); err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}

	// Updating the location url, Self Link and Application Self Link in AF
	afURL := updateQoSURL(afCtx.cfg, r, url.String())

	self, err := updateQoSSelfLink(afCtx.cfg, r, qsResp)
	if err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}
	qsResp.Self = Link(self)
	qsRespJSON, err = json.Marshal(qsResp)
	if err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", afURL)
	w.WriteHeader(resp.StatusCode)

	if _, err = w.Write(qsRespJSON); err != nil {
		errRspHeader(&w, "POST", err.Error(), http.StatusInternalServerError)
		return
	}
}
