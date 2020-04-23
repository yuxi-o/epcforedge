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

// QoSSubscriptionDeleteAPIService type
type QoSSubscriptionDeleteAPIService service

func (a *QoSSubscriptionDeleteAPIService) handleQoSDeleteResponse(
	r *http.Response, body []byte) error {

	newErr := GenericError{
		body:  body,
		error: r.Status,
	}

	switch r.StatusCode {

	case 400, 401, 403, 404, 429, 500, 503:

		var v ProblemDetails
		if r.StatusCode == 401 {

			if fetchNEFAuthorizationToken() != nil {
				log.Infoln("Token refresh failed")
			}
		}

		err := json.Unmarshal(body, &v)
		if err != nil {
			newErr.error = err.Error()
			return newErr
		}
		newErr.model = v
		return newErr

	default:
		b, _ := ioutil.ReadAll(r.Body)
		err := fmt.Errorf("NEF returned error - %s, %s", r.Status, string(b))
		return err
	}
}

/*
QoSSubscriptionDelete Deletes an already
existing subscription
Deletes an already existing subscription
 * @param ctx context.Context - for authentication, logging, cancellation,
 * deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param afID Identifier of the AF
 * @param subscriptionID Identifier of the subscription resource
*/
func (a *QoSSubscriptionDeleteAPIService) QoSSubscriptionDelete(
	ctx context.Context, afID string, subscriptionID string) (*http.Response,
	error) {
	var (
		method     = strings.ToUpper("Delete")
		deleteBody interface{}
	)

	// create path and map variables
	path := a.client.cfg.Protocol + "://" + a.client.cfg.NEFHostname +
		a.client.cfg.NEFPort + a.client.cfg.NEFQoSBasePath + "/" + afID +
		"/subscriptions/" + subscriptionID 

	headerParams := make(map[string]string)

	headerParams["Content-Type"] = contentType
	headerParams["Accept"] = contentType

	r, err := a.client.prepareRequest(ctx, path, method,
		deleteBody, headerParams)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.callAPI(r)
	if err != nil || resp == nil {
		return resp, err
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
		return resp, err
	}

	if resp.StatusCode > 300 {
		if err = a.handleQoSDeleteResponse(resp,
			respBody); err != nil {
			log.Errf("Handle delete response")
			return resp, err
		}
	}

	return resp, nil
}
