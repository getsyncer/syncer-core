package config

type Decoder[T TemplateConfig] func(Dynamic) (T, error)

func DefaultDecoder[T TemplateConfig]() func(runConfig Dynamic) (T, error) {
	return func(runConfig Dynamic) (T, error) {
		var cfg T
		if err := runConfig.Decode(&cfg); err != nil {
			return cfg, err
		}
		return cfg, nil
	}
}
