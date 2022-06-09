// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http_service

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/nitrictech/nitric/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/pkg/plugins/gateway/base_http"
	"github.com/nitrictech/nitric/pkg/providers/azure/core"
	"github.com/nitrictech/nitric/pkg/triggers"
	"github.com/nitrictech/nitric/pkg/worker"
)

type azMiddleware struct {
	provider core.AzProvider
}

// func (a *azMiddleware) handleSubscriptionValidation(ctx *fasthttp.RequestCtx, events []eventgrid.Event) {
// 	subPayload := events[0]
// 	var validateData eventgrid.SubscriptionValidationEventData
// 	if err := mapstructure.Decode(subPayload.Data, &validateData); err != nil {
// 		ctx.Error("Invalid subscription event data", 400)
// 		return
// 	}

// 	response := eventgrid.SubscriptionValidationResponse{
// 		ValidationResponse: validateData.ValidationCode,
// 	}

// 	responseBody, _ := json.Marshal(response)
// 	ctx.Success("application/json", responseBody)
// }

func (a *azMiddleware) handleSubscription(process worker.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if strings.ToUpper(string(ctx.Request.Header.Method())) == "OPTIONS" {
			ctx.SuccessString("text/plain", "success")
			return
		}
		var eventgridEvents []eventgrid.Event
		bytes := ctx.Request.Body()
		// TODO: verify topic for validity
		if err := json.Unmarshal(bytes, &eventgridEvents); err != nil {
			ctx.Error("Invalid event grid types", 400)
		}
		topicName := ctx.UserValue("name").(string)

		for _, event := range eventgridEvents {
			// XXX: Assume we have a nitric event for now
			// We have a valid nitric event
			// Decode and pass to our function
			var payloadBytes []byte
			if stringData, ok := event.Data.(string); ok {
				payloadBytes = []byte(stringData)
			} else if byteData, ok := event.Data.([]byte); ok {
				payloadBytes = byteData
			} else {
				// Assume a json serializable struct for now...
				payloadBytes, _ = json.Marshal(event.Data)
			}

			var evt *triggers.Event
			// Just extract the payload from the event type (payload from nitric event is directly mapped)
			evt = &triggers.Event{
				ID:      *event.ID,
				Topic:   topicName,
				Payload: payloadBytes,
			}

			wrkr, err := process.GetWorker(&worker.GetWorkerOptions{
				Event: evt,
				Filter: func(w worker.Worker) bool {
					_, isSubscription := w.(*worker.SubscriptionWorker)
					return isSubscription
				},
			})

			if err != nil {
				log.Default().Println("could not get worker for topic: ", topicName)
				// TODO: Handle error
				continue
			}

			err = wrkr.HandleEvent(evt)
			if err != nil {
				log.Default().Println("could not handle event: ", evt)
				// TODO: Handle error
				continue
			}
		}

		// TODO: event handling failure???
		ctx.SuccessString("text/plain", "success")
	}
}

func (a *azMiddleware) handleSchedule(process worker.WorkerPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		scheduleName := ctx.UserValue("name").(string)

		evt := &triggers.Event{
			ID:    "TODO",
			Topic: scheduleName,
		}

		wrkr, err := process.GetWorker(&worker.GetWorkerOptions{
			Event: evt,
			Filter: func(w worker.Worker) bool {
				_, isSchedule := w.(*worker.ScheduleWorker)
				return isSchedule
			},
		})
		if err != nil {
			log.Default().Println("could not get worker for schedule: ", scheduleName)
		}
		err = wrkr.HandleEvent(evt)
		if err != nil {
			log.Default().Println("could not handle event: ", evt)
		}

		ctx.SuccessString("text/plain", "success")
	}
}

func (a *azMiddleware) router(r *router.Router, pool worker.WorkerPool) {
	r.ANY(base_http.DefaultTopicRoute, a.handleSubscription(pool))
}

// Create a new HTTP Gateway plugin
func New(provider core.AzProvider) (gateway.GatewayService, error) {
	mw := &azMiddleware{
		provider: provider,
	}

	return base_http.New(base_http.BaseHttpGatewayOptions{
		Router: mw.router,
	})
}
