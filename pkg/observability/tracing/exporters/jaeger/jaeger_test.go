/*
 * Copyright 2018 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package jaeger provides a Jager Tracer
package jaeger

import (
	"testing"

	errs "github.com/tricksterproxy/trickster/pkg/observability/tracing/errors"
	"github.com/tricksterproxy/trickster/pkg/observability/tracing/options"
)

func TestNewTracer(t *testing.T) {

	_, err := NewTracer(nil)
	if err != errs.ErrNoTracerOptions {
		t.Error("expected error for no tracer options")
	}

	opt := options.New()
	opt.Tags = map[string]string{"test": "test"}
	opt.CollectorURL = "1.2.3.4:8000"
	opt.CollectorUser = "abc"
	opt.CollectorPass = "123"

	_, err = NewTracer(opt)
	if err != nil {
		t.Error(err)
	}

	opt.SampleRate = 1
	_, err = NewTracer(opt)
	if err != nil {
		t.Error(err)
	}

	opt.SampleRate = 0.5
	_, err = NewTracer(opt)
	if err != nil {
		t.Error(err)
	}

	opt.JaegerOptions.EndpointType = "agent"
	_, err = NewTracer(opt)
	if err != nil {
		t.Error(err)
	}

}
