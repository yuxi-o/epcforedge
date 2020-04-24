// CertusNet Copyright  

package af

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// Linger please
var (
	_ context.Context
)

// QoSSubscriptionPostAPIService type
type QoSSubscriptionPostAPIService service

func (a *QoSSubscriptionPostAPIService) handleQoSPostResponse(
	qs *AsSessionWithQoSSub, r *http.Response,
	body []byte) error {

	if r.StatusCode == 201 {

		err := json.Unmarshal(body, qs)
		if err != nil {
			log.Errf("Error decoding response body %s, ", err.Error())
		}
		return err
	}

	return handlePostPutPatchErrorResp(r, body)
}

/*
SubscriptionPost Creates a new subscription resource
Creates a new subscription resource
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param afID Identifier of the AF
 * @param body Request to create a new subscription resource

@return AsSessionWithQoSSub
*/
func (a *QoSSubscriptionPostAPIService) QoSSubscriptionPost(
	ctx context.Context, afID string, body AsSessionWithQoSSub) (AsSessionWithQoSSub,
	*http.Response, error) {

	var (
		method   = strings.ToUpper("Post")
		postBody interface{}
		ret      AsSessionWithQoSSub
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFQoSBasePath + "/" + afID +
		"/subscriptions"

	headerParams := make(map[string]string)

	headerParams["Content-Type"] = contentType
	headerParams["Accept"] = contentType

	// body params
	postBody = &body
	r, err := a.client.prepareRequest(ctx, path, method,
		postBody, headerParams)
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

	if err = a.handleQoSPostResponse(&ret, resp,
		respBody); err != nil {
		log.Errf("Handle Post response")
		return ret, resp, err
	}

	return ret, resp, nil
}
