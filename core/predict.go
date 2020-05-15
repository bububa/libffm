package core

import (
	"math"

	//"github.com/bububa/libffm"
	"github.com/bububa/libffm/models"
	"github.com/bububa/libffm/utils"
)

func Predict(aModel *models.Model, nodes []*models.Node) float64 {
	var r float64 = 1
	if aModel.Normalization {
		r = 0
		for _, node := range nodes {
			r += node.Value * node.Value
		}
		r = 1 / r
	}
	t := predictTx(aModel, nodes, r)
	return 1 / (1 + math.Exp(-1*t))
}

func predictTx(ffmModel *models.Model, nodes []*models.Node, r float64) float64 {
	align0 := utils.GetLatentFactorsNumberAligned(ffmModel.LatentFactors)
	align1 := ffmModel.Fields * align0

	var t float64
	for i, node := range nodes {
		feature1 := node.Feature
		field1 := node.Field
		value1 := node.Value
		if feature1 >= ffmModel.Features || field1 >= ffmModel.Fields {
			continue
		}
		for j := i + 1; j < len(nodes); j++ {
			feature2 := nodes[j].Feature
			field2 := nodes[j].Field
			value2 := nodes[j].Value
			if feature2 >= ffmModel.Features || field2 >= ffmModel.Fields {
				continue
			}
			w1 := feature1*align1 + field2*align0
			w2 := feature2*align1 + field1*align0
			v := value1 * value2 * r
			for d := 0; d < align0; d += 1 {
				t += ffmModel.W[w1+d] * ffmModel.W[w2+d] * v
			}
		}
	}
	return t
}
