package storage

type Options func(*options) error

type options struct {
	log     Logger
	storage Storage
}

func getDefaultOptions() options {
	return options{
		log:     nil,
		storage: nil,
	}
}

func validate(opts options) error {
	if opts.log == nil {
		return ErrBadLogger
	}

	return nil
}

func WithLogger(logger Logger) Options {
	return func(o *options) error {
		o.log = logger

		return nil
	}
}

func WithStorage(storage Storage) Options {
	return func(o *options) error {
		o.storage = storage

		return nil
	}
}
