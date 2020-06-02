package sqs

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func createResourceDescription(nd *b.NamespaceDescription, qu *string) (*b.ResourceDescription, error) {
	rd := b.ResourceDescription{}

	parts := strings.Split(*qu, "/")
	queueName := h.GetLastStringElement(parts)

	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("QueueName"),
			Value: queueName,
		},
	}
	if err := rd.BuildDimensions(dd); err != nil {
		return nil, err
	}

	rd.ID = queueName
	rd.Name = queueName
	rd.Type = aws.String("sqs")
	rd.Parent = nd

	return &rd, nil
}

// CreateResourceList fetches a list of all SQS Queues in the region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("Creating SQS resource list ...")
	session := sqs.New(nd.Parent.Session)
	input := sqs.ListQueuesInput{}
	result, err := session.ListQueues(&input)
	h.LogIfError(err)

	var w sync.WaitGroup
	w.Add(len(result.QueueUrls))
	ch := make(chan *b.ResourceDescription, len(result.QueueUrls))
	for _, qu := range result.QueueUrls {
		go func(qu *string, wg *sync.WaitGroup) {
			defer wg.Done()
			input := sqs.ListQueueTagsInput{
				QueueUrl: qu,
			}
			tags, err := session.ListQueueTags(&input)
			h.LogIfError(err)

			if nd.Parent.TagsFound(tags) {
				if r, err := createResourceDescription(nd, qu); err == nil {
					ch <- r
				}
				h.LogIfError(err)
			}
		}(qu, &w)
	}
	w.Wait()
	close(ch)

	resources := []*b.ResourceDescription{}
	for r := range ch {
		resources = append(resources, r)
	}
	nd.Mutex.Lock()
	nd.Resources = resources
	nd.Mutex.Unlock()
}
