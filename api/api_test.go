// Copyright (c) 2015 Spinpunch, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

var Client *model.Client

func Setup() {
	if Srv == nil {
		utils.LoadConfig("config.json")
		NewServer()
		StartServer()
		InitApi()
		Client = model.NewClient("http://localhost:" + utils.Cfg.ServiceSettings.Port + "/api/v1")
	}
}

func SetupBenchmark() (*model.Team, *model.User, *model.Channel) {
	Setup()

	team := &model.Team{Name: "Benchmark Team", Domain: "z-z-" + model.NewId() + "a", Email: "benchmark@nowhere.com", Type: model.TEAM_OPEN}
	team = Client.Must(Client.CreateTeam(team)).Data.(*model.Team)
	user := &model.User{TeamId: team.Id, Email: model.NewId() + "benchmark@test.com", FullName: "Mr. Benchmarker", Password: "pwd"}
	user = Client.Must(Client.CreateUser(user, "")).Data.(*model.User)
	Srv.Store.User().VerifyEmail(user.Id)
	Client.LoginByEmail(team.Domain, user.Email, "pwd")
	channel := &model.Channel{DisplayName: "Benchmark Channel", Name: "a" + model.NewId() + "a", Type: model.CHANNEL_OPEN, TeamId: team.Id}
	channel = Client.Must(Client.CreateChannel(channel)).Data.(*model.Channel)

	return team, user, channel
}

func TearDown() {
	if Srv != nil {
		StopServer()
	}
}