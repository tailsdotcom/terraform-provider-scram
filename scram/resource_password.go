package scram

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePasswordCreate,
		UpdateContext: resourcePasswordCreate,
		ReadContext:   resourcePasswordNoOp,
		DeleteContext: resourcePasswordNoOp,
		Schema: map[string]*schema.Schema{
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"scram_mech": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SCRAM-SHA-256",
			},
			"iter_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4096,
			},
			"salt": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"stored_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"server_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourcePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var salt = make([]byte, 16)
	_, err := rand.Read(salt[:])
	if err != nil {
		return diag.FromErr(err)
	}
	password := d.Get("password").(string)
	iterCount := d.Get("iter_count").(int)
	hasher := hmac.New(sha256.New, []byte(password))
	hasher.Write(salt)
	hasher.Write([]byte("\x00\x00\x00\x01"))
	ui := hasher.Sum(nil)
	var saltedPassword = make([]byte, 32)
	copy(saltedPassword, ui)
	for i := 1; i < iterCount; i++ {
		hasher := hmac.New(sha256.New, []byte(password))
		hasher.Write(ui)
		ui = hasher.Sum(nil)
		for j, b := range ui {
			saltedPassword[j] ^= b
		}
	}
	hasher = hmac.New(sha256.New, saltedPassword)
	hasher.Write([]byte("Client Key"))
	clientKey := hasher.Sum(nil)
	storedKey := sha256.Sum256(clientKey)
	hasher = hmac.New(sha256.New, saltedPassword)
	hasher.Write([]byte("Server Key"))
	serverKey := hasher.Sum(nil)

	saltB64 := base64.StdEncoding.EncodeToString(salt)
	storedKeyB64 := base64.StdEncoding.EncodeToString(storedKey[:])
	serverKeyB64 := base64.StdEncoding.EncodeToString(serverKey)

	if d.SetId(saltB64); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("salt", saltB64); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("stored_key", storedKeyB64); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("server_key", serverKeyB64); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePasswordNoOp(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
