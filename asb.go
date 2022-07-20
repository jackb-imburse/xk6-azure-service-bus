package asb

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"go.k6.io/k6/js/modules"
)

func init() {
    modules.Register("k6/x/azure-service-bus", new(Asb))
}

type Asb struct{

}

func (a *Asb) GetClient(connectionString string) *azservicebus.Client {

	var client, err = azservicebus.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		panic(err)
	}

	return client
}

func (a *Asb) SendMessage(queue string, message string, client *azservicebus.Client) {
	sender, err := client.NewSender(queue, nil)
	if err != nil {
		panic(err)
	}
	defer sender.Close(context.TODO())

	sbMessage := &azservicebus.Message{
		Body: []byte(message),
	}
	err = sender.SendMessage(context.TODO(), sbMessage, nil)
	if err != nil {
		panic(err)
	}
}