package tool

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bububa/libffm/models"
)

func LoadData(dataPath string) ([][]*models.Node, error) {
	fn, err := os.Open(dataPath)
	if err != nil {
		return nil, err
	}
	defer fn.Close()
	buf := bufio.NewReader(fn)
	var instances [][]*models.Node
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line = strings.TrimSpace(line)
		instances = append(instances, ParseDataToNodes(line))
	}
	return instances, nil
}

func ParseDataToNodes(line string) []*models.Node {
	re := regexp.MustCompile("((\\d+)\\:(\\d+):(\\-?\\d+(\\.\\d+)?))")
	matches := re.FindAllStringSubmatch(line, -1)
	var nodes []*models.Node
	for _, m := range matches {
		if len(m) != 6 {
			continue
		}
		field, err := strconv.ParseInt(m[2], 10, 64)
		if err != nil {
			continue
		}
		feature, err := strconv.ParseInt(m[3], 10, 64)
		if err != nil {
			continue
		}
		value, err := strconv.ParseFloat(m[4], 64)
		if err != nil {
			continue
		}
		nodes = append(nodes, &models.Node{
			Field:   int(field),
			Feature: int(feature),
			Value:   value,
		})
	}

	return nodes
}
