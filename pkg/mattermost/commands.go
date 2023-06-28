package mattermost

// func handlMsg(event *model.WebSocketEvent) {
// 	var post model.Post
// 	err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)
// 	if err != nil {
// 		printError(err)
// 		return
// 	}

// 	if post.UserId == botUser.Id { //Чтобы бот не обрабатывал свои сообщения
// 		return
// 	}
// }
