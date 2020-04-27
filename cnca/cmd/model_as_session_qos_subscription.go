// CertusNet Copyright  
package cnca

// AsSessionWithQoSSub is  As Session with Qos Subscription structure
type AsSessionWithQoSSub struct {
	// URL of created subscription resource
	Self Link `json:"self,omitempty"`
	// String identifying supported features per 
	SuppFeat SupportedFeatures `json:"supportedFeatures,omitempty"`
	// URL where notifications shall be sent
	NotificationDestination Link `json:"notificationDestination,omitempty"`
	// Identifies IP packet filters.
	FlowInfo []FlowInfo `json:"flowInfo,omitempty"`
	// Identifies Ethernet packet filters.
	EthFlowInfo []EthFlowDescription `json:"ethFlowInfo,omitempty"`
	// Identifies a pre-defined QoS information
	QosReference string `json:"qosReference,omitempty"`
	// string identifying a Ipv4 address formatted in the \"dotted decimal\"
	//notation as defined in IETF RFC 1166.
	UeIPv4Addr IPv4Addr `json:"ueIpv4Addr,omitempty"`
	// string identifying a Ipv6 address formatted according to clause 4
	// in IETF RFC 5952.
	UeIPv6Addr IPv6Addr `json:"ueIpv6Addr,omitempty"`
	// string identifying mac address of UE
	MacAddr MacAddr `json:"macAddr,omitempty"`
	UsageThreshold UsageThreshold `json:"usageThreshold,omitempty"`
	SponsorInfo SponsorInformation `json:"sponsorInfo,omitempty"`

	// Set to true by the AF to request the NEF to send a test notification.
	//Set to false or omitted otherwise.
	RequestTestNotification bool `json:"requestTestNotification,omitempty"`
	// Configuration used for sending notifications though web sockets
	WebsockNotifConfig WebsockNotifConfig `json:"websockNotifConfig,omitempty"`
}

type SponsorInformation struct{
	SponsorId string `json:"sponsorId"`
	// Indicates Application Service Provider ID
	AspId string `json:"aspId"`
}

// Volume is unsigned integer identifying a volume in units of bytes
type Volume uint64

type UsageThreshold struct {
	// Identify a period of time in units of seconds
	Duration DurationSec `json:"duration,omitempty"`
	TotalVolume Volume `json:"totalVolume,omitempty"`
	DownlinkVolume Volume `json:"downlinkVolume,omitempty"`
	UplinkVolume Volume `json:"uplinkVolume,omitempty"`
}

// AsSessionWithQoSSubPatch is  As Session with Qos Subscription patch structure
type AsSessionWithQoSSubPatch struct {
	// Identifies IP packet filters.
	FlowInfo []FlowInfo `json:"flowInfo,omitempty"`
	// Identifies Ethernet packet filters.
	EthFlowInfo []EthFlowDescription `json:"ethFlowInfo,omitempty"`
	// Identifies a pre-defined QoS information
	QosReference string `json:"qosReference,omitempty"`
	UsageThreshold UsageThreshold `json:"usageThreshold,omitempty"`
}

type AFAsSessionWithQoSSub struct {
	H		Header
	Policy	struct{
		// URL of created subscription resource
		Self string `yaml:"self,omitempty"`
		// String identifying supported features per 
		SuppFeat string `yaml:"supportedFeatures,omitempty"`
		// URL where notifications shall be sent
		NotificationDestination string `yaml:"notificationDestination,omitempty"`
		// Identifies IP packet filters.
		FlowInfo []struct {
			// Indicates the IP flow.
			FlowID int32 `yaml:"flowId"`
			// Indicates the packet filters of the IP flow. Refer to subclause 5.3.8 of 3GPP TS 29.214 for encoding.
			// It shall contain UL and/or DL IP flow description.
			FlowDescriptions []string `yaml:"flowDescriptions,omitempty"`
		} `yaml:"flowInfo,omitempty"`
		// Identifies Ethernet packet filters.
		EthFlowInfo []struct {
			DestMACAddr string `yaml:"destMacAddr,omitempty"`
			EthType     string `yaml:"ethType"`
			// Defines a packet filter of an IP flow.
			FDesc string `yaml:"fDesc,omitempty"`
			// Possible values are :
			// - DOWNLINK - The corresponding filter applies for traffic to the UE.
			// - UPLINK - The corresponding filter applies for traffic from the UE.
			// - BIDIRECTIONAL The corresponding filter applies for traffic both to and from the UE.
			// UNSPECIFIED - The corresponding filter applies for traffic to the UE (downlink), but has no specific
			//              direction declared. The service data flow detection shall apply the filter for uplink
			//              traffic as if the filter was bidirectional.
			FDir          string   `yaml:"fDir,omitempty"`
			SourceMACAddr string   `yaml:"sourceMacAddr,omitempty"`
			VLANTags      []string `yaml:"vlanTags,omitempty"`
		} `yaml:"ethFlowInfo,omitempty"`

		// Identifies a pre-defined QoS information
		QosReference string `yaml:"qosReference,omitempty"`
		// string identifying a Ipv4 address formatted in the \"dotted decimal\"
		//notation as defined in IETF RFC 1166.
		UeIPv4Addr string `yaml:"ueIpv4Addr,omitempty"`
		// string identifying a Ipv6 address formatted according to clause 4
		// in IETF RFC 5952.
		UeIPv6Addr string `yaml:"ueIpv6Addr,omitempty"`
		// string identifying mac address of UE
		MacAddr string `yaml:"macAddr,omitempty"`
		UsageThreshold struct {
			// Identify a period of time in units of seconds
			Duration DurationSec `yaml:"duration,omitempty"`
			TotalVolume Volume `yaml:"totalVolume,omitempty"`
			DownlinkVolume Volume `yaml:"downlinkVolume,omitempty"`
			UplinkVolume Volume `yaml:"uplinkVolume,omitempty"`
		} `yaml:"usageThreshold,omitempty"`

		SponsorInfo struct{
			SponsorId string `yaml:"sponsorId"`
			// Indicates Application Service Provider ID
			AspId string `yaml:"aspId"`
		} `yaml:"sponsorInfo,omitempty"`

		// Set to true by the AF to request the NEF to send a test notification.
		//Set to false or omitted otherwise.
		RequestTestNotification bool `yaml:"requestTestNotification,omitempty"`
		// Configuration used for sending notifications though web sockets
		WebsockNotifConfig struct {
			WebsocketURI        string `yaml:"websocketUri,omitempty"`
			RequestWebsocketURI bool   `yaml:"requestWebsocketUri,omitempty"`
		} `yaml:"websockNotifConfig,omitempty"`
	} `yaml:"policy"`
}
