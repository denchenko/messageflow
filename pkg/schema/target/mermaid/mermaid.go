// Package mermaid provides functionality for generating and rendering Mermaid diagrams.
package mermaid

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/holydocs/messageflow/pkg/messageflow"
)

// targetType defines the schema format type for Mermaid diagrams
const targetType = messageflow.TargetType("mermaid")

var (
	//go:embed templates/service_channels.tmpl
	serviceChannelsTemplateFS embed.FS

	//go:embed templates/channel_services.tmpl
	channelServicesTemplateFS embed.FS

	//go:embed templates/context_services.tmpl
	contextServicesTemplateFS embed.FS

	//go:embed templates/service_services.tmpl
	serviceServicesTemplateFS embed.FS
)

// Ensure Target implements messageflow interfaces.
var (
	_ messageflow.Target = (*Target)(nil)
)

// Target handles the generation and rendering of Mermaid diagrams from message flow schemas.
type Target struct {
	serviceChannelsTemplate *template.Template
	channelServicesTemplate *template.Template
	contextServicesTemplate *template.Template
	serviceServicesTemplate *template.Template
}

// NewTarget creates a new Mermaid diagram formatter instance.
// It initializes the templates from the embedded template files.
func NewTarget() (*Target, error) {
	funcMap := template.FuncMap{
		"sanitizeNodeID": sanitizeNodeID,
		"hasOperation": func(ops []messageflow.Operation, actionStr string, hasReply bool) bool {
			action := messageflow.Action(actionStr)
			for _, op := range ops {
				if op.Action == action {
					if hasReply {
						if op.Reply != nil {
							return true
						}
					} else {
						if op.Reply == nil {
							return true
						}
					}
				}
			}
			return false
		},
	}

	serviceChannelsTemplate, err := template.New("service_channels.tmpl").Funcs(funcMap).ParseFS(serviceChannelsTemplateFS, "templates/service_channels.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing service channels template: %w", err)
	}

	channelServicesTemplate, err := template.New("channel_services.tmpl").Funcs(funcMap).ParseFS(channelServicesTemplateFS, "templates/channel_services.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing channel services template: %w", err)
	}

	contextServicesTemplate, err := template.New("context_services.tmpl").Funcs(funcMap).ParseFS(contextServicesTemplateFS, "templates/context_services.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing context services template: %w", err)
	}

	serviceServicesTemplate, err := template.New("service_services.tmpl").Funcs(funcMap).ParseFS(serviceServicesTemplateFS, "templates/service_services.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing service services template: %w", err)
	}

	return &Target{
		serviceChannelsTemplate: serviceChannelsTemplate,
		channelServicesTemplate: channelServicesTemplate,
		contextServicesTemplate: contextServicesTemplate,
		serviceServicesTemplate: serviceServicesTemplate,
	}, nil
}

// Capabilities returns target capabilities.
func (t *Target) Capabilities() messageflow.TargetCapabilities {
	return messageflow.TargetCapabilities{
		Format: true,
		Render: true,
	}
}

type channelServicesPayload struct {
	Channel          string
	Message          string
	MessageName      string
	ReplyMessage     *string
	ReplyMessageName *string
	Senders          []string
	Receivers        []string
	OmitPayloads     bool
}

type contextServicesPayload struct {
	Services    []messageflow.Service
	Connections []connection
}

type serviceServicesPayload struct {
	MainService      messageflow.Service
	NeighborServices []messageflow.Service
}

type connection struct {
	From          string
	To            string
	Label         string
	Bidirectional bool
}

