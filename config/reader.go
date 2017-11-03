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

func (e *Reader) ReadConfigs(dir string) (HostsWithConfigs, error) {
	files, err := e.scanner.ScanDirectory(dir)
	if err != nil {
		return HostsWithConfigs{}, err
	}
	config := HostsWithConfigs{}
	for _, f := range files {
		c, err := e.ReadConfig(f)
		if err != nil {
			return HostsWithConfigs{}, err
		}
		config.Merge(c)
	}
	return config, nil
}

func (e *Reader) ReadConfig(file string) (HostsWithConfigs, error) {
	data, err := e.scanner.ReadFile(file)
	if err != nil {
		return HostsWithConfigs{}, err
	}
	c, err := e.decoder.Decode(data)
	if err != nil {
		return HostsWithConfigs{}, err
	}
	return c, nil
}
