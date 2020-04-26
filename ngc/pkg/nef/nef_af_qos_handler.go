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

var fLocation bool = false
var startSubID int = 1000
var maxSubID int = 9999
var fLoopbackSubID bool = false
var incSubID int = startSubID
var sID string
var loc string = "https://localhost:8060/3gpp-as-session-with-qos/v1/1/subscriptions/"
//var sliceSub []AsSessionWithQoSSub
var mapSub = make(map[string]AsSessionWithQoSSub, 1)

func (h QoSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	log.Infof("NEF receive QoS [" + r.Method +"] request: " + r.URL.String())

	var qsSub AsSessionWithQoSSub

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
	log.Infof("Preset SubID: " + strconv.Itoa(incSubID))
	log.Infof("request body: " + string(b))

	switch r.Method {
		case "DELETE":
			if _, ok := mapSub[sID]; ok {
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
				mdata, err := json.Marshal(sliceSub)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				} else { 
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.Write(mdata)
				}
			} else if v, ok := mapSub[sID]; ok { // get sId sub
				mdata, err := json.Marshal(v)
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
				sID, err = genSubID()
				if err != nil {
					log.Err(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else if (sID != "") && (r.Method == "PUT" || r.Method == "PATCH"){
				fLocation = false
				if  _, ok := mapSub[sID]; !ok {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err := json.Unmarshal(b, &qsSub)
			if err != nil {
				log.Err(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			qsSub.Self = Link(loc + sID)
			mapSub[sID] = qsSub
			mdata, err := json.Marshal(qsSub)
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

func genSubID() (string, error){
	var id = incSubID
	if (!fLoopbackSubID) && (incSubID <= maxSubID) {
		incSubID ++
		if incSubID > maxSubID {
			fLoopbackSubID = true
		}
		return strconv.Itoa(id), nil
	} else {
		lenSubID := len(mapSub)
		if lenSubID == (maxSubID - startSubID + 1) {
			return "", errors.New("All QoS IDs are located!")
		}
		if incSubID > maxSubID {
			incSubID = startSubID
		}
		for incSubID <= maxSubID {
			if _, ok := mapSub[strconv.Itoa(incSubID)]; !ok{
				id = incSubID 
				incSubID ++
				return strconv.Itoa(id), nil
			}
			incSubID ++
		} 

		return "", errors.New("All QoS IDs are located!")
	}
}

