// Code generated by protoc-gen-orion. DO NOT EDIT.
// source: Feed.proto

package Feed_proto

import (
	orion "github.com/carousell/Orion/orion"
)

// If you see error please update your orion-protoc-gen by running 'go get -u github.com/carousell/Orion/protoc-gen-orion'
var _ = orion.ProtoGenVersion1_0

// Encoders

// RegisterFeedUpperEncoder registers the encoder for Upper method in Feed
// it registers HTTP  path /api/1.0/upper/{msg} with "GET", "POST", "OPTIONS" methods
func RegisterFeedUpperEncoder(svr orion.Server, encoder orion.Encoder) {
	orion.RegisterEncoders(svr, "Feed", "Upper", []string{"GET", "POST", "OPTIONS"}, "/api/1.0/upper/{msg}", encoder)
}

// RegisterFeedUpperProxyEncoder registers the encoder for UpperProxy method in Feed
// it registers HTTP with "POST", "PUT" methods
func RegisterFeedUpperProxyEncoder(svr orion.Server, encoder orion.Encoder) {
	orion.RegisterEncoders(svr, "Feed", "UpperProxy", []string{"POST", "PUT"}, "", encoder)
}

// Handlers

// RegisterFeedUpperHandler registers the handler for Upper method in Feed
func RegisterFeedUpperHandler(svr orion.Server, handler orion.HTTPHandler) {
	orion.RegisterHandler(svr, "Feed", "Upper", "/api/1.0/upper/{msg}", handler)
}

// RegisterFeedUpperProxyHandler registers the handler for UpperProxy method in Feed
func RegisterFeedUpperProxyHandler(svr orion.Server, handler orion.HTTPHandler) {
	orion.RegisterHandler(svr, "Feed", "UpperProxy", "", handler)
}

// Decoders

// RegisterFeedUpperDecoder registers the decoder for Upper method in Feed
func RegisterFeedUpperDecoder(svr orion.Server, decoder orion.Decoder) {
	orion.RegisterDecoder(svr, "Feed", "Upper", decoder)
}

// RegisterFeedUpperProxyDecoder registers the decoder for UpperProxy method in Feed
func RegisterFeedUpperProxyDecoder(svr orion.Server, decoder orion.Decoder) {
	orion.RegisterDecoder(svr, "Feed", "UpperProxy", decoder)
}

//Streams

// RegisterFeedOrionServer registers Feed to Orion server
// Services need to pass either ServiceFactory or ServiceFactoryV2 implementation
func RegisterFeedOrionServer(sf interface{}, orionServer orion.Server) error {
	err := orionServer.RegisterService(&_Feed_serviceDesc, sf)
	if err != nil {
		return err
	}

	RegisterFeedUpperEncoder(orionServer, nil)
	RegisterFeedUpperProxyEncoder(orionServer, nil)
	return nil
}

// DefaultEncoder
func RegisterFeedDefaultEncoder(svr orion.Server, encoder orion.Encoder) {
	orion.RegisterDefaultEncoder(svr, "Feed", encoder)
}

// DefaultDecoder
func RegisterFeedDefaultDecoder(svr orion.Server, decoder orion.Decoder) {
	orion.RegisterDefaultDecoder(svr, "Feed", decoder)
}

