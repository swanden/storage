package storage

type Options func(*options) error

type options struct {
	log            Logger
	storageUseCase StorageUseCase
}

func getDefaultOptions() options {
	return options{
		log:            nil,
		storageUseCase: nil,
	}
}

func validate(opts options) error {
	if opts.log == nil {
		return ErrBadLogger
	}

	if opts.storageUseCase == nil {
		return ErrBadStorageUseCase
	}

	return nil
}

func WithLogger(logger Logger) Options {
	return func(o *options) error {
		o.log = logger

		return nil
	}
}

func WithStorageUseCase(storageUseCase StorageUseCase) Options {
	return func(o *options) error {
		o.storageUseCase = storageUseCase

		return nil
	}
}
