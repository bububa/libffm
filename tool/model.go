package tool

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bububa/libffm/models"
	"github.com/bububa/libffm/utils"
)

func LoadModel(modelPath string) (*models.Model, error) {
	var ffmModel models.Model
	fn, err := os.Open(modelPath)
	if err != nil {
		return nil, err
	}
	defer fn.Close()
	buf := bufio.NewReader(fn)
	var idx int
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line = strings.TrimSpace(line)
		if idx == 0 {
			re := regexp.MustCompile("^n (\\d+)$")
			matches := re.FindStringSubmatch(line)
			if len(matches) != 2 {
				return nil, errors.New("wrong model")
			}
			features, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ffmModel.Features = int(features)
		} else if idx == 1 {
			re := regexp.MustCompile("^m (\\d+)$")
			matches := re.FindStringSubmatch(line)
			if len(matches) != 2 {
				return nil, errors.New("wrong model")
			}
			fields, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ffmModel.Fields = int(fields)
		} else if idx == 2 {
			re := regexp.MustCompile("^k (\\d+)$")
			matches := re.FindStringSubmatch(line)
			if len(matches) != 2 {
				return nil, errors.New("wrong model")
			}
			latentFactors, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ffmModel.LatentFactors = int(latentFactors)
		} else if idx == 3 {
			re := regexp.MustCompile("^normalization (\\d+)$")
			matches := re.FindStringSubmatch(line)
			if len(matches) != 2 {
				return nil, errors.New("wrong model")
			}
			if matches[1] == "1" {
				ffmModel.Normalization = true
			}
		} else {
			if idx == 4 {
				wSize := ffmModel.Features * ffmModel.Fields * utils.GetLatentFactorsNumberAligned(ffmModel.LatentFactors)
				ffmModel.W = make([]float64, wSize)
			}
			arr := strings.Split(line, " ")
			for i, v := range arr {
				if i == 0 {
					continue
				}
				if i > ffmModel.LatentFactors {
					break
				}
				w, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, errors.New("wrong model")
				}
				wIdx := ffmModel.LatentFactors*(idx-4) + (i - 1)
				ffmModel.W[wIdx] = w
			}
		}
		idx += 1
	}
	return &ffmModel, nil
}
