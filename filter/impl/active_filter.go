// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @author yiji@apache.org
package impl

import (
	"github.com/feiyuw/dubbo-go/common/extension"
	"github.com/feiyuw/dubbo-go/common/logger"
	"github.com/feiyuw/dubbo-go/filter"
	"github.com/feiyuw/dubbo-go/protocol"
)

const active = "active"

func init() {
	extension.SetFilter(active, GetActiveFilter)
}

type ActiveFilter struct {
}

func (ef *ActiveFilter) Invoke(invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	logger.Infof("invoking active filter. %v,%v", invocation.MethodName(), len(invocation.Arguments()))

	protocol.BeginCount(invoker.GetUrl(), invocation.MethodName())
	return invoker.Invoke(invocation)
}

func (ef *ActiveFilter) OnResponse(result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {

	protocol.EndCount(invoker.GetUrl(), invocation.MethodName())
	return result
}

func GetActiveFilter() filter.Filter {
	return &ActiveFilter{}
}
