# Log

This package is an opinionated logger and logging http middleware for logging structured (JSON) logs.

## Usage

First you wrap your handlers using the logging middleware like so:

```
// Change this
http.HandleFunc("/test", TestHandleFunc)

// Into this
http.HandleFunc("/test", log.Middleware(TestHandleFunc)
```

Now you are logging the request/response and the request ID has been placed into the context for future use. So you're ready to log anywhere in your code (preferably where context is accessible so that you have request ID in the logs). Logging itself is pretty straightforward:

```
log.Info(ctx, "You can log like this")

log.Info(ctx, "Or", "like", "this")

logger := log.WithContext(ctx)
logger.Info("Or even like this")

log.Error(ctx, "There is also error level")
```

 For now there are only 2 levels you can log with: Info and Error. More may be added in hte future. You should probably also redirect all logs from the native log package into this logger by doing this in your main:

 ```
 import nativeLog "log"
 ...

// All logs from native log package will be sent to this logger with level "unknown"
nativeLog.SetOutput(log.GetWriter(log.LevelUnknown))
 ```

Last but not least, you may want to define custom serializing for your some types. You can either implement json.Marshaler on them or you can implement a "Data Serialization Converter" - basically transform data into a marshalable form. It is done via the `log.AddDataSerializationConverter` function like this:

```
log.AddDataSerializationConverter("error", func(i interface{}) (interface{}, error) {
    if err, ok := i.(error); ok {
        return "!error!:" + err.Error(), nil
    }

    // returning nil will not add the converted data to log
    return nil, nil
})
```

When logging an error with a serialization converter like this, a new entry will be added under the "Converted" key like so:

```
...
"Facts":[
  {
     "Type":"*json.SyntaxError",
     "Marshaled":{
        "Offset":9
     },
     "Printed":"\u0026json.SyntaxError{msg:\"invalid character 'd' looking for beginning of value\", Offset:9}",
     "Converted":{
        "error": "!error!:invalid character 'd' looking for beginning of value"
     }
  }
]
...
```

## Example logs

This example contains request, json unmarshaling error and response.

```
{
   "Time":"2018-05-28T18:31:19.98617+02:00",
   "Level":"info",
   "RequestID":"A1VPsEFcEB3fw0z7HQPFzo1TsIYKMFSd",
   "Message":"request received",
   "Source":{
      "FilePath":"/Users/dusan/Development/golang/src/github.com/DusanKasan/log/middleware.go",
      "LineNumber":82,
      "FunctionName":"github.com/DusanKasan/log.Middleware.func1"
   },
   "Facts":[
      {
         "Type":"string",
         "Marshaled":"request received",
         "Printed":"\"request received\"",
         "Converted":{

         }
      },
      {
         "Type":"github.com/DusanKasan/log.log.request",
         "Marshaled":{
            "Method":"POST",
            "Header":{
               "Accept":[
                  "*/*"
               ],
               "Content-Length":[
                  "24"
               ],
               "Content-Type":[
                  "application/x-www-form-urlencoded"
               ],
               "User-Agent":[
                  "curl/7.54.0"
               ]
            },
            "URL":{
               "Scheme":"",
               "Opaque":"",
               "User":null,
               "Host":"",
               "Path":"/request-one-time-password",
               "RawPath":"",
               "ForceQuery":false,
               "RawQuery":"",
               "Fragment":""
            },
            "Proto":"HTTP/1.1",
            "BodyPreview":"{\"emai\":dusan@kasan.sk\"}"
         },
         "Printed":"log.request{Method:\"POST\", Header:http.Header{\"Accept\":[]string{\"*/*\"}, \"Content-Length\":[]string{\"24\"}, \"Content-Type\":[]string{\"application/x-www-form-urlencoded\"}, \"User-Agent\":[]string{\"curl/7.54.0\"}}, URL:url.URL{Scheme:\"\", Opaque:\"\", User:(*url.Userinfo)(nil), Host:\"\", Path:\"/request-one-time-password\", RawPath:\"\", ForceQuery:false, RawQuery:\"\", Fragment:\"\"}, Proto:\"HTTP/1.1\", BodyPreview:\"{\\\"emai\\\":dusan@kasan.sk\\\"}\"}",
         "Converted":{

         }
      }
   ]
}{
   "Time":"2018-05-28T18:31:19.986842+02:00",
   "Level":"info",
   "RequestID":"A1VPsEFcEB3fw0z7HQPFzo1TsIYKMFSd",
   "Message":"invalid character 'd' looking for beginning of value",
   "Source":{
      "FilePath":"/Users/dusan/Development/golang/src/github.com/DusanKasan/log/log.go",
      "LineNumber":78,
      "FunctionName":"github.com/DusanKasan/log.(*contextLogger).Info"
   },
   "Facts":[
      {
         "Type":"*json.SyntaxError",
         "Marshaled":{
            "Offset":9
         },
         "Printed":"\u0026json.SyntaxError{msg:\"invalid character 'd' looking for beginning of value\", Offset:9}",
         "Converted":{

         }
      }
   ]
}{
   "Time":"2018-05-28T18:31:19.986991+02:00",
   "Level":"info",
   "RequestID":"A1VPsEFcEB3fw0z7HQPFzo1TsIYKMFSd",
   "Message":"response sent",
   "Source":{
      "FilePath":"/Users/dusan/Development/golang/src/github.com/DusanKasan/log/middleware.go",
      "LineNumber":87,
      "FunctionName":"github.com/DusanKasan/log.Middleware.func1"
   },
   "Facts":[
      {
         "Type":"string",
         "Marshaled":"response sent",
         "Printed":"\"response sent\"",
         "Converted":{

         }
      },
      {
         "Type":"github.com/DusanKasan/log.response",
         "Marshaled":{
            "Header":{
               "Request-Id":[
                  "A1VPsEFcEB3fw0z7HQPFzo1TsIYKMFSd"
               ]
            },
            "Body":"{\"Message\":\"invalid input\"}",
            "Status":400
         },
         "Printed":"log.response{Header:http.Header{\"Request-Id\":[]string{\"A1VPsEFcEB3fw0z7HQPFzo1TsIYKMFSd\"}}, Body:\"{\\\"Message\\\":\\\"invalid input\\\"}\", Status:400}",
         "Converted":{

         }
      }
   ]
}
```

## TODO

- tests
- document