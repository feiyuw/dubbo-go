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
	"github.com/feiyuw/dubbo-go/common"
	"github.com/feiyuw/dubbo-go/common/extension"
	"github.com/feiyuw/dubbo-go/common/logger"
	"github.com/feiyuw/dubbo-go/config"
	"github.com/feiyuw/dubbo-go/protocol"
)

const JSONRPC = "jsonrpc"

func init() {
	extension.SetProtocol(JSONRPC, GetProtocol)
}

var jsonrpcProtocol *JsonrpcProtocol

type JsonrpcProtocol struct {
	protocol.BaseProtocol
	serverMap map[string]*Server
}

func NewDubboProtocol() *JsonrpcProtocol {
	return &JsonrpcProtocol{
		BaseProtocol: protocol.NewBaseProtocol(),
		serverMap:    make(map[string]*Server),
	}
}

func (jp *JsonrpcProtocol) Export(invoker protocol.Invoker) protocol.Exporter {
	url := invoker.GetUrl()
	serviceKey := url.Key()
	exporter := NewJsonrpcExporter(serviceKey, invoker, jp.ExporterMap())
	jp.SetExporterMap(serviceKey, exporter)
	logger.Infof("Export service: %s", url.String())

	// start server
	jp.openServer(url)

	return exporter
}

func (jp *JsonrpcProtocol) Refer(url common.URL) protocol.Invoker {
	invoker := NewJsonrpcInvoker(url, NewHTTPClient(&HTTPOptions{
		HandshakeTimeout: config.GetConsumerConfig().ConnectTimeout,
		HTTPTimeout:      config.GetConsumerConfig().RequestTimeout,
	}))
	jp.SetInvokers(invoker)
	logger.Infof("Refer service: %s", url.String())
	return invoker
}

func (jp *JsonrpcProtocol) Destroy() {
	logger.Infof("jsonrpcProtocol destroy.")

	jp.BaseProtocol.Destroy()

	// stop server
	for key, server := range jp.serverMap {
		delete(jp.serverMap, key)
		server.Stop()
	}
}

func (jp *JsonrpcProtocol) openServer(url common.URL) {
	exporter, ok := jp.ExporterMap().Load(url.Key())
	if !ok {
		panic("[JsonrpcProtocol]" + url.Key() + "is not existing")
	}
	srv := NewServer(exporter.(protocol.Exporter))
	jp.serverMap[url.Location] = srv
	srv.Start(url)
}

func GetProtocol() protocol.Protocol {
	if jsonrpcProtocol != nil {
		return jsonrpcProtocol
	}
	return NewDubboProtocol()
}
