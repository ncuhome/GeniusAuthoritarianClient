package geniusAuth

import (
	"container/list"
	"context"
	"errors"
	"github.com/ncuhome/GeniusAuthoritarianClient/jwt"
	"github.com/ncuhome/GeniusAuthoritarianClient/jwt/jwtClaims"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"sync/atomic"
	"time"
)

type RpcJwtParser struct {
	Rpc *RpcClient
	Jwt *jwt.Parser

	// UserOperationIDTable uint64 => uint64 uid => UserOperationID
	UserOperationIDTable sync.Map
	// CanceledTokenTable uint64 => time.Time TokenID => ValidBefore
	CanceledTokenTable sync.Map
	Connected          atomic.Bool

	// OnError process error produced in rpc watch stream
	OnError func(err error)
}

var (
	RpcJwtNotConnected            = errors.New("jwt parser rpc not connected")
	RpcJwtUserOperationIDNotFound = errors.New("user operation id not exist")
	RpcJwtTokenInvalid            = errors.New("token canceled")
)

func (p *RpcJwtParser) init() error {
	if p.OnError == nil {
		p.OnError = func(_ error) {}
	}
	go p._RpcStream()
	go p._TableClean()

	return nil
}

func (p *RpcJwtParser) _RpcStream() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		rpcCtx, rpcCancel := context.WithCancel(ctx)
		srv, err := p.Rpc.WatchTokenOperation(rpcCtx, &emptypb.Empty{})
		if err != nil {
			rpcCancel()
			p.OnError(err)
			time.Sleep(time.Second * 5)
			continue
		}

		p.Connected.Store(true)
		for {
			msg, err := srv.Recv()
			if err != nil {
				p.OnError(err)
				break
			}
			for _, userOperation := range msg.UserOperation {
				p.UserOperationIDTable.Store(userOperation.Uid, userOperation.OperationId)
			}
			for _, canceledToken := range msg.CanceledToken {
				p.CanceledTokenTable.Store(canceledToken.Id, time.Unix(canceledToken.ValidBefore, 0))
			}
		}
		p.Connected.Store(false)
		rpcCancel()
	}
}

func (p *RpcJwtParser) _TableClean() {
	for {
		time.Sleep(time.Hour * 12)

		expiredCanceledToken := list.New()
		now := time.Now()
		p.CanceledTokenTable.Range(func(key, value any) bool {
			if value.(time.Time).Before(now) {
				expiredCanceledToken.PushBack(key)
			}
			return true
		})
		for el := expiredCanceledToken.Front(); el != nil; el = el.Next() {
			p.CanceledTokenTable.Delete(el.Value)
		}
	}
}

func (p *RpcJwtParser) TokenStatCheck(claims jwtClaims.ClaimsStandard) error {
	if claims.GetAppCode() != p.Rpc.Api.AppCode {
		return RpcJwtTokenInvalid
	}

	if !p.Connected.Load() {
		return RpcJwtNotConnected
	}

	currentUserOperationID, ok := p.UserOperationIDTable.Load(claims.GetUID())
	if !ok {
		return RpcJwtUserOperationIDNotFound
	} else if currentUserOperationID != claims.GetUserOperateID() {
		return RpcJwtTokenInvalid
	}

	_, ok = p.CanceledTokenTable.Load(claims.GetID())
	if ok {
		return RpcJwtTokenInvalid
	}

	return nil
}

func (p *RpcJwtParser) ParseRefreshToken(token string) (*jwtClaims.RefreshToken, bool, error) {
	claims, valid, err := p.Jwt.ParseRefreshToken(token)
	if err != nil || !valid {
		return claims, valid, err
	}
	return claims, valid, p.TokenStatCheck(claims)
}

func (p *RpcJwtParser) ParseAccessToken(token string) (*jwtClaims.AccessToken, bool, error) {
	claims, valid, err := p.Jwt.ParseAccessToken(token)
	if err != nil || !valid {
		return claims, valid, err
	}
	return claims, valid, p.TokenStatCheck(claims)
}
