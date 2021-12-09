/*
Copyright 2021 The Netease Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ctxlog

import (
	"context"
	"fmt"

	"k8s.io/klog/v2"
)

type contextKey string

// CtxKey for context based logging
var CtxKey = contextKey("ID")

// ReqID for logging request ID
var ReqID = contextKey("Req-ID")

func contextFormat(ctx context.Context, format string) string {
	id := ctx.Value(CtxKey)
	if id == nil {
		return format
	}
	a := fmt.Sprintf("ID: %v ", id)
	reqID := ctx.Value(ReqID)
	if reqID == nil {
		return a + format
	}
	a += fmt.Sprintf("Req-ID: %v ", reqID)
	return a + format
}

type CtxLog struct {
	level klog.Level
}

func V(level klog.Level) CtxLog {
	return CtxLog{level: level}
}

func (l CtxLog) Infof(ctx context.Context, format string, args ...interface{}) {
	if klog.V(l.level).Enabled() {
		msg := fmt.Sprintf(contextFormat(ctx, format), args...)
		klog.InfoDepth(1, msg)
	}
}

func ErrorS(ctx context.Context, err error, msg string, keysAndValues ...interface{}) {
	ctxMsg := contextFormat(ctx, msg)
	klog.ErrorSDepth(1, err, ctxMsg, keysAndValues...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(contextFormat(ctx, format), args...)
	klog.ErrorDepth(1, msg)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(contextFormat(ctx, format), args...)
	klog.InfoDepth(1, msg)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(contextFormat(ctx, format), args...)
	klog.WarningDepth(1, msg)
}
