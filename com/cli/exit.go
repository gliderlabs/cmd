package cli

// http://stackoverflow.com/questions/1101957/are-there-any-standard-exit-status-codes-in-linux

const (
	StatusOK    = 0
	StatusError = 1

	StatusUsageError    = 64 // cli usage error
	StatusDataError     = 65 // data format error
	StatusNoInput       = 66 // cannot open input
	StatusNoUser        = 67 // addressee unknown
	StatusNoHost        = 68 // host name unknown
	StatusUnavailable   = 69 // service unavailable
	StatusInternalError = 70 // internal software error
	StatusOSError       = 71 // system error (eg can't fork)
	StatusOSFile        = 72 // critical OS file missing
	StatusCreateError   = 73 // can't create (user) output file
	StatusIOError       = 74 // input/output error
	StatusTempFail      = 75 // temp failure; user invited to retry
	StatusProtocolError = 76 // remote error in protocol
	StatusNoPerm        = 77 // permission denied
	StatusConfigError   = 78 // configuration error

	StatusExecError      = 126
	StatusUnknownCommand = 127
	StatusInvalidCommand = 128
)
