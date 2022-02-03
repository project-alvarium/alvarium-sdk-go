package annotators

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
)

type TlsAnnotator struct {
	hash contracts.HashType
	kind contracts.AnnotationType
	sign config.SignatureInfo
}

func NewTlsAnnotator(cfg config.SdkInfo) interfaces.Annotator {
	a := TlsAnnotator{}
	a.hash = cfg.Hash.Type
	a.kind = contracts.AnnotationTLS
	a.sign = cfg.Signature
	return &a
}

func (a *TlsAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := DeriveHash(a.hash, data)
	hostname, _ := os.Hostname()
	isSatisfied := false

	// Currently this annotator should only be used in the context of HTTP. TLS is also applicable to pub/sub but
	// a different approach may be required in that scenario to annotate based on the connection rather than per
	// message. More thought required.
	//
	// The methodology below is also very suitable to HTTP requests given that the tls.ConnectionState is readily
	// available off the incoming request whereas pub/sub connection providers may only expose the tls.Config (see
	// https://pkg.go.dev/crypto/tls#Config) requiring a function implementation for
	// VerifyConnection func(ConnectionState) error
	val := ctx.Value(contracts.AnnotationTLS)
	if val != nil {
		tls, ok := val.(*tls.ConnectionState)
		if !ok {
			return contracts.Annotation{}, errors.New(fmt.Sprintf("unexpected type %T", tls))
		}
		if tls != nil {
			isSatisfied = tls.HandshakeComplete
		}
	}
	annotation := contracts.NewAnnotation(key, a.hash, hostname, a.kind, isSatisfied)
	sig, err := SignAnnotation(a.sign.PrivateKey, annotation)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation.Signature = string(sig)
	return annotation, nil
}
