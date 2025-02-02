package exporter

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	ctx = context.Background()

	// Create a grpc Exporer
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test"),
			attribute.String("environment", "production"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create a trace provider with the exporter and a resource
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Register the trace provider with the global trace provider
	otel.SetTracerProvider(tp)
	return tp, nil
}

/*

func main() {
	ctx := context.Background()

	// TracerProvider 초기화
	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	// 프로그램 종료 시 모든 스팬이 플러시될 수 있도록 Shutdown 호출
	defer func() {
		// Shutdown은 남은 스팬을 플러시하고 Exporter를 종료합니다.
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	// "example-tracer"라는 이름의 트레이서를 가져옵니다.
	tracer := otel.Tracer("example-tracer")

	// "main-span" 스팬 시작 (속성도 추가)
	ctx, span := tracer.Start(ctx, "main-span",
		trace.WithAttributes(attribute.String("foo", "bar")),
	)
	// 작업 시뮬레이션 (100ms 대기)
	time.Sleep(100 * time.Millisecond)
	// 이벤트 추가 (예: 추가 정보를 남김)
	span.AddEvent("An event occurred", trace.WithAttributes(attribute.Int("event.attr", 123)))
	span.End() // 스팬 종료

	log.Println("Tracing example finished. Waiting for export...")
	// 스팬이 내보내질 시간을 확보하기 위해 잠시 대기 (실제 환경에서는 Shutdown 시 플러시됨)
	time.Sleep(2 * time.Second)
}
*/
