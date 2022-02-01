package test

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"golang.org/x/crypto/ssh"
)

var folder = flag.String("folder", "", "Folder ID in Yandex.Cloud")
var svckeyfile = flag.String("svckeyfile", "", "Key file for service account in Yandex.Cloud")
var sshkeypath = flag.String("sshkeypath", "", "Private ssh key for access to virtual machines")

func TestEndToEndDeploymentScenario(t *testing.T) {
	fixtureFolder := "../"

	test_structure.RunTestStage(t, "setup", func() {
		terraformOptions := &terraform.Options{
			TerraformDir: fixtureFolder,

			Vars: map[string]interface{}{
				"yc_folder":                   *folder,
				"yc_service_account_key_file": *svckeyfile,
			},
		}

		test_structure.SaveTerraformOptions(t, fixtureFolder, terraformOptions)

		terraform.InitAndApply(t, terraformOptions)
	})

	test_structure.RunTestStage(t, "validate", func() {
		fmt.Println("Run some tests...")
		terraformOptions := test_structure.LoadTerraformOptions(t, fixtureFolder)

		// test load balancer ip existing
		loadbalancerIPAddress := terraform.Output(t, terraformOptions, "load_balancer_public_ip")

		if loadbalancerIPAddress == "" {
			t.Fatal("Cannot retrieve the public IP address value for the load balancer.")
		}

		// test ssh connect
		vmLinuxPublicIPAddress := terraform.Output(t, terraformOptions, "vm_linux_public_ip_address")

		key, err := ioutil.ReadFile(*sshkeypath)
		if err != nil {
			t.Fatalf("Unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			t.Fatalf("Unable to parse private key: %v", err)
		}

		sshConfig := &ssh.ClientConfig{
			User: "ubuntu",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		sshConnection, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", vmLinuxPublicIPAddress), sshConfig)
		if err != nil {
			t.Fatalf("Cannot establish SSH connection to vm-linux public IP address: %v", err)
		}

		defer sshConnection.Close()

		sshSession, err := sshConnection.NewSession()
		if err != nil {
			t.Fatalf("Cannot create SSH session to vm-linux public IP address: %v", err)
		}

		defer sshSession.Close()

		err = sshSession.Run("ping -c 1 8.8.8.8")
		if err != nil {
			t.Fatalf("Cannot ping 8.8.8.8: %v", err)
		}

		// test db connect
		rootCertPool := x509.NewCertPool()

		pem, err := ioutil.ReadFile("./dbrootca.crt")
		if err != nil {
			panic(err)
		}

		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			panic("Failed to append PEM.")
		}

		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs: rootCertPool,
		})

		dbclusterfqdn := terraform.Output(t, terraformOptions, "database_cluster_fqdn")
		dbuser := terraform.Output(t, terraformOptions, "database_user")
		dbpassword := terraform.Output(t, terraformOptions, "database_password")
		dbname := terraform.Output(t, terraformOptions, "database_name")

		mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=custom", dbuser, dbpassword, dbclusterfqdn, dbname)
		fmt.Println(mysqlInfo)

		db, err := sql.Open("mysql", mysqlInfo)
		if err != nil {
			t.Fatalf("Error connection to DB cluster: %v", err)
		}

		defer db.Close()

		err = db.Ping()
		if err != nil {
			t.Fatalf("DB is not available: %v", err)
		}
	})

	test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, fixtureFolder)
		terraform.Destroy(t, terraformOptions)
	})
}
