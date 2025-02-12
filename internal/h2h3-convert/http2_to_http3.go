package h2h3_convert

import (
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/http2"
)

var h2ErrCodeToH3 = map[http2.ErrCode]http3.ErrCode{
	http2.ErrCodeNo:       http3.ErrCodeNoError,
	http2.ErrCodeProtocol: http3.ErrCodeGeneralProtocolError,
	http2.ErrCodeInternal: http3.ErrCodeInternalError,
	//http2.ErrCodeFlowControl:  http3.ErrCodeExcessiveLoad,
	//http2.ErrCodeStreamClosed: http3.ErrCodeClosedCriticalStream,
	http2.ErrCodeFrameSize:     http3.ErrCodeFrameError,
	http2.ErrCodeRefusedStream: http3.ErrCodeRequestRejected,
	http2.ErrCodeCancel:        http3.ErrCodeRequestCanceled,
	//http2.ErrCodeCompression:        http3.ErrCodeMessageError,
	http2.ErrCodeConnect:         http3.ErrCodeConnectError,
	http2.ErrCodeEnhanceYourCalm: http3.ErrCodeExcessiveLoad,
	//http2.ErrCodeInadequateSecurity: http3.ErrCodeConnectError,
	http2.ErrCodeHTTP11Required: http3.ErrCodeVersionFallback,
}
