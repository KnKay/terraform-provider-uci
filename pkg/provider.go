package pkg

// This stuff is from the boilerplate example repo!
// I hope this is becoming clear during coding

import (
	"context"
	"time"

	"github.com/digineo/go-uci"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				// DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_HOST", nil),
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: false,
				// DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  false,
				Sensitive: true,
				// DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_PASSWORD", nil),
			},
		},
		ResourcesMap:         map[string]*schema.Resource{},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (tree interface{}, diags diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var host *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	if (username != "") || (password != "") || (*host != "") {
		// We have an error! There is no chance to get a connection without this!
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create HashiCups client",
			Detail:   "Unable to authenticate user for authenticated HashiCups client",
		})
		return nil, diags
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// Non-production only
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	tree, err := uci.NewSshTree(config, *host)
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  err.Error(),
		Detail:   err.Error(),
	})

	return tree, diags

}
