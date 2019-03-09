package judger

import (
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v2"

	"github.com/syzoj/syzoj-ng-go/app/model"
	"github.com/syzoj/syzoj-ng-go/judger/rpc"
)

type JudgerConfig struct {
	RpcUrl      string `yaml:"rpc_url"`
	JudgerId    string `yaml:"judger_id"`
	JudgerToken string `yaml:"judger_token"`
}

func (j *Judger) LoadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var config JudgerConfig
	if err = yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	j.rpcUrl = config.RpcUrl
	j.judgeRequest = new(rpc.JudgeRequest)
	var judgerId primitive.ObjectID
	if judgerId, err = model.DecodeObjectID(config.JudgerId); err != nil {
		return err
	}
	j.judgeRequest.JudgerId = model.ObjectIDProto(judgerId)
	j.judgeRequest.JudgerToken = proto.String(config.JudgerToken)
	return nil
}
