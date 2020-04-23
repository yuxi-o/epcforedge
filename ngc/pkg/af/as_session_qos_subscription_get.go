// CertusNet Copyright  

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func getQoSSubscription(cliCtx context.Context, afCtx *Context,
	subscriptionID string) (AsSessionWithQoSSub, *http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	qs, resp, err := cli.QoSSubGetAPI.QoSSubscriptionGet(
		cliCtx, afCtx.cfg.AfID, subscriptionID)

	if err != nil {
		return AsSessionWithQoSSub{}, resp, err
	}
	return qs, resp, nil
}

// GetSubscription function
func GetQoSSubscription(w http.ResponseWriter, r *http.Request) {

	var (
		err            error
		qsResp         AsSessionWithQoSSub
		resp           *http.Response
		subscriptionID string
		qsRespJSON     []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("AsSessionWithQoS Subscription get: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	subscriptionID = getQoSSubsIDFromURL(r)

	qsResp, resp, err = getQoSSubscription(cliCtx, afCtx, subscriptionID)
	if err != nil {
		log.Errf("AsSessionWithQoS Subscription get : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	qsRespJSON, err = json.Marshal(qsResp)
	if err != nil {
		log.Errf("AsSessionWithQoS Subscription get : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	if _, err = w.Write(qsRespJSON); err != nil {
		log.Errf("AsSessionWithQoS Subscription get: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
