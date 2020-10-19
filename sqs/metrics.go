package sqs

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Metrics is a map of default MetricDescriptions for this namespace
var Metrics = map[string]*b.MetricDescription{
	"ApproximateAgeOfOldestMessage": {
		Help:       aws.String("The approximate age of the oldest non-deleted message in the queue"),
		OutputName: aws.String("sqs_approximate_age_of_oldest_message"),
		Statistic:  h.StringPointers("Average", "Maximum"),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ApproximateNumberOfMessagesDelayed": {
		Help:       aws.String("The number of messages in the queue that are delayed and not available for reading immediately. This can happen when the queue is configured as a delay queue or when a message has been sent with a delay parameter"),
		OutputName: aws.String("sqs_approximate_number_of_messages_delayed"),
		Statistic:  h.StringPointers("Average", "Sum"),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ApproximateNumberOfMessagesNotVisible": {
		Help:       aws.String("The number of messages that are in flight. Messages are considered to be in flight if they have been sent to a client but have not yet been deleted or have not yet reached the end of their visibility window"),
		OutputName: aws.String("sqs_approximate_number_of_messages_not_visible"),
		Statistic:  h.StringPointers("Average", "Sum"),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"ApproximateNumberOfMessagesVisible": {
		Help:       aws.String("The number of messages available for retrieval from the queue"),
		OutputName: aws.String("sqs_approximate_number_of_messages_visible"),
		Statistic:  h.StringPointers("Average", "Sum"),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfEmptyReceives": {
		Help:       aws.String("The number of ReceiveMessage API calls that did not return a message"),
		OutputName: aws.String("sqs_number_of_empty_receives"),
		Statistic:  h.StringPointers("Average", "Sum"),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfMessagesDeleted": {
		Help:       aws.String("The number of messages deleted from the queue"),
		OutputName: aws.String("sqs_number_of_messages_deleted"),
		Statistic:  h.StringPointers("Average", "Sum"),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfMessagesReceived": {
		Help:       aws.String("The number of messages returned by calls to the ReceiveMessage action"),
		OutputName: aws.String("sqs_number_of_messages_received"),
		Statistic:  h.StringPointers("Average", "Sum"),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfMessagesSent": {
		Help:       aws.String("The number of messages added to a queue"),
		OutputName: aws.String("sqs_number_of_messages_sent"),
		Statistic:  h.StringPointers("Average", "Sum"),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"SentMessageSize": {
		Help:       aws.String("The size of messages added to a queue"),
		OutputName: aws.String("sqs_sent_message_size"),
		Statistic:  h.StringPointers("Average", "Maximum"),

		Dimensions: []*cloudwatch.Dimension{},
	},
}
