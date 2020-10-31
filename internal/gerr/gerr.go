package gerr

import "strings"

type multiErr struct {
	errs []error
	text string
}

func (err multiErr) Unwrap() error {
	if len(err.errs) == 0 {
		return nil
	}

	return err.errs[0]
}

func (err multiErr) Error() string {
	return err.text
}

func NewMultiError(errs ...error) error {
	mErr := multiErr{
		errs: make([]error, 0, len(errs)),
	}

	errTexts := make([]string, 0, len(errs))

	for _, err := range errs {
		if err == nil {
			continue
		}

		mErr.errs = append(mErr.errs, err)
		errTexts = append(errTexts, err.Error())
	}

	switch len(mErr.errs) {
	case 0:
		return nil
	case 1:
		return mErr.errs[0]
	default:
		mErr.text = strings.Join(errTexts, "; ")
		return mErr
	}
}
