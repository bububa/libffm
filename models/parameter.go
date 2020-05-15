package models

type Parameter struct {
	Eta           float64 // eta used for per-coordinate learning rate
	Lambda        float64 // used for l2-regularization
	NIters        int     // max iterations
	LatentFactors int     // latent factor dim
	Normalization bool    // instance-wise normalization
	Random        bool    //
	AutoStop      bool    // randomization training order of samples
}

func NewDefaultParameter() *Parameter {
	return &Parameter{
		Eta:           0.2,
		Lambda:        0.00002,
		NIters:        15,
		LatentFactors: 4,
		Normalization: true,
		Random:        true,
	}
}
