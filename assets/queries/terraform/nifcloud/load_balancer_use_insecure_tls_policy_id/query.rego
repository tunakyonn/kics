package Cx

import data.generic.terraform as tf_lib
import data.generic.common as common_lib

outdatedSSLPolicies := {
	"1",
	"2",
	"3",
	"5",
	"8"
}

CxPolicy[result] {

	lb := input.document[i].resource.nifcloud_load_balancer[name]
    lb.ssl_policy_id == outdatedSSLPolicies[_]

	result := {
		"documentId": input.document[i].id,
		"resourceType": "nifcloud_load_balancer",
		"resourceName": tf_lib.get_resource_name(lb, name),
		"searchKey": sprintf("nifcloud_load_balancer[%s]", [name]),
		"issueType": "IncorrectValue",
		"keyExpectedValue": sprintf("'nifcloud_load_balancer[%s]' should not use outdated/insecure TLS versions for encryption. You should be using TLS v1.2+.", [name]),
		"keyActualValue": sprintf("'nifcloud_load_balancer[%s]' using outdated SSL policy.", [name]),
	}
}
