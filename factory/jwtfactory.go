package factory

import (
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	corefactory "github.com/calvine/goauth/core/factory"
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/goauth/core/nullable"
	"github.com/calvine/richerror/errors"
)

type jwtFactory struct {
	issuer             string
	signers            []jwt.Signer
	currentSignerIndex int
}

func NewJWTFactory(issuer string, signers []jwt.Signer) (corefactory.JWTFactory, errors.RichError) {
	if len(signers) == 0 {
		return nil, errors.NewRichError("NoSignersProvidedForJWTFactory", "no signers were provided in jwt factory constructor").WithStack(0)
	}
	return &jwtFactory{
		issuer:             issuer,
		signers:            signers,
		currentSignerIndex: 0,
	}, nil
}

func (jf *jwtFactory) NewSignedJWT(sub string, aud []string, exp nullable.NullableDuration) (jwt.JWT, string, errors.RichError) {
	numSigners := len(jf.signers)
	if jf.currentSignerIndex >= numSigners {
		jf.currentSignerIndex = 0
	}
	signer := jf.signers[jf.currentSignerIndex]
	jf.currentSignerIndex++
	return jf.NewSignedJWTWithSigner(sub, aud, exp, signer)
}

func (jf *jwtFactory) NewSignedJWTWithSigner(sub string, aud []string, exp nullable.NullableDuration, signer jwt.Signer) (jwt.JWT, string, errors.RichError) {
	algFam := signer.GetAlgorithmFamily()
	var alg jwt.JWTSigningAlgorithm
	switch algFam {
	case jwt.HMAC:
		alg = jwt.HS384 // hardcoding this for now. should this be a parameter for the factory constructor for each implemented algorithm?
	default:
		err := coreerrors.NewJWTSigningMaterialAlgorithmTypeNotSupportedError(string(alg), true)
		return jwt.JWT{}, "", err
	}
	token := jwt.NewUnsignedJWT(alg, jf.issuer, aud, sub, exp, time.Time{})
	encodedJWT, err := token.SignAndEncode(signer)
	if err != nil {
		return jwt.JWT{}, "", err
	}
	return token, encodedJWT, nil
}
