package mocks

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type MockSqsOptions struct {
	Queues   []string
	Messages map[string][]*Message
}

type Message struct {
	Id *string
	ReceiptHandle *string
	Body          *string
}

type MockSqs struct {
	sqsiface.SQSAPI
	queues   []string
	messages map[string][]*Message
}

func (s *MockSqs) ListQueues(in *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	queueUrls := make([]*string, 0)

	for _, queue := range s.queues {
		queueUrls = append(queueUrls, &queue)
	}

	return &sqs.ListQueuesOutput{
		QueueUrls: queueUrls,
	}, nil
}

func (s *MockSqs) ReceiveMessage(in *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	for _, q := range s.queues {
		if *in.QueueUrl == q {
			mockMessages := s.messages[q]

			if mockMessages == nil || len(mockMessages) < 1 {
				return &sqs.ReceiveMessageOutput{}, nil
			}

			var messages []*sqs.Message

			for i, m := range mockMessages {
				// Only return up to the max number of messages requested.
				if int64(i) >= *in.MaxNumberOfMessages {
					break
				}
				messages = append(messages, &sqs.Message{
					Body:                   m.Body,
					ReceiptHandle:          m.ReceiptHandle,
				})
				mockMessages[i] = nil
			}

			res := &sqs.ReceiveMessageOutput{
				Messages: messages,
			}

			return res, nil
		}
	}

	return nil, fmt.Errorf("queue not found")
}

func (s *MockSqs) SendMessageBatch(in *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	for _, q := range s.queues {
		if *in.QueueUrl == q {
			if s.messages[q] == nil {
				s.messages[q] = make([]*Message, 0)
			}

			successfulMessages := make([]*sqs.SendMessageBatchResultEntry, 0)
			failedMessages := make([]*sqs.BatchResultErrorEntry, 0)
			for i, e := range in.Entries {
				mockReceiptHandle := fmt.Sprintf("%s%s", string(rune(i)), time.Now())

				s.messages[q] = append(s.messages[q], &Message{
					Id:    e.Id,
					Body: e.MessageBody,
					ReceiptHandle: &mockReceiptHandle,
				})

				successfulMessages = append(successfulMessages, &sqs.SendMessageBatchResultEntry{
					Id: e.Id,
				})
			}

			// TODO: Add a configurable failure mechanism here...
			return &sqs.SendMessageBatchOutput{
				Successful: successfulMessages,
				Failed:     failedMessages,
			}, nil
		}
	}

	return nil, fmt.Errorf("Queue: %s does not exist", *in.QueueUrl)
}

func NewMockSqs(opts *MockSqsOptions) *MockSqs {
	if opts.Messages == nil {
		opts.Messages = make(map[string][]*Message)
	}
	return &MockSqs{
		queues:   opts.Queues,
		messages: opts.Messages,
	}
}
