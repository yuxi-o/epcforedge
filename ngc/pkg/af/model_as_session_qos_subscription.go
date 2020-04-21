package af

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
	QosReference string `json:"qosReference, omitempty"`
	// string identifying a Ipv4 address formatted in the \"dotted decimal\"
	//notation as defined in IETF RFC 1166.
	UeIPv4Addr IPv4Addr `json:"ueIpv4Addr,omitempty"`
	// string identifying a Ipv6 address formatted according to clause 4
	// in IETF RFC 5952.
	UeIPv6Addr IPv6Addr `json:"ueIpv6Addr,omitempty"`
	// string identifying mac address of UE
	MacAddr MacAddr `json:"macAddr,omitempty"`
	UsageThreshold UsageThreshold `json:"usageThreshold"`
	SponsorInfo SponsorInfo `json:"sponsorInfo,omitempty"`

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
	QosReference string `json:"qosReference, omitempty"`
	UsageThreshold UsageThreshold `json:"usageThreshold"`
}
