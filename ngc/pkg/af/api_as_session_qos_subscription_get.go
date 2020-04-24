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

// QoSSubscriptionGetAPIService type
type QoSSubscriptionGetAPIService service

func (a *QoSSubscriptionGetAPIService) handleQoSGetResponse(
	qs *AsSessionWithQoSSub, r *http.Response,
	body []byte) error {

	if r.StatusCode == 200 {
		err := json.Unmarshal(body, qs)
		if err != nil {
			log.Errf("Error decoding response body %s: ", err.Error())
		}
		return err
	}
	return handleGetErrorResp(r, body)
}

/*
QoSSubscriptionGet Read an active subscriptions
for the AF and the subscription Id
Read an active QoS subscriptions for the AF and the subscription Id
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param afID Identifier of the AF
 * @param subscriptionID Identifier of the subscription resource

@return AsSessionWithQoSSub
*/
func (a *QoSSubscriptionGetAPIService) QoSSubscriptionGet(
	ctx context.Context, afID string, subscriptionID string) (AsSessionWithQoSSub,
	*http.Response, error) {
	var (
		method  = strings.ToUpper("Get")
		getBody interface{}
		ret     AsSessionWithQoSSub
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFQoSBasePath + "/" + afID +
		"/subscriptions/" + subscriptionID

	headerParams := make(map[string]string)

	headerParams["Content-Type"] = contentType
	headerParams["Accept"] = contentType

	r, err := a.client.prepareRequest(ctx, path, method,
		getBody, headerParams)
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

	if err = a.handleQoSGetResponse(&ret, resp,
		respBody); err != nil {

		return ret, resp, err
	}

	return ret, resp, nil
}
