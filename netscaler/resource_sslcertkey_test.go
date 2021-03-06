/*
Copyright 2016 Citrix Systems, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package netscaler

import (
	"fmt"
	"github.com/chiradeep/go-nitro/netscaler"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccSslcertkey_basic(t *testing.T) {
	if os.Getenv("TF_TEST_SSLCERTKEY") == "" {
		t.Skip("skipping test; $TF_TEST_SSLCERTKEY not set")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslcertkeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSslcertkey_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslcertkeyExist("netscaler_sslcertkey.foo", nil),

					resource.TestCheckResourceAttr(
						"netscaler_sslcertkey.foo", "cert", "/var/certs/server.crt"),
					resource.TestCheckResourceAttr(
						"netscaler_sslcertkey.foo", "certkey", "sample_ssl_cert"),
					resource.TestCheckResourceAttr(
						"netscaler_sslcertkey.foo", "key", "/var/certs/server.key"),
				),
			},
		},
	})
}

func testAccCheckSslcertkeyExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ssl cert name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource(netscaler.Sslcertkey.Type(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("SSL cert %s not found", n)
		}

		return nil
	}
}

func testAccCheckSslcertkeyDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netscaler_sslcertkey" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource(netscaler.Sslcertkey.Type(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("SSL certkey %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

const testAccSslcertkey_basic = `


resource "netscaler_sslcertkey" "foo" {
  certkey = "sample_ssl_cert"
  cert = "/var/certs/server.crt"
  key = "/var/certs/server.key"
  notificationperiod = 40
  expirymonitor = "ENABLED"
}
`
