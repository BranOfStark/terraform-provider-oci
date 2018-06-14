// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	DhcpOptionsRequiredOnlyResource = DhcpOptionsResourceDependencies + `
resource "oci_core_dhcp_options" "test_dhcp_options" {
	#Required
	compartment_id = "${var.compartment_id}"
	options {
		#Required
		type = "${var.dhcp_options_options_type}"
	}
	vcn_id = "${oci_core_vcn.test_vcn.id}"
}
`

	DhcpOptionsResourceConfig = DhcpOptionsResourceDependencies + `
resource "oci_core_dhcp_options" "test_dhcp_options" {
	#Required
	compartment_id = "${var.compartment_id}"
	options {
		#Required
		type = "${var.dhcp_options_options_type}"
	}
	vcn_id = "${oci_core_vcn.test_vcn.id}"

	#Optional
	display_name = "${var.dhcp_options_display_name}"
}
`
	DhcpOptionsPropertyVariables = `
variable "dhcp_options_display_name" { default = "MyDhcpOptions" }
variable "dhcp_options_options_type" { default = "DomainNameServer" }
variable "dhcp_options_state" { default = "AVAILABLE" }

`
	DhcpOptionsResourceDependencies = VcnPropertyVariables + VcnResourceConfig
)

func TestCoreDhcpOptionsResource_basic(t *testing.T) {
	t.Skip("Skipping generated test for now as it has not been worked on.")
	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getRequiredEnvSetting("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_core_dhcp_options.test_dhcp_options"
	datasourceName := "data.oci_core_dhcp_options.test_dhcp_options"

	var resId, resId2 string

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify create
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            config + DhcpOptionsPropertyVariables + compartmentIdVariableStr + DhcpOptionsRequiredOnlyResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "options.0.type", "DomainNameServer"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + DhcpOptionsResourceDependencies,
			},
			// verify create with optionals
			{
				Config: config + DhcpOptionsPropertyVariables + compartmentIdVariableStr + DhcpOptionsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "display_name", "MyDhcpOptions"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "options.0.type", "DomainNameServer"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// verify updates to updatable parameters
			{
				Config: config + `
variable "dhcp_options_display_name" { default = "displayName2" }
variable "dhcp_options_options_type" { default = "DomainNameServer" }
variable "dhcp_options_state" { default = "AVAILABLE" }

                ` + compartmentIdVariableStr + DhcpOptionsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "options.0.type", "DomainNameServer"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config + `
variable "dhcp_options_display_name" { default = "displayName2" }
variable "dhcp_options_options_type" { default = "DomainNameServer" }
variable "dhcp_options_state" { default = "AVAILABLE" }

data "oci_core_dhcp_options" "test_dhcp_options" {
	#Required
	compartment_id = "${var.compartment_id}"
	vcn_id = "${oci_core_vcn.test_vcn.id}"

	#Optional
	display_name = "${var.dhcp_options_display_name}"
	state = "${var.dhcp_options_state}"

    filter {
    	name = "id"
    	values = ["${oci_core_dhcp_options.test_dhcp_options.id}"]
    }
}
                ` + compartmentIdVariableStr + DhcpOptionsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),
					resource.TestCheckResourceAttrSet(datasourceName, "vcn_id"),

					resource.TestCheckResourceAttr(datasourceName, "dhcp_options.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "dhcp_options.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "dhcp_options.0.display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(datasourceName, "dhcp_options.0.id"),
					resource.TestCheckResourceAttr(datasourceName, "dhcp_options.0.options.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "dhcp_options.0.options.0.type", "DomainNameServer"),
					resource.TestCheckResourceAttrSet(datasourceName, "dhcp_options.0.state"),
					resource.TestCheckResourceAttrSet(datasourceName, "dhcp_options.0.time_created"),
					resource.TestCheckResourceAttrSet(datasourceName, "dhcp_options.0.vcn_id"),
				),
			},
		},
	})
}