func (t *Target) FormatSchema(
	_ context.Context,
	s messageflow.Schema,
	opts messageflow.FormatOptions,
) (messageflow.FormattedSchema, error) {
	fs := messageflow.FormattedSchema{
		Type: targetType,
	}

	var buf bytes.Buffer

	switch opts.Mode {
	case messageflow.FormatModeContextServices:
		payload := prepareContextServicesPayload(s)

		err := t.contextServicesTemplate.Execute(&buf, payload)
		if err != nil {
			return messageflow.FormattedSchema{}, fmt.Errorf("executing context services template: %w", err)
		}
	case messageflow.FormatModeServiceChannels:
		payload := prepareServiceChannelsPayload(s, opts.Service)

		err := t.serviceChannelsTemplate.Execute(&buf, payload)
		if err != nil {
			return messageflow.FormattedSchema{}, fmt.Errorf("executing service channels template: %w", err)
		}
	case messageflow.FormatModeChannelServices:
		payload := prepareChannelServicesPayload(s, opts.Channel, opts.OmitPayloads)

		err := t.channelServicesTemplate.Execute(&buf, payload)
		if err != nil {
			return messageflow.FormattedSchema{}, fmt.Errorf("executing channel services template: %w", err)
		}
	case messageflow.FormatModeServiceServices:
		payload := prepareServiceServicesPayload(s, opts.Service)

		err := t.serviceServicesTemplate.Execute(&buf, payload)
		if err != nil {
			return messageflow.FormattedSchema{}, fmt.Errorf("executing service services template: %w", err)
		}
	default:
		return messageflow.FormattedSchema{}, messageflow.NewUnsupportedFormatModeError(opts.Mode, []messageflow.FormatMode{
			messageflow.FormatModeServiceChannels,
			messageflow.FormatModeChannelServices,
			messageflow.FormatModeContextServices,
			messageflow.FormatModeServiceServices,
		})
	}

	fs.Data = buf.Bytes()

	return fs, nil
}

// RenderSchema renders a formatted Mermaid diagram to Mermaid code.
// For Mermaid, this simply returns the formatted code as-is since it will be embedded in markdown.
func (t *Target) RenderSchema(_ context.Context, s messageflow.FormattedSchema) ([]byte, error) {
	if s.Type != targetType {
		return nil, messageflow.NewUnsupportedFormatError(s.Type, targetType)
	}

	return s.Data, nil
}

// sanitizeNodeID converts a string to a valid Mermaid node ID (no spaces, special chars).
func sanitizeNodeID(name string) string {
	// Replace spaces and special characters with underscores
	result := strings.Builder{}
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		} else if r == ' ' || r == '-' {
			result.WriteRune('_')
		}
	}
	return result.String()
}

func prepareServiceChannelsPayload(s messageflow.Schema, serviceName string) messageflow.Service {
	if serviceName == "" && len(s.Services) == 1 {
		return s.Services[0]
	}

	for _, service := range s.Services {
		if service.Name == serviceName {
			return service
		}
	}

	return messageflow.Service{}
}

func prepareChannelServicesPayload(s messageflow.Schema, channel string, omitPayloads bool) channelServicesPayload {
	payload := channelServicesPayload{
		Channel:      channel,
		OmitPayloads: omitPayloads,
	}

	for _, service := range s.Services {
		for _, op := range service.Operation {
			if op.Channel.Name == channel {
				switch op.Action {
				case messageflow.ActionSend:
					payload.Senders = append(payload.Senders, service.Name)
				case messageflow.ActionReceive:
					payload.Receivers = append(payload.Receivers, service.Name)
				}

				if len(op.Channel.Messages) > 0 {
					firstMessage := op.Channel.Messages[0]
					if len(payload.Message) < len(firstMessage.Payload) {
						payload.Message = firstMessage.Payload
						payload.MessageName = firstMessage.Name
					}
				}

				if op.Reply != nil && len(op.Reply.Messages) > 0 {
					firstReplyMessage := op.Reply.Messages[0]
					if payload.ReplyMessage == nil ||
						(len(*payload.ReplyMessage) < len(firstReplyMessage.Payload)) {
						payload.ReplyMessage = &firstReplyMessage.Payload
						payload.ReplyMessageName = &firstReplyMessage.Name
					}
				}
			}
		}
	}

	return payload
}

