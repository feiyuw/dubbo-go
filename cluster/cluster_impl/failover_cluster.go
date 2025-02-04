/*
 * Licensed to the feiyuw Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the feiyuw License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.feiyuw.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cluster_impl

import (
	"github.com/feiyuw/dubbo-go/cluster"
	"github.com/feiyuw/dubbo-go/common/extension"
	"github.com/feiyuw/dubbo-go/protocol"
)

type failoverCluster struct{}

const name = "failover"

func init() {
	extension.SetCluster(name, NewFailoverCluster)
}

func NewFailoverCluster() cluster.Cluster {
	return &failoverCluster{}
}

func (cluster *failoverCluster) Join(directory cluster.Directory) protocol.Invoker {
	return newFailoverClusterInvoker(directory)
}
