package h2h3_convert

import (
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
)

var h3ErrCodeToH2 = map[http3.ErrCode]http2.ErrCode{
	http3.ErrCodeNoError:              http2.ErrCodeNo,
	http3.ErrCodeGeneralProtocolError: http2.ErrCodeProtocol,
	http3.ErrCodeInternalError:        http2.ErrCodeInternal,
	//http3.ErrCodeStreamCreationError:  http2.ErrCodeStreamClosed,
	http3.ErrCodeClosedCriticalStream: http2.ErrCodeProtocol,
	http3.ErrCodeFrameUnexpected:      http2.ErrCodeProtocol,
	http3.ErrCodeFrameError:           http2.ErrCodeFrameSize,
	http3.ErrCodeExcessiveLoad:        http2.ErrCodeEnhanceYourCalm,
	http3.ErrCodeIDError:              http2.ErrCodeProtocol,
	http3.ErrCodeSettingsError:        http2.ErrCodeProtocol,
	//http3.ErrCodeMissingSettings:      http2.ErrCodeProtocol,
	http3.ErrCodeRequestRejected:   http2.ErrCodeRefusedStream,
	http3.ErrCodeRequestCanceled:   http2.ErrCodeCancel,
	http3.ErrCodeRequestIncomplete: http2.ErrCodeProtocol,
	http3.ErrCodeMessageError:      http2.ErrCodeProtocol,
	http3.ErrCodeConnectError:      http2.ErrCodeConnect,
	http3.ErrCodeVersionFallback:   http2.ErrCodeHTTP11Required,
}
