package asb

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"go.k6.io/k6/js/modules"
	"google.golang.org/protobuf/proto"

	pb "github.com/jackb-imburse/xk6-azure-service-bus/contracts/order_management"
)

func init() {
	modules.Register("k6/x/azure-service-bus", new(Asb))
}

func convertDateToCSDateTime(date time.Time) *pb.DateTime {
	datetime := &pb.DateTime{
		Value: date.Unix() / 60 / 60 / 24,
		Scale: 0,
		Kind:  0,
	}
	return datetime
}

func convertUintToDecimal(value uint64) *pb.Decimal {
	decimal := &pb.Decimal{
		Lo:        value,
		Hi:        0,
		SignScale: 0,
	}
	return decimal
}

func convertStringToGuid(guidStr string) *pb.Guid {
	s := strings.Replace(guidStr, "-", "", -1)
	lo := s[12:16] + s[8:12] + s[0:8]
	hi := s[30:] + s[28:30] + s[26:28] + s[24:26] + s[22:24] + s[20:22] + s[18:20] + s[16:18]
	a, _ := strconv.ParseUint(lo, 16, 64)
	b, _ := strconv.ParseUint(hi, 16, 64)
	guid := &pb.Guid{
		Lo: a,
		Hi: b,
	}
	return guid
}

type Asb struct {
}

func (a *Asb) GetClient(connectionString string) *azservicebus.Client {

	var client, err = azservicebus.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		panic(err)
	}

	return client
}

func (a *Asb) EncodeStringBodyAndSendMessage(queue string, correlationId string, subject string, bodyDecoded string, client *azservicebus.Client) {
	sender, err := client.NewSender(queue, nil)
	if err != nil {
		panic(err)
	}
	defer sender.Close(context.TODO())

	sbMessage := &azservicebus.Message{
		CorrelationID: &correlationId,
		Subject:       &subject,
		Body:          []byte(bodyDecoded),
	}
	err = sender.SendMessage(context.TODO(), sbMessage, nil)
	if err != nil {
		panic(err)
	}
}

func (a *Asb) Order_CalculateScheduleStep3_V2(
	correlationId string,
	tenantId string,
	customerRef string,
	financialInstrumentId string,
	workflowId string,
	orderRef string,
	instructionRef string,
	country string,
	currency string,
	amount uint64,
	settledByDate string,
	appId string,
	appQueue string,
	paymentMethod string,
	client *azservicebus.Client) {

	intYear, err := strconv.Atoi(settledByDate[0:4])
	if err != nil {
		panic(err)
	}
	intMonth, err := strconv.Atoi(settledByDate[5:7])
	if err != nil {
		panic(err)
	}
	intDay, err := strconv.Atoi(settledByDate[8:10])
	if err != nil {
		panic(err)
	}

	msg := &pb.ScheduleInstructionStep3V2{
		CorrelationId:               correlationId,
		RequestId:                   convertStringToGuid("30f6432c-ae25-4fbf-bd80-3b537e8ef0f2"),
		TenantId:                    convertStringToGuid(tenantId),
		CustomerRef:                 customerRef,
		FinancialInstrumentId:       convertStringToGuid(financialInstrumentId),
		WorkflowId:                  workflowId,
		OrderRef:                    orderRef,
		InstructionRef:              instructionRef,
		Country:                     country,
		Currency:                    currency,
		Amount:                      convertUintToDecimal(amount),
		RequestedSettlementDate:     convertDateToCSDateTime(time.Date(intYear, time.Month(intMonth), intDay, 0, 0, 0, 0, time.UTC)),
		ForecastedSettlementDate:    convertDateToCSDateTime(time.Date(intYear, time.Month(intMonth), intDay, 0, 0, 0, 0, time.UTC)),
		ScheduledExecutionTimestamp: convertDateToCSDateTime(time.Date(intYear, time.Month(intMonth), intDay, 1, 0, 0, 0, time.UTC)),
		AppId:                       appId,
		AppQueue:                    appQueue,
		PaymentMethod:               paymentMethod,
	}

	out, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalln("Failed to encode msg:", err)
	}

	targetQueue := "transaction"

	sender, err := client.NewSender(targetQueue, nil)
	if err != nil {
		panic(err)
	}
	defer sender.Close(context.TODO())

	subject := "ORDER:CALCULATE_SCHEDULE_STEP_3:V2"

	sbMessage := &azservicebus.Message{
		CorrelationID: &correlationId,
		Subject:       &subject,
		Body:          out,
	}

	err = sender.SendMessage(context.TODO(), sbMessage, nil)
	if err != nil {
		panic(err)
	}
}
