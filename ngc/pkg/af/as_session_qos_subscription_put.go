
package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func modifyQoSSubscriptionByPut(cliCtx context.Context, qs AsSessionWithQoSSub,
	afCtx *Context, sID string) (AsSessionWithQoSSub,
	*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	qsResp, resp, err := cli.QoSSubPutAPI.QoSSubscriptionPut(cliCtx,
		afCtx.cfg.AfID, sID, qs)

	if err != nil {
		return AsSessionWithQoSSub{}, resp, err
	}
	return qsResp, resp, nil
}

// ModifyQoSSubscriptionPut function
func ModifyQoSSubscriptionPut(w http.ResponseWriter, r *http.Request) {

	var (
		err     error
		qs      AsSessionWithQoSSub
		qsResp  AsSessionWithQoSSub
		resp    *http.Response
		sID     string
		qsRespJSON      []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("AsSessionWithQoS Subscription put: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&qs); err != nil {
		log.Errf("AsSessionWithQoS Subscription modify: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sID = getQoSSubsIDFromURL(r)
	qsResp, resp, err = modifyQoSSubscriptionByPut(cliCtx, qs, afCtx, sID)
	if err != nil {
		log.Errf("AsSessionWithQos Subscription modify : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	// Updating the Self Link in AF
	self, err := updateQoSSelfLink(afCtx.cfg, r, qsResp)
	if err != nil {
		errRspHeader(&w, "PUT", err.Error(), http.StatusInternalServerError)
		return
	}
	qsResp.Self = Link(self)
	qsRespJSON, err = json.Marshal(qsResp)
	if err != nil {
		errRspHeader(&w, "PUT", err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	if _, err = w.Write(qsRespJSON); err != nil {
		errRspHeader(&w, "PUT", err.Error(), http.StatusInternalServerError)
		return
	}
}
