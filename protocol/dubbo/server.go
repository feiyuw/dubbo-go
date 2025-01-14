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
	"fmt"
	"net"
)

import (
	"github.com/dubbogo/getty"
	"gopkg.in/yaml.v2"
)

import (
	"github.com/feiyuw/dubbo-go/common"
	"github.com/feiyuw/dubbo-go/common/logger"
	"github.com/feiyuw/dubbo-go/config"
	"github.com/feiyuw/dubbo-go/protocol"
)

var srvConf *ServerConfig

func InitServer() {

	// load clientconfig from provider_config
	protocolConf := config.GetProviderConfig().ProtocolConf
	if protocolConf == nil {
		logger.Warnf("protocol_conf is nil")
		return
	}
	dubboConf := protocolConf.(map[interface{}]interface{})[DUBBO]
	if protocolConf == nil {
		logger.Warnf("dubboConf is nil")
		return
	}

	dubboConfByte, err := yaml.Marshal(dubboConf)
	if err != nil {
		panic(err)
	}
	conf := &ServerConfig{}
	err = yaml.Unmarshal(dubboConfByte, conf)
	if err != nil {
		panic(err)
	}

	if err := conf.CheckValidity(); err != nil {
		panic(err)
	}

	srvConf = conf
}

func SetServerConfig(s ServerConfig) {
	srvConf = &s
}

func GetServerConfig() ServerConfig {
	return *srvConf
}

type Server struct {
	conf      ServerConfig
	tcpServer getty.Server
	exporter  protocol.Exporter

	rpcHandler *RpcServerHandler
}

func NewServer(exporter protocol.Exporter) *Server {

	s := &Server{
		exporter: exporter,
		conf:     *srvConf,
	}

	s.rpcHandler = NewRpcServerHandler(s.exporter, s.conf.SessionNumber, s.conf.sessionTimeout)

	return s
}

func (s *Server) newSession(session getty.Session) error {
	var (
		ok      bool
		tcpConn *net.TCPConn
	)
	conf := s.conf

	if conf.GettySessionParam.CompressEncoding {
		session.SetCompressType(getty.CompressZip)
	}

	if tcpConn, ok = session.Conn().(*net.TCPConn); !ok {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection\n", session.Stat(), session.Conn()))
	}

	tcpConn.SetNoDelay(conf.GettySessionParam.TcpNoDelay)
	tcpConn.SetKeepAlive(conf.GettySessionParam.TcpKeepAlive)
	if conf.GettySessionParam.TcpKeepAlive {
		tcpConn.SetKeepAlivePeriod(conf.GettySessionParam.keepAlivePeriod)
	}
	tcpConn.SetReadBuffer(conf.GettySessionParam.TcpRBufSize)
	tcpConn.SetWriteBuffer(conf.GettySessionParam.TcpWBufSize)

	session.SetName(conf.GettySessionParam.SessionName)
	session.SetMaxMsgLen(conf.GettySessionParam.MaxMsgLen)
	session.SetPkgHandler(rpcServerPkgHandler)
	session.SetEventListener(s.rpcHandler)
	session.SetRQLen(conf.GettySessionParam.PkgRQSize)
	session.SetWQLen(conf.GettySessionParam.PkgWQSize)
	session.SetReadTimeout(conf.GettySessionParam.tcpReadTimeout)
	session.SetWriteTimeout(conf.GettySessionParam.tcpWriteTimeout)
	session.SetCronPeriod((int)(conf.sessionTimeout.Nanoseconds() / 1e6))
	session.SetWaitTime(conf.GettySessionParam.waitTimeout)
	logger.Debugf("app accepts new session:%s\n", session.Stat())

	return nil
}

func (s *Server) Start(url common.URL) {
	var (
		addr      string
		tcpServer getty.Server
	)

	addr = url.Location
	tcpServer = getty.NewTCPServer(
		getty.WithLocalAddress(addr),
	)
	tcpServer.RunEventLoop(s.newSession)
	logger.Debugf("s bind addr{%s} ok!", addr)
	s.tcpServer = tcpServer

}

func (s *Server) Stop() {
	s.tcpServer.Close()
}
