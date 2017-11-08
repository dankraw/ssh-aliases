package config

type Reader struct {
	decoder *Decoder
	scanner *Scanner
}

func NewReader() *Reader {
	return &Reader{
		decoder: NewDecoder(),
		scanner: NewScanner(),
	}
}

func (e *Reader) ReadConfigs(dir string) (RawConfigContext, error) {
	files, err := e.scanner.ScanDirectory(dir)
	if err != nil {
		return RawConfigContext{}, err
	}
	config := RawConfigContext{}
	for _, f := range files {
		c, err := e.ReadConfig(f)
		if err != nil {
			return RawConfigContext{}, err
		}
		config.Merge(c)
	}
	return config, nil
}

func (e *Reader) ReadConfig(file string) (RawConfigContext, error) {
	data, err := e.scanner.ReadFile(file)
	if err != nil {
		return RawConfigContext{}, err
	}
	c, err := e.decoder.Decode(data)
	if err != nil {
		return RawConfigContext{}, err
	}
	return c, nil
}
