package batches

import "errors"

var (
	errBatchesRepositoryCantGetBatches        = errors.New("err tax report service cant get batches")
	errBatchesServiceSanityCheckFailed        = errors.New("err sanity check invalid")
	errBatchesServiceCantCreateBatches        = errors.New("err cant create batches")
	errBatchesServiceCantCountTransactions    = errors.New("err cant count transactions")
	errBatchesServiceCantCountDWHTransactions = errors.New("err cant count dwh transactions")
	errBatchesServiceCantGetBatches           = errors.New("err cant get batches")
	errBatchesServiceNoBatchesToSend          = errors.New("err no batches to send")
	errBatchesServiceNoTransactions           = errors.New("err no transactions")
	errCantUnmarshalResponse                  = errors.New("err cant unmarshal response")
	errCantUnmarshalEncodedToken              = errors.New("err cant unmarshal encoded token")
	errCantMarshalRequest                     = errors.New("err cant marshal request")
	errClientCallFailed                       = errors.New("err client call failed")
	errCantCountBatches                       = errors.New("err can't count batches")
	errCantMarshalGetSeedRequest              = errors.New("err can't marshal get seed request")
	errCantMarshalGetTokenRequest             = errors.New("err can't marshal get token request")
	errCantUnmarshalEncodedSeed               = errors.New("err can't unmarshal encoded seed")
	errCantCanonicalizeTokenObject            = errors.New("err can't canonicalize token object")
	errCantCanonicalizeSignedInfo             = errors.New("err can't canonicalize signed info")
	errCantMarshalTokenObject                 = errors.New("err can't marshal token object")
	errGetSeedReturnsError                    = errors.New("err get seed returns error")
	errGetTokenReturnsError                   = errors.New("err get token returns error")
	errBatchesServiceNoBatches                = errors.New("err no batches")
)

type InvalidInputError struct {
	errDetail string
	Err       error
}

func (e *InvalidInputError) Error() string {
	return e.errDetail
}
