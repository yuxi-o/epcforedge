// CertusNet Copyright  

package af

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Linger please
var (
	_ context.Context
)

// QoSSubscriptionPatchAPIService type
type QoSSubscriptionPatchAPIService service

func (a *QoSSubscriptionPatchAPIService) handleQoSPatchResponse(
	qs *AsSessionWithQoSSub, r *http.Response,
	body []byte) error {

	if r.StatusCode == 200 {
		err := json.Unmarshal(body, qs)
		if err != nil {
			log.Errf("Error decoding response body %s, ", err.Error())
		}
		return err
	}

	return handlePostPutPatchErrorResp(r, body)

}

/*
QoSSubscriptionPatch Updates an existing subscriptionre
source
Updates an existing subscription resource
 * @param ctx context.Context - for authentication, logging, cancellation,
 * 	deadlines, tracing, etc. Passed from http.Request or
 *	context.Background().
 * @param afID Identifier of the AF
 * @param subscriptionID Identifier of the subscription resource
 * @param body Provides a patch for traffic subscription identified by
 *	subscription ID

@return AsSessionWithQoSSub
*/
func (a *QoSSubscriptionPatchAPIService) QoSSubscriptionPatch(
	ctx context.Context, afID string, subscriptionID string,
	body AsSessionWithQoSSubPatch) (AsSessionWithQoSSub, *http.Response, error) {

	var (
		method    = strings.ToUpper("Patch")
		patchBody interface{}
		ret       TrafficInfluSub
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFQoSBasePath + "/" + afID +
		"/subscriptions/" + subscriptionID 

	headerParams := make(map[string]string)

	headerParams["Content-Type"] = contentType
	headerParams["Accept"] = contentType

	// body params
	patchBody = &body
	r, err := a.client.prepareRequest(ctx, path, method,
		patchBody, headerParams)
	if err != nil {
		return ret, nil, err
	}

	resp, err := a.client.callAPI(r)
	if err != nil || resp == nil {
		return ret, resp, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Errf("response body was not closed properly")
		}
	}()

	if err != nil {
		log.Errf("http response body could not be read")
		return ret, resp, err
	}

	if err = a.handleQoSPatchResponse(&ret, resp,
		respBody); err != nil {
		log.Errf("Handle Patch response")
		return ret, resp, err
	}

	return ret, resp, nil
}
