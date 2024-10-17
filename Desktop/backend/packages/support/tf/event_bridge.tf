### Start Schema Registry ###
resource "aws_schemas_registry" "support_schema_registry" {
  name        = "support_schema_registry"
  description = "Schema registry for events that come from the support service."
}

resource "aws_schemas_schema" "message_sent_event" {
  name          = "message_sent_event"
  registry_name = aws_schemas_registry.support_schema_registry.name
  type          = "OpenApi3"
  description   = "Event fired when a message is sent from the support service"

  content = jsonencode({
    "openapi" : "3.0.0",
    "info" : {
      "version" : "1.0.0",
      "title" : "Event"
    },
    "paths" : {},
    "components" : {
      "schemas" : {
        "Event" : {
          "type" : "object",
          "required" : ["ReceiverId", "SenderType", "SenderId", "SessionId", "SentAt", "MessageId"],
          "properties" : {
            "ReceiverId" : {
              "type" : "string",
            },
            "SenderId" : {
              "type" : "string"
            },
            "SenderType" : {
              "type" : "string",
              "enum" : ["internal", "loop"]
            },
            "SentAt" : {
              "type" : "number",
              "format" : "epoch"
            },
            "SessionId" : {
              "type" : "string"
            },
            "MessageId" : {
              "type" : "string"
            }
          }
        }
      }
    }
  })
}
### End Schema Registry ###