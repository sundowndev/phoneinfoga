package phonenumbers

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
	fmt "fmt"
	"io/ioutil"
	"strings"
)

// intStringMap is our data structure for maps from prefixes to a single string
// this is used for our carrier and geocoding maps
type intStringMap struct {
	Map       map[int]string
	MaxLength int
}

func loadPrefixMap(data string) (*intStringMap, error) {
	rawBytes, err := decodeUnzipString(data)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(rawBytes)

	// ok, first read in our number of values
	var valueSize uint32
	err = binary.Read(reader, binary.LittleEndian, &valueSize)
	if err != nil {
		return nil, err
	}

	// then our values
	valueBytes := make([]byte, valueSize)
	n, err := reader.Read(valueBytes)
	if uint32(n) < valueSize {
		return nil, fmt.Errorf("unable to read all values: %v", err)
	}

	values := strings.Split(string(valueBytes), "\n")

	// read our # of mappings
	var mappingCount uint32
	err = binary.Read(reader, binary.LittleEndian, &mappingCount)
	if err != nil {
		return nil, err
	}

	maxLength := 0
	mappings := make(map[int]string, mappingCount)
	prefix := 0
	for i := 0; i < int(mappingCount); i++ {
		// first read our diff
		diff, err := binary.ReadUvarint(reader)
		if err != nil {
			return nil, err
		}

		prefix += int(diff)

		// then our map
		var valueIntern uint16
		err = binary.Read(reader, binary.LittleEndian, &valueIntern)
		if err != nil || int(valueIntern) >= len(values) {
			return nil, fmt.Errorf("unable to read interned value: %v", err)
		}

		mappings[prefix] = values[valueIntern]

		strPrefix := fmt.Sprintf("%d", prefix)
		if len(strPrefix) > maxLength {
			maxLength = len(strPrefix)
		}
	}

	// return our values
	return &intStringMap{
		Map:       mappings,
		MaxLength: maxLength,
	}, nil
}

// intStringArrayMap is our map from an int to an array of strings
// this is used for our timezone and region maps
type intStringArrayMap struct {
	Map       map[int][]string
	MaxLength int
}

func loadIntStringArrayMap(data string) (*intStringArrayMap, error) {
	rawBytes, err := decodeUnzipString(data)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(rawBytes)

	// ok, first read in our number of values
	var valueSize uint32
	err = binary.Read(reader, binary.LittleEndian, &valueSize)
	if err != nil {
		return nil, err
	}

	// then our values
	valueBytes := make([]byte, valueSize)
	n, err := reader.Read(valueBytes)
	if uint32(n) < valueSize {
		return nil, fmt.Errorf("unable to read all values: %v", err)
	}

	values := strings.Split(string(valueBytes), "\n")

	// read our # of mappings
	var mappingCount uint32
	err = binary.Read(reader, binary.LittleEndian, &mappingCount)
	if err != nil {
		return nil, err
	}

	maxLength := 0
	mappings := make(map[int][]string, mappingCount)
	key := 0
	for i := 0; i < int(mappingCount); i++ {
		// first read our diff
		diff, err := binary.ReadUvarint(reader)
		if err != nil {
			return nil, err
		}

		key += int(diff)

		// then our values
		var valueCount uint8
		if err = binary.Read(reader, binary.LittleEndian, &valueCount); err != nil {
			return nil, err
		}

		keyValues := make([]string, valueCount)
		for i := 0; i < int(valueCount); i++ {
			var valueIntern uint16
			err = binary.Read(reader, binary.LittleEndian, &valueIntern)
			if err != nil || int(valueIntern) >= len(values) {
				return nil, fmt.Errorf("unable to read interned value: %v", err)
			}
			keyValues[i] = values[valueIntern]
		}
		mappings[key] = keyValues

		strPrefix := fmt.Sprintf("%d", key)
		if len(strPrefix) > maxLength {
			maxLength = len(strPrefix)
		}
	}

	// return our values
	return &intStringArrayMap{
		Map:       mappings,
		MaxLength: maxLength,
	}, nil
}

func decodeUnzipString(data string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	zipReader, err := gzip.NewReader(bytes.NewReader(decodedBytes))
	if err != nil {
		return nil, err
	}

	rawBytes, err := ioutil.ReadAll(zipReader)
	if err != nil {
		return nil, err
	}

	return rawBytes, nil
}
