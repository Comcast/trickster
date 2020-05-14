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

package span

import (
	"net/http"
	"testing"

	"github.com/tricksterproxy/trickster/pkg/proxy/context"
	"github.com/tricksterproxy/trickster/pkg/tracing/exporters/stdout"
	"github.com/tricksterproxy/trickster/pkg/tracing/options"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
)

func TestNewChildSpan(t *testing.T) {

	// test with nil context and tracer:
	_, span := NewChildSpan(nil, nil, "test")

	if _, ok := span.(trace.NoopSpan); !ok {
		t.Error("expected NoopSpan")
	}

	// test with nil context but non-nil tracer
	tr, _ := stdout.NewTracer(nil)
	tr.Options.Tags = map[string]string{"testTagName": "testTagValue"}

	// force coverage of tags attachment
	tr.Options.TracerType = "zipkin"
	options.ProcessTracingOptions(map[string]*options.Options{"default": tr.Options}, nil)

	ctx, span := NewChildSpan(nil, tr, "test")
	if ctx == nil {
		t.Error("expected non-nil context")
	}

	if span == nil {
		t.Error("expected non-nil span")
	}

	ctx = context.WithHealthCheckFlag(ctx, true)
	_, span = NewChildSpan(ctx, tr, "test")
	if span != nil {
		t.Error("expected nil span")
	}
}

func TestPrepareRequest(t *testing.T) {

	r, _ := http.NewRequest("GET", "http://example.com", nil)

	_, sp := PrepareRequest(r, nil)
	if sp != nil {
		t.Error("expected nil")
	}

	tr, _ := stdout.NewTracer(nil)

	r = r.WithContext(context.WithHealthCheckFlag(r.Context(), true))
	_, sp = PrepareRequest(r, tr)
	if sp != nil {
		t.Error("expected nil")
	}
	r = r.WithContext(context.WithHealthCheckFlag(r.Context(), false))

	tr.Options.Tags = map[string]string{"testTagName": "testTagValue"}
	_, sp = PrepareRequest(r, tr)
	if sp == nil {
		t.Error("expected non-nill span")
	}

	tr.Options.TracerType = "zipkin"
	options.ProcessTracingOptions(map[string]*options.Options{"default": tr.Options}, nil)

	_, sp = PrepareRequest(r, tr)
	if sp == nil {
		t.Error("expected non-nill span")
	}
}

func TestFilterAttributes(t *testing.T) {
	SetAttributes(nil, nil)
	tr, _ := stdout.NewTracer(nil)
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	_, sp := PrepareRequest(r, tr)
	kvs := []kv.KeyValue{kv.String("testKey", "testValue")}
	tr.Options.OmitTagsList = []string{"testKey2"}
	tr.Options.OmitTags = map[string]bool{"testKey2": true}
	SetAttributes(tr, sp, kvs...)
	kvs = filterAttributes(tr, kvs)

	if len(kvs) != 1 {
		t.Errorf("expected %d got %d", 1, len(kvs))
	}

}
