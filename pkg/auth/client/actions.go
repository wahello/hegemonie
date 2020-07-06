// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package authclient

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-openapi/strfmt"
	"github.com/jfsmig/hegemonie/pkg/discovery"
	keto "github.com/ory/keto-client-go/client"
	keto_engines "github.com/ory/keto-client-go/client/engines"
	kratos "github.com/ory/kratos-client-go/client"
	kratos_admin "github.com/ory/kratos-client-go/client/admin"
	"os"
)

// ClientCLI gathers the authentication-related client actions available at the command line.
type ClientCLI struct{}

func (cfg *ClientCLI) doKratos(action func(cli *kratos.OryKratos) error) error {
	endpoint, err := discovery.DefaultDiscovery.Kratos()
	if err != nil {
		return err
	}
	cli := kratos.NewHTTPClientWithConfig(strfmt.Default, &kratos.TransportConfig{
		Host:     endpoint,
		BasePath: "/",
		Schemes:  []string{"http"},
	})
	return action(cli)
}

func (cfg *ClientCLI) doKeto(action func(cli *keto.OryKeto) error) error {
	endpoint, err := discovery.DefaultDiscovery.Keto()
	if err != nil {
		return err
	}
	cli := keto.NewHTTPClientWithConfig(strfmt.Default, &keto.TransportConfig{
		Host:     endpoint,
		BasePath: "/",
		Schemes:  []string{"http"},
	})
	return action(cli)
}

// DoShow dumps to os.Stdout a JSON object that summarizes the identity of the
// user with the given ID, plus the list of its ACL.
func (cfg *ClientCLI) DoShow(ctx context.Context, args []string) error {
	var anyError bool
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	for _, id := range args {
		out := make(map[string]interface{})
		cfg.doKratos(func(cli *kratos.OryKratos) error {
			inShow := kratos_admin.GetIdentityParams{
				Context: ctx,
				ID:      id,
			}

			reply, err := cli.Admin.GetIdentity(&inShow)
			if err != nil {
				out["user"] = err
				anyError = true
			} else {
				out["user"] = reply
			}
			return nil
		})
		cfg.doKeto(func(cli *keto.OryKeto) error {
			inShow := keto_engines.GetOryAccessControlPolicyParams{
				Context: ctx,
				ID:      id,
			}
			reply, err := cli.Engines.GetOryAccessControlPolicy(&inShow)
			if err != nil {
				out["acl"] = err
				anyError = true
			} else {
				out["acl"] = reply
			}
			return nil
		})
		encoder.Encode(out)
	}

	if anyError {
		return errors.New("Errors occured (cf. above)")
	}
	return nil
}

// DoList dumps to os.Stdout a JSON stream of objects, separated by CRLF characters,
// with all the known user identities.
func (cfg *ClientCLI) DoList(ctx context.Context, args []string) error {
	return cfg.doKratos(func(cli *kratos.OryKratos) error {
		var page, perPage int64 = 1, 100
		inList := kratos_admin.ListIdentitiesParams{
			Context: context.Background(),
			PerPage: &perPage,
			Page:    &page,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "")
		for {
			reply, err := cli.Admin.ListIdentities(&inList)
			if err != nil {
				return err
			}
			if len(reply.Payload) <= 0 {
				break
			}
			for _, identity := range reply.Payload {
				encoder.Encode(identity)
			}
			page++
		}
		return nil
	})
}

// DoCreate forces the insertion in the Kratos service of a user with the given characteristics.
func (cfg *ClientCLI) DoCreate(ctx context.Context, args []string) error {
	return errors.New("NYI")
}

// DoInvite initiate the process of inviting a user by email, waiting for an action from him/her.
func (cfg *ClientCLI) DoInvite(ctx context.Context, args []string) error {
	return errors.New("NYI")
}