func prepareContextServicesPayload(s messageflow.Schema) contextServicesPayload {
	formattedServices := make([]messageflow.Service, len(s.Services))
	for i, service := range s.Services {
		formattedServices[i] = messageflow.Service{
			Name:        service.Name,
			Description: formatDescription(service.Description),
			Operation:   service.Operation,
		}
	}

	payload := contextServicesPayload{
		Services:    formattedServices,
		Connections: []connection{},
	}

	servicePairs := make(map[string]map[string]bool) // service1->service2 -> hasSendOperation

	// First pass: collect all send operations between service pairs
	for _, service := range s.Services {
		for _, op := range service.Operation {
			if op.Action == messageflow.ActionSend {
				for _, otherService := range s.Services {
					if otherService.Name == service.Name {
						continue
					}

					for _, otherOp := range otherService.Operation {
						if otherOp.Channel.Name == op.Channel.Name && otherOp.Action == messageflow.ActionReceive {
							if servicePairs[service.Name] == nil {
								servicePairs[service.Name] = make(map[string]bool)
							}
							servicePairs[service.Name][otherService.Name] = true
							break
						}
					}
				}
			}
		}
	}

	// Second pass: create connections and detect bidirectional communication
	connectionMap := make(map[string]connection)

	for service1, receivers := range servicePairs {
		for service2 := range receivers {
			bidirectional := servicePairs[service2] != nil && servicePairs[service2][service1]

			var from, to string
			switch {
			case bidirectional && service1 < service2:
				from, to = service1, service2
			case bidirectional && service1 >= service2:
				from, to = service2, service1
			default:
				from, to = service1, service2
			}

			key := fmt.Sprintf("%s->%s", from, to)

			label := determineConnectionLabel(s, from, to)

			conn := connection{
				From:          from,
				To:            to,
				Label:         label,
				Bidirectional: bidirectional,
			}

			connectionMap[key] = conn
		}
	}

	keys := make([]string, 0, len(connectionMap))
	for key := range connectionMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		payload.Connections = append(payload.Connections, connectionMap[key])
	}

	return payload
}

// formatDescription formats a description string for better readability in Mermaid diagrams.
func formatDescription(desc string) string {
	if desc == "" {
		return ""
	}
	// For Mermaid, we can use line breaks with <br/> or keep it simple
	return desc
}

func determineConnectionLabel(s messageflow.Schema, service1, service2 string) string {
	var hasPub, hasReq bool

	svc1 := findServiceByName(s, service1)
	svc2 := findServiceByName(s, service2)

	for _, op1 := range svc1.Operation {
		for _, op2 := range svc2.Operation {
			if op1.Channel.Name != op2.Channel.Name {
				continue
			}

			switch {
			case op1.Action == messageflow.ActionSend && op2.Action == messageflow.ActionReceive:
				if op1.Reply != nil {
					hasReq = true
					continue
				}

				hasPub = true
			case op1.Action == messageflow.ActionReceive && op2.Action == messageflow.ActionSend:
				if op2.Reply != nil {
					hasReq = true
					continue
				}

				hasPub = true
			}
		}
	}

	switch {
	case hasPub && hasReq:
		return "Pub/Req"
	case hasReq:
		return "Req"
	default:
		return "Pub"
	}
}

func findServiceByName(s messageflow.Schema, name string) messageflow.Service {
	for _, service := range s.Services {
		if service.Name == name {
			return service
		}
	}
	return messageflow.Service{}
}

func prepareServiceServicesPayload(s messageflow.Schema, serviceName string) serviceServicesPayload {
	var mainService messageflow.Service
	if serviceName == "" && len(s.Services) == 1 {
		mainService = s.Services[0]
	} else {
		for _, service := range s.Services {
			if service.Name == serviceName {
				mainService = service
				break
			}
		}
	}

	var (
		neighborServices           = make([]messageflow.Service, 0)
		neighborServiceMap         = make(map[string]bool)
		mainServiceSendChannels    = make(map[string]bool)
		mainServiceReceiveChannels = make(map[string]bool)
	)

	for _, op := range mainService.Operation {
		switch op.Action {
		case messageflow.ActionSend:
			mainServiceSendChannels[op.Channel.Name] = true
		case messageflow.ActionReceive:
			mainServiceReceiveChannels[op.Channel.Name] = true
		}
	}

	for _, service := range s.Services {
		if service.Name == mainService.Name {
			continue
		}

		isNeighbor := false

		// Check if this service sends to channels that main service receives from
		for _, op := range service.Operation {
			if op.Action == messageflow.ActionSend && mainServiceReceiveChannels[op.Channel.Name] {
				isNeighbor = true
				break
			}
		}

		// Check if this service receives from channels that main service sends to
		if !isNeighbor {
			for _, op := range service.Operation {
				if op.Action == messageflow.ActionReceive && mainServiceSendChannels[op.Channel.Name] {
					isNeighbor = true
					break
				}
			}
		}

		if isNeighbor && !neighborServiceMap[service.Name] {
			neighborServices = append(neighborServices, service)
			neighborServiceMap[service.Name] = true
		}
	}

	return serviceServicesPayload{
		MainService:      mainService,
		NeighborServices: neighborServices,
	}
}
