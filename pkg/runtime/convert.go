package runtime

import "github.com/mitchellh/mapstructure"

func DecodeDict[T any](dict RTDict) (T, error) {
	mp := dict.Map()
	return DecodeMap[T](mp)
}

func DecodeMap[T any](mp map[interface{}]interface{}) (T, error) {
	var data T

	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &data,
		TagName:          "json",
		ErrorUnused:      false,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return data, err
	}

	if err = decoder.Decode(mp); err != nil {
		return data, err
	}

	return data, nil
}
