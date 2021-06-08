## go-gin-tracer

Хэлпер для работы с jaeger для rest-запросов
### Необходимые конфигурации

TRACE_HEADER=uber-trace-id // устанавливается, если нужен кастомный заголовок

JAEGER_SERVICE_NAME=service_control
JAEGER_AGENT_HOST=localhost
JAEGER_AGENT_PORT=6831
JAEGER_SAMPLER_TYPE=const
JAEGER_REPORTER_LOG_SPANS=true
JAEGER_SAMPLER_PARAM=1

Можно установить остальные параметры из github.com/uber/jaeger-client-go/config

### Методы

* SetJaegerTracer Получение объекта трассировщика. Пример использования:

```  
func main() {
 tracer, closer, err := gintraceser.SetJaegerTracer(env.TraceHeader)
 if err == nil {
   opentracing.SetGlobalTracer(tracer)
   defer closer.Close()
 }
}
  ```

* OpenTracingMiddleware
  Middleware для установки корневого span для входящих запросов. Пример использования:

```  
func (p Router) Router() *gin.Engine {
  r := gin.New()
  r.Use(gintraceser.OpenTracingMiddleware())
}
  ```

* AddTraceToRequest
 Добавление заголовка к исходящему запросу с текущим span. Пример использования:

```  
client := http.Client{}
req, _ := http.NewRequest("GET", "google.com", nil)
AddTraceToRequest(c,req)
resp,_ := client.Do(req)
  ```
