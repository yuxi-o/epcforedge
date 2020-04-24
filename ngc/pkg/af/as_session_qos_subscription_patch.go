// CertusNet Copyright  

package af

import (
	"context"
	"encoding/json"
	"net/http"
)

func modifyQoSSubscriptionByPatch(cliCtx context.Context, qs AsSessionWithQoSSubPatch,
	afCtx *Context, sID string) (AsSessionWithQoSSub,
	*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	qsRet, resp, err := cli.QoSSubPatchAPI.QoSSubscriptionPatch(cliCtx,
		afCtx.cfg.AfID, sID, qs)

	if err != nil {
		return AsSessionWithQoSSub{}, resp, err
	}
	return qsRet, resp, nil
}

// ModifySubscriptionPatch function
func ModifyQoSSubscriptionPatch(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		qsPatch        AsSessionWithQoSSubPatch
		qsResp         AsSessionWithQoSSub
		resp           *http.Response
		subscriptionID string
		qsRespJSON      []byte
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*Context)
	if afCtx == nil {
		log.Errf("AsSessionWithQoS Subscription patch: " +
			"af-ctx retrieved from request is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewDecoder(r.Body).Decode(&qsPatch); err != nil {
		log.Errf("AsSessionWithQoS Subscription modify: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subscriptionID = getQoSSubsIDFromURL(r)

	qsResp, resp, err = modifyQoSSubscriptionByPatch(cliCtx, qsPatch, afCtx,
		subscriptionID)
	if err != nil {
		log.Errf("AsSessionWithQoS Subscription modify : %s", err.Error())
		w.WriteHeader(getStatusCode(resp))
		return
	}

	// Updating the Self Link in AF
	self, err := updateQoSSelfLink(afCtx.cfg, r, qsResp)
	if err != nil {
		errRspHeader(&w, "PATCH", err.Error(), http.StatusInternalServerError)
		return
	}
	qsResp.Self = Link(self)
	qsRespJSON, err = json.Marshal(qsResp)
	if err != nil {
		errRspHeader(&w, "PATCH", err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	if _, err = w.Write(qsRespJSON); err != nil {
		errRspHeader(&w, "PATCH", err.Error(), http.StatusInternalServerError)
		return
	}
}
