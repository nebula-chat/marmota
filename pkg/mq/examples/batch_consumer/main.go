// Copyright 2022 Teamgram Authors
//  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: teamgramio (teamgram.io@gmail.com)
//

package main

import (
	"context"
	"flag"

	"github.com/teamgram/marmota/pkg/commands"
	kafka "github.com/teamgram/marmota/pkg/mq"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

var configFile = flag.String("f", "c.yaml", "the config file")

type Config struct {
	service.ServiceConf
	TestConsumer kafka.KafkaConsumerConf
}

type Server struct {
	grpcSrv *zrpc.RpcServer
	mq      *kafka.BatchConsumerGroup
}

func New() *Server {
	return new(Server)
}

func (s *Server) Initialize() error {
	var c Config
	conf.MustLoad(*configFile, &c)
	logx.SetUp(c.Log)

	logx.Infov(c)
	mq := kafka.MustKafkaBatchConsumer(&c.TestConsumer)

	mq.RegisterHandlers(
		func(triggerID string, idList []string) {
			logx.Debug("triggerID: ", triggerID, ", idList: ", len(idList))
		},
		func(ctx context.Context, value kafka.MsgChannelValue) {
			for _, msg := range value.MsgList {
				logx.Debug("AggregationID: ", value.AggregationID, ", TriggerID: ", value.TriggerID, ", Msg: ", string(msg.MsgData))
			}
		})

	s.mq = mq
	go s.mq.Start()

	return nil
}

func (s *Server) RunLoop() {
}

func (s *Server) Destroy() {
	s.mq.Stop()
	s.grpcSrv.Stop()
}

func main() {
	commands.Run(New())
}
