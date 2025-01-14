/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dubbo

import (
	"strconv"
	"sync"
)

import (
	perrors "github.com/pkg/errors"
)

import (
	"github.com/feiyuw/dubbo-go/common"
	"github.com/feiyuw/dubbo-go/common/constant"
	"github.com/feiyuw/dubbo-go/common/logger"
	"github.com/feiyuw/dubbo-go/protocol"
	invocation_impl "github.com/feiyuw/dubbo-go/protocol/invocation"
)

var Err_No_Reply = perrors.New("request need @reply")

type DubboInvoker struct {
	protocol.BaseInvoker
	client      *Client
	destroyLock sync.Mutex
}

func NewDubboInvoker(url common.URL, client *Client) *DubboInvoker {
	return &DubboInvoker{
		BaseInvoker: *protocol.NewBaseInvoker(url),
		client:      client,
	}
}

func (di *DubboInvoker) Invoke(invocation protocol.Invocation) protocol.Result {

	var (
		err    error
		result protocol.RPCResult
	)

	inv := invocation.(*invocation_impl.RPCInvocation)
	url := di.GetUrl()
	// async
	async, err := strconv.ParseBool(inv.AttachmentsByKey(constant.ASYNC_KEY, "false"))
	if err != nil {
		logger.Errorf("ParseBool - error: %v", err)
		async = false
	}
	if async {
		if callBack, ok := inv.CallBack().(func(response CallResponse)); ok {
			result.Err = di.client.AsyncCall(url.Location, url, inv.MethodName(), inv.Arguments(), callBack, inv.Reply())
		} else {
			result.Err = di.client.CallOneway(url.Location, url, inv.MethodName(), inv.Arguments())
		}
	} else {
		if inv.Reply() == nil {
			result.Err = Err_No_Reply
		} else {
			result.Err = di.client.Call(url.Location, url, inv.MethodName(), inv.Arguments(), inv.Reply())
		}
	}
	if result.Err == nil {
		result.Rest = inv.Reply()
	}
	logger.Debugf("result.Err: %v, result.Rest: %v", result.Err, result.Rest)

	return &result
}

func (di *DubboInvoker) Destroy() {
	if di.IsDestroyed() {
		return
	}
	di.destroyLock.Lock()
	defer di.destroyLock.Unlock()

	if di.IsDestroyed() {
		return
	}

	di.BaseInvoker.Destroy()

	if di.client != nil {
		di.client.Close() // close client
	}
}
