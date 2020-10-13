package getter

// Request represents the path and body parameters coming in to this endpoint.
//
//     GET /imds/instance/{subscription_id}/{resource_group}/{vmss_name}/{instance_id}/
//
type Request struct {
	PathParams PathParams
}

// PathParams hold the parameters we get from the URL.
type PathParams struct {
	SubscriptionID    string
	ResourceGroupName string
	VMSSName          string
	InstanceID        string
}
