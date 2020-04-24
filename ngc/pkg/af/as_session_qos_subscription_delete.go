// CertusNet Copyright  

package af

import (
	"context"
	"net/http"
)

func deleteQoSSubscription(cliCtx context.Context, afCtx *Context,
	sID string) (*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	resp, err := cli.QoSSubDeleteAPI.QoSSubscriptionDelete(cliCtx,
		afCtx.cfg.AfID, sID)

	if err != nil {
		return resp, err
	}
	return resp, nil

}

// DeleteQoSSubscription function
func DeleteQoSSubscription(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		resp           *http.Response
		subscriptionID string
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

	resp, err = deleteQoSSubscription(cliCtx, afCtx, subscriptionID)
	if err != nil {
		if resp != nil {
			errRspHeader(&w, "DELETE", err.Error(), resp.StatusCode)
		} else {
			errRspHeader(&w, "DELETE", err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(resp.StatusCode)
}
