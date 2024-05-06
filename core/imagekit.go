package core

import (
	"github.com/imagekit-developer/imagekit-go"
)

type ImageKit struct {
	*imagekit.ImageKit
}

func NewImageKit(env *Env) *ImageKit {
	return &ImageKit{
		imagekit.NewFromParams(imagekit.NewParams{
			PublicKey:   env.IkPublicKey,
			PrivateKey:  env.IkPrivateKey,
			UrlEndpoint: env.IkUrlEndpoint,
		}),
	}
}
