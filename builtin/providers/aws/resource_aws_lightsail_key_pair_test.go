package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSLightsailKeyPair_basic(t *testing.T) {
	var conf lightsail.KeyPair
	lightsailName := fmt.Sprintf("tf-test-lightsail-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_lightsail_key_pair.lightsail_key_pair_test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLightsailKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLightsailKeyPairConfig_basic(lightsailName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLightsailKeyPairExists("aws_lightsail_key_pair.lightsail_key_pair_test", &conf),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "arn"),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "fingerprint"),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "public_key"),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "private_key"),
				),
			},
		},
	})
}

func TestAccAWSLightsailKeyPair_imported(t *testing.T) {
	var conf lightsail.KeyPair
	lightsailName := fmt.Sprintf("tf-test-lightsail-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_lightsail_key_pair.lightsail_key_pair_test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLightsailKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLightsailKeyPairConfig_imported(lightsailName, testLightsailKeyPairPubKey1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLightsailKeyPairExists("aws_lightsail_key_pair.lightsail_key_pair_test", &conf),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "arn"),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "fingerprint"),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "public_key"),
					resource.TestCheckResourceAttr("aws_lightsail_key_pair.lightsail_key_pair_test", "encrypted_fingerprint", ""),
					resource.TestCheckResourceAttr("aws_lightsail_key_pair.lightsail_key_pair_test", "encrypted_private_key", ""),
					resource.TestCheckResourceAttr("aws_lightsail_key_pair.lightsail_key_pair_test", "private_key", ""),
				),
			},
		},
	})
}

func TestAccAWSLightsailKeyPair_encrypted(t *testing.T) {
	var conf lightsail.KeyPair
	lightsailName := fmt.Sprintf("tf-test-lightsail-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_lightsail_key_pair.lightsail_key_pair_test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLightsailKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLightsailKeyPairConfig_encrypted(lightsailName, testLightsailKeyPairPubKey1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLightsailKeyPairExists("aws_lightsail_key_pair.lightsail_key_pair_test", &conf),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "arn"),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "fingerprint"),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "encrypted_fingerprint"),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "encrypted_private_key"),
					resource.TestCheckResourceAttrSet("aws_lightsail_key_pair.lightsail_key_pair_test", "public_key"),
					resource.TestCheckResourceAttr("aws_lightsail_key_pair.lightsail_key_pair_test", "private_key", ""),
				),
			},
		},
	})
}

func testAccCheckAWSLightsailKeyPairExists(n string, res *lightsail.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No LightsailKeyPair set")
		}

		conn := testAccProvider.Meta().(*AWSClient).lightsailconn

		respKeyPair, err := conn.GetKeyPair(&lightsail.GetKeyPairInput{
			KeyPairName: aws.String(rs.Primary.Attributes["name"]),
		})

		if err != nil {
			return err
		}

		if respKeyPair == nil || respKeyPair.KeyPair == nil {
			return fmt.Errorf("KeyPair (%s) not found", rs.Primary.Attributes["name"])
		}
		*res = *respKeyPair.KeyPair
		return nil
	}
}

func testAccCheckAWSLightsailKeyPairDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_lightsail_key_pair" {
			continue
		}

		conn := testAccProvider.Meta().(*AWSClient).lightsailconn

		respKeyPair, err := conn.GetKeyPair(&lightsail.GetKeyPairInput{
			KeyPairName: aws.String(rs.Primary.Attributes["name"]),
		})

		if err == nil {
			if respKeyPair.KeyPair != nil {
				return fmt.Errorf("LightsailKeyPair %q still exists", rs.Primary.ID)
			}
		}

		// Verify the error
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NotFoundException" {
				return nil
			}
		}
		return err
	}

	return nil
}

func testAccAWSLightsailKeyPairConfig_basic(lightsailName string) string {
	return fmt.Sprintf(`
provider "aws" {
  region = "us-east-1"
}
resource "aws_lightsail_key_pair" "lightsail_key_pair_test" {
  name = "%s"
}
`, lightsailName)
}

func testAccAWSLightsailKeyPairConfig_imported(lightsailName, key string) string {
	return fmt.Sprintf(`
provider "aws" {
  region = "us-east-1"
}
resource "aws_lightsail_key_pair" "lightsail_key_pair_test" {
  name = "%s"
	
	public_key = "%s"
}
`, lightsailName, lightsailPubKey)
}

func testAccAWSLightsailKeyPairConfig_encrypted(lightsailName, key string) string {
	return fmt.Sprintf(`
provider "aws" {
  region = "us-east-1"
}
resource "aws_lightsail_key_pair" "lightsail_key_pair_test" {
  name = "%s"
	
	pgp_key = <<EOF
%s
EOF
}
`, lightsailName, key)
}

const lightsailPubKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 phodgson@thoughtworks.com`
const testLightsailKeyPairPubKey1 = `mQENBFXbjPUBCADjNjCUQwfxKL+RR2GA6pv/1K+zJZ8UWIF9S0lk7cVIEfJiprzzwiMwBS5cD0da
rGin1FHvIWOZxujA7oW0O2TUuatqI3aAYDTfRYurh6iKLC+VS+F7H+/mhfFvKmgr0Y5kDCF1j0T/
063QZ84IRGucR/X43IY7kAtmxGXH0dYOCzOe5UBX1fTn3mXGe2ImCDWBH7gOViynXmb6XNvXkP0f
sF5St9jhO7mbZU9EFkv9O3t3EaURfHopsCVDOlCkFCw5ArY+DUORHRzoMX0PnkyQb5OzibkChzpg
8hQssKeVGpuskTdz5Q7PtdW71jXd4fFVzoNH8fYwRpziD2xNvi6HABEBAAG0EFZhdWx0IFRlc3Qg
S2V5IDGJATgEEwECACIFAlXbjPUCGy8GCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEOfLr44B
HbeTo+sH/i7bapIgPnZsJ81hmxPj4W12uvunksGJiC7d4hIHsG7kmJRTJfjECi+AuTGeDwBy84TD
cRaOB6e79fj65Fg6HgSahDUtKJbGxj/lWzmaBuTzlN3CEe8cMwIPqPT2kajJVdOyrvkyuFOdPFOE
A7bdCH0MqgIdM2SdF8t40k/ATfuD2K1ZmumJ508I3gF39jgTnPzD4C8quswrMQ3bzfvKC3klXRlB
C0yoArn+0QA3cf2B9T4zJ2qnvgotVbeK/b1OJRNj6Poeo+SsWNc/A5mw7lGScnDgL3yfwCm1gQXa
QKfOt5x+7GqhWDw10q+bJpJlI10FfzAnhMF9etSqSeURBRW5AQ0EVduM9QEIAL53hJ5bZJ7oEDCn
aY+SCzt9QsAfnFTAnZJQrvkvusJzrTQ088eUQmAjvxkfRqnv981fFwGnh2+I1Ktm698UAZS9Jt8y
jak9wWUICKQO5QUt5k8cHwldQXNXVXFa+TpQWQR5yW1a9okjh5o/3d4cBt1yZPUJJyLKY43Wvptb
6EuEsScO2DnRkh5wSMDQ7dTooddJCmaq3LTjOleRFQbu9ij386Do6jzK69mJU56TfdcydkxkWF5N
ZLGnED3lq+hQNbe+8UI5tD2oP/3r5tXKgMy1R/XPvR/zbfwvx4FAKFOP01awLq4P3d/2xOkMu4Lu
9p315E87DOleYwxk+FoTqXEAEQEAAYkCPgQYAQIACQUCVduM9QIbLgEpCRDny6+OAR23k8BdIAQZ
AQIABgUCVduM9QAKCRAID0JGyHtSGmqYB/4m4rJbbWa7dBJ8VqRU7ZKnNRDR9CVhEGipBmpDGRYu
lEimOPzLUX/ZXZmTZzgemeXLBaJJlWnopVUWuAsyjQuZAfdd8nHkGRHG0/DGum0l4sKTta3OPGHN
C1z1dAcQ1RCr9bTD3PxjLBczdGqhzw71trkQRBRdtPiUchltPMIyjUHqVJ0xmg0hPqFic0fICsr0
YwKoz3h9+QEcZHvsjSZjgydKvfLYcm+4DDMCCqcHuJrbXJKUWmJcXR0y/+HQONGrGJ5xWdO+6eJi
oPn2jVMnXCm4EKc7fcLFrz/LKmJ8seXhxjM3EdFtylBGCrx3xdK0f+JDNQaC/rhUb5V2XuX6VwoH
/AtY+XsKVYRfNIupLOUcf/srsm3IXT4SXWVomOc9hjGQiJ3rraIbADsc+6bCAr4XNZS7moViAAcI
PXFv3m3WfUlnG/om78UjQqyVACRZqqAGmuPq+TSkRUCpt9h+A39LQWkojHqyob3cyLgy6z9Q557O
9uK3lQozbw2gH9zC0RqnePl+rsWIUU/ga16fH6pWc1uJiEBt8UZGypQ/E56/343epmYAe0a87sHx
8iDV+dNtDVKfPRENiLOOc19MmS+phmUyrbHqI91c0pmysYcJZCD3a502X1gpjFbPZcRtiTmGnUKd
OIu60YPNE4+h7u2CfYyFPu3AlUaGNMBlvy6PEpU=`

const testLigthsailKeyPairPrivKey1 = `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
7u2CfYyFPu3AlUaGNMBlvy6PEpU=`
