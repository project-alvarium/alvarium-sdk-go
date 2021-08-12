package annotators

import (
	"context"
	"encoding/json"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/test"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTlsAnnotator_Base(t *testing.T) {
	b, err := ioutil.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	badHashType := cfg
	badHashType.Hash.Type = "invalid"

	badKeyType := cfg
	badKeyType.Signature.PrivateKey.Type = "invalid"

	keyNotFound := cfg
	keyNotFound.Signature.PrivateKey.Path = "/dev/null/private.key"

	rndString := test.FactoryRandomFixedLengthString(1024, test.AlphanumericCharset)
	tests := []struct {
		name        string
		data        string
		cfg         config.SdkInfo
		expectError bool
	}{
		{"tls annotation OK", rndString, cfg, false},
		{"tls bad hash type", rndString, badHashType, false}, // returns "none" hash type
		{"tls bad key type", rndString, badKeyType, true},
		{"tls key not found", rndString, keyNotFound, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tls := NewTlsAnnotator(tt.cfg)
			anno, err := tls.Do(context.Background(), []byte(tt.data))
			test.CheckError(err, tt.expectError, tt.name, t)
			if err == nil {
				result, err := verifySignature(tt.cfg.Signature.PublicKey, anno)
				if err != nil {
					t.Error(err.Error())
				} else if !result {
					t.Error("signature not verified")
				}
			}
		})
	}
}

func TestTlsAnnotator_ServeTLS(t *testing.T) {
	b, err := ioutil.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	tls := NewTlsAnnotator(cfg)

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contracts.AnnotationTLS, r.TLS)
		anno, err := tls.Do(ctx, []byte(test.FactoryRandomFixedLengthString(1024, test.AlphanumericCharset)))
		if err != nil {
			t.Error(err)
		}
		b, _ := json.Marshal(anno)
		w.Write(b)
	}))
	defer ts.Close()

	client := ts.Client()
	result, err := client.Get(ts.URL)
	if err != nil {
		t.Error(err)
	} else {
		defer result.Body.Close()
		var a contracts.Annotation
		body, _ := io.ReadAll(result.Body)
		json.Unmarshal(body, &a)
		if !a.IsSatisfied {
			t.Error("annotation.IsSatisfied should be true")
		}
	}
}

func TestTlsAnnotator_Serve(t *testing.T) {
	b, err := ioutil.ReadFile("../../test/res/config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	var cfg config.SdkInfo
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	tls := NewTlsAnnotator(cfg)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contracts.AnnotationTLS, r.TLS)
		anno, err := tls.Do(ctx, []byte(test.FactoryRandomFixedLengthString(1024, test.AlphanumericCharset)))
		if err != nil {
			t.Error(err)
		}
		b, _ := json.Marshal(anno)
		w.Write(b)
	}))
	defer ts.Close()

	client := ts.Client()
	result, err := client.Get(ts.URL)
	if err != nil {
		t.Error(err)
	} else {
		defer result.Body.Close()
		var a contracts.Annotation
		body, _ := io.ReadAll(result.Body)
		json.Unmarshal(body, &a)
		if a.IsSatisfied {
			t.Error("annotation.IsSatisfied should be false")
		}
	}
}
