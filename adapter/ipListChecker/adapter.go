// Copyright 2016 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package ipListChecker defines an adapter that checks the existence of an IP address in an online IP address list.
package ipListChecker

import (
	"net/url"

	"github.com/golang/protobuf/proto"

	"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/aspect/listChecker"
	"istio.io/mixer/pkg/registry"
)

// Register records the existence of this adapter
func Register(r registry.Registrar) error {
	return r.RegisterCheckList(newAdapter())
}

type adapterState struct{}

func newAdapter() listChecker.Adapter { return &adapterState{} }
func (a *adapterState) Name() string  { return "istio/ipListChecker" }

func (a *adapterState) Description() string {
	return "Checks whether an IP address is present in an IP address list."
}

func (a *adapterState) Close() error { return nil }

func (a *adapterState) ValidateConfig(cfg proto.Message) (ce *aspect.ConfigErrors) {
	c := cfg.(*Config)

	u, err := url.Parse(c.ProviderUrl)
	if err != nil {
		ce = ce.Append("ProviderUrl", err)
	} else {
		if u.Scheme == "" || u.Host == "" {
			ce = ce.Appendf("ProviderUrl", "URL scheme and host cannot be empty")
		}
	}

	return
}

func (a *adapterState) DefaultConfig() proto.Message {
	return &Config{
		ProviderUrl:     "http://localhost",
		RefreshInterval: 60,
		Ttl:             120,
	}
}

func (a *adapterState) NewAspect(env aspect.Env, cfg proto.Message) (listChecker.Aspect, error) {
	return newAspect(env, cfg.(*Config))
}
