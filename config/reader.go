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

func (e *Reader) ReadConfigs(dir string) (Config, error) {
	files, err := e.scanner.ScanDirectory(dir)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	for _, f := range files {
		c, err := e.ReadConfig(f)
		if err != nil {
			return Config{}, err
		}
		config.Merge(c)
	}
	return config, nil
}

func (e *Reader) ReadConfig(file string) (Config, error) {
	data, err := e.scanner.ReadFile(file)
	if err != nil {
		return Config{}, err
	}
	c, err := e.decoder.Decode(data)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
