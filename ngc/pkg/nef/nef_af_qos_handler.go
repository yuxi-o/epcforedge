// CertusNet Copyright  

package ngcnef

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"errors"
	"strconv"
)

// deal with QoS request
type QoSHandler struct{}

const startSubID int = 1000
const maxSubID int = 9999
const numSubID int = maxSubID - startSubID + 1
var fLocation bool = false
var fLoopbackSubID bool = false
var incSubID int = startSubID
var sID string
var loc string = "https://localhost:8060/3gpp-as-session-with-qos/v1/1/subscriptions/"
//var sliceSub []AsSessionWithQoSSub
var mapSub = make(map[string]AsSessionWithQoSSub, 1)

func (h QoSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	log.Infof("NEF receive QoS [" + r.Method +"] request: " + r.URL.String())

	var qsSub AsSessionWithQoSSub
	var qsSubPatch AsSessionWithQoSSubPatch
	var err error
	var ok bool
	var mdata []byte

	surl := strings.Split(r.URL.String(), "subscriptions")
	if surl[1] == "" { // endwith "subscriptions"
		sID = ""
	} else {
		sID = surl[1][1:]
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Infof("NEF QoS read body error")
		return
	}
	log.Infof("request body: " + string(b))

	switch r.Method {
		case "DELETE":
			if _, ok = mapSub[sID]; ok {
				delete(mapSub, sID)
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		case "GET":
			if sID == "" { // get all
				var sliceSub []AsSessionWithQoSSub
				for _, v := range mapSub{
					sliceSub = append(sliceSub, v)
				}
				mdata, err = json.Marshal(sliceSub)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				} else { 
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.Write(mdata)
				}
			} else if v, ok := mapSub[sID]; ok { // get sId sub
				mdata, err = json.Marshal(v)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				} else { 
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.Write(mdata)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		case "POST", "PUT", "PATCH":
			if sID == "" && r.Method == "POST" {
				fLocation = true
				sID, err = genQoSSubID()
				if err != nil {
					log.Err(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else if (sID != "") && (r.Method == "PUT" || r.Method == "PATCH"){
				fLocation = false
				if  qsSub, ok = mapSub[sID]; !ok {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if r.Method == "PATCH" {
				err = json.Unmarshal(b, &qsSubPatch)
			} else {
				err = json.Unmarshal(b, &qsSub)
			}
			if err != nil {
				log.Err(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if r.Method == "PATCH" {
				patchQoSSub(&qsSub, qsSubPatch)
			}

			qsSub.Self = Link(loc + sID)
			mapSub[sID] = qsSub
			mdata, err = json.Marshal(qsSub)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else { 
				if fLocation {
					w.Header().Set("Location", string(qsSub.Self))
					w.WriteHeader(http.StatusCreated)
				} else {
					w.WriteHeader(http.StatusOK)
				}
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.Write(mdata)
			}
	}
}

func genQoSSubID() (string, error){

	if (incSubID > maxSubID) || (incSubID < startSubID) {
		incSubID = startSubID
	}

	var id = incSubID
	log.Infof("Preset SubID: " + strconv.Itoa(incSubID))

	lenSubID := len(mapSub)
	if lenSubID >= (maxSubID - startSubID + 1) {
		return "", errors.New("All QoS IDs are located!")
	}

	//for incSubID <= maxSubID {
	for {
		if _, ok := mapSub[strconv.Itoa(incSubID)]; !ok{
			id = incSubID 
			incSubID ++
			return strconv.Itoa(id), nil
		}

		incSubID ++
		if incSubID > maxSubID {
			incSubID = startSubID
		}
		if incSubID == id { // loopback
			return "", errors.New("No QoS ID is located!")
		}
	}
}

func patchQoSSub(q *AsSessionWithQoSSub, qp AsSessionWithQoSSubPatch){
	if qp.FlowInfo != nil {
		q.FlowInfo = qp.FlowInfo
	}

	if qp.EthFlowInfo != nil {
		q.EthFlowInfo = qp.EthFlowInfo
	}

	if qp.QosReference != "" {
		q.QosReference = qp.QosReference
	}

	var ut UsageThreshold = UsageThreshold{}
	if qp.UsageThreshold != ut {
		q.UsageThreshold = qp.UsageThreshold
	}
}

