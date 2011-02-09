package cobranding

const DEFAULT_VERSION = "2009-01-09"
const ENDPOINT_COBRANDED_SANDBOX = "https://authorize.payments-sandbox.amazon.com/cobranded-ui/actions/start"
const ENDPOINT_COBRANDED_PRODUCTION = "https://authorize.payments.amazon.com/cobranded-ui/actions/start"

// Changing these two will NOT currently change the underlying operations
const DEFAULT_SIGNATURE_VERSION = "2"
const DEFAULT_SIGNATURE_METHOD  = "HmacSHA256"