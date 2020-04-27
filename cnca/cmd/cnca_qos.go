// CertusNet Copyright  
package cnca

import (
	"errors"
	"io/ioutil"
	"net/http"
	"fmt"
	"bytes"

	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

const (

)

var qosCmd = &cobra.Command {
	Use:	"qos",
	Short:	"Applies/Create and Manages AsSessionWithQoS subscriptions",
	Args:	cobra.MaximumNArgs(3),
	SilenceUsage:	true,
}

func init(){
	const help = `
  Applies/Creates and Manages AsSessionWithQoS subsriptions

Usage:
  Create AsSessionWithQoS subscription:
		cnca qos apply -f <config.yaml>
  Get All AsSessionWithQoS subscriptions:
		cnca qos get subscriptions
  Get single AsSessionWithQoS subscription:
		cnca qos get subscription <subscription-id>
  Update single AsSessionWithQoS subscription:
		cnca qos patch subscription <subscription-id> -f <config.yaml>
  Delete single AsSessionWithQoS subscription:
		cnca qos delete subscription <subscription-id>

Flags:
  -h, --help			help
  -f, --filename		YAML configuration file
`

	const qosGetHelp = `
  Get active NGC AF AsSessionWithQos subscription(s)

Usage:
  cnca qos get { subscriptions |
				 subscription <subscription-id> }

Flags:
  -h, --help	help

`

	cncaCmd.AddCommand(qosCmd)
	qosCmd.SetHelpTemplate(help)

	qosCmd.AddCommand(qosApplyCmd)
	qosApplyCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = qosApplyCmd.MarkFlagRequired("filename")

	qosCmd.AddCommand(qosGetCmd)
	qosGetCmd.SetHelpTemplate(qosGetHelp)

	qosCmd.AddCommand(qosDeleteCmd)

	qosCmd.AddCommand(qosPatchCmd)
	qosPatchCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	_ = qosPatchCmd.MarkFlagRequired("filename")
}

// ---------------Patch--------------------
var qosPatchCmd = &cobra.Command{
	Use:	"patch",
	Short:	"Update active NGC AF AsSessionWithQoS subscription",
	Args:	cobra.MaximumNArgs(2),
	Run: func (cmd *cobra.Command, args[]string){

		if args[0] != "subscription" {
			fmt.Println("Unsupported intput")
			return
		}

		data, err := readInputData(cmd)
		if err != nil {
			fmt.Println(err)
		}

		var c Header
		if err = yaml.Unmarshal(data, &c); err != nil {
			fmt.Println(err)
			return
		}

		if c.Kind != "ngc_qos" {
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
			return
		}
		
		var s AFAsSessionWithQoSSub
		if err = yaml.Unmarshal(data, &s); err != nil {
			fmt.Println(err)
			return
		}

		var sub[] byte
		sub, err = yaml.Marshal(s.Policy)
		if err != nil {
			fmt.Println(err)
			return
		}

		sub, err = y2j.YAMLToJSON(sub)
		if err != nil {
			fmt.Println(err)
			return
		}

		sub, err = qosPatchSub(args[1], sub)
		if err != nil {
			klog.Info(err)
			return
		}

		sub, err = y2j.JSONToYAML(sub)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[PATCH] Qos Subscription %s patched\n%s", args[0], string(sub))
	},
}

func qosPatchSub(sID string, sub []byte) ([]byte, error){
	url := getNgcAFQoSServiceURL() + "/" + sID

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(sub))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	return res, nil
}
// ---------------DELETE----------------------------
var qosDeleteCmd = &cobra.Command{
	Use:	"delete",
	Short:	"Delete active NGC AF AsSessionWithQoS subscription",
	Args:	cobra.ExactArgs(2),
	Run: func (cmd *cobra.Command, args[]string){
		if args[0] != "subscription" {
			fmt.Println("Unsupported intput")
			return
		}
		
		err := qosDeleteSub(args[1])
		if err != nil {
			klog.Info(err)
			return
		}
		
		fmt.Printf("QoS subscription %s deleted\n", args[1])
	},
}

func qosDeleteSub(sID string) error {
	url := getNgcAFQoSServiceURL() + "/" + sID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}
	
	return nil
}

// ------------------GET-----------------------
var qosGetCmd = &cobra.Command {
	Use:	"get",
	Short:	"Get active NGC AF AsSessionWithQoS subscription(s)",
	Args:	cobra.MaximumNArgs(2),
	Run: func (cmd *cobra.Command, args[]string){
		
		if len(args) < 1 {
			fmt.Println(errors.New("Missing input"))
			return
		}

		if args[0] == "subscription" {
			if len(args) < 2 {
				fmt.Println(errors.New("Missing input"))
				return
			}

			sub, err := qosGetSub(args[1])
			if err != nil {
				klog.Info(err)
				return
			}

			sub, err = y2j.JSONToYAML(sub)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("[GET] Active AF QoS subscription:\n%s", string(sub))
			return
		} else if args[0] == "subscriptions" {
			sub, err := qosGetSub("all")
			if err != nil {
				klog.Info(err)
				return
			}

			if string(sub) == "[]" {
				sub = []byte("none")
			}

			sub, err = y2j.JSONToYAML(sub)
			if err != nil {
				fmt.Println(err)
				return
			}
			
			fmt.Printf("[GET] Active AF Qos Subscriptions:\n%s", string(sub))
			return
		} else {
			fmt.Println(errors.New("Unsupported command"))
			return
		}
	},
}

func qosGetSub(sID string) ([]byte, error){
	var sub []byte

	url := getNgcAFQoSServiceURL()
	if sID != "all" {
		url = string(url + "/" + sID)
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	sub, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// ----------POST------------------------
var qosApplyCmd = &cobra.Command{
	Use:	"apply",
	Short:	"Apply NGC AF AsSessionWithQoS subscription using YAML configuration file",
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string){
		
		data, err := readInputData(cmd)
		if err != nil {
			fmt.Println(err)
		}

		var c Header
		if err = yaml.Unmarshal(data, &c); err != nil {
			fmt.Println(err)
			return
		}

		if c.Kind != "ngc_qos" {
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
			return
		}

		var s AFAsSessionWithQoSSub
		if err = yaml.Unmarshal(data, &s); err != nil {
			fmt.Println(err)
			return
		}

		var sub []byte
		sub, err = yaml.Marshal(s.Policy)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		sub, err = y2j.YAMLToJSON(sub)
		if err != nil {
			fmt.Println(err)
			return
		}

		sub, self, err := qosPostSub(sub)
		if err != nil {
			klog.Info(err)
			return
		}

		sub, err = y2j.JSONToYAML(sub)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[POST] QoS subscription URI: %s\n%s", self, string(sub))
	},
}

func qosPostSub(sub []byte) ([]byte, string, error){
	url := getNgcAFQoSServiceURL()

	req, err := http.NewRequest("POST", url, bytes.NewReader(sub))
	if err != nil {
		return nil, "", err
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	self := resp.Header.Get("Location")
	if self == "" {
		return nil, "", fmt.Errorf("Empty QoS subscription URI returned from AF")
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return res, self, nil
}

func getNgcAFQoSServiceURL() string{
	if UseHTTPProtocol == HTTP2 {
		return NgcAFServiceHTTP2Endpoint + "/qos/subscriptions"
	}
	return NgcAFServiceEndpoint + "/qos/subscriptions"
}

