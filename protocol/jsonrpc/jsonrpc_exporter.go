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

package jsonrpc

import (
	"sync"
)

import (
	"github.com/feiyuw/dubbo-go/common"
	"github.com/feiyuw/dubbo-go/common/constant"
	"github.com/feiyuw/dubbo-go/common/logger"
	"github.com/feiyuw/dubbo-go/protocol"
)

type JsonrpcExporter struct {
	protocol.BaseExporter
}

func NewJsonrpcExporter(key string, invoker protocol.Invoker, exporterMap *sync.Map) *JsonrpcExporter {
	return &JsonrpcExporter{
		BaseExporter: *protocol.NewBaseExporter(key, invoker, exporterMap),
	}
}

func (je *JsonrpcExporter) Unexport() {
	service := je.GetInvoker().GetUrl().GetParam(constant.INTERFACE_KEY, "")
	je.BaseExporter.Unexport()
	err := common.ServiceMap.UnRegister(JSONRPC, service)
	if err != nil {
		logger.Errorf("[JsonrpcExporter.Unexport] error: %v", err)
	}
}
