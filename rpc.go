package geniusAuth

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/ncuhome/GeniusAuthoritarianClient/rpc/appProto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sync"
	"time"
)

func (c Client) NewRpcClient(addr string) (*RpcClient, error) {
	client := RpcClient{
		api:  &c,
		addr: addr,
	}
	return &client, client.initConnection()
}

type RpcClient struct {
	api *Client

	addr    string
	keypair RpcClientKeypair
	appProto.AppClient
}

type RpcClientKeypair struct {
	sync.RWMutex
	Cert *tls.Certificate
	Cred *RpcClientCredential
}

func (k *RpcClientKeypair) Valid() bool {
	return k.Cred != nil && k.Cred.ValidBefore > time.Now().Add(time.Minute*5).Unix()
}

func (rpc *RpcClient) loadKeypair() (*tls.Certificate, *RpcClientCredential, error) {
	rpc.keypair.RLock()
	if rpc.keypair.Valid() {
		defer rpc.keypair.RUnlock()
		return rpc.keypair.Cert, rpc.keypair.Cred, nil
	}
	rpc.keypair.RUnlock()
	// upgrade lock and try to create keypair
	rpc.keypair.Lock()
	defer rpc.keypair.Unlock()
	if !rpc.keypair.Valid() {
		cred, err := rpc.api.CreateRpcClientCredential()
		if err != nil {
			return nil, nil, err
		}
		cert, err := tls.X509KeyPair(cred.Cert, cred.Key)
		if err != nil {
			return nil, nil, err
		}
		rpc.keypair.Cert, rpc.keypair.Cred = &cert, cred
	}
	return rpc.keypair.Cert, rpc.keypair.Cred, nil
}

func (rpc *RpcClient) initConnection() error {
	pubKeys, err := rpc.api.GetServerPublicKeys()
	if err != nil {
		return err
	}
	block, _ := pem.Decode(pubKeys.Ca)
	if block == nil || block.Type != "CERTIFICATE" {
		return errors.New("decode server ca cert failed")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)

	conn, err := grpc.Dial(rpc.addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			cert, _, err := rpc.loadKeypair()
			return cert, err
		},
		RootCAs: caPool,
	})))
	if err != nil {
		return err
	}
	rpc.AppClient = appProto.NewAppClient(conn)
	return nil
}
