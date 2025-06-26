package jaeger

type amqpHeadersCarrier map[string]any

func (c amqpHeadersCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, val := range c {
		v, ok := val.(string)
		if !ok {
			continue
		}

		if err := handler(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (c amqpHeadersCarrier) Set(key, val string) {
	c[key] = val
}
