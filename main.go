package main

import (
	"context"
	"github.com/duysmile/go-pubsub/pubsub"
	"github.com/duysmile/go-pubsub/pubsub/plocal"
	"time"
)

func main() {
	localP := plocal.NewPubSubLocal()

	var sendMailTopic pubsub.Topic = "send-email"

	ctx := context.Background()

	runner := NewPubSubRunner(localP)
	runner.Init(ctx)

	_ = localP.Publish(ctx, sendMailTopic, pubsub.NewMessage("duy1@gmail.com"))
	_ = localP.Publish(ctx, sendMailTopic, pubsub.NewMessage("duy2@gmail.com"))
	_ = localP.Publish(ctx, sendMailTopic, pubsub.NewMessage("duy3@gmail.com"))
	_ = localP.Publish(ctx, sendMailTopic, pubsub.NewMessage("duy4@gmail.com"))
	time.Sleep(10 * time.Second)
}
