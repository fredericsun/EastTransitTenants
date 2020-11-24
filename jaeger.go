package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jaegertracing/jaeger/model"

	"github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"google.golang.org/grpc"
)

type status int

const (
	Before status = iota + 1
	After
)

type MySpan struct {
	TraceID       string
	SpanID        string
	OperationName string
	References    []map[string]string
	StartTime     time.Time
	Duration      time.Duration
	Process       *model.Process
	ProcessID     string
}

func ToMySpan(span model.Span) MySpan {
	mySpan := MySpan{
		TraceID:       span.TraceID.String(),
		SpanID:        span.SpanID.String(),
		OperationName: span.OperationName,
		References:    make([]map[string]string, 0),
		StartTime:     span.StartTime,
		Duration:      span.Duration,
		Process:       span.Process,
		ProcessID:     span.ProcessID,
	}
	for _, ref := range span.References {
		refMap := make(map[string]string)
		refMap["TraceID"] = ref.TraceID.String()
		refMap["SpanID"] = ref.SpanID.String()
		refMap["RefType"] = ref.RefType.String()
		mySpan.References = append(mySpan.References, refMap)
	}
	return mySpan
}

func (s status) GoString() string {
	switch s {
	case Before:
		return "before"
	default:
		return "after"
	}
}

// JaegerClient is a wrapper for grpc JaegerClient for jaeger query
type JaegerClient struct {
	cc *grpc.ClientConn
}

// NewJaegerClient creates a new JaegerClient
func NewJaegerClient(addr string) *JaegerClient {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &JaegerClient{cc: cc}
}

// QueryServices queries jaeger for all available services
func (c *JaegerClient) QueryServices() (*api_v2.GetServicesResponse, error) {
	client := api_v2.NewQueryServiceClient(c.cc)
	res, err := client.GetServices(context.Background(), &api_v2.GetServicesRequest{})
	if err != nil {
		return nil, err
	}

	d, err := res.Marshal()
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(filepath.Join("data", "services"), d, 0644); err != nil {
		return nil, err
	}

	return res, nil
}

// QueryOperations queries jaeger for all operations of a service
func (c *JaegerClient) QueryOperations(svc string) (*api_v2.GetOperationsResponse, error) {
	client := api_v2.NewQueryServiceClient(c.cc)
	return client.GetOperations(context.Background(), &api_v2.GetOperationsRequest{
		Service:  svc,
		SpanKind: "",
	})
}

// QueryTraces queries Jaeger for last 20 traces of a service's operation
func (c *JaegerClient) QueryTraces(svc, op string, since time.Time, depth int32) (map[string][]MySpan, error) {
	client := api_v2.NewQueryServiceClient(c.cc)
	stream, err := client.FindTraces(context.Background(), &api_v2.FindTracesRequest{
		Query: &api_v2.TraceQueryParameters{
			ServiceName:   svc,
			OperationName: op,
			StartTimeMin:  since,
			StartTimeMax:  time.Now(),
			SearchDepth:   depth,
		},
	})
	if err != nil {
		return nil, err
	}

	traces := make(map[string][]MySpan)

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		for _, s := range msg.GetSpans() {
			traceID := s.TraceID.String()
			mySpan := ToMySpan(s)
			traces[traceID] = append(traces[traceID], mySpan)
		}
	}

	return traces, nil
}

// QueryChunks queries jaeger for spans from inputted services since the inputted time
func (c *JaegerClient) QueryChunks(path string, status status, services []string, since time.Time) (map[string]*api_v2.SpansResponseChunk, error) {
	// Set data folder for saving chunks
	chunksDir := filepath.Join(path, status.GoString())
	if err := os.MkdirAll(chunksDir, 0755); err != nil {
		return nil, err
	}

	client := api_v2.NewQueryServiceClient(c.cc)
	result := make(map[string]*api_v2.SpansResponseChunk, 0)

	for _, svc := range services {
		if !strings.Contains(svc, ".default") {
			svc += ".default"
		}
		// Find all traces for this svc in the past hour with search depth 50
		res, err := client.FindTraces(context.Background(), &api_v2.FindTracesRequest{
			Query: &api_v2.TraceQueryParameters{
				ServiceName:  svc,
				StartTimeMin: since,
				StartTimeMax: time.Now(),
				SearchDepth:  30,
			},
		})
		if err != nil {
			return nil, err
		}

		// Populate spans
		var spans []model.Span

		for {
			c, err := res.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			spans = append(spans, c.GetSpans()...)
		}

		chunk := &api_v2.SpansResponseChunk{Spans: spans}

		// Write chunks to file
		if err := writeChunksToFile(chunk, filepath.Join(chunksDir, svc)); err != nil {
			return nil, err
		}

		// update map
		result[svc] = chunk
	}

	return result, nil
}

func writeChunksToFile(chunk *api_v2.SpansResponseChunk, path string) error {
	b, err := chunk.Marshal()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0644)
}

func writeTraceToFile(trace []model.Span, path string) error {
	return writeChunksToFile(&api_v2.SpansResponseChunk{
		Spans: trace,
	}, path)
}

func main() {
	// flag.Parse()

	// defer zLogger.Sync()

	// Create jaeger service client
	jc := NewJaegerClient("localhost:16686")

	// Quert all services
	// res, err := jc.QueryServices()

	// Query traces
	_, err := jc.QueryTraces("ts-ui-dashboard.default", "", time.Now().Add(time.Minute*-(10)), 0)

	if err != nil {
		fmt.Println(err)
		return
	}

	// For print purpose
	// parsed, err := json.Marshal(res)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(parsed))
}
